package repository

import (
	"fmt"

	"bradobrei/backend/internal/models"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) GetAll() ([]models.EmployeeProfile, error) {
	var profiles []models.EmployeeProfile
	err := r.db.
		Preload("User").
		Preload("Salons").
		Preload("Services").
		Find(&profiles).Error
	return profiles, err
}

func (r *EmployeeRepository) GetByID(id uint) (*models.EmployeeProfile, error) {
	var p models.EmployeeProfile
	err := r.db.
		Preload("User").
		Preload("Salons").
		Preload("Services").
		First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *EmployeeRepository) GetByUserID(userID uint) (*models.EmployeeProfile, error) {
	var p models.EmployeeProfile
	err := r.db.
		Where("user_id = ?", userID).
		Preload("Salons").
		Preload("Services").
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *EmployeeRepository) Create(p *models.EmployeeProfile) error {
	return r.db.Create(p).Error
}

func (r *EmployeeRepository) Update(p *models.EmployeeProfile) error {
	return r.db.Save(p).Error
}

func (r *EmployeeRepository) ReplaceSalonAssignments(profileID uint, salonIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			"DELETE FROM employee_salons WHERE employee_profile_id = ?",
			profileID,
		).Error; err != nil {
			return err
		}

		for _, salonID := range salonIDs {
			if err := tx.Exec(
				"INSERT INTO employee_salons (employee_profile_id, salon_id) VALUES (?, ?)",
				profileID, salonID,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *EmployeeRepository) UpdateSchedule(profileID uint, schedule string) error {
	return r.db.Model(&models.EmployeeProfile{}).
		Where("id = ?", profileID).
		Update("work_schedule", schedule).Error
}

func (r *EmployeeRepository) AssignToSalon(profileID uint, salonID uint) error {
	result := r.db.Exec(
		"INSERT INTO employee_salons (employee_profile_id, salon_id) VALUES (?, ?)",
		profileID, salonID,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("сотрудник уже закреплён за этим салоном")
	}
	return nil
}

func (r *EmployeeRepository) RemoveFromSalon(profileID uint, salonID uint) error {
	return r.db.Exec(
		"DELETE FROM employee_salons WHERE employee_profile_id = ? AND salon_id = ?",
		profileID, salonID,
	).Error
}

func (r *EmployeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmployeeProfile{}, id).Error
}

func (r *EmployeeRepository) Fire(profileID uint, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			"DELETE FROM employee_salons WHERE employee_profile_id = ?",
			profileID,
		).Error; err != nil {
			return err
		}
		if err := tx.Exec(
			"DELETE FROM employee_services WHERE employee_profile_id = ?",
			profileID,
		).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.EmployeeProfile{}, profileID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.User{}, userID).Error; err != nil {
			return err
		}
		return nil
	})
}
