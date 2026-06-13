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

	"lessonsgo/testpaymentgeteway/internal/api"
	"lessonsgo/testpaymentgeteway/internal/bank"
	"lessonsgo/testpaymentgeteway/internal/domain/payment"
	"lessonsgo/testpaymentgeteway/internal/storage/memory"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	bankURL := flag.String("bank-url", "http://localhost:8081", "bank simulator base URL")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	store := memory.New()
	bankClient := bank.NewHTTPClient(*bankURL)
	svc := payment.NewService(store, bankClient)
	handler := api.NewHandler(svc, logger)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		logger.Info("payment gateway started", "addr", *addr, "bank_url", *bankURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("shutdown complete")
}
