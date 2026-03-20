module github.com/rado31/rabbit/storage

go 1.23

require (
	github.com/jackc/pgx/v5            v5.8.0
	github.com/rabbitmq/amqp091-go     v1.10.0
	github.com/rado31/rabbit/proto     v0.0.0
	google.golang.org/grpc             v1.79.3
)

replace github.com/rado31/rabbit/proto => ../../proto
