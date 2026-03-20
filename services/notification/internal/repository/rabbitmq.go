package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ClientCreatedEvent struct {
	ID      int32  `json:"id"`
	Surname string `json:"surname"`
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Email   string `json:"email"`
}

type EventHandler func(e ClientCreatedEvent)

type RabbitMQRepository struct {
	ch *amqp.Channel
}

func New(ch *amqp.Channel) (*RabbitMQRepository, error) {
	_, err := ch.QueueDeclare("clients", true, false, false, false, nil)

	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &RabbitMQRepository{ch: ch}, nil
}

func (r *RabbitMQRepository) Consume(ctx context.Context, handler EventHandler) error {
	msgs, err := r.ch.Consume("clients", "", true, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var e ClientCreatedEvent

				if err := json.Unmarshal(msg.Body, &e); err != nil {
					log.Printf("notification: unmarshal event: %v", err)
					continue
				}

				handler(e)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
