package model

import (
	"testing"

	"google.golang.org/protobuf/encoding/prototext"

	gpb "github.com/ilnar/wf/gen/pb/graph"
)

func TestLoadGraphFromProtoWideGraph(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" name: "node a" }
		nodes { id: "b" name: "node b" }
		nodes { id: "c" name: "node c" }
		edges { from_id: "a" to_id: "b" }
		edges { from_id: "a" to_id: "c" }
	`
	pb := &gpb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), pb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &Graph{}
	if err := g.FromProto(pb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}
	if g.Head() != "a" {
		t.Errorf("unexpected head: %s; want a", g.Head())
	}
	gotNodes, err := g.NextNodeIDs("a")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	wantNodes := []string{"b", "c"}
	if len(gotNodes) != len(wantNodes) {
		t.Errorf("unexpected next nodes: %v; want %v", gotNodes, wantNodes)
	}
	for i, n := range gotNodes {
		if n != wantNodes[i] {
			t.Errorf("unexpected next node: %s; want %s", n, wantNodes[i])
		}
	}
	gotNodes, err = g.NextNodeIDs("b")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	if len(gotNodes) != 0 {
		t.Errorf("unexpected next node count: %d; want 0", len(gotNodes))
	}
}

func TestLoadGraphFromProtoLongGraph(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" name: "node a" }
		nodes { id: "b" name: "node b" }
		nodes { id: "c" name: "node c" }
		nodes { id: "d" name: "node d" }
		nodes { id: "e" name: "node e" }
		edges { from_id: "a" to_id: "b" }
		edges { from_id: "b" to_id: "c" }
		edges { from_id: "c" to_id: "d" }
		edges { from_id: "d" to_id: "e" }
	`
	pb := &gpb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), pb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &Graph{}
	if err := g.FromProto(pb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}
	if g.Head() != "a" {
		t.Errorf("unexpected head: %s; want a", g.Head())
	}
	gotNodes, err := g.NextNodeIDs("a")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	if len(gotNodes) != 1 {
		t.Errorf("unexpected next node count: %d; want 1", len(gotNodes))
	}
	if gotNodes[0] != "b" {
		t.Errorf("unexpected next node: %s; want b", gotNodes[0])
	}

	gotNodes, err = g.NextNodeIDs("b")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	if len(gotNodes) != 1 {
		t.Errorf("unexpected next node count: %d; want 1", len(gotNodes))
	}
	if gotNodes[0] != "c" {
		t.Errorf("unexpected next node: %s; want c", gotNodes[0])
	}

	gotNodes, err = g.NextNodeIDs("e")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	if len(gotNodes) != 0 {
		t.Errorf("unexpected next node count: %d; want 0", len(gotNodes))
	}
}

func TestLoadGraphFromProtoCorruptedGraph(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" name: "node a" }
		nodes { id: "b" name: "node b" }
		nodes { id: "c" name: "node c" }
		edges { from_id: "a" to_id: "b" }
		edges { from_id: "b" to_id: "c" }
		edges { from_id: "c" to_id: "d" }
	`
	pb := &gpb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), pb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &Graph{}
	if err := g.FromProto(pb); err == nil {
		t.Fatalf("unexpected success loading corrupted graph")
	}
}

func TestLoadGraphFromProtoUnitialised(t *testing.T) {
	var g *Graph
	if err := g.FromProto(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Head() != "" {
		t.Errorf("unexpected head: %s; want empty", g.Head())
	}
}
