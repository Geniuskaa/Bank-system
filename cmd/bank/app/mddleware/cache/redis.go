package cache

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

const cacheTimeout time.Duration = 55 * time.Millisecond

type Service struct {
	pool *redis.Pool
}

func NewService(cachePool *redis.Pool) *Service {
	return &Service{pool: cachePool}
}

func (s *Service) FromCache(ctx context.Context, key string) ([]byte, error) {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	reply, err := redis.DoWithTimeout(conn, cacheTimeout, "GET", key)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	value, err := redis.Bytes(reply, err)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return value, err
}

func (s *Service) ToCache(ctx context.Context, key string, value []byte) error {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	_, err = redis.DoWithTimeout(conn, cacheTimeout, "SET", key, value)
	if err != nil {
		log.Print(err)
	}
	return err
}
