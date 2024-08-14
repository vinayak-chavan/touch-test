package main

import (
    "log"
    "touch-test/config"
    "touch-test/routes"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "os"
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

    // Register Routes
    routes.RegisterRoutes(r, db, rdb)

    // Start the server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Could not start the server: %v", err)
    }
}
