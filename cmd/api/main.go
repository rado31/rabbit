package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"  binding:"required"`
	Age       int    `json:"age"        binding:"required,min=1,max=150"`
	Email     string `json:"email"      binding:"required,email"`
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

type rpcClient struct {
	ch      *amqp.Channel
	replyQ  string
	pending map[string]chan []byte
	mu      sync.Mutex
}

func newRPCClient(conn *amqp.Connection) (*rpcClient, error) {
	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	for _, q := range []string{queueCreate, queueGet} {
		if _, err := ch.QueueDeclare(q, true, false, false, false, nil); err != nil {
			return nil, err
		}
	}

	q, err := ch.QueueDeclare("", false, true, true, false, nil)

	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(q.Name, "", true, true, false, false, nil)

	if err != nil {
		return nil, err
	}

	c := &rpcClient{
		ch:      ch,
		replyQ:  q.Name,
		pending: make(map[string]chan []byte),
	}

	go func() {
		for msg := range msgs {
			c.mu.Lock()

			ch, ok := c.pending[msg.CorrelationId]

			if ok {
				delete(c.pending, msg.CorrelationId)
			}

			c.mu.Unlock()

			if ok {
				ch <- msg.Body
			}
		}
	}()

	return c, nil
}

func (c *rpcClient) call(queue string, payload any) (*Reply, error) {
	body, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	corrID := corrID()
	replyCh := make(chan []byte, 1)

	c.mu.Lock()
	c.pending[corrID] = replyCh
	c.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = c.ch.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: corrID,
		ReplyTo:       c.replyQ,
		Body:          body,
	})

	if err != nil {
		c.mu.Lock()

		delete(c.pending, corrID)

		c.mu.Unlock()

		return nil, err
	}

	select {
	case raw := <-replyCh:
		var reply Reply

		return &reply, json.Unmarshal(raw, &reply)
	case <-ctx.Done():
		c.mu.Lock()

		delete(c.pending, corrID)

		c.mu.Unlock()

		return nil, fmt.Errorf("timeout waiting for reply from %s", queue)
	}
}

func corrID() string {
	b := make([]byte, 16)

	rand.Read(b)

	return fmt.Sprintf("%x", b)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func main() {
	conn, err := amqp.Dial(getenv("AMQP_URL", "amqp://guest:guest@localhost:5672/"))

	if err != nil {
		log.Fatal("rabbitmq:", err)
	}

	defer conn.Close()

	rpc, err := newRPCClient(conn)

	if err != nil {
		log.Fatal("rpc client:", err)
	}

	r := gin.Default()

	r.POST("/clients", func(c *gin.Context) {
		var req CreateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		reply, err := rpc.call(queueCreate, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if reply.Error != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": reply.Error})
			return
		}

		c.JSON(http.StatusCreated, reply.Client)
	})

	r.GET("/clients/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
			return
		}

		reply, err := rpc.call(queueGet, map[string]int{"id": id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if reply.Error != "" {
			c.JSON(http.StatusNotFound, gin.H{"error": reply.Error})
			return
		}

		c.JSON(http.StatusOK, reply.Client)
	})

	addr := getenv("ADDR", ":8080")
	log.Println("listening on", addr)
	log.Fatal(r.Run(addr))
}
