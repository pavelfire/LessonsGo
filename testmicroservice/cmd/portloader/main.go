package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"lessonsgo/testmicroservice/internal/domain/port"
	"lessonsgo/testmicroservice/internal/parser"
	"lessonsgo/testmicroservice/internal/storage/memory"
)

func main() {
	filePath := flag.String("file", "ports.json", "path to ports JSON file")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	store := memory.New()
	svc := port.NewService(store)

	processed, err := loadPorts(ctx, logger, *filePath, svc)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Info("interrupted during load",
				"processed", processed,
				"stored", svc.Count(),
			)
			os.Exit(0)
		}
		logger.Error("failed to load ports", "error", err)
		os.Exit(1)
	}

	logger.Info("ports loaded",
		"processed", processed,
		"stored", svc.Count(),
	)

	<-ctx.Done()
	logger.Info("shutdown complete", "stored", svc.Count())
}

func loadPorts(ctx context.Context, logger *slog.Logger, path string, svc *port.Service) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	logger.Info("loading ports", "file", path)

	return parser.StreamPorts(ctx, f, func(ctx context.Context, p port.Port) error {
		return svc.Upsert(ctx, p)
	})
}
