package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"scheduler/database"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func buildReq(task database.Task) (*http.Request, error) {
	var headers map[string]string
	json.Unmarshal(task.ActionHeaders, &headers)
	var payload []byte
	json.Unmarshal(task.ActionPayload, &payload)

	req, err := http.NewRequest(task.ActionMethod, task.ActionUrl, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func getResponse(req *http.Request) (*http.Response, time.Duration, error) {
	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	return resp, duration, err
}

func (wp *WorkerPool) saveResult(ctx context.Context, task database.Task, resp *http.Response, duration time.Duration, taskErr error) {
	var statusCode int32
	var success bool
	var responseHeaders json.RawMessage
	var responseBody json.RawMessage
	var errorMessage string

	if taskErr != nil {
		statusCode = 0
		success = false
		errorMessage = taskErr.Error()
	} else {
		statusCode = int32(resp.StatusCode)
		success = resp.StatusCode >= 200 && resp.StatusCode < 300

		hdrs, _ := json.Marshal(resp.Header)
		responseHeaders = hdrs

		rawBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var tmp interface{}
		if err := json.Unmarshal(rawBody, &tmp); err != nil {
			responseBody, _ = json.Marshal(string(rawBody))
		} else {
			responseBody = rawBody
		}

	}

	_, dbErr := wp.db.CreateTaskResult(ctx, database.CreateTaskResultParams{
		TaskID:          task.ID,
		RunAt:           pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		StatusCode:      statusCode,
		Success:         success,
		ResponseHeaders: []byte(string(responseHeaders)),
		ResponseBody:    []byte(string(responseBody)),
		ErrorMessage:    pgtype.Text{String: errorMessage, Valid: errorMessage != ""},
		DurationMs:      int32(duration.Milliseconds()),
	})

	if dbErr != nil {
		log.Printf("Failed to save task result for task %s: %v", task.Name, dbErr)
		return
	}

	_, err := wp.db.UpdateTaskStatus(ctx, database.UpdateTaskStatusParams{
		ID:     task.ID,
		Status: "completed",
	})
	if err != nil {
		log.Printf("Failed to update task %s status: %v", task.Name, err)
		return
	}

	log.Printf("Task %s completed (success=%v)", task.Name, success)
}

func (wp *WorkerPool) executeTask(ctx context.Context, task database.Task) {
	log.Printf("Executing task: %s [%s %s]", task.Name, task.ActionMethod, task.ActionUrl)

	req, err := buildReq(task)
	if err != nil {
		log.Printf("Failed to build request for task %s: %v", task.Name, err)
		wp.saveResult(ctx, task, nil, 0, err)
		return
	}
	resp, duration, err := getResponse(req)

	wp.saveResult(ctx, task, resp, duration, err)
}
