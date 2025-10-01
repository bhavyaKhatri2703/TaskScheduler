package main

import (
	"context"
	"fmt"
	"log"
	"scheduler/application/api"
	"scheduler/database"
	"scheduler/scheduler"
	"scheduler/workers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*database.Queries, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "user=postgres password=postgres dbname=scheduler sslmode=disable")

	if err != nil {
		return nil, err
	}

	db := database.New(pool)

	return db, nil
}

func main() {

	ctx := context.Background()
	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}

	taskChan := make(chan database.Task, 100)

	schedulerEngine := scheduler.NewScheduler(db, taskChan)
	schedulerCtx, schedulerCancel := context.WithCancel(ctx)
	go schedulerEngine.StartScheduler(schedulerCtx)

	workerPool := workers.NewWorkerPool(db, taskChan, 5)
	workerCtx, workerCancel := context.WithCancel(ctx)
	workerPool.Start(workerCtx)

	server := api.NewServer(db)
	server.SetRoutes()

	fmt.Println("Server is running on http://localhost:8080")
	if err := server.Router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	schedulerCancel()

	workerCancel()
	workerPool.Stop()
}
