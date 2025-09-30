package workers

import (
	"context"
	"log"
	"scheduler/database"
	"time"
)

func (wp *WorkerPool) executeTask(ctx context.Context, task database.Task) {

	log.Printf("Executing task: %s [%s %s]", task.Name, task.ActionMethod, task.ActionUrl)

	// do work

	time.Sleep(2 * time.Second)

	log.Printf("Task %s completed", task.Name)
}
