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

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// Create godoc
// @Summary Create a new category
// @Description Create a new income or expense category for the authenticated user
// @Tags categories
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category details"
// @Success 201 {object} dto.CategoryDetailResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrCategoryExists) {
			code = http.StatusConflict
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "CREATE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CategoryDetailResponse{
		Message: "Category created successfully",
		Data:    *res,
	})
}

// GetAll godoc
// @Summary Get all categories
// @Description Get all active categories for the authenticated user
// @Tags categories
// @Produce json
// @Success 200 {object} dto.CategoriesListResponse
// @Router /categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	categories, err := h.svc.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "FETCH_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.CategoriesListResponse{
		Message: "Categories fetched successfully",
		Data:    categories,
	})
}

// GetByID godoc
// @Summary Get category by ID
// @Description Get details of a specific category
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.CategoryDetailResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid category ID"},
		})
		return
	}

	res, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrCategoryNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "GET_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.CategoryDetailResponse{
		Message: "Category fetched successfully",
		Data:    *res,
	})
}

// Update godoc
// @Summary Update a category
// @Description Update details of an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category details"
// @Success 200 {object} dto.CategoryDetailResponse
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid category ID"},
		})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrCategoryNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UPDATE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.CategoryDetailResponse{
		Message: "Category updated successfully",
		Data:    *res,
	})
}

// Delete godoc
// @Summary Delete a category
// @Description Soft delete a category
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.LogoutResponse
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid category ID"},
		})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "DELETE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}
