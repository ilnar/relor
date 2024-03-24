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
	jq     chan model.Job
}

func New(l Logger, s *storage.JobStorage) *Server {
	return &Server{
		logger: l,
		jobs:   s,
		jq:     make(chan model.Job, 1000),
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
	s.jq <- *j
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

func (s *Server) Listen(in *pb.ListenRequest, stream pb.JobService_ListenServer) error {
	s.logger.InfoContext(stream.Context(), "Listening for jobs", "worker_id", in.WorkerId)
	for {
		select {
		case j := <-s.jq:
			s.logger.InfoContext(stream.Context(), "Sending job",
				"id", j.ID().ID,
				"workflow_id", j.ID().WorkflowID,
				"workflow_action", j.ID().WorkflowAction,
				"worker_id", in.WorkerId,
			)
			if err := stream.Send(&pb.Job{
				Id: j.ID().ID.String(),
				Reference: &pb.Reference{
					WorkflowId:     j.ID().WorkflowID.String(),
					WorkflowAction: j.ID().WorkflowAction,
				},
			}); err != nil {
				return status.Errorf(codes.Internal, "failed to send job: %v", err)
			}
		case <-stream.Context().Done():
			s.logger.InfoContext(stream.Context(), "Stopped listening for jobs", "worker_id", in.WorkerId)
			return nil
		}
	}
}
