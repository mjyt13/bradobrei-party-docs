package handlers

import (
	"net/http"
	"strconv"

	"bradobrei/backend/internal/dto"
	"bradobrei/backend/internal/middleware"
	"bradobrei/backend/internal/models"
	"bradobrei/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceService *services.ServiceService
}

func NewServiceHandler(serviceService *services.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceService: serviceService}
}

// GET /api/v1/services
func (h *ServiceHandler) GetAll(c *gin.Context) {
	list, err := h.serviceService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/v1/services/:id
func (h *ServiceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Code: 400})
		return
	}
	svc, err := h.serviceService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "not_found", Code: 404, Message: "Услуга не найдена",
		})
		return
	}
	c.JSON(http.StatusOK, svc)
}

// GET /api/v1/services/my
func (h *ServiceHandler) GetMy(c *gin.Context) {
	claims, _ := middleware.GetCurrentClaims(c)
	list, err := h.serviceService.GetByMaster(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, list)
}

// POST /api/v1/services
func (h *ServiceHandler) Create(c *gin.Context) {
	var svc models.Service
	if err := c.ShouldBindJSON(&svc); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: err.Error(),
		})
		return
	}
	if err := h.serviceService.Create(&svc); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "validation_error", Code: 400, Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, svc)
}

// PUT /api/v1/services/:id
func (h *ServiceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Code: 400})
		return
	}
	existing, err := h.serviceService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "not_found", Code: 404, Message: "Услуга не найдена",
		})
		return
	}
	if err := c.ShouldBindJSON(existing); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: err.Error(),
		})
		return
	}
	existing.ID = uint(id)
	if err := h.serviceService.Update(existing); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "validation_error", Code: 400, Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, existing)
}

// DELETE /api/v1/services/:id
func (h *ServiceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Code: 400})
		return
	}
	if err := h.serviceService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Услуга удалена"})
}

// POST /api/v1/services/:id/assign-master
func (h *ServiceHandler) AssignToMaster(c *gin.Context) {
	claims, _ := middleware.GetCurrentClaims(c)
	serviceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request", Code: 400})
		return
	}

	var req dto.AssignServiceToMasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "bad_request", Code: 400, Message: err.Error(),
		})
		return
	}

	if err := h.serviceService.AddToMaster(
		claims.UserID, claims.Role, req.TargetUserID, uint(serviceID),
	); err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: "forbidden", Code: 403, Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Услуга добавлена мастеру"})
}

// DELETE /api/v1/services/:id/assign-master/:profileId
func (h *ServiceHandler) RemoveFromMaster(c *gin.Context) {
	serviceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	profileID, _ := strconv.ParseUint(c.Param("profileId"), 10, 64)

	if err := h.serviceService.RemoveFromMaster(uint(profileID), uint(serviceID)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal", Code: 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Услуга убрана у мастера"})
}
