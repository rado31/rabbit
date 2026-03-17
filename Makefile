api:
	@go run ./cmd/api

storage:
	@go run ./cmd/storage

migrate:
	@psql "postgres://postgres:postgres@localhost:5432/clients_db" -f migrations/001_init.sql

.PHONY: api storage migrate
