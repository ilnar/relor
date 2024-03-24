package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ilnar/wf/gen/pb/api"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type Notify func(c chan<- os.Signal, sig ...os.Signal)

type Server struct {
	logger Logger
	port   int
	notify Notify
	wfs    pb.WorkflowServiceServer
	js     pb.JobServiceServer
}

func New(port int, logger Logger, wfs pb.WorkflowServiceServer, js pb.JobServiceServer) *Server {
	return &Server{
		logger: logger,
		port:   port,
		notify: signal.Notify,
		wfs:    wfs,
		js:     js,
	}
}

func (s Server) Serve(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Starting server", "port", s.port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to listen", "err", err)
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()

	gs := grpc.NewServer()
	defer gs.GracefulStop()

	pb.RegisterWorkflowServiceServer(gs, s.wfs)
	pb.RegisterJobServiceServer(gs, s.js)
	reflection.Register(gs)

	stopChan := make(chan os.Signal, 1)
	s.notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error)
	go func() {
		if err := gs.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		s.logger.ErrorContext(ctx, "Error serving", "err", err)
		return err
	case <-stopChan:
		s.logger.InfoContext(ctx, "Received stop signal")
	case <-ctx.Done():
	}
	s.logger.InfoContext(ctx, "Stopping server")
	return nil
}
