package server

import (
	"Bank-system/cmd/bank/app/handler"
	"Bank-system/cmd/bank/app/mddleware"
	"Bank-system/cmd/bank/app/mddleware/cache"
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gomodule/redigo/redis"
	"net/http"
)

var ErrNoIdentifier = errors.New("no identifier")
var ErrNoAuthentication = errors.New("no authentication")
var ErrContextEmpty = errors.New("Context is empty!")

type IdentifierFunc func(ctx context.Context) (*string, error)
type UserDetails func(ctx context.Context, id *string) (interface{}, error)

var AuthenticationContextKey = &contextKey{"authentication context"}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

type Server struct {
	mux            *chi.Mux
	handlerStorage *handler.Handler
}

func NewServer(mux *chi.Mux, handlerStorage *handler.Handler) *Server {
	return &Server{mux: mux, handlerStorage: handlerStorage}
}

func (s *Server) Init() {
	cacheMd := cache.Cache(func(ctx context.Context, path string) ([]byte, error) {
		value, err := s.handlerStorage.CacheSvc.FromCache(ctx, path)
		if err != nil && errors.Is(err, redis.ErrNil) {
			return nil, cache.ErrNotInCache
		}
		return value, err
	}, func(ctx context.Context, path string, data []byte) error {
		return s.handlerStorage.CacheSvc.ToCache(context.Background(), path, data)
	})

	s.mux.Use(mddleware.Recoverer) // Добавление middleWare который ловит панику и обрабатывает её
	//s.mux.Post("/cards", s.handlerStorage.GetAllCards)

	//Запросы с использование MongoDB
	s.mux.With(middleware.Logger).Get("/orders", s.handlerStorage.FindAll)
	s.mux.With(middleware.Logger).Get("/orders/{id:[0-9a-f]+}", s.handlerStorage.FindByID)
	s.mux.With(middleware.Logger).Get("/orders/search", s.handlerStorage.Search)
	s.mux.With(middleware.Logger).Post("/orders", s.handlerStorage.Save)

	//Кешируемые запросы
	s.mux.With(middleware.Logger, cacheMd).Get("/cached/films", s.handlerStorage.All)
	s.mux.With(middleware.Logger, cacheMd).Get("/cached/films/{id}", s.handlerStorage.ByID)
	s.mux.With(middleware.Logger).Post("/cached/films/upload", s.handlerStorage.Upload)

	s.mux.Get("/getUserCards", s.handlerStorage.GetUserCards)
	s.mux.Get("/getCardTransactions", s.handlerStorage.GetCardTransactions)

	s.mux.Post("/api/users", s.handlerStorage.RegisterUser)
	s.mux.Post("/token", s.handlerStorage.TokenGenerator)

	s.mux.Get("/returnPanic", s.handlerStorage.ReturnPanic)

	s.mux.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
