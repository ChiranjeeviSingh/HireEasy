package main

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/database"
	"fmt"
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load Configuration
	config.LoadConfig()

	// Connect to Database
	database.Connect()

	// Create Router
	router := gin.Default()

	// Get CORS Config
	corsConfig := config.GetConfig().CORSConfig

	// ✅ Apply CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowedOrigins,
		AllowMethods:     corsConfig.AllowedMethods,
		AllowHeaders:     corsConfig.AllowedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
	}))

	log.Println("✅ CORS Middleware Applied")
	log.Printf("🌍 Allowed Origins: %s", strings.Join(corsConfig.AllowedOrigins, ", "))

	// Setup API Routes
	api.SetupRoutes(router)

	// Log registered routes
	for _, r := range router.Routes() {
		fmt.Println("🔍 Route registered:", r.Method, r.Path)
	}

	// Start Server
	log.Println("🚀 Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("❌ Could not start server: %s\n", err)
	}
}
