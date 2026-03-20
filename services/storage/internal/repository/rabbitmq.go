package repository

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rado31/rabbit/storage/internal/model"
)

type RabbitMQRepository struct {
	ch *amqp.Channel
}

func NewRabbitMQRepository(ch *amqp.Channel) (*RabbitMQRepository, error) {
	_, err := ch.QueueDeclare("clients", true, false, false, false, nil)

	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &RabbitMQRepository{ch: ch}, nil
}

func (r *RabbitMQRepository) PublishClientCreated(ctx context.Context, c *model.Client) error {
	body, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return r.ch.PublishWithContext(ctx,
		"",
		"clients",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
