package user

import (
	"Bank-system/cmd/bank/app/dto"
	"Bank-system/pkg/card"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Role string

var ErrRegistrationFailed = errors.New("Registration of new user is failed.")
var ErrPasswordWrong = errors.New("Your password is incorrect.")

const (
	CLIENT    Role = "CLIENT"
	OPERATOR  Role = "OPERATOR"
	VIPCLIENT Role = "VIPCLIENT"
	ADMIN     Role = "ADMIN"
)

type User struct {
	Id         int64
	FirstName  string
	SecondName string
	Login      string
	Password   []byte
	Email      string
	Cards      []*card.Card
}

type Service struct {
	pool *pgxpool.Pool
}

type DbError struct {
	Err error
}

func NewDbError(err error) *DbError {
	return &DbError{Err: err}
}

func (e DbError) Error() string {
	return fmt.Sprintf("db error: %s", e.Err.Error())
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) AuthorizeUser(ctx context.Context, userDto *dto.UserRegAndLognDTO) (string, error) {

	row := s.pool.QueryRow(ctx, `SELECT id,login,password from users where login = $1`, userDto.Login)
	userID := new(int64)
	realUser := dto.UserRegAndLognDTO{}
	err := row.Scan(userID, &realUser.Login, &realUser.Password)
	if err != nil {
		log.Println(ErrPasswordWrong)
		return "", ErrPasswordWrong
	}

	if userDto.Login == realUser.Login {
		err = bcrypt.CompareHashAndPassword([]byte(realUser.Password), []byte(userDto.Password))
		if err != nil {
			log.Println(err)
			return "", err
		}
		log.Println("User is athorized")
	}

	token, err := uuid.NewUUID()
	_, err = s.pool.Exec(ctx, `INSERT INTO user_to_token (user_id, token) VALUES ($1, $2)`, userID, token.String())
	if err != nil {
		log.Println(err)
		return "", err
	}

	return token.String(), nil
}

func (s *Service) RegisterUser(context context.Context, dto *dto.UserRegAndLognDTO) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 10)
	if err != nil {
		log.Println(err)
		return -2, err
	}
	log.Println(hashedPassword)

	row := s.pool.QueryRow(context, `INSERT INTO users(login, password, role) VALUES($1, $2, $3) 
                                         ON CONFLICT DO NOTHING RETURNING id`,
		dto.Login, hashedPassword, CLIENT)

	id := new(int64)
	err = row.Scan(id)
	if err != nil {
		log.Println(ErrRegistrationFailed)
		return -2, ErrRegistrationFailed
	}

	return *id, nil
}
