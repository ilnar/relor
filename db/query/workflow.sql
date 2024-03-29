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
SET next_action_at = $2, current_node = $3
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
