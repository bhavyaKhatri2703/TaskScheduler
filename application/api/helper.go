package api

import (
	"encoding/json"
	"log"
	entity "scheduler/application/entity"
	"scheduler/database"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robfig/cron/v3"
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

func NextCronTime(expr string) (*time.Time, error) {
	schedule, err := cron.ParseStandard(expr)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	next := schedule.Next(now)
	return &next, nil
}

func taskResultToResponse(result database.TaskResult) (entity.TaskResultResponse, error) {
	response := entity.TaskResultResponse{
		ID:         result.ID,
		TaskID:     result.TaskID,
		RunAt:      result.RunAt.Time,
		StatusCode: result.StatusCode,
		Success:    result.Success,
		DurationMs: result.DurationMs,
		CreatedAt:  result.CreatedAt.Time,
	}

	if result.ResponseHeaders != nil {
		var headers map[string]interface{}
		if err := json.Unmarshal(result.ResponseHeaders, &headers); err == nil {
			response.ResponseHeaders = headers
		}
	}
	if result.ResponseBody != nil {
		var body interface{}
		if err := json.Unmarshal(result.ResponseBody, &body); err == nil {
			response.ResponseBody = body
		}
	}

	if result.ErrorMessage.Valid {
		response.ErrorMessage = result.ErrorMessage.String
	}

	return response, nil
}
