-- name: CreateWorkflow :one
INSERT INTO workflows (
  current_node,
  status,
  graph
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetWorkflow :one
SELECT * FROM workflows
WHERE id = $1 LIMIT 1;
