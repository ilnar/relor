-- name: CreateWorkflow :one
INSERT INTO workflows (
  current_node,
  status,
  graph
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetWorflow :one
SELECT * FROM workflows
WHERE id = $1 LIMIT 1;
