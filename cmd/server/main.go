package main

import (
	"context"
	"log/slog"

	"github.com/ilnar/wf/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := 8080
	logger := slog.Default()

	srv := server.New(port, logger)
	if err := srv.Serve(ctx); err != nil {
		logger.ErrorContext(ctx, "Error serving", "err", err)
	}
}
