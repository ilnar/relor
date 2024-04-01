package graphviz

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/ilnar/wf/internal/model"
)

type edgeKey struct {
	from, to, label string
}

const dotTemplate = `digraph G {

node [fontname="Helvetica,Arial,sans-serif" shape="box" style="rounded,filled" fillcolor="white"];
edge [fontname="Helvetica,Arial,sans-serif" fontsize="10"];

Start{{- range .VisitedNodes }}, {{ . }} {{- end }} [fillcolor="grey"];

"{{ .CurrentNode }}" [fillcolor="magenta" color="magenta" fontcolor="white"];

Start -> "{{ .StartNode }}";
{{- range .Edges }}
"{{ .From }}" -> "{{ .To }}" [label="{{ .Label }}" weight={{ .Weight }}{{ if eq .Weight 0 }} color="grey"{{ end }}];
{{- end }}

}`

type edgeData struct {
	From, To, Label string
	Weight          int
}

type dotData struct {
	VisitedNodes           []string
	StartNode, CurrentNode string
	Edges                  []edgeData
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

	// Get seen nodes.
	seenNodes := make(map[string]struct{})
	seenNodes[w.Graph.Head()] = struct{}{}
	for k, cnt := range counts {
		if cnt > 0 {
			seenNodes[k.from] = struct{}{}
			seenNodes[k.to] = struct{}{}
		}
	}
	delete(seenNodes, w.CurrentNode)

	//Prepare template data.
	data := dotData{
		VisitedNodes: make([]string, 0, len(seenNodes)),
		StartNode:    w.Graph.Head(),
		CurrentNode:  w.CurrentNode,
		Edges:        make([]edgeData, 0, len(counts)),
	}
	for _, key := range keySec {
		label := key.label
		if counts[key] > 0 {
			label += fmt.Sprintf(" (%d)", counts[key])
		}
		data.Edges = append(data.Edges, edgeData{
			From:   key.from,
			To:     key.to,
			Label:  label,
			Weight: counts[key],
		})
	}
	for k := range seenNodes {
		data.VisitedNodes = append(data.VisitedNodes, k)
	}

	// Render the template.
	tmpl, err := template.New("dot").Parse(dotTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return out.String(), nil
}
