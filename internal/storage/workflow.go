package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gemlab-dev/relor/gen/sqlc"
	"github.com/gemlab-dev/relor/internal/model"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/gemlab-dev/relor/gen/pb/graph"
)

type DBQuery interface {
	CreateWorkflow(ctx context.Context, db sqlc.DBTX, arg sqlc.CreateWorkflowParams) (sqlc.Workflow, error)
	GetWorkflow(ctx context.Context, db sqlc.DBTX, id uuid.UUID) (sqlc.Workflow, error)
	GetNextWorkflows(ctx context.Context, db sqlc.DBTX) ([]sqlc.Workflow, error)
	UpdateWorkflowNextAction(ctx context.Context, db sqlc.DBTX, arg sqlc.UpdateWorkflowNextActionParams) (sqlc.Workflow, error)
	UpdateWorkflowStatus(ctx context.Context, db sqlc.DBTX, arg sqlc.UpdateWorkflowStatusParams) (sqlc.Workflow, error)
	UpdateWorkflowNextActionAt(ctx context.Context, db sqlc.DBTX, arg sqlc.UpdateWorkflowNextActionAtParams) (sqlc.Workflow, error)
	CreateTransition(ctx context.Context, db sqlc.DBTX, arg sqlc.CreateTransitionParams) (sqlc.Transition, error)
	GetLatestTransition(ctx context.Context, db sqlc.DBTX, workflowID uuid.UUID) ([]sqlc.Transition, error)
	GetFirstTransition(ctx context.Context, db sqlc.DBTX, workflowID uuid.UUID) ([]sqlc.Transition, error)
	UpdateTransitionNext(ctx context.Context, db sqlc.DBTX, arg sqlc.UpdateTransitionNextParams) (sqlc.Transition, error)
	GetTransitions(ctx context.Context, db sqlc.DBTX, workflowID uuid.UUID) ([]sqlc.Transition, error)
}

