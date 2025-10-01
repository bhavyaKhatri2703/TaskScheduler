package entity

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTaskReq struct {
	Name    string      `json:"name"`
	Trigger TriggerData `json:"trigger"`
	Action  ActionData  `json:"action"`
}

type TaskResponse struct {
	ID        pgtype.UUID `json:"id"`
	Name      string      `json:"name"`
	Trigger   TriggerData `json:"trigger"`
	Action    ActionData  `json:"action"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	NextRun   *time.Time  `json:"next_run"`
}

type UpdateTaskRequest struct {
	Name    *string      `json:"name"`
	Trigger *TriggerData `json:"trigger"`
	Action  *ActionData  `json:"action"`
}

type ListTasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

type TaskResultResponse struct {
	ID              pgtype.UUID            `json:"id"`
	TaskID          pgtype.UUID            `json:"task_id"`
	RunAt           time.Time              `json:"run_at"`
	StatusCode      int32                  `json:"status_code"`
	Success         bool                   `json:"success"`
	ResponseHeaders map[string]interface{} `json:"response_headers,omitempty"`
	ResponseBody    interface{}            `json:"response_body,omitempty"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	DurationMs      int32                  `json:"duration_ms"`
	CreatedAt       time.Time              `json:"created_at"`
}
