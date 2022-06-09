package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Client struct {
	Id       int64
	FullName string
	Birthday time.Time
}

func main() {
	dsn := "postgres://app:pass@localhost:5200/emilDB"
	ctx := context.Background() // про контексты почитай,посмотри подробнее
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Println(err)
		return
	}
	defer pool.Close()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Release()

	tag, err := conn.Exec(ctx, `
INSERT INTO cards (number, balance, issuer, holder, owner_id, status) VALUES ($1, $2, $3, $4, $5, $6)`,
		"0420 3020", 1500_00, "Visa", "PETR IVANOV", 1, "ACTIVE")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(tag.RowsAffected())

	//Если нам нужно сделать только один запрос к бд, тогда создадим метод для простоты работы
	//	tag, err = pool.Exec(ctx, `
	//INSERT INTO cards (number, balance, issuer, holder, owner_id, status) VALUES ($1, $2, $3, $4, $5, $6)`,
	//		"0420 3020", 1500_00, "VISA", "PETR IVANOV", 1, "ACTIVE")
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//

	client := &Client{}
	err = conn.QueryRow(ctx, `
		SELECT id, full_name, birthday FROM clients WHERE id = $1
	`, 1).Scan(&client.Id, &client.FullName, &client.Birthday)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(client)

	client = &Client{}
	err = conn.QueryRow(ctx, `
		SELECT id, full_name, birthday FROM clients WHERE id = $1
	`, 1).Scan(&client.Id, &client.FullName, &client.Birthday)
	if err != nil {
		if err != pgx.ErrNoRows {
			log.Println(err)
			return
		}
	}
	log.Println(client)

	clients := make([]*Client, 0)
	rows, err := conn.Query(ctx, `
		SELECT id, full_name, birthday FROM clients WHERE status = $1
	`, "ACTIVE")
	if err != nil {
		// TODO: отдельная обработка для ErrNoRows
		log.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		client := &Client{}
		err = rows.Scan(&client.Id, &client.FullName, &client.Birthday)
		if err != nil {
			log.Println(err)
			return
		}
		clients = append(clients, client)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return
	}

	var id int64
	err = conn.QueryRow(ctx, `
		INSERT INTO cards(number, balance, issuer, holder, owner_id, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, "**** 0002", 100_00, "Visa", "PETR IVANOV", 2, "ACTIVE").Scan(&id)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(id)
}
