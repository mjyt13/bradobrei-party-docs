package repository

import (
	"bradobrei/backend/internal/models"

	"gorm.io/gorm"
)

type MaterialExpenseRepository struct {
	db *gorm.DB
}

func NewMaterialExpenseRepository(db *gorm.DB) *MaterialExpenseRepository {
	return &MaterialExpenseRepository{db: db}
}

func (r *MaterialExpenseRepository) GetAll() ([]models.MaterialExpense, error) {
	var expenses []models.MaterialExpense
	err := r.db.Preload("Material").Preload("Salon").Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

func (r *MaterialExpenseRepository) GetByID(id uint) (*models.MaterialExpense, error) {
	var expense models.MaterialExpense
	if err := r.db.Preload("Material").Preload("Salon").First(&expense, id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *MaterialExpenseRepository) Create(expense *models.MaterialExpense) error {
	return r.db.Create(expense).Error
}

func (r *MaterialExpenseRepository) Update(expense *models.MaterialExpense) error {
	return r.db.Save(expense).Error
}

func (r *MaterialExpenseRepository) Delete(id uint) error {
	return r.db.Delete(&models.MaterialExpense{}, id).Error
}
