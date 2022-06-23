package main

import (
	"Bank-system/cmd/bank/app"
	"context"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort        = "9999"
	defaultHost        = "0.0.0.0"
	defaultPostgresDSN = "postgres://app:pass@localhost:5200/emilDB"
	defaultMongoDSN    = "mongodb://app:pass@localhost:27017/" + defaultMongoDB
	defaultMongoDB     = "emil"
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

	postgresDsn, ok := os.LookupEnv("APP_POSTGRES_DSN")
	if !ok {
		postgresDsn = defaultPostgresDSN
	}

	mongoDsn, ok := os.LookupEnv("APP_MONGO_DSN")
	if !ok {
		mongoDsn = defaultMongoDSN
	}

	mongoDB, ok := os.LookupEnv("APP_MONGO_DB")
	if !ok {
		mongoDB = defaultMongoDB
	}

	if err := execute(net.JoinHostPort(host, port), postgresDsn, mongoDsn, mongoDB); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, dsn string, mongoDsn string, mongoDBname string) (err error) {
	ctx := context.Background()
	pstgrCtx, cancelPstgr := context.WithCancel(ctx)
	mongoCtx, cancelMongo := context.WithCancel(ctx)

	pool, err := pgxpool.Connect(pstgrCtx, dsn)
	if err != nil {
		log.Print(err)
		cancelPstgr()
		return err
	}

	client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(mongoDsn))
	if err != nil {
		log.Print(err)
		cancelMongo()
		return err
	}

	mongoDatabase := client.Database(mongoDBname)

	handlerStorage := app.NewHandler(pool, mongoDatabase)

	mux := chi.NewRouter()
	application := app.NewServer(mux, handlerStorage)
	application.Init()

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
