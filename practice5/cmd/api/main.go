package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "time"

    _ "github.com/lib/pq"

    "practice5_ready/internal/config"
    "practice5_ready/internal/handlers"
    "practice5_ready/internal/repository"
)

func main() {
    cfg := config.Load()

    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    if err := db.Ping(); err != nil {
        log.Fatalf("failed to ping db: %v", err)
    }

    repo := repository.NewRepository(db)
    handler := handlers.NewHandler(repo)

    mux := http.NewServeMux()
    mux.HandleFunc("/health", handler.Health)
    mux.HandleFunc("/users", handler.GetUsers)
    mux.HandleFunc("/friends/common", handler.GetCommonFriends)

    server := &http.Server{
        Addr:         ":8080",
        Handler:      loggingMiddleware(mux),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  30 * time.Second,
    }

    log.Println("Server started on :8080")
    log.Printf("Try in Postman: http://localhost:8080/users?limit=5&offset=0&order_by=name")

    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server failed: %v", err)
    }
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.String())
        next.ServeHTTP(w, r)
    })
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
