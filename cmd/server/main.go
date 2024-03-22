package main

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/ilnar/wf/gen/sqlc"
	"github.com/ilnar/wf/internal/server"
	"github.com/ilnar/wf/internal/storage"
	"github.com/ilnar/wf/internal/workflow"

	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.Default()

	conn, err := sql.Open("postgres", "user=root dbname=workflow password=secret sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	store := storage.NewWorkflowStorage(sqlc.New(conn))

	wfs := workflow.New(logger, store)

	port := 50051
	srv := server.New(port, logger, wfs)
	if err := srv.Serve(ctx); err != nil {
		logger.ErrorContext(ctx, "Error serving", "err", err)
	}
}
