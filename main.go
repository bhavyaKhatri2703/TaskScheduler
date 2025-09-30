package main

import (
	"context"
	"fmt"
	"log"
	"scheduler/application/api"
	"scheduler/database"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*database.Queries, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "user=postgres password=postgres dbname=scheduler sslmode=disable")

	if err != nil {
		return nil, err
	}

	db := database.New(conn)

	return db, nil
}

func main() {

	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}
	server := api.NewServer(db)
	server.SetRoutes()

	fmt.Println("Server is running on http://localhost:8080")
	if err := server.Router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
