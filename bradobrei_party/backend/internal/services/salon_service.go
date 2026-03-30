package services

import (
	"fmt"
	"strconv"
	"strings"

	"bradobrei/backend/internal/models"
	"bradobrei/backend/internal/repository"
)

type SalonService struct {
	salonRepo *repository.SalonRepository
}

func NewSalonService(salonRepo *repository.SalonRepository) *SalonService {
	return &SalonService{salonRepo: salonRepo}
}

func (s *SalonService) GetAll() ([]models.Salon, error) {
	return s.salonRepo.GetAll()
}

func (s *SalonService) GetByID(id uint) (*models.Salon, error) {
	return s.salonRepo.GetByID(id)
}

func (s *SalonService) Create(salon *models.Salon) error {
	if err := normalizeSalonLocation(salon); err != nil {
		return err
	}
	return s.salonRepo.Create(salon)
}

func (s *SalonService) Update(salon *models.Salon) error {
	if err := normalizeSalonLocation(salon); err != nil {
		return err
	}
	return s.salonRepo.Update(salon)
}

func (s *SalonService) Delete(id uint) error {
	return s.salonRepo.Delete(id)
}

func (s *SalonService) GetMasters(salonID uint) ([]models.EmployeeProfile, error) {
	return s.salonRepo.GetMasters(salonID)
}

func normalizeSalonLocation(salon *models.Salon) error {
	if salon.Location == nil {
		return nil
	}

	raw := strings.TrimSpace(*salon.Location)
	if raw == "" {
		salon.Location = nil
		return nil
	}

	normalized, err := normalizePoint(raw)
	if err != nil {
		return err
	}

	salon.Location = &normalized
	return nil
}

func normalizePoint(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	upper := strings.ToUpper(trimmed)

	if strings.HasPrefix(upper, "POINT(") && strings.HasSuffix(trimmed, ")") {
		inside := strings.TrimSpace(trimmed[len("POINT(") : len(trimmed)-1])
		if strings.Contains(inside, ",") {
			return normalizeLatLonPair(inside)
		}

		parts := strings.Fields(inside)
		if len(parts) != 2 {
			return "", fmt.Errorf("координаты должны быть в формате \"широта, долгота\" или \"POINT(долгота широта)\"")
		}

		lon, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return "", fmt.Errorf("не удалось разобрать долготу")
		}
		lat, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return "", fmt.Errorf("не удалось разобрать широту")
		}
		if err := validateLatLon(lat, lon); err != nil {
			return "", err
		}

		return fmt.Sprintf("POINT(%g %g)", lon, lat), nil
	}

	return normalizeLatLonPair(trimmed)
}

func normalizeLatLonPair(raw string) (string, error) {
	cleaned := strings.NewReplacer(",", " ", ";", " ").Replace(raw)
	parts := strings.Fields(cleaned)
	if len(parts) != 2 {
		return "", fmt.Errorf("Координаты должны содержать широту и долготу")
	}

	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "", fmt.Errorf("Не удалось разобрать широту")
	}
	lon, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return "", fmt.Errorf("Не удалось разобрать долготу")
	}

	if err := validateLatLon(lat, lon); err != nil {
		return "", err
	}

	return fmt.Sprintf("POINT(%g %g)", lon, lat), nil
}

func validateLatLon(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("широта должна быть в диапазоне от -90 до 90")
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("долгота должна быть в диапазоне от -180 до 180")
	}
	return nil
}
