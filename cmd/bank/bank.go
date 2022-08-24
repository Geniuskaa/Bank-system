package main

import (
	"context"
	"github.com/Geniuskaa/Bank-system/pkg/app/handlr"
	"github.com/Geniuskaa/Bank-system/pkg/app/server"
	"github.com/go-chi/chi"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort        = "9998"
	defaultHost        = "0.0.0.0"
	defaultPostgresDSN = "postgres://app:pass@localhost:5200/emilDB"
	defaultMongoDSN    = "mongodb://app:pass@localhost:27017/" + defaultMongoDB
	defaultMongoDB     = "emil"
	defaultCacheDSN    = "redis://localhost:6379/0"
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

	cacheDSN, ok := os.LookupEnv("APP_CACHE_DSN")
	if !ok {
		cacheDSN = defaultCacheDSN
	}

	if err := execute(net.JoinHostPort(host, port), postgresDsn, mongoDsn, mongoDB, cacheDSN); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, dsn string, mongoDsn string, mongoDBname string, cacheDSN string) (err error) {

	ctx := context.Background()
	pstgrCtx, cancelPstgr := context.WithCancel(ctx)
	mongoCtx, cancelMongo := context.WithCancel(ctx)

	defer func() {
		cancelPstgr()
		cancelMongo()
	}()

	//Создание пула подключений к PostgreSQL
	pool, err := pgxpool.Connect(pstgrCtx, dsn)
	if err != nil {
		log.Print(err)
		cancelPstgr()
		return err
	}

	//Создание подключения к MongoDB
	client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(mongoDsn))
	if err != nil {
		log.Print(err)
		cancelMongo()
		return err
	}

	mongoDatabase := client.Database(mongoDBname)

	//Создание пула подключений к Redis
	cache := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialURL(cacheDSN)
		},
	}

	handlerStorage := handlr.NewHandler(pool, mongoDatabase, cache)

	mux := chi.NewRouter()
	application := server.NewServer(mux, handlerStorage)
	application.Init()

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
