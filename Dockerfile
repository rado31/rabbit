FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/service ./cmd/${SERVICE}

FROM alpine:3.21

RUN addgroup -S app && adduser -S -G app app

USER app

COPY --from=builder /bin/service /bin/service

ENTRYPOINT ["/bin/service"]
