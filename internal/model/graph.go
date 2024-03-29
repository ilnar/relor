package model

import (
	"fmt"

	gpb "github.com/ilnar/wf/gen/pb/graph"
)

type node struct {
	id, name string
	in, out  map[string]*edge
}

type edge struct {
	from, to *node
	label    string
}

type Graph struct {
	root *node
	idx  map[string]*node
}

func (g *Graph) FromProto(pb *gpb.Graph) error {
	if g == nil {
		return nil
	}

	var err error
	g.idx, err = indexNodes(pb)
	if err != nil {
		return fmt.Errorf("failed to index nodes: %w", err)
	}

	if s, ok := g.idx[pb.Start]; ok {
		g.root = s
	} else {
		return fmt.Errorf("root node is not found: %s", pb.Start)
	}

	for _, n := range pb.Edges {
		from, ok := g.idx[n.FromId]
		if !ok {
			return fmt.Errorf("corrupted graph; node not found: from %s; edge %v", n.FromId, n)
		}
		to, ok := g.idx[n.ToId]
		if !ok {
			return fmt.Errorf("corrupted graph; node not found: to %s; edge %v", n.ToId, n)
		}
		if n.Condition == nil {
			return fmt.Errorf("corrupted graph; edge condition is not set: %v", n)
		}
		if n.Condition.OperationResult == "" {
			return fmt.Errorf("corrupted graph; edge condition operation result is not set: %v", n)
		}
		if _, found := from.out[n.Condition.OperationResult]; found {
			return fmt.Errorf("corrupted graph; duplicate edge: %v", n)
		}
		e := &edge{from: from, to: to, label: n.Condition.OperationResult}
		from.out[n.Condition.OperationResult] = e
		to.in[n.Condition.OperationResult] = e
	}

	return nil
}
func (g *Graph) ToProto() (*gpb.Graph, error) {
	if g == nil || g.idx == nil {
		return nil, fmt.Errorf("graph is not initialized")
	}
	pb := &gpb.Graph{Start: g.root.id}
	for _, n := range g.idx {
		pb.Nodes = append(pb.Nodes, &gpb.Node{Id: n.id, Name: n.name})
		for label, e := range n.out {
			pb.Edges = append(pb.Edges, &gpb.Edge{
				FromId:    e.from.id,
				ToId:      e.to.id,
				Condition: &gpb.TransitionCondition{OperationResult: label},
			},
			)
		}
	}
	return pb, nil
}

func (g *Graph) NextNodeID(nodeID, label string) (string, error) {
	if g == nil || g.idx == nil {
		return "", fmt.Errorf("graph is not initialized")
	}
	n, ok := g.idx[nodeID]
	if !ok {
		return "", fmt.Errorf("node not found: %s", nodeID)
	}
	e, found := n.out[label]
	if !found {
		return "", fmt.Errorf("edge not found: node %q; out label %q", nodeID, label)
	}
	return e.to.id, nil
}

func (g *Graph) Head() string {
	if g == nil || g.root == nil {
		return ""
	}
	return g.root.id
}

func indexNodes(g *gpb.Graph) (map[string]*node, error) {
	i := make(map[string]*node, len(g.Nodes))
	for _, n := range g.Nodes {
		if _, found := i[n.Id]; found {
			return nil, fmt.Errorf("duplicate node id: %s", n.Id)
		}
		i[n.Id] = &node{
			id:   n.Id,
			name: n.Name,
			in:   make(map[string]*edge),
			out:  make(map[string]*edge),
		}
	}
	return i, nil
}
