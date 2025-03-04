package api

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// add cors middleware
	router.Use(cors.Default())
	api := router.Group("/api", middleware.AuthMiddleware())

	// Authentication routes
	api.POST("/login", handlers.LoginH)
	api.POST("/register", handlers.RegisterH)

	//job routes
	jobs := api.Group("/jobs")
	{
		jobs.POST("", handlers.CreateJobH)                        // Create job
		jobs.PUT("/:jobId", handlers.UpdateJobH)                  // Update job
		jobs.GET("/:jobId", handlers.GetJobByIdH)                 // Get specific job by id
		jobs.GET("/jobtitle/:jobtitle", handlers.GetJobsByTitleH) // Get jobs by jobtitle - Has to include the jobtitle(could be a subset)
		jobs.GET("/status/:status", handlers.GetJobsByStatusH)    // Get jobs by status
		jobs.GET("", handlers.ListUserJobsH)                      // List all jobs for user
		jobs.DELETE("/:jobId", handlers.DeleteJobH)               // Delete job
	}

}
