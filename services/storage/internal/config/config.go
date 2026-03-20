package config

import "os"

type Config struct {
	GRPCAddr string
	PgURL    string
	AMQPURL  string
}

func Load() Config {
	return Config{
		GRPCAddr: getenv("GRPC_ADDR", ":50051"),
		PgURL:    getenv("PG_URL", "postgres://postgres:postgres@localhost:5432/clients_db"),
		AMQPURL:  getenv("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
