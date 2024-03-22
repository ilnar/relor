package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ilnar/wf/gen/sqlc"
	"github.com/ilnar/wf/internal/model"
	"google.golang.org/protobuf/encoding/prototext"

	pb "github.com/ilnar/wf/gen/pb/graph"
)

type fakeDBQuery struct {
	wfs map[uuid.UUID]sqlc.Workflow
}

func newFakeDBQuery() *fakeDBQuery {
	return &fakeDBQuery{
		wfs: make(map[uuid.UUID]sqlc.Workflow),
	}
}

func (f *fakeDBQuery) CreateWorkflow(_ context.Context, arg sqlc.CreateWorkflowParams) (sqlc.Workflow, error) {
	w := sqlc.Workflow{
		ID:          arg.ID,
		CurrentNode: arg.CurrentNode,
		Status:      arg.Status,
		Graph:       arg.Graph,
	}
	f.wfs[arg.ID] = w
	return w, nil
}

func (f *fakeDBQuery) GetWorkflow(_ context.Context, id uuid.UUID) (sqlc.Workflow, error) {
	w, ok := f.wfs[id]
	if !ok {
		return sqlc.Workflow{}, errors.New("not found")
	}
	return w, nil
}

func TestCreateWorkflow(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" name: "node a" }
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
	ws := NewWorkflowStorage(newFakeDBQuery())
	want := &model.Workflow{
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
	ws := NewWorkflowStorage(newFakeDBQuery())
	_, err := ws.GetWorkflow(ctx, id)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
