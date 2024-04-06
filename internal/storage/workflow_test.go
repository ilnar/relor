package storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ilnar/wf/gen/sqlc"
	"github.com/ilnar/wf/internal/model"
	"google.golang.org/protobuf/encoding/prototext"

	pb "github.com/ilnar/wf/gen/pb/graph"
)

type fakeDBTX struct {
}

func (f *fakeDBTX) BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error) {
	return &sql.Tx{}, nil
}
func (f *fakeDBTX) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return nil, nil
}
func (f *fakeDBTX) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, nil
}
func (f *fakeDBTX) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, nil
}
func (f *fakeDBTX) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return nil
}

type fakeDBQuery struct {
	wfs map[uuid.UUID]sqlc.Workflow
	ts  map[uuid.UUID][]sqlc.Transition
}

func newFakeDBQuery() *fakeDBQuery {
	return &fakeDBQuery{
		wfs: make(map[uuid.UUID]sqlc.Workflow),
		ts:  make(map[uuid.UUID][]sqlc.Transition),
	}
}

func (f *fakeDBQuery) CreateWorkflow(_ context.Context, _ sqlc.DBTX, arg sqlc.CreateWorkflowParams) (sqlc.Workflow, error) {
	w := sqlc.Workflow{
		ID:          arg.ID,
		CurrentNode: arg.CurrentNode,
		Status:      arg.Status,
		Graph:       arg.Graph,
	}
	f.wfs[arg.ID] = w
	return w, nil
}

func (f *fakeDBQuery) GetWorkflow(_ context.Context, _ sqlc.DBTX, id uuid.UUID) (sqlc.Workflow, error) {
	w, ok := f.wfs[id]
	if !ok {
		return sqlc.Workflow{}, errors.New("not found")
	}
	return w, nil
}

func (f *fakeDBQuery) GetNextWorkflows(_ context.Context, _ sqlc.DBTX) ([]sqlc.Workflow, error) {
	return nil, nil
}

func (f *fakeDBQuery) UpdateWorkflowNextAction(_ context.Context, _ sqlc.DBTX, arg sqlc.UpdateWorkflowNextActionParams) (sqlc.Workflow, error) {
	w, ok := f.wfs[arg.ID]
	if !ok {
		return sqlc.Workflow{}, errors.New("not found")
	}
	w.CurrentNode = arg.CurrentNode
	f.wfs[arg.ID] = w

	return w, nil
}

func (f *fakeDBQuery) UpdateWorkflowStatus(_ context.Context, _ sqlc.DBTX, arg sqlc.UpdateWorkflowStatusParams) (sqlc.Workflow, error) {
	w, ok := f.wfs[arg.ID]
	if !ok {
		return sqlc.Workflow{}, errors.New("not found")
	}
	w.Status = arg.Status
	f.wfs[arg.ID] = w

	return w, nil
}

func (f *fakeDBQuery) UpdateWorkflowNextActionAt(_ context.Context, _ sqlc.DBTX, arg sqlc.UpdateWorkflowNextActionAtParams) (sqlc.Workflow, error) {
	return sqlc.Workflow{}, nil
}

func (f *fakeDBQuery) CreateTransition(_ context.Context, _ sqlc.DBTX, arg sqlc.CreateTransitionParams) (sqlc.Transition, error) {
	tn := sqlc.Transition{
		ID:         uuid.New(),
		WorkflowID: arg.WorkflowID,
		FromNode:   arg.FromNode,
		ToNode:     arg.ToNode,
		Label:      arg.Label,
	}
	if len(f.ts[arg.WorkflowID]) > 0 {
		tn.Previous = uuid.NullUUID{
			UUID:  f.ts[arg.WorkflowID][len(f.ts[arg.WorkflowID])-1].ID,
			Valid: true,
		}
		f.ts[arg.WorkflowID][len(f.ts[arg.WorkflowID])-1].Next = uuid.NullUUID{
			UUID:  tn.ID,
			Valid: true,
		}
	}
	f.ts[arg.WorkflowID] = append(f.ts[arg.WorkflowID], tn)
	return tn, nil
}

