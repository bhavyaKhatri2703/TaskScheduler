package entity

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskReq struct {
	Name    string      `json:"name"`
	Trigger TriggerData `json:"trigger"`
	Action  ActionData  `json:"action"`
}

type TaskResponse struct {
	ID        uuid.UUID   `json:"id"`
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
