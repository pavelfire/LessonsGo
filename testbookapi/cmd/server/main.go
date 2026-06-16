package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lessonsgo/testbookapi/internal/api"
	"lessonsgo/testbookapi/internal/auth"
	"lessonsgo/testbookapi/internal/book"
	"lessonsgo/testbookapi/internal/cart"
	"lessonsgo/testbookapi/internal/category"
	"lessonsgo/testbookapi/internal/domain"
	"lessonsgo/testbookapi/internal/storage/memory"
	"lessonsgo/testbookapi/internal/storage/postgres"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	databaseURL := flag.String("database-url", os.Getenv("DATABASE_URL"), "PostgreSQL connection URL")
	jwtSecret := flag.String("jwt-secret", envOrDefault("JWT_SECRET", "dev-secret-change-me"), "JWT signing secret")
	useMemory := flag.Bool("memory", false, "use in-memory storage instead of PostgreSQL")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx := context.Background()
	var repo domain.Repository
	var cleanup func()

	if *useMemory || *databaseURL == "" {
		logger.Warn("using in-memory storage; data will not persist")
		repo = memory.New()
		cleanup = func() {}
	} else {
		store, err := postgres.New(ctx, *databaseURL)
		if err != nil {
			logger.Error("failed to connect to database", "error", err)
			os.Exit(1)
		}
		repo = store
		cleanup = store.Close
	}
	defer cleanup()

	authSvc := auth.NewService(repo, *jwtSecret, 24*time.Hour)
	categorySvc := category.NewService(repo)
	bookSvc := book.NewService(repo)
	cartSvc := cart.NewService(repo)

	handler := api.NewHandler(authSvc, categorySvc, bookSvc, cartSvc, logger)

	mux := http.NewServeMux()
	handler.Register(mux)

	serverCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	api.StartCartCleanup(serverCtx, cartSvc, logger, time.Minute)

	server := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("book shop API started", "addr", *addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-serverCtx.Done()
	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("shutdown complete")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
