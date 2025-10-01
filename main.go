package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"scheduler/application/api"
	"scheduler/database"
	_ "scheduler/docs"
	"scheduler/scheduler"
	"scheduler/workers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*database.Queries, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	pool, err := pgxpool.New(ctx, dsn)

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
