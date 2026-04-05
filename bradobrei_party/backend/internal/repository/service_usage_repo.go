package repository

import (
	"bradobrei/backend/internal/models"

	"gorm.io/gorm"
)

type ServiceUsageRepository struct {
	db *gorm.DB
}

func NewServiceUsageRepository(db *gorm.DB) *ServiceUsageRepository {
	return &ServiceUsageRepository{db: db}
}

func (r *ServiceUsageRepository) Create(usage *models.ServiceUsage) error {
	return r.db.Create(usage).Error
}
