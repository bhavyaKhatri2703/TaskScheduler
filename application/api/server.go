package api

import (
	"scheduler/database"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
	DB     *database.Queries
}

func NewServer(db *database.Queries) *Server {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	s := &Server{
		Router: router,
		DB:     db,
	}
	return s
}

func (s *Server) SetRoutes() {
	r := s.Router
	r.POST("/tasks", s.CreateTask)
	r.GET("/tasks", s.ListTasks)
	r.GET("/tasks/:id", s.GetTask)
	r.PUT("/tasks/:id", s.UpdateTask)
	r.DELETE("/tasks/:id", s.CancelTask)
	r.GET("/tasks/:id/results", s.ListTaskResults)
}
