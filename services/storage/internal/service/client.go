package service

import (
	"context"
	"fmt"

	"github.com/rado31/rabbit/storage/internal/model"
	"github.com/rado31/rabbit/storage/internal/repository"
)

type ClientService struct {
	db *repository.PostgresRepository
	mq *repository.RabbitMQRepository
}

func New(db *repository.PostgresRepository, mq *repository.RabbitMQRepository) *ClientService {
	return &ClientService{db: db, mq: mq}
}

func (s *ClientService) Create(ctx context.Context, c model.Client) (*model.Client, error) {
	tx, err := s.db.BeginTx(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	saved, err := s.db.SaveTx(ctx, tx, c)

	if err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	if err := s.mq.PublishClientCreated(ctx, saved); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return saved, nil
}

func (s *ClientService) GetByID(ctx context.Context, id int32) (*model.Client, error) {
	return s.db.FindByID(ctx, id)
}
