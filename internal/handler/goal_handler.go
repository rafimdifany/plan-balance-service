package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/service"
	"plan-balance-service/pkg/utils"
)

type GoalHandler struct {
	svc service.GoalService
}

func NewGoalHandler(svc service.GoalService) *GoalHandler {
	return &GoalHandler{svc: svc}
}

// Create godoc
// @Summary Create a new goal
// @Description Create a new savings or budget goal
// @Tags goals
// @Accept json
// @Produce json
// @Param request body dto.CreateGoalRequest true "Goal details"
// @Success 201 {object} dto.GoalDetailResponse
// @Router /goals [post]
func (h *GoalHandler) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	var req dto.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "CREATE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.GoalDetailResponse{
		Message: "Goal created successfully",
		Data:    *res,
	})
}

// List godoc
// @Summary Get all goals
// @Description Get a list of goals with optional filters for type and active status
// @Tags goals
// @Produce json
// @Param type query string false "Goal Type (SAVINGS, BUDGET)"
// @Param is_active query bool false "Active Status"
// @Success 200 {object} dto.GoalsListResponse
// @Router /goals [get]
func (h *GoalHandler) List(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	filters := make(map[string]interface{})
	if val := c.Query("type"); val != "" {
		filters["type"] = val
	}
	if val := c.Query("is_active"); val != "" {
		filters["is_active"] = val == "true"
	}

	res, err := h.svc.List(c.Request.Context(), userID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "FETCH_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary Get goal by ID
// @Description Get detailed information for a single goal
// @Tags goals
// @Produce json
// @Param id path string true "Goal ID"
// @Success 200 {object} dto.GoalDetailResponse
// @Router /goals/{id} [get]
func (h *GoalHandler) GetByID(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid goal ID"},
		})
		return
	}

	res, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrGoalNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "GET_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.GoalDetailResponse{
		Message: "Goal fetched successfully",
		Data:    *res,
	})
}

// Update godoc
// @Summary Update goal
// @Description Update all fields of a goal
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Param request body dto.UpdateGoalRequest true "Update details"
// @Success 200 {object} dto.GoalDetailResponse
// @Router /goals/{id} [put]
func (h *GoalHandler) Update(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid goal ID"},
		})
		return
	}

	var req dto.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrGoalNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UPDATE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.GoalDetailResponse{
		Message: "Goal updated successfully",
		Data:    *res,
	})
}

// Delete godoc
// @Summary Delete goal
// @Description Soft delete a goal
// @Tags goals
// @Produce json
// @Param id path string true "Goal ID"
// @Success 200 {object} gin.H
// @Router /goals/{id} [delete]
func (h *GoalHandler) Delete(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid goal ID"},
		})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrGoalNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "DELETE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Goal deleted successfully",
	})
}
