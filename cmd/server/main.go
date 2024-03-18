package main

import (
	"context"
	"log/slog"

	"github.com/ilnar/wf/internal/server"
	"github.com/ilnar/wf/internal/workflow"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := 50051
	logger := slog.Default()
	wfs := workflow.New(logger)

	srv := server.New(port, logger, wfs)
	if err := srv.Serve(ctx); err != nil {
		logger.ErrorContext(ctx, "Error serving", "err", err)
	}
}
