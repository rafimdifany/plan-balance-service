package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/service"
	"plan-balance-service/pkg/utils"
)

type DashboardHandler struct {
	svc service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// GetSummary godoc
// @Summary Get dashboard summary
// @Description Get a comprehensive overview of the user's financial status, including balance, monthly totals, recent transactions, goals, and todos
// @Tags dashboard
// @Produce json
// @Success 200 {object} dto.DashboardResponse
// @Router /dashboard/summary [get]
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	res, err := h.svc.GetSummary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "FETCH_FAILED", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
