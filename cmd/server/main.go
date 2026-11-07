package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/todo-microservice/internal/handler"
	"github.com/yourusername/todo-microservice/internal/middleware"
	"github.com/yourusername/todo-microservice/internal/repository"
	"github.com/yourusername/todo-microservice/internal/service"
)

func main() {
	// ── Configuration (env-first with sensible defaults) ──────────────────────
	port := getEnv("PORT", "8080")
	dsn := getEnv("DATABASE_DSN", "./todos.db")

	// ── Dependency wiring ─────────────────────────────────────────────────────
	repo, err := repository.NewSQLiteRepository(dsn)
	if err != nil {
		log.Fatalf("failed to initialise repository: %v", err)
	}

	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	// ── Router ────────────────────────────────────────────────────────────────
	mux := http.NewServeMux()

	// Health-check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	h.RegisterRoutes(mux)

	// ── Middleware stack ──────────────────────────────────────────────────────
	root := middleware.Chain(mux, middleware.Logger, middleware.CORS)

	// ── HTTP server ───────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      root,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start in a goroutine so the main thread can listen for signals.
	go func() {
		log.Printf("Todo microservice listening on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("Server stopped.")
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
