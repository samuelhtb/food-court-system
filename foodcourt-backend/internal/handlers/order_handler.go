package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/dto"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/services"
)

type OrderHandler struct {
	service         services.OrderService
	midtransService services.MidtransService
}

func NewOrderHandler(service services.OrderService, midtransService services.MidtransService) *OrderHandler {
	return &OrderHandler{service, midtransService}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid", "details": err.Error()})
		return
	}

	// Get OrderID from service
	orderID, err := h.service.CreateOrder(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add order_id to JSON response
	var midtransToken string
	if req.PaymentMethod == "qris" {
		order, err := h.service.GetOrderWithDetails(orderID)
		if err != nil {
			log.Printf("Gagal mengambil detail pesanan untuk Midtrans: %v", err)
		} else {
			token, err := h.midtransService.CreateSnapTransaction(order)
			if err != nil {
				log.Printf("Gagal membuat transaksi Midtrans: %v", err)
			} else {
				midtransToken = token
			}
		}
	}

	c.JSON(http.StatusCreated, dto.CreateOrderResponse{
		Message:       "Pesanan berhasil dibuat!",
		OrderID:       orderID,
		MidtransToken: midtransToken,
	})
}

func (h *OrderHandler) GetTenantOrders(c *gin.Context) {
	// Get userID from JWT token (saved by middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tenantID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID format"})
		return
	}

	subOrders, err := h.service.GetTenantOrders(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subOrders})
}

func (h *OrderHandler) UpdateTenantOrderStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := uuid.Parse(userID.(string))

	subOrderIDParam := c.Param("id")
	subOrderID, err := uuid.Parse(subOrderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	var req dto.UpdateSubOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid", "details": err.Error()})
		return
	}

	if err := h.service.UpdateTenantOrderStatus(subOrderID, tenantID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pesanan berhasil diperbarui menjadi " + req.Status})
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pesanan: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *OrderHandler) MarkOrderAsPaid(c *gin.Context) {
	orderIDParam := c.Param("id")
	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID Pesanan tidak valid"})
		return
	}

	if err := h.service.MarkOrderAsPaid(orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses pembayaran: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil! Pesanan telah diteruskan ke dapur tenant."})
}

func (h *OrderHandler) GetTenantEarnings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tenantID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID format"})
		return
	}

	totalEarnings, totalOrders, err := h.service.GetTenantEarnings(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung laporan pendapatan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Laporan pendapatan berhasil ditarik",
		"data": gin.H{
			"total_orders_completed": totalOrders,
			"total_earnings":         totalEarnings,
		},
	})
}

func (h *OrderHandler) GetOrderDetails(c *gin.Context) {
	orderIDParam := c.Param("id")
	orderID, err := uuid.Parse(orderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID Pesanan tidak valid"})
		return
	}

	order, err := h.service.GetOrderWithDetails(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	// Transformasi ke DTO untuk meratakan rincian item (flattening)
	var items []dto.TrackOrderItemResponse
	for _, so := range order.SubOrders {
		for _, it := range so.OrderItems {
			items = append(items, dto.TrackOrderItemResponse{
				MenuName: it.Menu.Name,
				Quantity: it.Quantity,
				Price:    it.PriceAtOrder,
			})
		}
	}

	response := dto.TrackOrderResponse{
		ID:            order.ID,
		CustomerName:  order.CustomerName,
		PaymentMethod: order.PaymentMethod,
		PaymentStatus: order.PaymentStatus,
		OrderStatus:   order.OrderStatus,
		TotalAmount:   order.TotalAmount,
		Items:         items,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail pesanan berhasil diambil",
		"data":    response,
	})
}

// Laporan Pemasukan (Reports)

func (h *OrderHandler) GetAdminIncomeReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	report, err := h.service.GetAdminIncomeReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil laporan pendapatan sistem"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *OrderHandler) GetTenantIncomeReport(c *gin.Context) {
	tenantIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak"})
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	report, err := h.service.GetTenantIncomeReport(tenantID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil laporan pendapatan tenant"})
		return
	}

	c.JSON(http.StatusOK, report)
}

