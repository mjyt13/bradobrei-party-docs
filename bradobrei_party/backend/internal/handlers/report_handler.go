package handlers

import (
	"net/http"
	"strconv"
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

func parseOptionalUintQuery(c *gin.Context, key string) (uint, error) {
	value := c.Query(key)
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
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

// Reviews godoc
// @Summary Отчёт 2.2.5 Отзывы
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало периода YYYY-MM-DD"
// @Param to query string false "Конец периода YYYY-MM-DD"
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
	c.JSON(http.StatusOK, gin.H{"report": "reviews", "period": gin.H{"from": from, "to": to}, "data": data})
}

// InventoryMovement godoc
// @Summary Отчёт 2.2.6 Ведомость движения ТМЦ
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало периода YYYY-MM-DD"
// @Param to query string false "Конец периода YYYY-MM-DD"
// @Param salon_id query int false "ID салона"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/inventory-movement [get]
func (h *ReportHandler) InventoryMovement(c *gin.Context) {
	from, to, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "Формат даты: YYYY-MM-DD",
		})
		return
	}

	salonID, err := parseOptionalUintQuery(c, "salon_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "salon_id должен быть целым числом",
		})
		return
	}

	data, err := h.reportService.GetInventoryMovement(from, to, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "inventory_movement", "period": gin.H{"from": from, "to": to}, "data": data})
}

// ClientLoyalty godoc
// @Summary Отчёт 2.2.7 Анализ клиентской лояльности
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало периода YYYY-MM-DD"
// @Param to query string false "Конец периода YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/client-loyalty [get]
func (h *ReportHandler) ClientLoyalty(c *gin.Context) {
	from, to, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "Формат даты: YYYY-MM-DD",
		})
		return
	}

	data, err := h.reportService.GetClientLoyalty(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "client_loyalty", "period": gin.H{"from": from, "to": to}, "data": data})
}

// CancelledBookings godoc
// @Summary Отчёт 2.2.8 Отменённые и нереализованные бронирования
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало периода YYYY-MM-DD"
// @Param to query string false "Конец периода YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/cancelled-bookings [get]
func (h *ReportHandler) CancelledBookings(c *gin.Context) {
	from, to, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "Формат даты: YYYY-MM-DD",
		})
		return
	}

	data, err := h.reportService.GetCancelledBookings(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "cancelled_bookings", "period": gin.H{"from": from, "to": to}, "data": data})
}

// FinancialSummary godoc
// @Summary Отчёт 2.2.9 Финансовый отчёт по транзакциям
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param from query string false "Начало периода YYYY-MM-DD"
// @Param to query string false "Конец периода YYYY-MM-DD"
// @Param salon_id query int false "ID салона"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /reports/financial-summary [get]
func (h *ReportHandler) FinancialSummary(c *gin.Context) {
	from, to, err := parsePeriod(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "Формат даты: YYYY-MM-DD",
		})
		return
	}

	salonID, err := parseOptionalUintQuery(c, "salon_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: "salon_id должен быть целым числом",
		})
		return
	}

	data, err := h.reportService.GetFinancialSummary(from, to, salonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": "financial_summary", "period": gin.H{"from": from, "to": to}, "data": data})
}
