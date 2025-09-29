-- name: CreateTask :one
INSERT INTO tasks (name, trigger_type, trigger_datetime, trigger_cron, action_method, action_url, action_headers, action_payload, status, next_run)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
    trigger_type = COALESCE($3, trigger_type),
    trigger_datetime = COALESCE($4, trigger_datetime),
    trigger_cron = COALESCE($5, trigger_cron),
    action_method = COALESCE($6, action_method),
    action_url = COALESCE($7, action_url),
    action_headers = COALESCE($8, action_headers),
    action_payload = COALESCE($9, action_payload),
    status = COALESCE($10, status),
    next_run = COALESCE($11, next_run),
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: CancelTask :one
UPDATE tasks
SET status = 'cancelled',
    updated_at = now()
WHERE id = $1
RETURNING *;
