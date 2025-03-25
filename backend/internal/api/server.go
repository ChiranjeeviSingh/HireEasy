package api

import (
	"backend/internal/api/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	router := gin.Default()
	
	// Setup routes
	router.POST("/forms/:form_uuid/submit", handlers.HandleFormSubmission)
	router.GET("/forms/:form_uuid/submissions", handlers.GetFormSubmissions)

	return &Server{
		router: router,
	}
}

func (s *Server) Run() error {
	log.Println("Starting server on :8080")
	return s.router.Run(":8080")
} 