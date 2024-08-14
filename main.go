package main

import (
	"log"
	"os"
	"touch-test/config"
	"touch-test/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	r := gin.Default()

	// Initialize Database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	rdb := config.NewRedisClient()

	// Health check endpoint
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Application is working fine!!",
		})
	})

	// Create a new group with `/api` prefix
	api := r.Group("/api/v1")

	// Register routes with the `/api` prefix
	routes.RegisterRoutes(api, db, rdb)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Could not start the server: %v", err)
	}
}
