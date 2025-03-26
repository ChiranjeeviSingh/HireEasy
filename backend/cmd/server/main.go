package main

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

    // Load the configuration and connect to the database
    config.LoadConfig()
    database.Connect()

    // Initialize Gin router
    router := gin.Default()

    // CORS configuration
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // Allow frontend origin
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow all methods
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"}, // Allow necessary headers
        ExposeHeaders:    []string{"Content-Length", "Authorization"}, // Expose necessary headers
        AllowCredentials: true,
        MaxAge:           12 * time.Hour, // Cache preflight response for 12 hours
    }))

    // Handle preflight OPTIONS requests (CORS preflight check)
    router.OPTIONS("/*any", func(c *gin.Context) {
        c.Status(http.StatusOK)
    })

    // Set up routes
    api.SetupRoutes(router)

    // Start the server
    log.Println("Starting server on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}