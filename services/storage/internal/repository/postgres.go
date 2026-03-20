package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rado31/rabbit/storage/internal/model"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.pool.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	return tx, nil
}

func (r *PostgresRepository) SaveTx(
	ctx context.Context,
	tx pgx.Tx,
	c model.Client,
) (*model.Client, error) {
	var saved model.Client

	err := tx.QueryRow(ctx,
		`INSERT INTO clients (surname, name, age, email)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, surname, name, age, email`,
		c.Surname, c.Name, c.Age, c.Email,
	).Scan(&saved.ID, &saved.Surname, &saved.Name, &saved.Age, &saved.Email)

	if err != nil {
		return nil, fmt.Errorf("insert client: %w", err)
	}

	return &saved, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id int32) (*model.Client, error) {
	var c model.Client

	err := r.pool.QueryRow(ctx,
		`SELECT id, surname, name, age, email FROM clients WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.Surname, &c.Name, &c.Age, &c.Email)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("client %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("find client: %w", err)
	}

	return &c, nil
}
