package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

type Client struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

type Reply struct {
	Client *Client `json:"client,omitempty"`
	Error  string  `json:"error,omitempty"`
}

const (
	queueCreate = "clients.create"
	queueGet    = "clients.get"
)

func createClient(ctx context.Context, pool *pgxpool.Pool, req CreateRequest) (*Client, error) {
	var c Client

	err := pool.QueryRow(ctx,
		`INSERT INTO clients (first_name, last_name, age, email)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, first_name, last_name, age, email`,
		req.FirstName, req.LastName, req.Age, req.Email,
	).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Age, &c.Email)

	return &c, err
}

func getClient(ctx context.Context, pool *pgxpool.Pool, id int) (*Client, error) {
	var c Client

	err := pool.QueryRow(ctx,
		`SELECT id, first_name, last_name, age, email FROM clients WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Age, &c.Email)

	return &c, err
}

func reply(ch *amqp.Channel, d amqp.Delivery, client *Client, errMsg string) {
	r := Reply{Client: client, Error: errMsg}

	body, _ := json.Marshal(r)

	ch.PublishWithContext(context.Background(), "", d.ReplyTo, false, false, amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: d.CorrelationId,
		Body:          body,
	})

	d.Ack(false)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, getenv("PG_URL", "postgres://postgres:postgres@localhost:5432/clients_db"))

	if err != nil {
		log.Fatal("postgres:", err)
	}

	defer pool.Close()

	conn, err := amqp.Dial(getenv("AMQP_URL", "amqp://guest:guest@localhost:5672/"))

	if err != nil {
		log.Fatal("rabbitmq:", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal("channel:", err)
	}

	defer ch.Close()

	for _, q := range []string{queueCreate, queueGet} {
		if _, err := ch.QueueDeclare(q, true, false, false, false, nil); err != nil {
			log.Fatal("declare queue:", err)
		}
	}

	ch.Qos(1, 0, false)

	createMsgs, _ := ch.Consume(queueCreate, "", false, false, false, false, nil)
	getMsgs, _ := ch.Consume(queueGet, "", false, false, false, false, nil)

	log.Println("storage listening...")

	go func() {
		for d := range createMsgs {
			var req CreateRequest

			if err := json.Unmarshal(d.Body, &req); err != nil {
				reply(ch, d, nil, "invalid payload")
				continue
			}

			client, err := createClient(ctx, pool, req)

			if err != nil {
				reply(ch, d, nil, err.Error())
				continue
			}

			log.Printf("created client id=%d", client.ID)
			reply(ch, d, client, "")
		}
	}()

	go func() {
		for d := range getMsgs {
			var payload struct {
				ID int `json:"id"`
			}

			if err := json.Unmarshal(d.Body, &payload); err != nil {
				reply(ch, d, nil, "invalid payload")
				continue
			}

			client, err := getClient(ctx, pool, payload.ID)

			if err != nil {
				reply(ch, d, nil, "client not found")
				continue
			}

			log.Printf("fetched client id=%d", client.ID)
			reply(ch, d, client, "")
		}
	}()

	select {}
}
