package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"

	clientpb "github.com/rado31/rabbit/proto/gen"
	"github.com/rado31/rabbit/storage/internal/config"
	"github.com/rado31/rabbit/storage/internal/handler"
	"github.com/rado31/rabbit/storage/internal/repository"
	"github.com/rado31/rabbit/storage/internal/service"
)

func must(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.PgURL)
	must(err, "postgres connect")
	defer pool.Close()

	must(pool.Ping(ctx), "postgres ping")

	amqpConn, err := amqp.Dial(cfg.AMQPURL)
	must(err, "rabbitmq connect")
	defer amqpConn.Close()

	amqpCh, err := amqpConn.Channel()
	must(err, "rabbitmq channel")
	defer amqpCh.Close()

	db := repository.NewPostgresRepository(pool)

	mq, err := repository.NewRabbitMQRepository(amqpCh)
	must(err, "rabbitmq repository")

	svc := service.New(db, mq)
	grpcHandler := handler.New(svc)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	must(err, "tcp listen")

	grpcServer := grpc.NewServer()
	clientpb.RegisterClientServiceServer(grpcServer, grpcHandler)

	log.Printf("storage gRPC server listening on %s", cfg.GRPCAddr)
	log.Fatal(grpcServer.Serve(lis))
}
