package services

import (
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
)

type MenuService interface {
	CreateMenu(menu *models.Menu) error
	GetAllMenus() ([]models.Menu, error)
}

type menuService struct {
	repo repositories.MenuRepository
}

func NewMenuService(repo repositories.MenuRepository) MenuService {
	return &menuService{repo}
}

func (s *menuService) CreateMenu(menu *models.Menu) error {
	//  Harga tidak boleh 0
	if menu.Price <= 0 {
		return interface{}(nil).(error)
	}
	return s.repo.CreateMenu(menu)
}

func (s *menuService) GetAllMenus() ([]models.Menu, error) {
	return s.repo.GetAllMenus()
}