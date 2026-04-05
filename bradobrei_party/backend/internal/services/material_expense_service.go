package services

import (
	"errors"

	"bradobrei/backend/internal/dto"
	"bradobrei/backend/internal/models"
	"bradobrei/backend/internal/repository"

	"gorm.io/gorm"
)

type MaterialExpenseService struct {
	db *gorm.DB
}

func NewMaterialExpenseService(db *gorm.DB) *MaterialExpenseService {
	return &MaterialExpenseService{db: db}
}

func (s *MaterialExpenseService) GetAll() ([]models.MaterialExpense, error) {
	return repository.NewMaterialExpenseRepository(s.db).GetAll()
}

func (s *MaterialExpenseService) GetByID(id uint) (*models.MaterialExpense, error) {
	return repository.NewMaterialExpenseRepository(s.db).GetByID(id)
}

func (s *MaterialExpenseService) Create(req dto.CreateMaterialExpenseRequest) (*models.MaterialExpense, error) {
	var created *models.MaterialExpense

	err := s.db.Transaction(func(tx *gorm.DB) error {
		expenseRepo := repository.NewMaterialExpenseRepository(tx)
		inventoryRepo := repository.NewInventoryRepository(tx)

		expense := &models.MaterialExpense{
			MaterialID:    req.MaterialID,
			SalonID:       req.SalonID,
			PurchasePrice: req.PurchasePrice,
			Quantity:      req.Quantity,
		}

		if err := expenseRepo.Create(expense); err != nil {
			return err
		}

		if err := inventoryRepo.AdjustQuantity(req.SalonID, req.MaterialID, req.Quantity); err != nil {
			return err
		}

		loaded, err := expenseRepo.GetByID(expense.ID)
		if err != nil {
			return err
		}
		created = loaded
		return nil
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *MaterialExpenseService) Update(id uint, req dto.UpdateMaterialExpenseRequest) (*models.MaterialExpense, error) {
	var updated *models.MaterialExpense

	err := s.db.Transaction(func(tx *gorm.DB) error {
		expenseRepo := repository.NewMaterialExpenseRepository(tx)
		inventoryRepo := repository.NewInventoryRepository(tx)

		existing, err := expenseRepo.GetByID(id)
		if err != nil {
			return err
		}

		if err := inventoryRepo.AdjustQuantity(existing.SalonID, existing.MaterialID, -existing.Quantity); err != nil {
			return errors.New("не удалось откатить предыдущий остаток: " + err.Error())
		}

		if err := inventoryRepo.AdjustQuantity(req.SalonID, req.MaterialID, req.Quantity); err != nil {
			return errors.New("не удалось применить новый остаток: " + err.Error())
		}

		existing.MaterialID = req.MaterialID
		existing.SalonID = req.SalonID
		existing.PurchasePrice = req.PurchasePrice
		existing.Quantity = req.Quantity

		if err := expenseRepo.Update(existing); err != nil {
			return err
		}

		loaded, err := expenseRepo.GetByID(existing.ID)
		if err != nil {
			return err
		}
		updated = loaded
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *MaterialExpenseService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		expenseRepo := repository.NewMaterialExpenseRepository(tx)
		inventoryRepo := repository.NewInventoryRepository(tx)

		existing, err := expenseRepo.GetByID(id)
		if err != nil {
			return err
		}

		if err := inventoryRepo.AdjustQuantity(existing.SalonID, existing.MaterialID, -existing.Quantity); err != nil {
			return errors.New("не удалось скорректировать остаток при удалении закупки: " + err.Error())
		}

		return expenseRepo.Delete(id)
	})
}
