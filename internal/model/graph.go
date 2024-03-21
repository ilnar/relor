package model

import (
	"fmt"

	gpb "github.com/ilnar/wf/gen/pb/graph"
)

type Node struct {
	ID, Name string
	In, Out  []*Edge
}

type Edge struct {
	From, To *Node
}

type Graph struct {
	root *Node
	idx  map[string]*Node
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
		e := &Edge{From: from, To: to}
		from.Out = append(from.Out, e)
		to.In = append(to.In, e)
	}

	return nil
}

func (g *Graph) NextNodeIDs(id string) ([]string, error) {
	if g == nil || g.idx == nil {
		return nil, fmt.Errorf("graph is not initialized")
	}
	n, ok := g.idx[id]
	if !ok {
		return nil, fmt.Errorf("node not found: %s", id)
	}
	next := make([]string, 0, len(n.Out))
	for _, e := range n.Out {
		next = append(next, e.To.ID)
	}
	return next, nil
}

func (g *Graph) Head() string {
	if g == nil || g.root == nil {
		return ""
	}
	return g.root.ID
}

func indexNodes(g *gpb.Graph) (map[string]*Node, error) {
	i := make(map[string]*Node, len(g.Nodes))
	for _, n := range g.Nodes {
		i[n.Id] = &Node{ID: n.Id, Name: n.Name}
	}
	return i, nil
}
