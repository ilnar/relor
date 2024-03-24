package job

import (
	"context"
	"time"

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
	pb.UnimplementedJobServiceServer

	logger Logger
	jobs   *storage.JobStorage
}

func New(l Logger, s *storage.JobStorage) *Server {
	return &Server{
		logger: l,
		jobs:   s,
	}
}

func (s *Server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	if in.Reference == nil {
		return nil, status.Errorf(codes.InvalidArgument, "reference is empty")
	}
	r := in.Reference
	if r.WorkflowId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "workflow id is empty")
	}
	if r.WorkflowAction == "" {
		return nil, status.Errorf(codes.InvalidArgument, "workflow action is empty")
	}
	wid, err := uuid.Parse(r.WorkflowId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse workflow id: %v", err)
	}
	jid := model.JobID{
		ID:             id,
		WorkflowID:     wid,
		WorkflowAction: r.WorkflowAction,
	}
	j := model.NewJob(jid, time.Now())
	if err := s.jobs.Save(j); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save job: %v", err)
	}
	return &pb.CreateResponse{}, nil
}

func (s *Server) Claim(ctx context.Context, in *pb.ClaimRequest) (*pb.ClaimResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	j, err := s.jobs.Get(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to get job: %v", err)
	}

	if err := j.ClaimAt(time.Now()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to claim job: %v", err)
	}

	if err := s.jobs.Save(j); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save job: %v", err)
	}

	return &pb.ClaimResponse{ActionId: j.ID().WorkflowAction}, nil
}

func (s *Server) Release(ctx context.Context, in *pb.ReleaseRequest) (*pb.ReleaseResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	j, err := s.jobs.Get(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to get job: %v", err)
	}
	if err := j.CloseAt(time.Now(), "released"); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to release job: %v", err)
	}
	if err := s.jobs.Save(j); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save job: %v", err)
	}
	return &pb.ReleaseResponse{}, nil
}

func (s *Server) Complete(ctx context.Context, in *pb.CompleteRequest) (*pb.CompleteResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	j, err := s.jobs.Get(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to get job: %v", err)
	}
	if err := j.CloseAt(time.Now(), "completed"); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to complete job: %v", err)
	}
	if err := s.jobs.Save(j); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save job: %v", err)
	}

	return &pb.CompleteResponse{}, nil
}
