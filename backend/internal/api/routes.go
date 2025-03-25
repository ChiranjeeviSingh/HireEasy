package api

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes defines all API routes
func SetupRoutes(router *gin.Engine) {
	// Add CORS middleware
	router.Use(cors.Default())

	// Public routes (no auth required)
	public := router.Group("/api")
	{
		// Authentication routes
		public.POST("/login", handlers.LoginH)
		public.POST("/register", handlers.RegisterH)

		// Form submission routes (unauthenticated)
		public.POST("/jobs/:job_id/apply", handlers.HandleFormSubmission)    // Submit job application
		public.GET("/jobs/:job_id/submissions", handlers.GetFormSubmissions) // Get job submissions
	}

	// Protected routes (auth required)
	api := router.Group("/api", middleware.AuthMiddleware())
	{
		// Job routes
		jobs := api.Group("/jobs")
		{
			jobs.POST("", handlers.CreateJobH)                        // Create job
			jobs.PUT("/:jobId", handlers.UpdateJobH)                  // Update job
			jobs.GET("/:jobId", handlers.GetJobByIdH)                // Get specific job by id
			jobs.GET("/jobtitle/:jobtitle", handlers.GetJobsByTitleH) // Get jobs by jobtitle
			jobs.GET("/status/:status", handlers.GetJobsByStatusH)    // Get jobs by status
			jobs.GET("", handlers.ListUserJobsH)                      // List all jobs for user
			jobs.DELETE("/:jobId", handlers.DeleteJobH)               // Delete job
		}

		// Form template routes
		formTemplates := api.Group("/forms/templates")
		{
			formTemplates.POST("", handlers.CreateFormTemplateH)                     // Create form template
			formTemplates.GET("/:form_template_id", handlers.GetFormTemplateH)       // Get specific template
			formTemplates.GET("", handlers.ListFormTemplatesH)                       // List all templates
			formTemplates.DELETE("/:form_template_id", handlers.DeleteFormTemplateH) // Delete template
		}

		// Application form routes
		applicationForms := api.Group("")
		{
			applicationForms.POST("/jobs/:job_id/forms", handlers.LinkJobToFormTemplateH)   // Link job to form template
			applicationForms.PATCH("/forms/:form_uuid/status", handlers.UpdateFormStatusH)  // Update form status
			applicationForms.GET("/forms/:form_uuid", handlers.GetFormDetailsH)             // Get form details
			applicationForms.DELETE("/forms/:form_uuid", handlers.DeleteFormH)              // Delete form
		}
	}
}
