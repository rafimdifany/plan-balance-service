package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": "EMAIL_ALREADY_EXISTS"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{
		Message: "User created successfully",
		Data:    *data,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			h.sendError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email or password is incorrect")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Message: "Login successful",
		Data:    *data,
	})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req dto.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.authService.GoogleLogin(c.Request.Context(), req)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Message: "Login successful",
		Data:    *data,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.authService.Refresh(c.Request.Context(), req)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, "INVALID_TOKEN", "Session expired or invalid")
		return
	}

	c.JSON(http.StatusOK, dto.RefreshResponse{
		Message: "Request has been proceed successfully",
		Data:    *data,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.Logout(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LogoutResponse{
		Message: "Request has been proceed successfully",
	})
}

func (h *AuthHandler) sendError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
