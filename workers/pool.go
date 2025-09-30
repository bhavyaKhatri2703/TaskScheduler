package workers

import (
	"context"
	"log"
	"scheduler/database"
	"sync"
)

type WorkerPool struct {
	db       *database.Queries
	taskChan <-chan database.Task
	count    int
	wg       *sync.WaitGroup
}

func NewWorkerPool(db *database.Queries, taskChan <-chan database.Task, workerCount int) *WorkerPool {
	return &WorkerPool{
		db:       db,
		taskChan: taskChan,
		count:    workerCount,
		wg:       &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	log.Printf("starting %d workers", wp.count)

	for i := 1; i <= wp.count; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}

}

func (wp *WorkerPool) Stop() {
	log.Println("Waiting for workers to finish...")
	wp.wg.Wait()
	log.Println("All workers stopped")
}
