package workers

import (
	"context"
	"log"
)

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopped", id)
			return

		case task, ok := <-wp.taskChan:
			if !ok {
				log.Printf("Worker %d: task channel closed", id)
				return
			}

			log.Printf("Worker %d: processing task %s", id, task.Name)
			wp.executeTask(ctx, task)
		}
	}
}
