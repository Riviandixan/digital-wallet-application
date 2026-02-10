package handler

import (
	"digital-wallet-application/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	WalletUsecase domain.WalletUsecase
}

func NewWalletHandler(r *gin.Engine, wu domain.WalletUsecase) {
	handler := &WalletHandler{
		WalletUsecase: wu,
	}

	api := r.Group("/api")
	{
		api.GET("/balance/:user_id", handler.GetBalance)
		api.POST("/withdraw", handler.Withdraw)
	}
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}

	wallet, err := h.WalletUsecase.GetBalance(c.Request.Context(), userID)
	if err != nil {
		if err == domain.ErrWalletNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		// Log the actual error for debugging
		log.Printf("Error getting balance for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": wallet.UserID,
		"balance": wallet.Balance,
	})
}

type WithdrawRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Amount float64   `json:"amount" binding:"required,gt=0"`
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	wallet, err := h.WalletUsecase.Withdraw(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case domain.ErrInsufficientFunds, domain.ErrInvalidAmount:
			status = http.StatusBadRequest
		case domain.ErrWalletNotFound:
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "withdrawal successful",
		"user_id":     wallet.UserID,
		"new_balance": wallet.Balance,
	})
}
