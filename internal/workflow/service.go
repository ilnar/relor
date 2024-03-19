package workflow

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pb "github.com/ilnar/wf/gen/pb/api"
	"github.com/ilnar/wf/internal/model"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
}

type Server struct {
	pb.UnimplementedWorkflowServiceServer

	logger Logger
}

func New(logger Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Run(ctx context.Context, in *pb.RunRequest) (*pb.RunResponse, error) {
	g := model.Graph{}
	if err := g.FromProto(in.Graph); err != nil {
		return nil, fmt.Errorf("failed to parse graph: %w", err)
	}
	if in.Id == "" {
		return nil, fmt.Errorf("id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id: %w", err)
	}
	w := model.NewWorkflow(id, &g)
	s.logger.InfoContext(ctx, "Running workflow", "w", w)
	return &pb.RunResponse{}, nil
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	r := &pb.GetResponse{
		State: &pb.WorkflowState{
			Status:      string(model.WorkflowStatusPending),
			CurrentNode: "start",
		},
	}
	return r, nil
}