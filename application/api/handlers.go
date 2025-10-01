package api

import (
	"encoding/json"
	"log"
	"net/http"
	entity "scheduler/application/entity"
	"scheduler/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Server) GetTask(c *gin.Context) {
	idParam := c.Param("id")
	var pguuid pgtype.UUID
	err := pguuid.Scan(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := s.DB.GetTask(c, pguuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching task"})
		return
	}

	response, err := taskToResponse(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating response"})
		return
	}

	c.JSON(http.StatusAccepted, response)
}

func (s *Server) CreateTask(c *gin.Context) {
	var req entity.CreateTaskReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dateTime, err := StringToTimestamptz(req.Trigger.DateTime)
	cron := StringToPgText(req.Trigger.Cron)
	reqHeaders, _ := json.Marshal(req.Action.Headers)

	var nextRun pgtype.Timestamptz

	if req.Trigger.Type == "one-off" && req.Trigger.DateTime != "" {
		nextRun = dateTime
	} else if req.Trigger.Type == "cron" && req.Trigger.Cron != "" {
		if t, err := NextCronTime(req.Trigger.Cron); err == nil {
			nextRun = pgtype.Timestamptz{Time: *t, Valid: true}
		}
	}

	task, err := s.DB.CreateTask(c, database.CreateTaskParams{
		Name:            req.Name,
		TriggerType:     req.Trigger.Type,
		TriggerDatetime: dateTime,
		TriggerCron:     cron,
		ActionMethod:    req.Action.Method,
		ActionUrl:       req.Action.URL,
		ActionHeaders:   reqHeaders,
		ActionPayload:   req.Action.Payload,
		Status:          "scheduled",
		NextRun:         nextRun,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task: " + err.Error()})
		return
	}

	response := entity.TaskResponse{
		ID:     task.ID,
		Name:   task.Name,
		Status: task.Status,
		Trigger: entity.TriggerData{
			Type: task.TriggerType,
		},
		Action: entity.ActionData{
			Method: task.ActionMethod,
			URL:    task.ActionUrl,
		},
		CreatedAt: task.CreatedAt.Time,
		UpdatedAt: task.UpdatedAt.Time,
		NextRun:   &task.NextRun.Time,
	}

	if task.TriggerDatetime.Valid {
		response.Trigger.DateTime = task.TriggerDatetime.Time.Format(time.RFC3339)
	}
	if task.TriggerCron.Valid {
		response.Trigger.Cron = task.TriggerCron.String
	}

	var headers map[string]string
	if err := json.Unmarshal(task.ActionHeaders, &headers); err == nil {
		response.Action.Headers = headers
	}

	var payload json.RawMessage
	if err := json.Unmarshal(task.ActionPayload, &payload); err == nil {
		response.Action.Payload = payload
	}

	c.JSON(http.StatusCreated, response)

}

func (s *Server) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")

	offset := (page - 1) * size

	tasks, err := s.DB.ListTasks(c, database.ListTasksParams{
		Column1: status,
		Limit:   int32(size),
		Offset:  int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	var taskResponses []entity.TaskResponse
	for _, task := range tasks {
		response, err := taskToResponse(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format task"})
			return
		}
		taskResponses = append(taskResponses, response)
	}

	response := entity.ListTasksResponse{
		Tasks: taskResponses,
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) CancelTask(c *gin.Context) {
	idParam := c.Param("id")
	var pguuid pgtype.UUID
	err := pguuid.Scan(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := s.DB.CancelTask(c, pguuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel task"})
		return
	}

	response, err := taskToResponse(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format response"})
		return
	}

	c.JSON(http.StatusOK, response)
}
func (s *Server) ListTaskResults(c *gin.Context) {
	idParam := c.Param("id")
	var pguuid pgtype.UUID
	err := pguuid.Scan(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	results, err := s.DB.ListTaskResults(c, pguuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task results"})
		return
	}

	var response []entity.TaskResultResponse
	for _, result := range results {
		resultResponse, err := taskResultToResponse(result)
		if err != nil {
			log.Printf("Error converting result: %v", err)
			continue
		}
		response = append(response, resultResponse)
	}

	c.JSON(http.StatusOK, gin.H{"results": response})
}

func (s *Server) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	var pguuid pgtype.UUID
	err := pguuid.Scan(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req entity.UpdateTaskRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currTask, err := s.DB.GetTask(c, pguuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	params := database.UpdateTaskParams{
		ID:              pguuid,
		Name:            currTask.Name,
		TriggerType:     currTask.TriggerType,
		TriggerDatetime: currTask.TriggerDatetime,
		TriggerCron:     currTask.TriggerCron,
		ActionMethod:    currTask.ActionMethod,
		ActionUrl:       currTask.ActionUrl,
		ActionHeaders:   currTask.ActionHeaders,
		ActionPayload:   currTask.ActionPayload,
		Status:          currTask.Status,
		NextRun:         currTask.NextRun,
	}

	if req.Name != nil {
		params.Name = *req.Name
	}

	if req.Trigger != nil {
		params.TriggerType = req.Trigger.Type

		if req.Trigger.Type == "one-off" {
			if req.Trigger.DateTime == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "datetime is required for one-off tasks"})
				return
			}

			dateTime, err := StringToTimestamptz(req.Trigger.DateTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datetime format"})
				return
			}

			params.TriggerDatetime = dateTime
			params.NextRun = dateTime
			params.TriggerCron = pgtype.Text{Valid: false}
		} else if req.Trigger.Type == "cron" {
			if req.Trigger.Cron == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "cron expression is required for cron tasks"})
				return
			}

			nextTime, err := NextCronTime(req.Trigger.Cron)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cron expression: " + err.Error()})
				return
			}

			params.TriggerCron = pgtype.Text{String: req.Trigger.Cron, Valid: true}
			params.TriggerDatetime = pgtype.Timestamptz{Valid: false}
			params.NextRun = pgtype.Timestamptz{Time: *nextTime, Valid: true}
		}

	}

	if req.Action != nil {
		if req.Action.Method != "" {
			params.ActionMethod = req.Action.Method
		}

		if req.Action.URL != "" {
			params.ActionUrl = req.Action.URL
		}

		if req.Action.Headers != nil {
			headersJSON, err := json.Marshal(req.Action.Headers)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process headers"})
				return
			}
			params.ActionHeaders = headersJSON
		}

		if req.Action.Payload != nil {
			payloadJSON, err := json.Marshal(req.Action.Payload)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payload"})
				return
			}
			params.ActionPayload = payloadJSON
		}
	}

	updatedTask, err := s.DB.UpdateTask(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task: " + err.Error()})
		return
	}

	response, err := taskToResponse(updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to format response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (s *Server) ListAllTasksResults(c *gin.Context) {
	results, err := s.DB.ListAllTaskResults(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task results"})
		return
	}

	var response []entity.TaskResultResponse
	for _, result := range results {
		resultResponse, err := taskResultToResponse(result)
		if err != nil {
			log.Printf("Error converting result: %v", err)
			continue
		}
		response = append(response, resultResponse)
	}

	c.JSON(http.StatusOK, gin.H{"results": response})
}
