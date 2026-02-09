package main

import (
	"digital-wallet-application/internal/config"
	"digital-wallet-application/internal/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect do database: %v", err)
	}
	defer db.Close()

	// Initialize Gin Engine
	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start Server
	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
