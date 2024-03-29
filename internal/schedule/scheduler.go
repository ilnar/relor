package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	pb "github.com/ilnar/wf/gen/pb/api"
	"github.com/ilnar/wf/internal/model"
	"github.com/ilnar/wf/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const pollInterval = 30 * time.Second

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type Scheduler struct {
	wfStore   *storage.WorkflowStorage
	jobClient pb.JobServiceClient
	logger    Logger
}

func New(wfStore *storage.WorkflowStorage, l Logger) *Scheduler {
	return &Scheduler{
		wfStore: wfStore,
		logger:  l,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Starting scheduler")

	conn, err := grpc.DialContext(ctx, "localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	s.jobClient = pb.NewJobServiceClient(conn)

	if err := s.poll(ctx); err != nil {
		return fmt.Errorf("failed to poll: %w", err)
	}
	s.logger.InfoContext(ctx, "Stopping scheduler")
	return nil
}

func (s *Scheduler) poll(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(pollInterval):
			wfs, err := s.wfStore.GetNextWorkflows(ctx)
			if err != nil {
				s.logger.ErrorContext(ctx, "failed to list workflows", "err", err)
				continue
			}
			for _, w := range wfs {
				if err := s.schedule(ctx, w); err != nil {
					s.logger.ErrorContext(ctx, "failed to schedule", "err", err)
				}
			}
		}
	}
}

func (s *Scheduler) schedule(ctx context.Context, w model.Workflow) error {
	labels, err := w.Graph.OutLabels(w.CurrentNode)
	if err != nil {
		return fmt.Errorf("failed to get out labels: %w", err)
	}
	// No more actions to take; complete the workflow.
	if len(labels) == 0 {
		s.logger.InfoContext(ctx, "Completing workflow", "workflow", w)
		w.Status = model.WorkflowStatusCompleted
		if err := s.wfStore.UpdateStatus(ctx, w.ID, model.WorkflowStatusCompleted); err != nil {
			return fmt.Errorf("failed to update workflow: %w", err)
		}
		return nil
	}
	// Schedule the next action.
	jcp := &pb.CreateRequest{
		Id: w.ID.String(),
		Reference: &pb.Reference{
			WorkflowId:     uuid.NewString(),
			WorkflowAction: w.CurrentNode,
		},
		ResultLabels: labels,
	}
	s.logger.InfoContext(ctx, "Scheduling workflow", "request", jcp)
	_, err = s.jobClient.Create(ctx, jcp)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}
	return nil
}
