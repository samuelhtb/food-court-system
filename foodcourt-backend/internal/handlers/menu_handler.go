package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/dto"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/services"
)

type MenuHandler struct {
	service services.MenuService
}

func NewMenuHandler(service services.MenuService) *MenuHandler {
	return &MenuHandler{service}
}

// Helper untuk mengambil TenantID dari context JWT
func getTenantID(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("user_id") // Pastikan middleware JWT kamu set "user_id"
	if !exists {
		return uuid.Nil, errors.New("unauthorized")
	}
	return uuid.Parse(userIDStr.(string))
}

func (h *MenuHandler) Create(c *gin.Context) {
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak"})
		return
	}

	var req dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateMenu(tenantID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Menu berhasil ditambahkan"})
}

func (h *MenuHandler) GetAll(c *gin.Context) {
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak"})
		return
	}

	menus, err := h.service.GetMenus(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menus})
}

func (h *MenuHandler) Update(c *gin.Context) {
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak"})
		return
	}

	menuID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID menu tidak valid"})
		return
	}

	var req dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMenu(menuID, tenantID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil diperbarui"})
}

func (h *MenuHandler) Delete(c *gin.Context) {
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak"})
		return
	}

	menuID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID menu tidak valid"})
		return
	}

	if err := h.service.DeleteMenu(menuID, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil dihapus"})
}