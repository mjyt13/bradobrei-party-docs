package handlers

import (
	"net/http"
	"strconv"

	"bradobrei/backend/internal/dto"
	"bradobrei/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService *services.InventoryService
}

func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// GetBySalon godoc
// @Summary Остатки материалов по салону
// @Description Возвращает фактические остатки материалов на складе выбранного салона.
// @Tags inventories
// @Produce json
// @Security BearerAuth
// @Param salonId path int true "ID салона"
// @Success 200 {array} models.Inventory
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /inventories/salon/{salonId} [get]
func (h *InventoryHandler) GetBySalon(c *gin.Context) {
	salonID, err := strconv.ParseUint(c.Param("salonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Code: 400})
		return
	}

	items, err := h.inventoryService.GetBySalon(uint(salonID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}

	c.JSON(http.StatusOK, items)
}

// SetQuantity godoc
// @Summary Обновить остаток материала на складе
// @Description Устанавливает фактическое количество материала на складе конкретного салона по ID материала.
// @Tags inventories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param salonId path int true "ID салона"
// @Param materialId path int true "ID материала"
// @Param request body dto.SetInventoryQuantityRequest true "Новое количество материала"
// @Success 200 {object} models.Inventory
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /inventories/salon/{salonId}/material/{materialId} [put]
func (h *InventoryHandler) SetQuantity(c *gin.Context) {
	salonID, err := strconv.ParseUint(c.Param("salonId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Code: 400})
		return
	}

	materialID, err := strconv.ParseUint(c.Param("materialId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Code: 400})
		return
	}

	var req dto.SetInventoryQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	item, err := h.inventoryService.SetQuantity(uint(salonID), uint(materialID), req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "inventory_update_failed",
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, item)
}
