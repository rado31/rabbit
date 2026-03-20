package config

import "os"

type Config struct {
	Addr        string
	StorageAddr string
}

func Load() Config {
	return Config{
		Addr:        getenv("ADDR", ":8080"),
		StorageAddr: getenv("STORAGE_ADDR", "localhost:50051"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
