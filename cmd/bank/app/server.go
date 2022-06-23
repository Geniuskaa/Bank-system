package app

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
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
	mux            *chi.Mux
	handlerStorage *Handler
}

func NewServer(mux *chi.Mux, handlerStorage *Handler) *Server {
	return &Server{mux: mux, handlerStorage: handlerStorage}
}

func (s *Server) Init() {
	s.mux.Use(s.Recoverer) // Добавление middleWare который ловит панику и обрабатывает её
	s.mux.Post("/cards", s.handlerStorage.getAllCards)

	s.mux.With(middleware.Logger).Get("/orders}", s.handlerStorage.FindAll)
	s.mux.With(middleware.Logger).Get("/orders/{id:[0-9a-f]+}", s.handlerStorage.FindByID)
	s.mux.With(middleware.Logger).Get("/orders/search", s.handlerStorage.Search)
	s.mux.With(middleware.Logger).Post("/orders", s.handlerStorage.Save)

	s.mux.Get("/getUserCards", s.handlerStorage.getUserCards)
	s.mux.Get("/getCardTransactions", s.handlerStorage.getCardTransactions)

	s.mux.Post("/api/users", s.handlerStorage.registerUser)
	s.mux.Post("/token", s.handlerStorage.tokenGenerator)

	s.mux.Get("/returnPanic", s.handlerStorage.returnPanic)

}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
