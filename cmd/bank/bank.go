package main

import (
	"Bank-system/cmd/bank/app"
	"Bank-system/pkg/card"
	"context"
	mux2 "github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	pgx := newPool()
	defer pgx.Close()
	if pgx == nil {
		os.Exit(1)
	}

	if err := execute(net.JoinHostPort(host, port), pgx); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, pool *pgxpool.Pool) (err error) {
	cardSvc := card.NewService(pool)
	mux := mux2.NewRouter()
	application := app.NewServer(cardSvc, mux)
	application.Init()

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}

func newPool() *pgxpool.Pool {
	dsn := "postgres://app:pass@localhost:5200/emilDB"
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Println(err)
		return nil
	}

	return pool
}