type TxManager interface {
	sqlc.DBTX

	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
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

type NextAction struct {
	ID                uuid.UUID
	Label             string
	CurrentTransition uuid.UUID
}

func (s *WorkflowStorage) UpdateNextAction(ctx context.Context, na NextAction) error {
	tx, err := s.txm.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.updateNextActionTxn(ctx, tx, na); err != nil {
		return fmt.Errorf("failed to update next action: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *WorkflowStorage) updateNextActionTxn(ctx context.Context, tx sqlc.DBTX, na NextAction) error {
	// Get workflow.
	w, err := s.q.GetWorkflow(ctx, tx, na.ID)
	if err != nil {
		return fmt.Errorf("failed to get workflow: %w", err)
	}
	wf, err := workflowToModel(w)
	if err != nil {
		return fmt.Errorf("failed to convert workflow: %w", err)
	}

	// Validate optimistic lock.
	lastTns, err := s.q.GetLatestTransition(ctx, tx, na.ID)
	if err != nil {
		return fmt.Errorf("failed to get latest transition: %w", err)
	}
	if len(lastTns) > 1 {
		return fmt.Errorf("more than one last transitions found: %v", lastTns)
	}
	if len(lastTns) == 1 && lastTns[0].ID != na.CurrentTransition {
		return fmt.Errorf("optimistic lock failed: transition IDs don't match: %v != %v", lastTns[0].ID, na.CurrentTransition)
	}

	// Validate current action and get next node.
	nextNode, err := wf.Graph.NextNodeID(wf.CurrentNode, na.Label)
	if err != nil {
		return fmt.Errorf("failed to get next node: %w", err)
	}

	// Create transitions.
	if err := s.createTransition(ctx, tx, na.ID, wf.CurrentNode, nextNode, na.Label); err != nil {
		return fmt.Errorf("failed to update transitions: %w", err)
	}

	// Update next action.
	if _, err := s.q.UpdateWorkflowNextAction(ctx, tx, sqlc.UpdateWorkflowNextActionParams{
		ID:          na.ID,
		CurrentNode: nextNode,
	}); err != nil {
		return fmt.Errorf("failed to update workflow next action: %w", err)
	}

	// Update status if no next labels.
	nextLabels, err := wf.Graph.OutLabels(nextNode)
	if err != nil {
		return fmt.Errorf("failed to get out labels: %w", err)
	}
	if len(nextLabels) == 0 {
		if _, err := s.q.UpdateWorkflowStatus(ctx, tx, sqlc.UpdateWorkflowStatusParams{
			ID:     na.ID,
			Status: string(model.WorkflowStatusCompleted),
		}); err != nil {
			return fmt.Errorf("failed to update workflow status: %w", err)
		}
	}
	return nil
}

func (s *WorkflowStorage) createTransition(ctx context.Context, tx sqlc.DBTX, workflowID uuid.UUID, fromNode, toNode, label string) error {
	firstTnID, lastTnID, err := getFirstAndLastTranstitions(ctx, s.q, tx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get first and last transitions: %w", err)
	}

	newTn, err := s.q.CreateTransition(ctx, tx, sqlc.CreateTransitionParams{
		WorkflowID: workflowID,
		FromNode:   fromNode,
		ToNode:     toNode,
		Label:      label,
		Previous:   lastTnID,
		// This creates a loop, however we need this to meet these constraints:
		// 1. We can't have two transitions with empty Next field.
		// 2. We can't have two transitions with the same Next field.
		// 3. Next has to point to a real transition.
		// The only free and valid transition ID if the one at the beginning of the linked list.
		// We will set this to NULL in the next step.
		Next: firstTnID,
	})
	if err != nil {
		return fmt.Errorf("failed to stage transition: %w", err)
	}

	if lastTnID.Valid {
		if _, err := s.q.UpdateTransitionNext(ctx, tx, sqlc.UpdateTransitionNextParams{
			ID:   lastTnID.UUID,
			Next: uuid.NullUUID{UUID: newTn.ID, Valid: true},
		}); err != nil {
			return fmt.Errorf("failed to update transition next: %w", err)
		}
	}

	if firstTnID.Valid {
		if _, err := s.q.UpdateTransitionNext(ctx, tx, sqlc.UpdateTransitionNextParams{
			ID: newTn.ID,
			// Now as we set the previous transition's Next field, we can use Null again.
			Next: uuid.NullUUID{},
		}); err != nil {
			return fmt.Errorf("failed to update transition next: %w", err)
		}
	}
	return nil
}

func getFirstAndLastTranstitions(ctx context.Context, q DBQuery, tx sqlc.DBTX, id uuid.UUID) (first uuid.NullUUID, last uuid.NullUUID, err error) {
	firstTns, err := q.GetFirstTransition(ctx, tx, id)
	if err != nil {
		return
	}
	if len(firstTns) > 1 {
		err = fmt.Errorf("more than one first transitions found: %v", firstTns)
		return
	}
	if len(firstTns) == 1 {
		first = uuid.NullUUID{UUID: firstTns[0].ID, Valid: true}
	}

	lastTns, err := q.GetLatestTransition(ctx, tx, id)
	if err != nil {
		return
	}
	if len(lastTns) > 1 {
		err = fmt.Errorf("more than one last transitions found: %v", lastTns)
		return
	}
	if len(lastTns) == 1 {
		last = uuid.NullUUID{UUID: lastTns[0].ID, Valid: true}
	}
	return
}

func (s *WorkflowStorage) GetHistory(ctx context.Context, workflowID uuid.UUID) (*model.Transition, error) {
	w, err := s.q.GetWorkflow(ctx, s.txm, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}
	ts, err := s.q.GetTransitions(ctx, s.txm, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transitions: %w", err)
	}
	rts := make([]model.RawTransition, 0, len(ts))
	for _, t := range ts {
		rts = append(rts, model.RawTransition{
			ID:      t.ID,
			WID:     t.WorkflowID,
			From:    t.FromNode,
			To:      t.ToNode,
			Label:   t.Label,
			Created: t.CreatedAt,
			Prev:    t.Previous,
			Next:    t.Next,
		})
	}
	t, err := model.NewTransitionHistory(w.CreatedAt, rts)
	if err != nil {
		return nil, fmt.Errorf("failed to build transition history: %w", err)
	}
	return t, nil
}

func (s *WorkflowStorage) GetLatestTransition(ctx context.Context, workflowID uuid.UUID) (uuid.UUID, error) {
	t, err := s.q.GetLatestTransition(ctx, s.txm, workflowID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get latest transition: %w", err)
	}
	if len(t) > 1 {
		return uuid.Nil, fmt.Errorf("more than one last transitions found: %v", t)
	}
	if len(t) == 0 {
		return uuid.Nil, nil
	}
	return t[0].ID, nil
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
