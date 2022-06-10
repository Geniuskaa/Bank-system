package main

import (
	"Bank-system/cmd/bank/app"
	"Bank-system/pkg/card"
	"Bank-system/pkg/user"
	"context"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@localhost:5200/emilDB"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = defaultDSN
	}

	if err := execute(net.JoinHostPort(host, port), dsn); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, dsn string) (err error) {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Print(err)
		return err
	}

	userSvc := user.NewService(pool)
	cardSvc := card.NewService(pool)

	mux := chi.NewRouter()
	application := app.NewServer(cardSvc, mux, pool, userSvc)
	application.Init()

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
