# ── Stage 1: generate proto ───────────────────────────────────────────────────
FROM golang:alpine AS proto-gen

RUN apk add --no-cache protobuf && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /repo
COPY proto/ ./proto/

RUN mkdir -p proto/gen && \
    protoc \
      --proto_path=proto \
      --go_out=proto/gen      --go_opt=paths=source_relative \
      --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
      client.proto

# ── Stage 2: build binary ─────────────────────────────────────────────────────
FROM golang:alpine AS builder

WORKDIR /repo
COPY --from=proto-gen /repo/proto      ./proto
COPY services/api-gateway              ./services/api-gateway

WORKDIR /repo/services/api-gateway
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/service ./cmd/

# ── Stage 3: minimal runtime ──────────────────────────────────────────────────
FROM alpine:3.21

RUN addgroup -S app && adduser -S -G app app
USER app

COPY --from=builder /bin/service /bin/service
ENTRYPOINT ["/bin/service"]
