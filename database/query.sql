-- name: CreateTask :one
INSERT INTO tasks (name, "trigger", action, status, next_run)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE ($1::TEXT IS NULL OR status = $1)
LIMIT $2 OFFSET $3;

-- name: ListTaskResults :many
SELECT * FROM task_results
WHERE task_id = $1;

-- name: UpdateTask :one
UPDATE tasks
SET name = COALESCE($2, name),
    "trigger" = COALESCE($3, "trigger"),
    action = COALESCE($4, action),
    status = COALESCE($5, status),
    next_run = COALESCE($6, next_run),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CancelTask :one
UPDATE tasks
SET status = 'cancelled',
    updated_at = now()
WHERE id = $1
RETURNING *;
