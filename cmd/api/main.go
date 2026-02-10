package main

import (
	"digital-wallet-application/internal/config"
	"digital-wallet-application/internal/database"
	"digital-wallet-application/internal/handler"
	"digital-wallet-application/internal/repository"
	"digital-wallet-application/internal/usecase"
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
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Repositories
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize Usecases
	walletUsecase := usecase.NewWalletUsecase(walletRepo, transactionRepo)

	// Initialize Gin Engine
	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Setup Handlers
	handler.NewWalletHandler(r, walletUsecase)

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
