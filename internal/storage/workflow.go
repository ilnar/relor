package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ilnar/wf/gen/sqlc"
	"github.com/ilnar/wf/internal/model"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/ilnar/wf/gen/pb/graph"
)

type DBQuery interface {
	CreateWorkflow(ctx context.Context, arg sqlc.CreateWorkflowParams) (sqlc.Workflow, error)
	GetWorkflow(ctx context.Context, id uuid.UUID) (sqlc.Workflow, error)
}

type WorkflowStorage struct {
	q DBQuery
}

func NewWorkflowStorage(q DBQuery) *WorkflowStorage {
	return &WorkflowStorage{q: q}
}

func (s *WorkflowStorage) CreateWorkflow(ctx context.Context, w *model.Workflow) error {
	gpb, err := w.Graph.ToProto()
	if err != nil {
		return fmt.Errorf("failed to convert graph to proto: %w", err)
	}
	b, err := protojson.Marshal(gpb)
	if err != nil {
		return fmt.Errorf("failed to marshal graph to json: %w", err)
	}
	arg := sqlc.CreateWorkflowParams{
		ID:          w.ID,
		CurrentNode: w.CurrentNode,
		Status:      string(w.Status),
		Graph:       b,
	}
	if _, err := s.q.CreateWorkflow(ctx, arg); err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}
	return nil
}

func (s *WorkflowStorage) GetWorkflow(ctx context.Context, id uuid.UUID) (*model.Workflow, error) {
	w, err := s.q.GetWorkflow(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}
	gpb := &pb.Graph{}
	if err := protojson.Unmarshal(w.Graph, gpb); err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph: %w", err)
	}
	g := &model.Graph{}
	if err := g.FromProto(gpb); err != nil {
		return nil, fmt.Errorf("failed to build graph: %w", err)
	}
	return &model.Workflow{
		ID:          w.ID,
		CurrentNode: w.CurrentNode,
		Status:      model.WorkflowStatus(w.Status),
		Graph:       g,
	}, nil
}
