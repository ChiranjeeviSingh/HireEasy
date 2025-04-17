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
			jobs.PUT("/:job_id", handlers.UpdateJobH)                  // Update job
			jobs.GET("/:job_id", handlers.GetJobByIdH)                // Get specific job by id
			jobs.GET("/jobtitle/:jobtitle", handlers.GetJobsByTitleH) // Get jobs by jobtitle
			jobs.GET("/status/:status", handlers.GetJobsByStatusH)    // Get jobs by status
			jobs.GET("", handlers.ListUserJobsH)                      // List all jobs for user
			jobs.DELETE("/:job_id", handlers.DeleteJobH)               // Delete job
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
            applicationForms.POST("/jobs/:job_id/forms", handlers.LinkJobToFormTemplateH)   // Link job to form template, return unique URL and form_uuid
            applicationForms.PATCH("/forms/:form_uuid/status", handlers.UpdateFormStatusH)  // Update form status (active/inactive)
            applicationForms.GET("/forms/:form_uuid", handlers.GetFormDetailsH)             // Get job and form template details (unauthenticated)
            applicationForms.DELETE("/forms/:form_uuid", handlers.DeleteFormH)              // Delete form and unlink from job
        }


        // Profile(Interviwer) routes
        profiles := api.Group("/profiles")
        {
            profiles.POST("", handlers.CreateProfileH)          // Create new profile
            profiles.PUT("", handlers.UpdateMyProfileH)         // Update own profile
            profiles.GET("/me", handlers.GetMyProfileH)         // Get own profile
            profiles.GET("/user/:user_name", handlers.GetUserProfileH) // Get other users profile
        }

        // Availability routes
        availability := api.Group("/availability")
        {
            availability.POST("", handlers.CreateAvailabilityH)         // Create availability slot(Interviewer action)
            availability.DELETE("/:id", handlers.DeleteAvailabilityH)   // Delete availability slot(Interviewer action)
            availability.GET("/me", handlers.GetMyAvailabilityH)       // Get own availability using JWT token with optional date range
            availability.GET("/user/:user_name", handlers.GetUserAvailabilityH)  // Get specific user's availability(using user's username)
            availability.GET("", handlers.GetAllAvailabilityH)         // Get all available people with with optional date range, profile filters
        }
   }

}