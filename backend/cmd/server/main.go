package main

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/database"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.Connect()

	router := gin.Default()

	// ✅ Custom CORS setup to allow Authorization header
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api.SetupRoutes(router)

	for _, r := range router.Routes() {
		fmt.Println("🔍 Route registered:", r.Method, r.Path)
	}

	log.Println("Starting server on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
