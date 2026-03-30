package services

import (
	"errors"

	"bradobrei/backend/internal/dto"
	"bradobrei/backend/internal/models"
	"bradobrei/backend/internal/repository"
)

type EmployeeService struct {
	employeeRepo *repository.EmployeeRepository
	userRepo     *repository.UserRepository
}

func NewEmployeeService(
	employeeRepo *repository.EmployeeRepository,
	userRepo *repository.UserRepository,
) *EmployeeService {
	return &EmployeeService{employeeRepo: employeeRepo, userRepo: userRepo}
}

func (s *EmployeeService) GetAll() ([]models.EmployeeProfile, error) {
	return s.employeeRepo.GetAll()
}

func (s *EmployeeService) GetByID(id uint) (*models.EmployeeProfile, error) {
	return s.employeeRepo.GetByID(id)
}

func (s *EmployeeService) GetMyProfile(userID uint) (*models.EmployeeProfile, error) {
	return s.employeeRepo.GetByUserID(userID)
}

func (s *EmployeeService) HireEmployee(req dto.HireEmployeeRequest) (*models.EmployeeProfile, error) {
	if req.Role == models.RoleClient {
		return nil, errors.New("нельзя нанять пользователя с ролью CLIENT")
	}

	var emailPtr *string
	if req.Email != "" {
		emailPtr = &req.Email
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: req.PasswordHash,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Email:        emailPtr,
		Role:         req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("не удалось создать пользователя: " + err.Error())
	}

	var schedulePtr *string
	if req.WorkSchedule != "" {
		schedulePtr = &req.WorkSchedule
	}

	profile := &models.EmployeeProfile{
		UserID:         user.ID,
		Specialization: req.Specialization,
		ExpectedSalary: req.ExpectedSalary,
		WorkSchedule:   schedulePtr,
	}

	if err := s.employeeRepo.Create(profile); err != nil {
		return nil, errors.New("не удалось создать профиль: " + err.Error())
	}

	if req.SalonID > 0 {
		if err := s.employeeRepo.AssignToSalon(profile.ID, req.SalonID); err != nil {
			return nil, err
		}
	}

	return s.employeeRepo.GetByID(profile.ID)
}

func (s *EmployeeService) UpdateSchedule(userID uint, schedule string) error {
	profile, err := s.employeeRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("профиль сотрудника не найден")
	}
	return s.employeeRepo.UpdateSchedule(profile.ID, schedule)
}

func (s *EmployeeService) UpdateProfile(profile *models.EmployeeProfile) error {
	return s.employeeRepo.Update(profile)
}

func (s *EmployeeService) UpdateEmployee(profileID uint, req dto.UpdateEmployeeRequest) (*models.EmployeeProfile, error) {
	if req.Role == models.RoleClient {
		return nil, errors.New("сотруднику нельзя назначить роль CLIENT")
	}

	profile, err := s.employeeRepo.GetByID(profileID)
	if err != nil {
		return nil, errors.New("профиль сотрудника не найден")
	}

	user, err := s.userRepo.GetByID(profile.UserID)
	if err != nil {
		return nil, errors.New("учётная запись сотрудника не найдена")
	}

	var emailPtr *string
	if req.Email != "" {
		emailPtr = &req.Email
	}

	var schedulePtr *string
	if req.WorkSchedule != "" {
		schedulePtr = &req.WorkSchedule
	}

	user.Username = req.Username
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.Email = emailPtr
	user.Role = req.Role

	profile.Specialization = req.Specialization
	profile.ExpectedSalary = req.ExpectedSalary
	profile.WorkSchedule = schedulePtr

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("не удалось обновить пользователя: " + err.Error())
	}
	if err := s.employeeRepo.Update(profile); err != nil {
		return nil, errors.New("не удалось обновить профиль: " + err.Error())
	}
	if err := s.employeeRepo.ReplaceSalonAssignments(profile.ID, uniqueSalonIDs(req.SalonIDs)); err != nil {
		return nil, errors.New("не удалось обновить закрепление за салонами: " + err.Error())
	}

	return s.employeeRepo.GetByID(profile.ID)
}

func (s *EmployeeService) PatchEmployee(profileID uint, req dto.PatchEmployeeRequest) (*models.EmployeeProfile, error) {
	profile, err := s.employeeRepo.GetByID(profileID)
	if err != nil {
		return nil, errors.New("профиль сотрудника не найден")
	}

	user, err := s.userRepo.GetByID(profile.UserID)
	if err != nil {
		return nil, errors.New("учётная запись сотрудника не найдена")
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Email != nil {
		if *req.Email == "" {
			user.Email = nil
		} else {
			user.Email = req.Email
		}
	}
	if req.Role != nil {
		if *req.Role == models.RoleClient {
			return nil, errors.New("сотруднику нельзя назначить роль CLIENT")
		}
		user.Role = *req.Role
	}
	if req.Specialization != nil {
		profile.Specialization = *req.Specialization
	}
	if req.ExpectedSalary != nil {
		profile.ExpectedSalary = *req.ExpectedSalary
	}
	if req.WorkSchedule != nil {
		if *req.WorkSchedule == "" {
			profile.WorkSchedule = nil
		} else {
			profile.WorkSchedule = req.WorkSchedule
		}
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("не удалось обновить пользователя: " + err.Error())
	}
	if err := s.employeeRepo.Update(profile); err != nil {
		return nil, errors.New("не удалось обновить профиль: " + err.Error())
	}
	if req.SalonIDs != nil {
		if err := s.employeeRepo.ReplaceSalonAssignments(profile.ID, uniqueSalonIDs(*req.SalonIDs)); err != nil {
			return nil, errors.New("не удалось обновить закрепление за салонами: " + err.Error())
		}
	}

	return s.employeeRepo.GetByID(profile.ID)
}

func (s *EmployeeService) FireEmployee(profileID uint) error {
	profile, err := s.employeeRepo.GetByID(profileID)
	if err != nil {
		return errors.New("профиль сотрудника не найден")
	}

	return s.employeeRepo.Fire(profile.ID, profile.UserID)
}

func (s *EmployeeService) AssignToSalon(profileID, salonID uint) error {
	return s.employeeRepo.AssignToSalon(profileID, salonID)
}

func (s *EmployeeService) RemoveFromSalon(profileID, salonID uint) error {
	return s.employeeRepo.RemoveFromSalon(profileID, salonID)
}

func uniqueSalonIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))

	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}

	return result
}
