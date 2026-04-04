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

// Register
// @Summary      Register User Baru
// @Description  Membuat akun user dengan email & password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req body dto.RegisterRequest true "User Registration Request"
// @Success      201 {object} dto.AuthResponse
// @Failure      400,409 {object} dto.ErrorResponse
// @Router       /api/v1/auth/register [post]
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

// Login
// @Summary      Login User
// @Description  Autentikasi user menggunakan email & password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req body dto.LoginRequest true "Login Request"
// @Success      200 {object} dto.AuthResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /api/v1/auth/login [post]
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

// GoogleLogin
// @Summary      Google Login
// @Description  Autentikasi user menggunakan Google ID Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req body dto.GoogleLoginRequest true "Google Login Request"
// @Success      200 {object} dto.AuthResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /api/v1/auth/google [post]
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

// Refresh
// @Summary      Refresh Token
// @Description  Memperbarui Access Token menggunakan Refresh Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success      200 {object} dto.RefreshResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /api/v1/auth/refresh [post]
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

// Logout
// @Summary      Logout
// @Description  Mencabut Refresh Token dan mengakhiri sesi
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        req body dto.RefreshTokenRequest true "Logout Request"
// @Success      200 {object} dto.LogoutResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /api/v1/auth/logout [post]
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
