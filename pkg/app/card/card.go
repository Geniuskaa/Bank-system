package card

import (
	"context"
	"fmt"
	"github.com/Geniuskaa/Bank-system/pkg/app/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
)

type Card struct {
	Id      int64
	Number  string
	Issuer  string
	Type    dto.CardType
	Balance int64
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
	return fmt.Sprintf("db serviceError: %s", e.Err.Error())
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) All(ctx context.Context, login *dto.UserLoginDTO, role string) ([]*Card, error) {

	if strings.EqualFold(role, "CLIENT") {
		rows, err := s.pool.Query(ctx, `
		SELECT cards.id, number, balance FROM cards join clients ON clients.id=cards.owner_id
	join users u on clients.id = u.client_id where u.login = $1 LIMIT 50`, login.Login)

		if err != nil {
			if err != pgx.ErrNoRows {
				return nil, NewDbError(err)
			}
			return nil, nil
		}
		defer rows.Close()

		var result []*Card
		for rows.Next() {
			card := &Card{}
			err = rows.Scan(&card.Id, &card.Number, &card.Balance)
			if err != nil {
				return nil, NewDbError(err)
			}
			result = append(result, card)
		}
		err = rows.Err()
		if err != nil {
			return nil, NewDbError(err)
		}
		return result, nil

	} else {

		rows, err := s.pool.Query(ctx, `
		SELECT id, number, balance FROM cards
		LIMIT 50`)
		if err != nil {
			if err != pgx.ErrNoRows {
				return nil, NewDbError(err)
			}
			return nil, nil
		}
		defer rows.Close()

		var result []*Card
		for rows.Next() {
			card := &Card{}
			err = rows.Scan(&card.Id, &card.Number, &card.Balance)
			if err != nil {
				return nil, NewDbError(err)
			}
			result = append(result, card)
		}
		err = rows.Err()
		if err != nil {
			return nil, NewDbError(err)
		}
		return result, nil
	}
}

func (s *Service) FindUserCardsById(c context.Context, userId int64) ([]*Card, error) {
	rows, err := s.pool.Query(c, `SELECT id, number, balance, issuer, type from cards where owner_id = $1 limit 10`, userId)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, NewDbError(err)
		}
		return nil, nil
	}
	defer rows.Close()

	var cards []*Card
	for rows.Next() {
		card := &Card{}
		err = rows.Scan(&card.Id, &card.Number, &card.Balance, &card.Issuer, &card.Type)
		if err != nil {
			return nil, NewDbError(err)
		}
		cards = append(cards, card)
	}
	err = rows.Err()
	if err != nil {
		return nil, NewDbError(err)
	}
	return cards, nil
}

func (s *Service) FindCardTransactionsByCardId(ctx context.Context, id int64) ([]*dto.TransactionDto, error) {
	rows, err := s.pool.Query(ctx, `SELECT amount, mcc, status, date from transactions where sender_id = $1 limit 10`, id)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, NewDbError(err)
		}
		log.Println(err)
		return nil, nil
	}
	defer rows.Close()

	var transactions []*dto.TransactionDto
	for rows.Next() {
		transaction := &dto.TransactionDto{}
		err = rows.Scan(&transaction.Amount, &transaction.MCC, &transaction.Status, &transaction.Date)
		if err != nil {
			return nil, NewDbError(err)
		}
		transactions = append(transactions, transaction)
	}
	err = rows.Err()
	if err != nil {
		return nil, NewDbError(err)
	}
	return transactions, nil
}
