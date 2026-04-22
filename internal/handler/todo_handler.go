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

type TodoHandler struct {
	svc service.TodoService
}

func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// Create godoc
// @Summary Create a new todo
// @Description Create a new todo item with title, description, status, priority and due date
// @Tags todos
// @Accept json
// @Produce json
// @Param request body dto.CreateTodoRequest true "Todo details"
// @Success 201 {object} dto.TodoDetailResponse
// @Router /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	var req dto.CreateTodoRequest
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

	c.JSON(http.StatusCreated, dto.TodoDetailResponse{
		Message: "Todo created successfully",
		Data:    *res,
	})
}

// List godoc
// @Summary Get todos with filters
// @Description Get a paginated list of todos with optional filters for category, status, priority, and due date range
// @Tags todos
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param category_id query string false "Category ID"
// @Param status query string false "Status (TODO, IN_PROGRESS, DONE)"
// @Param priority query string false "Priority (LOW, MEDIUM, HIGH)"
// @Param start_due_date query string false "Start Due Date (ISO)"
// @Param end_due_date query string false "End Due Date (ISO)"
// @Success 200 {object} dto.TodosListResponse
// @Router /todos [get]
func (h *TodoHandler) List(c *gin.Context) {
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
	if val := c.Query("category_id"); val != "" {
		uid, _ := uuid.Parse(val)
		filters["category_id"] = uid
	}
	if val := c.Query("status"); val != "" {
		filters["status"] = val
	}
	if val := c.Query("priority"); val != "" {
		filters["priority"] = val
	}
	if val := c.Query("start_due_date"); val != "" {
		t, _ := time.Parse(time.RFC3339, val)
		filters["start_due_date"] = t
	}
	if val := c.Query("end_due_date"); val != "" {
		t, _ := time.Parse(time.RFC3339, val)
		filters["end_due_date"] = t
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

// GetByID godoc
// @Summary Get todo by ID
// @Description Get detailed information for a single todo item
// @Tags todos
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.TodoDetailResponse
// @Router /todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid todo ID"},
		})
		return
	}

	res, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrTodoNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "GET_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.TodoDetailResponse{
		Message: "Todo fetched successfully",
		Data:    *res,
	})
}

// Update godoc
// @Summary Update todo details
// @Description Update all fields of a todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param request body dto.UpdateTodoRequest true "Update details"
// @Success 200 {object} dto.TodoDetailResponse
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid todo ID"},
		})
		return
	}

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrTodoNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UPDATE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.TodoDetailResponse{
		Message: "Todo updated successfully",
		Data:    *res,
	})
}

// PatchStatus godoc
// @Summary Patch todo status
// @Description Update only the status of a todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param request body dto.PatchTodoStatusRequest true "Status update"
// @Success 200 {object} dto.TodoDetailResponse
// @Router /todos/{id}/status [patch]
func (h *TodoHandler) PatchStatus(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid todo ID"},
		})
		return
	}

	var req dto.PatchTodoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.PatchStatus(c.Request.Context(), id, userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrTodoNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "PATCH_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.TodoDetailResponse{
		Message: "Todo status updated successfully",
		Data:    *res,
	})
}

// Delete godoc
// @Summary Delete todo
// @Description Soft delete a todo item
// @Tags todos
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} gin.H
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
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
			Error: dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid todo ID"},
		})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrTodoNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "DELETE_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deleted successfully",
	})
}
