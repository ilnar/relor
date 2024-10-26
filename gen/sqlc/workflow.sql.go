// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workflow.sql

package sqlc

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

const createTransition = `-- name: CreateTransition :one
INSERT INTO transitions (
  workflow_id,
  from_node,
  to_node,
  label,
  previous,
  "next"
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id, workflow_id, from_node, to_node, label, created_at, previous, next
`

type CreateTransitionParams struct {
	WorkflowID uuid.UUID     `json:"workflow_id"`
	FromNode   string        `json:"from_node"`
	ToNode     string        `json:"to_node"`
	Label      string        `json:"label"`
	Previous   uuid.NullUUID `json:"previous"`
	Next       uuid.NullUUID `json:"next"`
}

func (q *Queries) CreateTransition(ctx context.Context, db DBTX, arg CreateTransitionParams) (Transition, error) {
	row := db.QueryRowContext(ctx, createTransition,
		arg.WorkflowID,
		arg.FromNode,
		arg.ToNode,
		arg.Label,
		arg.Previous,
		arg.Next,
	)
	var i Transition
	err := row.Scan(
		&i.ID,
		&i.WorkflowID,
		&i.FromNode,
		&i.ToNode,
		&i.Label,
		&i.CreatedAt,
		&i.Previous,
		&i.Next,
	)
	return i, err
}

const createWorkflow = `-- name: CreateWorkflow :one
INSERT INTO workflows (
  id,
  current_node,
  status,
  graph
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, current_node, status, graph, created_at, next_action_at
`

type CreateWorkflowParams struct {
	ID          uuid.UUID       `json:"id"`
	CurrentNode string          `json:"current_node"`
	Status      string          `json:"status"`
	Graph       json.RawMessage `json:"graph"`
}

func (q *Queries) CreateWorkflow(ctx context.Context, db DBTX, arg CreateWorkflowParams) (Workflow, error) {
	row := db.QueryRowContext(ctx, createWorkflow,
		arg.ID,
		arg.CurrentNode,
		arg.Status,
		arg.Graph,
	)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.CurrentNode,
		&i.Status,
		&i.Graph,
		&i.CreatedAt,
		&i.NextActionAt,
	)
	return i, err
}

const getFirstTransition = `-- name: GetFirstTransition :many
SELECT id, workflow_id, from_node, to_node, label, created_at, previous, next FROM transitions
WHERE workflow_id = $1 AND previous IS NULL
`

