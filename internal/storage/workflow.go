package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ilnar/wf/gen/sqlc"
	"github.com/ilnar/wf/internal/model"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/ilnar/wf/gen/pb/graph"
)

type DBQuery interface {
	CreateWorkflow(ctx context.Context, db sqlc.DBTX, arg sqlc.CreateWorkflowParams) (sqlc.Workflow, error)
	GetWorkflow(ctx context.Context, db sqlc.DBTX, id uuid.UUID) (sqlc.Workflow, error)
	GetNextWorkflows(ctx context.Context, db sqlc.DBTX) ([]sqlc.Workflow, error)
	UpdateWorkflowNextAction(ctx context.Context, db sqlc.DBTX, arg sqlc.UpdateWorkflowNextActionParams) (sqlc.Workflow, error)
	UpdateWorkflowStatus(ctx context.Context, db sqlc.DBTX, arg sqlc.UpdateWorkflowStatusParams) (sqlc.Workflow, error)
	UpdateWorkflowNextActionAt(ctx context.Context, db sqlc.DBTX, arg sqlc.UpdateWorkflowNextActionAtParams) (sqlc.Workflow, error)
}

type TxManager interface {
	sqlc.DBTX

	Begin() (*sql.Tx, error)
}

type WorkflowStorage struct {
	q   DBQuery
	txm TxManager
}

func NewWorkflowStorage(q DBQuery, txm TxManager) *WorkflowStorage {
	return &WorkflowStorage{q: q, txm: txm}
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
	if _, err := s.q.CreateWorkflow(ctx, s.txm, arg); err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}
	return nil
}

func (s *WorkflowStorage) UpdateNextAction(ctx context.Context, id uuid.UUID, label string) error {
	tx, err := s.txm.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get workflow.
	w, err := s.q.GetWorkflow(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get workflow: %w", err)
	}
	wf, err := workflowToModel(w)
	if err != nil {
		return fmt.Errorf("failed to convert workflow: %w", err)
	}

	// Update next action.
	nextNode, err := wf.Graph.NextNodeID(wf.CurrentNode, label)
	if err != nil {
		return fmt.Errorf("failed to get next node: %w", err)
	}
	if _, err := s.q.UpdateWorkflowNextAction(ctx, tx, sqlc.UpdateWorkflowNextActionParams{
		ID:          id,
		CurrentNode: nextNode,
	}); err != nil {
		return fmt.Errorf("failed to update workflow next action: %w", err)
	}
	nextLabels, err := wf.Graph.OutLabels(nextNode)
	if err != nil {
		return fmt.Errorf("failed to get out labels: %w", err)
	}
	if len(nextLabels) == 0 {
		if _, err := s.q.UpdateWorkflowStatus(ctx, tx, sqlc.UpdateWorkflowStatusParams{
			ID:     id,
			Status: string(model.WorkflowStatusCompleted),
		}); err != nil {
			return fmt.Errorf("failed to update workflow status: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *WorkflowStorage) UpdateTimeout(ctx context.Context, id uuid.UUID, timeout time.Duration) error {
	if _, err := s.q.UpdateWorkflowNextActionAt(ctx, s.txm, sqlc.UpdateWorkflowNextActionAtParams{
		ID:    id,
		Delay: int64(timeout.Seconds()),
	}); err != nil {
		return fmt.Errorf("failed to update workflow next action at: %w", err)
	}
	return nil
}

func (s *WorkflowStorage) GetWorkflow(ctx context.Context, id uuid.UUID) (model.Workflow, error) {
	w, err := s.q.GetWorkflow(ctx, s.txm, id)
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
	ws, err := s.q.GetNextWorkflows(ctx, s.txm)
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
