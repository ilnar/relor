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

-- name: GetWorkflow :one
SELECT * FROM workflows
WHERE id = $1 LIMIT 1;

-- name: GetNextWorkflows :many
SELECT * FROM workflows
WHERE status = 'running'
  AND next_action_at <= now()
LIMIT 10;
