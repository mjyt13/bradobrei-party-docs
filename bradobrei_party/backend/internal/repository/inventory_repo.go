package repository

import (
	"bradobrei/backend/internal/models"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetBySalon(salonID uint) ([]models.Inventory, error) {
	var inv []models.Inventory
	err := r.db.Where("salon_id = ?", salonID).Preload("Material").Find(&inv).Error
	return inv, err
}

func (r *InventoryRepository) GetItem(salonID, materialID uint) (*models.Inventory, error) {
	var inv models.Inventory
	err := r.db.
		Where("salon_id = ? AND material_id = ?", salonID, materialID).
		First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InventoryRepository) GetItemWithMaterial(salonID, materialID uint) (*models.Inventory, error) {
	var inv models.Inventory
	err := r.db.
		Where("salon_id = ? AND material_id = ?", salonID, materialID).
		Preload("Material").
		First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// WriteOff — списание материалов при подтверждении бронирования (транзакция)
func (r *InventoryRepository) WriteOff(salonID uint, items []models.ServiceMaterial, quantity int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, sm := range items {
			need := sm.QuantityPerUse * float64(quantity)
			result := tx.Model(&models.Inventory{}).
				Where("salon_id = ? AND material_id = ? AND quantity >= ?", salonID, sm.MaterialID, need).
				UpdateColumn("quantity", gorm.Expr("quantity - ?", need))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound // недостаточно материала
			}
		}
		return nil
	})
}

func (r *InventoryRepository) Upsert(inv *models.Inventory) error {
	return r.db.Save(inv).Error
}

func (r *InventoryRepository) AdjustQuantity(salonID, materialID uint, delta float64) error {
	inv, err := r.GetItem(salonID, materialID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if delta < 0 {
				return gorm.ErrRecordNotFound
			}

			return r.db.Create(&models.Inventory{
				SalonID:    salonID,
				MaterialID: materialID,
				Quantity:   delta,
			}).Error
		}
		return err
	}

	newQty := inv.Quantity + delta
	if newQty < 0 {
		return gorm.ErrRecordNotFound
	}

	inv.Quantity = newQty
	return r.Upsert(inv)
}

func (r *InventoryRepository) SetQuantity(salonID, materialID uint, quantity float64) error {
	inv, err := r.GetItem(salonID, materialID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(&models.Inventory{
				SalonID:    salonID,
				MaterialID: materialID,
				Quantity:   quantity,
			}).Error
		}
		return err
	}

	inv.Quantity = quantity
	return r.Upsert(inv)
}
