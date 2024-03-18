package model

import gpb "github.com/ilnar/wf/gen/pb/graph"

type Graph struct {
}

func (g *Graph) FromProto(pb *gpb.Graph) error {
	if g == nil {
		return nil
	}

	return nil
}

func (g *Graph) Head() string {
	return "start"
}