func (q *Queries) GetFirstTransition(ctx context.Context, db DBTX, workflowID uuid.UUID) ([]Transition, error) {
	rows, err := db.QueryContext(ctx, getFirstTransition, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transition
	for rows.Next() {
		var i Transition
		if err := rows.Scan(
			&i.ID,
			&i.WorkflowID,
			&i.FromNode,
			&i.ToNode,
			&i.Label,
			&i.CreatedAt,
			&i.Previous,
			&i.Next,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLatestTransition = `-- name: GetLatestTransition :many
SELECT id, workflow_id, from_node, to_node, label, created_at, previous, next FROM transitions
WHERE workflow_id = $1 AND "next" IS NULL
`

func (q *Queries) GetLatestTransition(ctx context.Context, db DBTX, workflowID uuid.UUID) ([]Transition, error) {
	rows, err := db.QueryContext(ctx, getLatestTransition, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transition
	for rows.Next() {
		var i Transition
		if err := rows.Scan(
			&i.ID,
			&i.WorkflowID,
			&i.FromNode,
			&i.ToNode,
			&i.Label,
			&i.CreatedAt,
			&i.Previous,
			&i.Next,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextWorkflows = `-- name: GetNextWorkflows :many
SELECT id, current_node, status, graph, created_at, next_action_at FROM workflows
WHERE status = 'running'
  AND next_action_at <= now()
LIMIT 10
`

func (q *Queries) GetNextWorkflows(ctx context.Context, db DBTX) ([]Workflow, error) {
	rows, err := db.QueryContext(ctx, getNextWorkflows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Workflow
	for rows.Next() {
		var i Workflow
		if err := rows.Scan(
			&i.ID,
			&i.CurrentNode,
			&i.Status,
			&i.Graph,
			&i.CreatedAt,
			&i.NextActionAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTransitions = `-- name: GetTransitions :many
SELECT id, workflow_id, from_node, to_node, label, created_at, previous, next FROM transitions
WHERE workflow_id = $1
`

func (q *Queries) GetTransitions(ctx context.Context, db DBTX, workflowID uuid.UUID) ([]Transition, error) {
	rows, err := db.QueryContext(ctx, getTransitions, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transition
	for rows.Next() {
		var i Transition
		if err := rows.Scan(
			&i.ID,
			&i.WorkflowID,
			&i.FromNode,
			&i.ToNode,
			&i.Label,
			&i.CreatedAt,
			&i.Previous,
			&i.Next,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkflow = `-- name: GetWorkflow :one
SELECT id, current_node, status, graph, created_at, next_action_at FROM workflows
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetWorkflow(ctx context.Context, db DBTX, id uuid.UUID) (Workflow, error) {
	row := db.QueryRowContext(ctx, getWorkflow, id)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.CurrentNode,
		&i.Status,
		&i.Graph,
		&i.CreatedAt,
		&i.NextActionAt,
	)
	return i, err
}

const updateTransitionNext = `-- name: UpdateTransitionNext :one
UPDATE transitions
SET "next" = $2
WHERE id = $1
RETURNING id, workflow_id, from_node, to_node, label, created_at, previous, next
`

type UpdateTransitionNextParams struct {
	ID   uuid.UUID     `json:"id"`
	Next uuid.NullUUID `json:"next"`
}

func (q *Queries) UpdateTransitionNext(ctx context.Context, db DBTX, arg UpdateTransitionNextParams) (Transition, error) {
	row := db.QueryRowContext(ctx, updateTransitionNext, arg.ID, arg.Next)
	var i Transition
	err := row.Scan(
		&i.ID,
		&i.WorkflowID,
		&i.FromNode,
		&i.ToNode,
		&i.Label,
		&i.CreatedAt,
		&i.Previous,
		&i.Next,
	)
	return i, err
}

const updateWorkflowNextAction = `-- name: UpdateWorkflowNextAction :one
UPDATE workflows
SET current_node = $2, next_action_at = now() + interval '2 seconds'
WHERE id = $1
RETURNING id, current_node, status, graph, created_at, next_action_at
`

type UpdateWorkflowNextActionParams struct {
	ID          uuid.UUID `json:"id"`
	CurrentNode string    `json:"current_node"`
}

func (q *Queries) UpdateWorkflowNextAction(ctx context.Context, db DBTX, arg UpdateWorkflowNextActionParams) (Workflow, error) {
	row := db.QueryRowContext(ctx, updateWorkflowNextAction, arg.ID, arg.CurrentNode)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.CurrentNode,
		&i.Status,
		&i.Graph,
		&i.CreatedAt,
		&i.NextActionAt,
	)
	return i, err
}

const updateWorkflowNextActionAt = `-- name: UpdateWorkflowNextActionAt :one
UPDATE workflows
SET next_action_at = now() + $2::bigint * interval '1 second'
WHERE id = $1
RETURNING id, current_node, status, graph, created_at, next_action_at
`

type UpdateWorkflowNextActionAtParams struct {
	ID    uuid.UUID `json:"id"`
	Delay int64     `json:"delay"`
}

func (q *Queries) UpdateWorkflowNextActionAt(ctx context.Context, db DBTX, arg UpdateWorkflowNextActionAtParams) (Workflow, error) {
	row := db.QueryRowContext(ctx, updateWorkflowNextActionAt, arg.ID, arg.Delay)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.CurrentNode,
		&i.Status,
		&i.Graph,
		&i.CreatedAt,
		&i.NextActionAt,
	)
	return i, err
}

const updateWorkflowStatus = `-- name: UpdateWorkflowStatus :one
UPDATE workflows
SET status = $2
WHERE id = $1
RETURNING id, current_node, status, graph, created_at, next_action_at
`

type UpdateWorkflowStatusParams struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func (q *Queries) UpdateWorkflowStatus(ctx context.Context, db DBTX, arg UpdateWorkflowStatusParams) (Workflow, error) {
	row := db.QueryRowContext(ctx, updateWorkflowStatus, arg.ID, arg.Status)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.CurrentNode,
		&i.Status,
		&i.Graph,
		&i.CreatedAt,
		&i.NextActionAt,
	)
	return i, err
}
