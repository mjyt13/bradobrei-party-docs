package handlers

import (
	"net/http"
	"time"

	"bradobrei/backend/internal/dto"
	"bradobrei/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService *services.ReportService
}

func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func parsePeriod(c *gin.Context) (time.Time, time.Time, error) {
	var from, to time.Time
	var err error

	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return from, to, err
		}
	} else {
		from = time.Now().AddDate(0, -1, 0)
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return from, to, err
		}
	} else {
		to = time.Now()
	}

	return from, to, nil
}

// Employees godoc
// @Summary Отчёт 2.2.1 Реестр персонала
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/employees [get]
func (h *ReportHandler) Employees(c *gin.Context) {
	data, err := h.reportService.GetEmployeeList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "employee_list", "data": data})
}

// SalonActivity godoc
// @Summary Отчёт 2.2.2 Активность филиалов
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало периода YYYY-MM-DD"
// @Param to query string false "Конец периода YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/salon-activity [get]
func (h *ReportHandler) SalonActivity(c *gin.Context) {
	from, to, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "Формат даты: YYYY-MM-DD",
		})
		return
	}

	data, err := h.reportService.GetSalonActivity(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "salon_activity", "period": gin.H{"from": from, "to": to}, "data": data})
}

// Доступ: ADMIN, ACCOUNTANT, NETWORK_MANAGER
// ServicePopularity godoc
// @Summary Отчёт 2.2.3 Популярность услуг
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало периода YYYY-MM-DD"
// @Param to query string false "Конец периода YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/service-popularity [get]
func (h *ReportHandler) ServicePopularity(c *gin.Context) {
	from, to, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "Формат даты: YYYY-MM-DD",
		})
		return
	}

	data, err := h.reportService.GetServicePopularity(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "service_popularity", "period": gin.H{"from": from, "to": to}, "data": data})
}

// Доступ: ADMIN, ACCOUNTANT, NETWORK_MANAGER
// MasterActivity godoc
// @Summary Отчёт 2.2.4 Активность мастеров
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало периода YYYY-MM-DD"
// @Param to query string false "Конец периода YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/master-activity [get]
func (h *ReportHandler) MasterActivity(c *gin.Context) {
	from, to, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "Формат даты: YYYY-MM-DD",
		})
		return
	}

	data, err := h.reportService.GetMasterActivity(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "master_activity", "period": gin.H{"from": from, "to": to}, "data": data})
}

// Доступ: ADMIN
// Reviews godoc
// @Summary Отчёт 2.2.5 Отзывы
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/reviews [get]
func (h *ReportHandler) Reviews(c *gin.Context) {
	from, to, _ := parsePeriod(c)

	data, err := h.reportService.GetReviews(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "reviews", "data": data})
}
