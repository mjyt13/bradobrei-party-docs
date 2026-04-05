package services

import (
	"errors"

	"bradobrei/backend/internal/models"
	"bradobrei/backend/internal/repository"
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
}

func NewInventoryService(inventoryRepo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

func (s *InventoryService) GetBySalon(salonID uint) ([]models.Inventory, error) {
	return s.inventoryRepo.GetBySalon(salonID)
}

func (s *InventoryService) SetQuantity(salonID, materialID uint, quantity float64) (*models.Inventory, error) {
	if quantity < 0 {
		return nil, errors.New("количество материала не может быть отрицательным")
	}

	if err := s.inventoryRepo.SetQuantity(salonID, materialID, quantity); err != nil {
		return nil, err
	}

	return s.inventoryRepo.GetItemWithMaterial(salonID, materialID)
}
