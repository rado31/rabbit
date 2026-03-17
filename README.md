# Rabbit example of usage

Two Go services communicating over RabbitMQ using the RPC (request-reply) pattern.

## Services

| Service | Entry point | Role |
|---------|-------------|------|
| `api` | `cmd/api` | Accepts HTTP requests, publishes tasks to a queue, waits for a reply |
| `storage` | `cmd/storage` | Listens on queues, reads/writes data in PostgreSQL, sends replies |

## Stack

- [gin](https://github.com/gin-gonic/gin) — HTTP server
- [amqp091-go](https://github.com/rabbitmq/amqp091-go) — RabbitMQ client
- [pgx](https://github.com/jackc/pgx) — PostgreSQL client

## How it works

```
POST /clients
     │
     ▼
  [ api ]  ── clients.create ──►  [ storage ]
     │             RabbitMQ             │
     │                                  ├── INSERT INTO clients
     │                                  │
     └──── reply queue ◄────────────────┘
```

1. `api` receives an HTTP request and publishes a message to `clients.create` or `clients.get`
2. The message carries `ReplyTo` (the name of a private reply queue) and `CorrelationId` (a unique call ID)
3. `api` blocks and waits for a reply (5 second timeout)
4. `storage` processes the message, talks to PostgreSQL, and publishes the result back to `ReplyTo`
5. `api` receives the reply, matches it by `CorrelationId`, and returns the HTTP response

## Running

### Docker (recommended)

```bash
docker compose up --build
```

That's it. Compose starts RabbitMQ, PostgreSQL, runs the migration automatically, then starts both services.

| URL | Description |
|-----|-------------|
| `http://localhost:8080` | API |
| `http://localhost:15672` | RabbitMQ management UI (guest / guest) |

### Locally

```bash
# Start dependencies
docker compose up rabbitmq postgres -d

# Apply migration
make migrate

# Terminal 1
make storage

# Terminal 2
make api
```

### Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AMQP_URL` | `amqp://guest:guest@localhost:5672/` | RabbitMQ address |
| `PG_URL` | `postgres://postgres:postgres@localhost:5432/clients_db` | PostgreSQL address |
| `ADDR` | `:8080` | HTTP server address (`api` only) |

## API

### Create a client

```
POST /clients
Content-Type: application/json

{
  "first_name": "Ra",
  "last_name":  "Do",
  "age":        26,
  "email":      "rado@example.com"
}
```

```json
{
  "id":         1,
  "first_name": "Ra",
  "last_name":  "Do",
  "age":        26,
  "email":      "rado@example.com"
}
```

### Get a client

```
GET /clients/:id
```

```json
{
  "id":         1,
  "first_name": "Ra",
  "last_name":  "Do",
  "age":        26,
  "email":      "rado@example.com"
}
```

## Project structure

```
cmd/
  api/main.go         — HTTP API + AMQP RPC client
  storage/main.go     — AMQP consumer + PostgreSQL
migrations/
  001_init.sql        — database schema
Dockerfile            — multi-stage build (shared by both services via ARG SERVICE)
docker-compose.yml    — full stack: api, storage, rabbitmq, postgres
```
