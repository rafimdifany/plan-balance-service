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

type AssetHandler struct {
	svc service.AssetService
}

func NewAssetHandler(svc service.AssetService) *AssetHandler {
	return &AssetHandler{svc: svc}
}

// Create godoc
// @Summary Create a new asset
// @Description Create a new financial asset (Bank, Cash, E-wallet) for the authenticated user
// @Tags assets
// @Accept json
// @Produce json
// @Param request body dto.CreateAssetRequest true "Asset details"
// @Success 201 {object} dto.AssetDetailResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /assets [post]
func (h *AssetHandler) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrAssetExists) {
			code = http.StatusConflict
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "CREATE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.AssetDetailResponse{
		Message: "Asset created successfully",
		Data:    *res,
	})
}

// GetAll godoc
// @Summary Get all assets
// @Description Get all active assets for the authenticated user with total balance summary
// @Tags assets
// @Produce json
// @Success 200 {object} dto.AssetsListResponse
// @Router /assets [get]
func (h *AssetHandler) GetAll(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "FETCH_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary Get asset by ID
// @Description Get details of a specific asset
// @Tags assets
// @Produce json
// @Param id path string true "Asset ID"
// @Success 200 {object} dto.AssetDetailResponse
// @Router /assets/{id} [get]
func (h *AssetHandler) GetByID(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid asset ID"},
		})
		return
	}

	res, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrAssetNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "GET_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.AssetDetailResponse{
		Message: "Asset fetched successfully",
		Data:    *res,
	})
}

// Update godoc
// @Summary Update an asset
// @Description Update details of an existing asset (balance cannot be updated via this endpoint)
// @Tags assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Param request body dto.UpdateAssetRequest true "Asset details"
// @Success 200 {object} dto.AssetDetailResponse
// @Router /assets/{id} [put]
func (h *AssetHandler) Update(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid asset ID"},
		})
		return
	}

	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrAssetNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UPDATE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.AssetDetailResponse{
		Message: "Asset updated successfully",
		Data:    *res,
	})
}

// Delete godoc
// @Summary Delete an asset
// @Description Soft delete an asset
// @Tags assets
// @Produce json
// @Param id path string true "Asset ID"
// @Success 200 {object} dto.LogoutResponse
// @Router /assets/{id} [delete]
func (h *AssetHandler) Delete(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid asset ID"},
		})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrAssetNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "DELETE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Asset deleted successfully",
	})
}
