package app

import (
	"Bank-system/cmd/bank/app/dto"
	"Bank-system/pkg/card"
	"Bank-system/pkg/user"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var ErrNoIdentifier = errors.New("no identifier")
var ErrNoAuthentication = errors.New("no authentication")
var ErrContextEmpty = errors.New("Context is empty!")

type IdentifierFunc func(ctx context.Context) (*string, error)
type UserDetails func(ctx context.Context, id *string) (interface{}, error)

var authenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

type Server struct {
	cardSvc *card.Service
	userSvc *user.Service
	mux     *chi.Mux
	pool    *pgxpool.Pool
}

func NewServer(cardSvc *card.Service, mux *chi.Mux, pool *pgxpool.Pool, userSvc *user.Service) *Server {
	return &Server{cardSvc: cardSvc, mux: mux, pool: pool, userSvc: userSvc}
}

func (s *Server) Init() {
	//s.mux.Use(s.recoverer) // Добавление middleWare который ловит панику и обрабатывает её
	s.mux.With(s.Authorization).Post("/cards", s.getAllCards)

	s.mux.Get("/getUserCards", s.getUserCards)
	s.mux.Get("/getCardTransactions", s.getCardTransactions)

	s.mux.Post("/api/users", s.registerUser)
	s.mux.Post("/token", s.tokenGenerator)

	s.mux.Get("/returnPanic", s.returnPanic)

}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func marshalDtos(anyObject interface{}) ([]byte, error) {
	respBody, err := json.Marshal(anyObject)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return respBody, nil
}

func internalServerErrorJson(writer http.ResponseWriter) (w http.ResponseWriter) {
	w = writer
	w.WriteHeader(http.StatusInternalServerError)

	intErr := dto.Error{Error: "err.internal_server_error"}
	respBody, err := json.Marshal(intErr)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		return
	}
	return
}

func userAlreadyUsedErrorJson(writer http.ResponseWriter) (w http.ResponseWriter) {
	w = writer
	w.WriteHeader(http.StatusConflict)

	intErr := dto.Error{Error: "err.username_already_used"}
	respBody, err := json.Marshal(intErr)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		return
	}
	return
}

func (s *Server) returnPanic(w http.ResponseWriter, r *http.Request) {
	panic("I will be caught by middleware!")
}

func (s *Server) Authorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("new request: %s %s", request.Method, request.URL.Path)

		token := request.Header.Get("Authorization")

		row := s.pool.QueryRow(request.Context(), `SELECT role from user_to_token join users ON
   	user_to_token.user_id = users.id where token = $1`, token)
		if row == nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		var userRole user.Role
		err := row.Scan(&userRole)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			handler.ServeHTTP(writer, request)
			return
		}

		ctx := context.WithValue(request.Context(), authenticationContextKey, userRole)
		request = request.WithContext(ctx)

		handler.ServeHTTP(writer, request)
	})
}

func (s *Server) recoverer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("new request: %s %s", request.Method, request.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				log.Println("panic occurred:", err)
			}
		}()
		handler.ServeHTTP(writer, request)
	})
}

func (s *Server) getUserCards(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	log.Println(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cards, err := s.cardSvc.FindUserCardsById(r.Context(), int64(id)) //FindUserCardsById(r.Context(), int64(id))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dtos := make([]*dto.CardDTO, len(cards))
	for i, c := range cards {

		if c == nil {
			continue
		}

		dtos[i] = &dto.CardDTO{
			Id:      c.Id,
			Number:  c.Number,
			Issuer:  c.Issuer,
			Type:    c.Type,
			Balance: c.Balance,
		}
	}

	respBody, err := marshalDtos(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	// по умолчанию статус 200 Ok
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getCardTransactions(writer http.ResponseWriter, request *http.Request) {
	cardID, err := strconv.Atoi(request.URL.Query().Get("cardID"))
	if err != nil {
		log.Println(err)
		return
	}

	transactions, err := s.cardSvc.FindCardTransactionsByCardId(request.Context(), int64(cardID))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	respondBody, err := marshalDtos(transactions)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	_, err = writer.Write(respondBody)
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) registerUser(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer = internalServerErrorJson(writer)
		return
	}

	err = request.Body.Close()
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	userDTO := &dto.UserRegAndLognDTO{}
	err = json.Unmarshal(body, userDTO)
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	id, err := s.userSvc.RegisterUser(request.Context(), userDTO)
	if err != nil {
		writer = userAlreadyUsedErrorJson(writer)
		return
	}

	responseBody, err := marshalDtos(id)
	if err != nil {
		writer = internalServerErrorJson(writer)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(responseBody)

}

func (s *Server) tokenGenerator(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer = internalServerErrorJson(writer)
		return
	}

	err = request.Body.Close()
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	userDTO := &dto.UserRegAndLognDTO{}
	err = json.Unmarshal(body, userDTO)
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	token, err := s.userSvc.AuthorizeUser(request.Context(), userDTO)
	if err != nil {
		if err == user.ErrPasswordWrong {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer = internalServerErrorJson(writer)
		return
	}

	responsBody, err := marshalDtos(dto.Token{Token: token})
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(responsBody)

}

func (s *Server) getAllCards(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	login := &dto.UserLoginDTO{}
	err = json.Unmarshal(body, login)
	if err != nil {
		writer = internalServerErrorJson(writer)
		return
	}

	role := request.Context().Value(authenticationContextKey).(user.Role)

	cards, err := s.cardSvc.All(request.Context(), login, string(role))
	if err != nil {
		return
	}

	responseBody, err := marshalDtos(cards)
	if err != nil {
		log.Println(err)
		writer = internalServerErrorJson(writer)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.Write(responseBody)
}
