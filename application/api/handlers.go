package api

import (
	"encoding/json"
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
		// add next_run
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

func (s *Server) ListTasts(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"results": results})
}
