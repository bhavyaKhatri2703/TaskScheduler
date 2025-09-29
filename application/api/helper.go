package api

import (
	"encoding/json"
	"log"
	entity "scheduler/application/entity"
	"scheduler/database"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func taskToResponse(task database.Task) (entity.TaskResponse, error) {

	trigger := entity.TriggerData{
		Type:     task.TriggerType,
		DateTime: task.TriggerDatetime.Time.Format(time.RFC3339),
		Cron:     task.TriggerCron.String,
	}
	var headers map[string]string

	err := json.Unmarshal(task.ActionHeaders, &headers)
	if err != nil {
		log.Printf("failed to unmarshal headers: %v", err)
		return entity.TaskResponse{}, err
	}

	action := entity.ActionData{
		Method:  task.ActionMethod,
		URL:     task.ActionUrl,
		Headers: headers,
		Payload: task.ActionPayload,
	}

	return entity.TaskResponse{
		ID:        task.ID,
		Name:      task.Name,
		Trigger:   trigger,
		Action:    action,
		Status:    task.Status,
		CreatedAt: task.CreatedAt.Time,
		UpdatedAt: task.UpdatedAt.Time,
		NextRun:   &task.NextRun.Time,
	}, nil
}

func StringToTimestamptz(s string) (pgtype.Timestamptz, error) {
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return pgtype.Timestamptz{}, err
	}
	return pgtype.Timestamptz{Time: parsedTime, Valid: true}, nil
}

func StringToPgText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}
