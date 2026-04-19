package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/service"
	"plan-balance-service/pkg/utils"
)

type TransactionHandler struct {
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// Create godoc
// @Summary Create a new transaction
// @Description Create a new income or expense transaction and update asset balance accordingly
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Transaction details"
// @Success 201 {object} dto.TransactionDetailResponse
// @Router /transactions [post]
func (h *TransactionHandler) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	var req dto.CreateTransactionRequest
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

	c.JSON(http.StatusCreated, dto.TransactionDetailResponse{
		Message: "Transaction created successfully",
		Data:    *res,
	})
}

// List godoc
// @Summary Get transactions with filters
// @Description Get a paginated list of transactions with optional filters for type, asset, category, and date range
// @Tags transactions
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param type query string false "Transaction Type"
// @Param asset_id query string false "Asset ID"
// @Param category_id query string false "Category ID"
// @Param start_date query string false "Start Date (ISO)"
// @Param end_date query string false "End Date (ISO)"
// @Success 200 {object} dto.TransactionsListResponse
// @Router /transactions [get]
func (h *TransactionHandler) List(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filters := make(map[string]interface{})
	if val := c.Query("type"); val != "" {
		filters["type"] = val
	}
	if val := c.Query("asset_id"); val != "" {
		uid, _ := uuid.Parse(val)
		filters["asset_id"] = uid
	}
	if val := c.Query("category_id"); val != "" {
		uid, _ := uuid.Parse(val)
		filters["category_id"] = uid
	}
	if val := c.Query("start_date"); val != "" {
		t, _ := time.Parse(time.RFC3339, val)
		filters["start_date"] = t
	}
	if val := c.Query("end_date"); val != "" {
		t, _ := time.Parse(time.RFC3339, val)
		filters["end_date"] = t
	}

	res, err := h.svc.List(c.Request.Context(), userID, filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "FETCH_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetSummary godoc
// @Summary Get monthly summary
// @Description Get total income and expense for a specific month and year
// @Tags transactions
// @Produce json
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year (e.g. 2024)"
// @Success 200 {object} dto.TransactionSummaryResponse
// @Router /transactions/summary [get]
func (h *TransactionHandler) GetSummary(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	res, err := h.svc.GetMonthlySummary(c.Request.Context(), userID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "FETCH_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary Get transaction by ID
// @Description Get detailed information for a single transaction
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} dto.TransactionDetailResponse
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid transaction ID"},
		})
		return
	}

	res, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrTransactionNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "GET_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.TransactionDetailResponse{
		Message: "Transaction fetched successfully",
		Data:    *res,
	})
}

// Update godoc
// @Summary Update transaction details
// @Description Update non-monetary fields (note, category, date) of a transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param request body dto.UpdateTransactionRequest true "Update details"
// @Success 200 {object} dto.TransactionDetailResponse
// @Router /transactions/{id} [put]
func (h *TransactionHandler) Update(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid transaction ID"},
		})
		return
	}

	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrTransactionNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UPDATE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.TransactionDetailResponse{
		Message: "Transaction updated successfully",
		Data:    *res,
	})
}

// Delete godoc
// @Summary Delete transaction and revert balance
// @Description Soft delete a transaction and atomically revert its effect on the asset balance
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} dto.LogoutResponse
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) Delete(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid transaction ID"},
		})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrTransactionNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "DELETE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted and balance reverted successfully",
	})
}
