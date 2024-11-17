package gossip

import (
	"context"
	"fmt"

	"github.com/hashicorp/memberlist"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
}

type Gossip struct {
	mlist *memberlist.Memberlist
	l     Logger
}

func NewGossip(ctx context.Context, name string, port int, seed []string, l Logger) (*Gossip, error) {
	eh := &eventHandler{l: l}
	config := memberlist.DefaultLocalConfig()
	config.BindPort = port
	config.Events = eh
	config.Name = name
	ml, err := memberlist.Create(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create memberlist: %w", err)
	}
	n, err := ml.Join(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to join memberlist: %w", err)
	}
	l.InfoContext(ctx, "Nodes joined", "count", n)
	return &Gossip{mlist: ml, l: l}, nil
}

func (g *Gossip) Members() []string {
	members := g.mlist.Members()
	addrs := make([]string, 0, len(members))
	for _, m := range members {
		s := fmt.Sprintf("%s:%d", m.Addr, m.Port)
		addrs = append(addrs, s)
	}
	return addrs
}
