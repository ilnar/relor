package gossip

import (
	"context"

	"github.com/hashicorp/memberlist"
)

type eventHandler struct {
	l Logger
}

func (e *eventHandler) NotifyJoin(n *memberlist.Node) {
	e.l.InfoContext(context.Background(), "Node joined", "hostname", n.Addr, "port", n.Port)
}

func (e *eventHandler) NotifyLeave(n *memberlist.Node) {
	e.l.InfoContext(context.Background(), "Node left", "hostname", n.Addr, "port", n.Port)
}

func (e *eventHandler) NotifyUpdate(n *memberlist.Node) {
	e.l.InfoContext(context.Background(), "Node updated", "hostname", n.Addr, "port", n.Port)
}
