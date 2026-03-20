package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rado31/rabbit/api-gateway/internal/config"
	"github.com/rado31/rabbit/api-gateway/internal/handler"
	"github.com/rado31/rabbit/api-gateway/internal/service"
	clientpb "github.com/rado31/rabbit/proto/gen"
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

	conn, err := grpc.NewClient(
		cfg.StorageAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	must(err, "grpc dial storage")
	defer conn.Close()

	grpcClient := clientpb.NewClientServiceClient(conn)
	svc := service.New(grpcClient)
	h := handler.New(svc)

	r := gin.Default()
	h.RegisterRoutes(r)

	srv := &http.Server{Addr: cfg.Addr, Handler: r}

	go func() {
		log.Printf("api-gateway listening on %s", cfg.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	must(srv.Shutdown(shutdownCtx), "graceful shutdown")
}
