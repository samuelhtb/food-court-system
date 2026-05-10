package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/services"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/pkg/midtrans"
)

type MidtransHandler struct {
	midtransService services.MidtransService
}

func NewMidtransHandler(midtransService services.MidtransService) *MidtransHandler {
	return &MidtransHandler{midtransService}
}

func (h *MidtransHandler) HandleNotification(c *gin.Context) {
	var notif midtrans.Notification
	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification body"})
		return
	}

	if err := h.midtransService.ProcessNotification(notif); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification processed"})
}

func (h *MidtransHandler) VerifyLocalPayment(c *gin.Context) {
	var req struct {
		OrderID uuid.UUID `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.midtransService.VerifyPayment(req.OrderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment verified locally"})
}
