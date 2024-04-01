package graphviz

import (
	"fmt"
	"strings"

	"github.com/ilnar/wf/internal/model"
)

type edgeKey struct {
	from, to, label string
}

// Dot returns a Graphviz representation of the workflow.
// TODO: Switch to a 3rd party labirary.
func Dot(w model.Workflow, t *model.Transition) (string, error) {
	keySec := make([]edgeKey, 0)
	counts := make(map[edgeKey]int)

	// TODO: Use proper BFS here to render edges in topolofical order.
	g, err := w.Graph.ToProto()
	if err != nil {
		return "", fmt.Errorf("failed to convert graph to proto: %w", err)
	}
	for _, e := range g.Edges {
		key := edgeKey{from: e.FromId, to: e.ToId, label: e.Condition.OperationResult}
		keySec = append(keySec, key)
		counts[key] = 0
	}

	// Increase counts of executed transitions.
	th := t
	for th != nil {
		key := edgeKey{label: th.Label()}
		key.from, key.to = th.FromTo()
		counts[key]++
		th = th.Next()
	}

	var sb strings.Builder
	sb.WriteString("digraph G { ")
	for _, e := range keySec {
		sb.WriteString(fmt.Sprintf("%s -> %s [label=\"%s %d\"]; ", e.from, e.to, e.label, counts[e]))
	}
	sb.WriteString("}")
	return sb.String(), nil
}
