# Rabbit

A Go microservices demo using gRPC, RabbitMQ, PostgreSQL, and SMTP.

## Architecture

```
Client
  │
  ▼
api-gateway   (HTTP :8080)
  │  REST (Gin)
  │  gRPC
  ▼
storage       (gRPC :50051)
  │  pgx
  │  PostgreSQL
  │
  └─ RabbitMQ ──► notification
                    SMTP (email on registration)
```

## Services

| Service      | Transport | Description                                      |
|--------------|-----------|--------------------------------------------------|
| api-gateway  | HTTP      | Accepts REST requests, proxies to storage via gRPC |
| storage      | gRPC      | Persists clients in PostgreSQL, publishes events to RabbitMQ |
| notification | —         | Consumes RabbitMQ events, sends welcome emails via SMTP |

## API

### Create client

```
POST /clients
Content-Type: application/json

{
  "surname": "Doe",
  "name":    "John",
  "age":     30,
  "email":   "john@example.com"
}
```

**Response** `201 Created`

```json
{
  "id":      1,
  "surname": "Doe",
  "name":    "John",
  "age":     30,
  "email":   "john@example.com"
}
```

### Get client

```
GET /clients/:id
```

**Response** `200 OK` or `404 Not Found`

## Running

```bash
make up
```

Starts all services via Docker Compose. On first run, images are built from source.

```bash
make down   # stop and remove containers + volumes
```

### Ports

| Port  | Service               |
|-------|-----------------------|
| 8080  | api-gateway (HTTP)    |
| 15672 | RabbitMQ management UI (guest/guest) |
| 8025  | Mailpit web UI (captured emails) |

## Development

### Regenerate proto

Requires `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc` installed locally.

```bash
make proto
```

Generated files are written to `proto/gen/`.

## Environment variables

### storage

| Variable   | Default                                    | Description        |
|------------|--------------------------------------------|--------------------|
| `GRPC_ADDR`| `:50051`                                   | gRPC listen address |
| `PG_URL`   | `postgres://postgres:postgres@localhost:5432/clients_db` | PostgreSQL DSN |
| `AMQP_URL` | `amqp://guest:guest@localhost:5672/`       | RabbitMQ URL       |

### api-gateway

| Variable       | Default        | Description          |
|----------------|----------------|----------------------|
| `ADDR`         | `:8080`        | HTTP listen address  |
| `STORAGE_ADDR` | `localhost:50051` | Storage gRPC address |

### notification

| Variable    | Default                              | Description        |
|-------------|--------------------------------------|--------------------|
| `AMQP_URL`  | `amqp://guest:guest@localhost:5672/` | RabbitMQ URL       |
| `SMTP_HOST` | `localhost`                          | SMTP server host   |
| `SMTP_PORT` | `25`                                 | SMTP server port   |
| `SMTP_USER` | —                                    | SMTP username      |
| `SMTP_PASS` | —                                    | SMTP password      |
| `SMTP_FROM` | `noreply@example.com`                | Sender address     |
