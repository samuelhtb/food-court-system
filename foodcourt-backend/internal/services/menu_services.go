package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/dto"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
)

type MenuService interface {
	CreateMenu(tenantID uuid.UUID, req dto.CreateMenuRequest) error
	GetMenus(tenantID uuid.UUID) ([]dto.MenuResponse, error)
	GetPublicMenus() ([]dto.MenuResponse, error)
	UpdateMenu(id uuid.UUID, tenantID uuid.UUID, req dto.UpdateMenuRequest) error
	DeleteMenu(id uuid.UUID, tenantID uuid.UUID) error
}

type menuService struct {
	repo repositories.MenuRepository
}

func NewMenuService(repo repositories.MenuRepository) MenuService {
	return &menuService{repo}
}

func (s *menuService) CreateMenu(tenantID uuid.UUID, req dto.CreateMenuRequest) error {
	newMenu := models.Menu{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsAvailable: true, // Default saat dibuat pertama kali
	}
	return s.repo.CreateMenu(&newMenu)
}

func (s *menuService) GetMenus(tenantID uuid.UUID) ([]dto.MenuResponse, error) {
	menus, err := s.repo.GetMenusByTenantID(tenantID)
	if err != nil {
		return nil, err
	}

	var responses []dto.MenuResponse
	for _, m := range menus {
		responses = append(responses, dto.MenuResponse{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			Price:       m.Price,
			Stock:       m.Stock,
			ImageURL:    m.ImageURL,
			IsAvailable: m.IsAvailable,
		})
	}
	return responses, nil
}

func (s *menuService) GetPublicMenus() ([]dto.MenuResponse, error) {
	menus, err := s.repo.GetAllMenus()
	if err != nil {
		return nil, err
	}

	var responses []dto.MenuResponse
	for _, m := range menus {
		responses = append(responses, dto.MenuResponse{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			Price:       m.Price,
			Stock:       m.Stock,
			ImageURL:    m.ImageURL,
			IsAvailable: m.IsAvailable,
			TenantName:  m.Tenant.TenantName,
		})
	}
	return responses, nil
}

func (s *menuService) UpdateMenu(id uuid.UUID, tenantID uuid.UUID, req dto.UpdateMenuRequest) error {
	// Pastikan menu ini milik tenant tersebut sebelum di-update
	menu, err := s.repo.GetMenuByIDAndTenantID(id, tenantID)
	if err != nil {
		return errors.New("menu tidak ditemukan atau bukan milik anda")
	}

	// Update data dari request
	menu.Name = req.Name
	menu.Description = req.Description
	menu.Price = req.Price
	menu.Stock = req.Stock
	menu.ImageURL = req.ImageURL
	menu.IsAvailable = req.IsAvailable

	return s.repo.UpdateMenu(menu)
}

func (s *menuService) DeleteMenu(id uuid.UUID, tenantID uuid.UUID) error {
	return s.repo.DeleteMenu(id, tenantID)
}