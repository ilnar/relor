package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ilnar/wf/gen/sqlc"
	"github.com/ilnar/wf/internal/model"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/ilnar/wf/gen/pb/graph"
)

type DBQuery interface {
	CreateWorkflow(ctx context.Context, arg sqlc.CreateWorkflowParams) (sqlc.Workflow, error)
	GetWorkflow(ctx context.Context, id uuid.UUID) (sqlc.Workflow, error)
	GetNextWorkflows(ctx context.Context) ([]sqlc.Workflow, error)
	UpdateWorkflowNextAction(ctx context.Context, arg sqlc.UpdateWorkflowNextActionParams) (sqlc.Workflow, error)
	UpdateWorkflowStatus(ctx context.Context, arg sqlc.UpdateWorkflowStatusParams) (sqlc.Workflow, error)
}

type WorkflowStorage struct {
	q DBQuery
}

func NewWorkflowStorage(q DBQuery) *WorkflowStorage {
	return &WorkflowStorage{q: q}
}

func (s *WorkflowStorage) CreateWorkflow(ctx context.Context, w model.Workflow) error {
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

func (s *WorkflowStorage) UpdateStatus(ctx context.Context, id uuid.UUID, status model.WorkflowStatus) error {
	if _, err := s.q.UpdateWorkflowStatus(ctx, sqlc.UpdateWorkflowStatusParams{
		ID:     id,
		Status: string(status),
	}); err != nil {
		return fmt.Errorf("failed to update workflow status: %w", err)
	}
	return nil
}

func (s *WorkflowStorage) UpdateNextAction(ctx context.Context, id uuid.UUID, node string, nextActionAt time.Time) error {
	if _, err := s.q.UpdateWorkflowNextAction(ctx, sqlc.UpdateWorkflowNextActionParams{
		ID:           id,
		CurrentNode:  node,
		NextActionAt: nextActionAt,
	}); err != nil {
		return fmt.Errorf("failed to update workflow next action: %w", err)
	}
	return nil
}

func (s *WorkflowStorage) GetWorkflow(ctx context.Context, id uuid.UUID) (model.Workflow, error) {
	w, err := s.q.GetWorkflow(ctx, id)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to get workflow: %w", err)
	}
	wm, err := workflowToModel(w)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to convert workflow: %w", err)
	}
	return wm, nil
}

func (s *WorkflowStorage) GetNextWorkflows(ctx context.Context) ([]model.Workflow, error) {
	ws, err := s.q.GetNextWorkflows(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get next workflow: %w", err)
	}
	wms := make([]model.Workflow, 0, len(ws))
	for _, w := range ws {
		wm, err := workflowToModel(w)
		if err != nil {
			return nil, fmt.Errorf("failed to convert workflow: %w", err)
		}
		wms = append(wms, wm)
	}
	return wms, nil
}

func workflowToModel(w sqlc.Workflow) (model.Workflow, error) {
	gpb := &pb.Graph{}
	if err := protojson.Unmarshal(w.Graph, gpb); err != nil {
		return model.Workflow{}, fmt.Errorf("failed to unmarshal graph: %w", err)
	}
	g := &model.Graph{}
	if err := g.FromProto(gpb); err != nil {
		return model.Workflow{}, fmt.Errorf("failed to build graph: %w", err)
	}
	return model.Workflow{
		ID:          w.ID,
		CurrentNode: w.CurrentNode,
		Status:      model.WorkflowStatus(w.Status),
		Graph:       g,
	}, nil
}
