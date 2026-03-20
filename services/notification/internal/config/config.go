package config

import (
	"os"
	"strconv"
)

type Config struct {
	AMQPURL  string
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string
}

func Load() Config {
	return Config{
		AMQPURL:  getenv("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
		SMTPHost: getenv("SMTP_HOST", "localhost"),
		SMTPPort: getenvInt("SMTP_PORT", 25),
		SMTPUser: getenv("SMTP_USER", ""),
		SMTPPass: getenv("SMTP_PASS", ""),
		SMTPFrom: getenv("SMTP_FROM", "noreply@example.com"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}

	return fallback
}
