package schedule

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	pb "github.com/ilnar/wf/gen/pb/api"
	"github.com/ilnar/wf/internal/model"
	"github.com/ilnar/wf/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	pollInterval  = 5 * time.Second
	defaulTimeout = 60 * time.Second
)

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type Scheduler struct {
	wfStore    *storage.WorkflowStorage
	jobClient  pb.JobServiceClient
	logger     Logger
	workerName string
}

func New(wfStore *storage.WorkflowStorage, l Logger) *Scheduler {
	return &Scheduler{
		wfStore:    wfStore,
		logger:     l,
		workerName: "scheduler-" + uuid.NewString(),
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

	listener, err := s.jobClient.Listen(ctx, &pb.ListenRequest{
		WorkerId: s.workerName,
		ListenTo: pb.ListenOperation_COMPLETE,
	})
	if err != nil {
		return fmt.Errorf("failed to listen for jobs: %w", err)
	}

	go func() {
		for {
			resp, err := listener.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				s.logger.ErrorContext(ctx, "failed to receive job", "err", err)
				return
			}
			s.logger.InfoContext(ctx, "Received job", "job", resp)
			if resp.ResultLabel == "" {
				s.logger.ErrorContext(ctx, "missing result label", "job", resp)
				continue
			}
			wfID, err := uuid.Parse(resp.Reference.WorkflowId)
			if err != nil {
				s.logger.ErrorContext(ctx, "failed to parse workflow ID", "err", err)
				continue
			}
			na := storage.NextAction{
				ID:            wfID,
				Label:         resp.ResultLabel,
				CurrentAction: resp.Reference.WorkflowAction,
			}
			if err := s.wfStore.UpdateNextAction(listener.Context(), na); err != nil {
				s.logger.ErrorContext(ctx, "failed to update next workflow action", "err", err)
				continue
			}
		}
	}()

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

	// Get the timeout for this job.
	// When the timeout is reached, a new job for the same action will be created.
	timeout, err := w.Graph.Timeout(w.CurrentNode)
	if err != nil {
		return fmt.Errorf("failed to get timeout: %w", err)
	}
	if timeout == 0 {
		timeout = time.Duration(defaulTimeout)
	}
	if err := s.wfStore.UpdateTimeout(ctx, w.ID, timeout); err != nil {
		return fmt.Errorf("failed to update timeout: %w", err)
	}

	// Schedule the next action.
	jcp := &pb.CreateRequest{
		Id: uuid.NewString(),
		Reference: &pb.Reference{
			WorkflowId:     w.ID.String(),
			WorkflowAction: w.CurrentNode,
		},
		ResultLabels: labels,
		Ttl:          durationpb.New(timeout),
	}
	s.logger.InfoContext(ctx, "Scheduling workflow", "request", jcp)
	_, err = s.jobClient.Create(ctx, jcp)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}
