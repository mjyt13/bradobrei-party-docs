package dto

import "bradobrei/backend/internal/models"

// ErrorResponse — единый формат ошибок HTTP API.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type RegisterRequest struct {
	Username string          `json:"username" binding:"required,min=3,max=50"`
	Password string          `json:"password" binding:"required,min=6"`
	FullName string          `json:"full_name" binding:"required"`
	Phone    string          `json:"phone"`
	Email    string          `json:"email" binding:"omitempty,email"`
	Role     models.UserRole `json:"role"`
}

type CreateBookingRequest struct {
	StartTime  string `json:"start_time" binding:"required"`
	SalonID    uint   `json:"salon_id" binding:"required"`
	MasterID   *uint  `json:"master_id"`
	ServiceIDs []uint `json:"service_ids" binding:"required,min=1"`
	Notes      string `json:"notes"`
}

type CreateReviewRequest struct {
	Text   string `json:"text" binding:"required"`
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
}

type HireEmployeeRequest struct {
	Username string          `json:"username" binding:"required,min=3,max=50"`
	Password string          `json:"password" binding:"required,min=6"`
	FullName string          `json:"full_name" binding:"required"`
	Phone    string          `json:"phone"`
	Email    string          `json:"email" binding:"omitempty,email"`
	Role     models.UserRole `json:"role" binding:"required"`

	Specialization string  `json:"specialization"`
	ExpectedSalary float64 `json:"expected_salary"`
	WorkSchedule   string  `json:"work_schedule"`
	SalonID        uint    `json:"salon_id"`

	PasswordHash string `json:"-"`
}

type UpdateScheduleRequest struct {
	Schedule string `json:"schedule" binding:"required"`
}

type AssignSalonRequest struct {
	SalonID uint `json:"salon_id" binding:"required"`
}

type AssignServiceToMasterRequest struct {
	TargetUserID uint `json:"target_user_id" binding:"required"`
}
