package workflow

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ilnar/wf/gen/pb/api"
	"github.com/ilnar/wf/internal/model"
	"github.com/ilnar/wf/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
}

type Server struct {
	pb.UnimplementedWorkflowServiceServer

	logger Logger
	store  *storage.WorkflowStorage
}

func New(l Logger, s *storage.WorkflowStorage) *Server {
	return &Server{
		logger: l,
		store:  s,
	}
}

func (s *Server) Run(ctx context.Context, in *pb.RunRequest) (*pb.RunResponse, error) {
	g := model.Graph{}
	if err := g.FromProto(in.Graph); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse graph: %v", err)
	}
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	w := model.NewWorkflow(id, &g)
	s.logger.InfoContext(ctx, "Running workflow", "w", w)
	if err := s.store.CreateWorkflow(ctx, *w); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create workflow: %v", err)
	}
	return &pb.RunResponse{}, nil
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	w, err := s.store.GetWorkflow(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get workflow: %v", err)
	}
	return &pb.GetResponse{
		State: &pb.WorkflowState{
			Status:      string(w.Status),
			CurrentNode: w.CurrentNode,
		},
	}, nil
}

func (s *Server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	if in.ResultLabel == "" {
		return nil, status.Errorf(codes.InvalidArgument, "result label is empty")
	}
	if err := s.store.UpdateNextAction(ctx, id, in.ResultLabel); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update next action: %v", err)
	}
	return &pb.UpdateResponse{}, nil
}
