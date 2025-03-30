package main

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/database"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()

	database.Connect()

	router := gin.Default()

	api.SetupRoutes(router)
	for _, r := range router.Routes() {
		fmt.Println("🔍 Route registered:", r.Method, r.Path)
	}

	log.Println("Starting server on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
