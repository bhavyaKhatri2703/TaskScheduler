package scheduler

import (
	"context"
	"log"
	"scheduler/database"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Scheduler struct {
	db       *database.Queries
	interval time.Duration
	taskChan chan<- database.Task
}

func (s *Scheduler) StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			s.pollAndQueueTasks(ctx)
		}
	}
}

func NewScheduler(db *database.Queries, taskChan chan<- database.Task) *Scheduler {
	return &Scheduler{
		db:       db,
		interval: 30 * time.Second,
		taskChan: taskChan,
	}
}

func (s *Scheduler) pollAndQueueTasks(ctx context.Context) {
	readyTasks, err := s.findReadyTasks(ctx)
	if err != nil {
		log.Printf("Error finding scheduled tasks: %v", err)
		return
	}

	if len(readyTasks) == 0 {
		log.Println("No scheduled tasks found")
		return
	}

	log.Printf("Found scheduled ready tasks", len(readyTasks))

	for _, task := range readyTasks {
		s.queueTask(task)

	}
}

func (s *Scheduler) findReadyTasks(ctx context.Context) ([]database.Task, error) {
	now := pgtype.Timestamptz{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	return s.db.GetTasksToRun(ctx, now)
}

func (s *Scheduler) queueTask(task database.Task) {
	select {
	case s.taskChan <- task:
		log.Printf("Queued task: %s", task.Name)
	default:
		log.Printf("Task queue full, dropping task: %s", task.Name)
	}
}
