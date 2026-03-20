package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rado31/rabbit/notification/internal/config"
	"github.com/rado31/rabbit/notification/internal/repository"
	"github.com/rado31/rabbit/notification/internal/service"
)

func must(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := amqp.Dial(cfg.AMQPURL)
	must(err, "rabbitmq connect")
	defer conn.Close()

	ch, err := conn.Channel()
	must(err, "rabbitmq channel")
	defer ch.Close()

	mq, err := repository.New(ch)
	must(err, "rabbitmq repository")

	svc := service.New(cfg, mq)
	must(svc.Start(ctx), "start consumer")

	log.Println("notification service started")
	<-ctx.Done()
	log.Println("notification service stopped")
}
