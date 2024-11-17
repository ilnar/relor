package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/google/uuid"
	"github.com/ilnar/wf/gen/sqlc"
	"github.com/ilnar/wf/internal/gossip"
	"github.com/ilnar/wf/internal/job"
	"github.com/ilnar/wf/internal/schedule"
	"github.com/ilnar/wf/internal/server"
	"github.com/ilnar/wf/internal/storage"
	"github.com/ilnar/wf/internal/workflow"

	configpb "github.com/ilnar/wf/gen/pb/config"

	_ "github.com/lib/pq"
)

var config = flag.String("config", "config.json", "Path to the config file")
var gossipPort = flag.Int("gossip-port", 7946, "Port for gossip protocol")

func loadConfig(path string) (*configpb.Config, error) {
	if path == "" {
		return nil, errors.New("config file is required")
	}

	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	cfg := configpb.Config{}
	if err := protojson.Unmarshal(jsonData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return &cfg, nil
}

func main() {
	flag.Parse()

	cfg, err := loadConfig(*config)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	conn, err := sql.Open("postgres", "user=root dbname=workflow password=secret sslmode=disable")
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}
	defer conn.Close()
	wfStore := storage.NewWorkflowStorage(sqlc.New(), conn)
	jobStore := storage.NewJobStorage()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.Default()

	if cfg.GetCluster() != nil {
		if *gossipPort == 0 {
			log.Fatal("gossip port is required")
		}
		name := uuid.New().String()
		seed := []string{}
		for _, m := range cfg.GetCluster().GetGossipSeed() {
			seed = append(seed, fmt.Sprintf("%s:%d", m.GetHostname(), m.GetPort()))
		}
		g, err := gossip.NewGossip(ctx, name, *gossipPort, seed, logger)
		if err != nil {
			log.Fatalf("failed to create gossip: %v", err)
		}
		logger.InfoContext(ctx, "Gossip started", "members", g.Members())
	}

	wfs := workflow.New(logger, wfStore)
	// TODO: Job service should be a separate service.
	js := job.New(logger, jobStore)

	jaddr := cfg.GetJobServiceAddr()
	if jaddr == nil {
		log.Fatal("job service address is required")
	}
	jobServiceAddr := fmt.Sprintf("%s:%d", jaddr.GetHostname(), jaddr.GetPort())
	sch := schedule.New(wfStore, logger, jobServiceAddr)
	go sch.Run(ctx)

	srv := server.New(int(cfg.GetApiPort()), logger, wfs, js)
	if err := srv.Serve(ctx); err != nil {
		logger.ErrorContext(ctx, "Error serving", "err", err)
	}
}
