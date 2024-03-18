package model

import "github.com/google/uuid"

type WorkflowStatus string

const (
	WorkflowStatusUnknown   WorkflowStatus = "unknown"
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
)

type Workflow struct {
	ID          uuid.UUID
	CurrentNode string
	Status      WorkflowStatus
	Graph       *Graph
}

func NewWorkflow(id uuid.UUID, g *Graph) *Workflow {
	return &Workflow{
		ID:          id,
		Status:      WorkflowStatusPending,
		CurrentNode: g.Head(),
		Graph:       g,
	}
}
