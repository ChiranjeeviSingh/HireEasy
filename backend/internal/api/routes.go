package api

import (
	"backend/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Authentication routes
	api.POST("/login", handlers.Login)
	api.POST("/register", handlers.Register)
	api.POST("/apply", HandleJobApplication)
}
