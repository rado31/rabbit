# ── Stage 1: build binary ─────────────────────────────────────────────────────
FROM golang:alpine AS builder

WORKDIR /repo
COPY services/notification ./services/notification

WORKDIR /repo/services/notification
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/service ./cmd/

# ── Stage 2: minimal runtime ──────────────────────────────────────────────────
FROM alpine:3.21

RUN addgroup -S app && adduser -S -G app app
USER app

COPY --from=builder /bin/service /bin/service
ENTRYPOINT ["/bin/service"]
