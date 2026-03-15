package config

import "os"

type Config struct {
    DatabaseURL string
}

func Load() Config {
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        databaseURL = "postgres://postgres:postgres@localhost:5432/practice5?sslmode=disable"
    }

    return Config{DatabaseURL: databaseURL}
}
