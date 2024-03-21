package storage

import (
	"context"
	"fmt"

	"github.com/ilnar/wf/gen/sqlc"
	"github.com/ilnar/wf/internal/model"
	"google.golang.org/protobuf/encoding/protojson"
)

type Storage struct {
	q *sqlc.Queries
}

func New(q *sqlc.Queries) *Storage {
	return &Storage{q: q}
}

// CreateWorkflow creates a new workflow
func (s *Storage) CreateWorkflow(ctx context.Context, w model.Workflow) error {
	gpb, err := w.Graph.ToProto()
	if err != nil {
		return fmt.Errorf("failed to convert graph to proto: %w", err)
	}
	b, err := protojson.Marshal(gpb)
	if err != nil {
		return fmt.Errorf("failed to marshal graph to json: %w", err)
	}
	arg := sqlc.CreateWorkflowParams{
		CurrentNode: w.CurrentNode,
		Status:      string(w.Status),
		Graph:       b,
	}
	if _, err := s.q.CreateWorkflow(ctx, arg); err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}
	return nil
}
