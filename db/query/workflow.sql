-- name: CreateWorkflow :one
INSERT INTO workflows (
  id,
  current_node,
  status,
  graph
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateWorkflowStatus :one
UPDATE workflows
SET status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateWorkflowNextAction :one
UPDATE workflows
SET current_node = $2, next_action_at = now() + interval '2 seconds'
WHERE id = $1
RETURNING *;

-- name: UpdateWorkflowNextActionAt :one
UPDATE workflows
SET next_action_at = now() + sqlc.arg(delay)::bigint * interval '1 second'
WHERE id = $1
RETURNING *;

-- name: GetWorkflow :one
SELECT * FROM workflows
WHERE id = $1 LIMIT 1;

-- name: GetNextWorkflows :many
SELECT * FROM workflows
WHERE status = 'running'
  AND next_action_at <= now()
LIMIT 10;

-- name: CreateTransition :one
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
RETURNING *;

-- name: GetLatestTransition :many
SELECT * FROM transitions
WHERE workflow_id = $1 AND "next" IS NULL;

-- name: GetFirstTransition :many
SELECT * FROM transitions
WHERE workflow_id = $1 AND previous IS NULL;

-- name: UpdateTransitionNext :one
UPDATE transitions
SET "next" = $2
WHERE id = $1
RETURNING *;