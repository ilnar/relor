package model

import (
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	gpb "github.com/ilnar/wf/gen/pb/graph"
)

func TestLoadGraphFromProtoWideGraph(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
		nodes { id: "b" }
		nodes { id: "c" }
		edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
		edges { from_id: "a" to_id: "c" condition { operation_result: "not_ok" } }
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
	got, err := g.NextNodeID("a", "ok")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	want := "b"
	if got != want {
		t.Errorf("wrong node, got %q; want %q", got, want)
	}
	got, err = g.NextNodeID("a", "not_ok")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	want = "c"
	if got != want {
		t.Errorf("wrong node, got %q; want %q", got, want)
	}
}

func TestLoadGraphFromProtoLongGraph(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
		nodes { id: "b" }
		nodes { id: "c" }
		nodes { id: "d" }
		nodes { id: "e" }
		edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
		edges { from_id: "b" to_id: "c" condition { operation_result: "ok" } }
		edges { from_id: "c" to_id: "d" condition { operation_result: "ok" } }
		edges { from_id: "d" to_id: "e" condition { operation_result: "ok" } }
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
	got, err := g.NextNodeID("a", "ok")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	if got != "b" {
		t.Errorf("unexpected next node: %s; want b", got)
	}

	got, err = g.NextNodeID("b", "ok")
	if err != nil {
		t.Fatalf("failed to get next nodes: %v", err)
	}
	if got != "c" {
		t.Errorf("unexpected next node: %s; want c", got)
	}
}

func TestLoadGraphFromProtoCorruptedGraph(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
		nodes { id: "b" }
		nodes { id: "c" }
		edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
		edges { from_id: "b" to_id: "c" condition { operation_result: "ok" } }
		edges { from_id: "c" to_id: "d" condition { operation_result: "ok" } }
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

func TestGraphToProto(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
		nodes { id: "b" }
		nodes { id: "c" }
		edges { from_id: "a" to_id: "b" condition { operation_result: "res1" } }
		edges { from_id: "a" to_id: "c" condition { operation_result: "res2" } }
	`
	pb := &gpb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), pb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &Graph{}
	if err := g.FromProto(pb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}
	got, err := g.ToProto()
	if err != nil {
		t.Fatalf("failed to convert to proto: %v", err)
	}
	if diff := cmp.Diff(pb, got, protocmp.Transform(), protocmp.SortRepeatedFields(pb, "nodes", "edges")); diff != "" {
		t.Errorf("unexpected difference: %v", diff)
	}

}

func TestGraphDuplicateNodes(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
		nodes { id: "a" }
	`
	pb := &gpb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), pb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &Graph{}
	if err := g.FromProto(pb); err == nil {
		t.Fatalf("unexpected success loading graph with duplicate nodes")
	}
}

func TestGraphDuplicateEdges(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
		nodes { id: "b" }
		edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
		edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
	`
	pb := &gpb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), pb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &Graph{}
	if err := g.FromProto(pb); err == nil {
		t.Fatalf("unexpected success loading graph with duplicate nodes")
	}
}

func TestGetOutLabels(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
		nodes { id: "b" }
		nodes { id: "c" }
		edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
		edges { from_id: "a" to_id: "c" condition { operation_result: "not_ok" } }
	`
	pb := &gpb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), pb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &Graph{}
	if err := g.FromProto(pb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}
	got, err := g.OutLabels("a")
	if err != nil {
		t.Fatalf("failed to get out labels: %v", err)
	}
	want := []string{"ok", "not_ok"}
	copts := cmpopts.SortSlices(func(a, b string) bool { return a < b })
	if diff := cmp.Diff(want, got, copts); diff != "" {
		t.Errorf("unexpected difference: %v", diff)
	}
}

func TestGetOutLabelsAtLeaf(t *testing.T) {
	txt := `
		start: "a"
		nodes { id: "a" }
	`
	pb := &gpb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), pb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &Graph{}
	if err := g.FromProto(pb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}
	got, err := g.OutLabels("a")
	if err != nil {
		t.Fatalf("failed to get out labels: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("unexpected out labels: %v; want empty", got)
	}
}