func (f *fakeDBQuery) UpdateTransitionNext(_ context.Context, _ sqlc.DBTX, arg sqlc.UpdateTransitionNextParams) (sqlc.Transition, error) {
	return sqlc.Transition{}, nil
}

func (f *fakeDBQuery) GetLatestTransition(_ context.Context, _ sqlc.DBTX, workflowID uuid.UUID) ([]sqlc.Transition, error) {
	ts, ok := f.ts[workflowID]
	if !ok {
		return nil, nil
	}
	if len(ts) == 0 {
		return nil, nil
	}

	return []sqlc.Transition{ts[len(ts)-1]}, nil
}

func (f *fakeDBQuery) GetFirstTransition(_ context.Context, _ sqlc.DBTX, workflowID uuid.UUID) ([]sqlc.Transition, error) {
	ts, ok := f.ts[workflowID]
	if !ok {
		return nil, nil
	}
	if len(ts) == 0 {
		return nil, nil
	}

	return []sqlc.Transition{ts[0]}, nil
}

func (f *fakeDBQuery) GetTransitions(_ context.Context, _ sqlc.DBTX, workflowID uuid.UUID) ([]sqlc.Transition, error) {
	return nil, nil
}

func TestCreateWorkflow(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
	`
	gpb := &pb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), gpb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &model.Graph{}
	if err := g.FromProto(gpb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}

	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	ctx := context.Background()
	ws := NewWorkflowStorage(newFakeDBQuery(), &fakeDBTX{})
	want := model.Workflow{
		ID:          id,
		CurrentNode: "a",
		Status:      model.WorkflowStatusPending,
		Graph:       g,
	}
	if err := ws.CreateWorkflow(ctx, want); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := ws.GetWorkflow(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID {
		t.Errorf("got id %v, want %v", got.ID, want.ID)
	}
}

func TestWorkflowNotFound(t *testing.T) {
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	ctx := context.Background()
	ws := NewWorkflowStorage(newFakeDBQuery(), &fakeDBTX{})
	_, err := ws.GetWorkflow(ctx, id)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateWorkflow(t *testing.T) {
	txt := `
	start: "a"
	nodes { id: "a" }
	nodes { id: "b" }
	nodes { id: "c" }
	edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
	edges { from_id: "b" to_id: "c" condition { operation_result: "ok" } }
	`
	gpb := &pb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), gpb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &model.Graph{}
	if err := g.FromProto(gpb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}

	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	ctx := context.Background()
	dbtx := &fakeDBTX{}
	wfstore := NewWorkflowStorage(newFakeDBQuery(), dbtx)
	w := model.Workflow{
		ID:          id,
		CurrentNode: "a",
		Status:      model.WorkflowStatusPending,
		Graph:       g,
	}
	if err := wfstore.CreateWorkflow(ctx, w); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tid1, err := wfstore.GetLatestTransition(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tid1 != uuid.Nil {
		t.Fatalf("unexpected transition id: %v", tid1)
	}

	err = wfstore.updateNextActionTxn(ctx, dbtx, NextAction{
		ID:                id,
		Label:             "ok",
		CurrentTransition: tid1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tid2, err := wfstore.GetLatestTransition(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tid2 == uuid.Nil {
		t.Fatalf("unexpected transition id: %v", tid2)
	}

	// Check that we can't update the same transition twice; optimistic lock rejection.
	err = wfstore.updateNextActionTxn(ctx, dbtx, NextAction{
		ID:                id,
		Label:             "ok",
		CurrentTransition: tid1,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	err = wfstore.updateNextActionTxn(ctx, dbtx, NextAction{
		ID:                id,
		Label:             "ok",
		CurrentTransition: tid2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	finalw, err := wfstore.GetWorkflow(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if finalw.CurrentNode != "c" {
		t.Errorf("unexpected current node: %v", finalw.CurrentNode)
	}
}
