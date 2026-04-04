package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"plan-balance-service/internal/config"
	"plan-balance-service/internal/db"
	"plan-balance-service/internal/handler"
	"plan-balance-service/internal/middleware"
	"plan-balance-service/internal/repository"
	"plan-balance-service/internal/service"
	"plan-balance-service/pkg/logger"
	"plan-balance-service/pkg/utils"

	_ "plan-balance-service/docs" // Import Swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Plan Balance API
// @version         1.0
// @description     REST API Service for Plan Balance Application.
// @host            localhost:8080
// @BasePath        /
func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Logger
	logger.InitLogger(cfg.Environment)
	defer logger.Sync()
	log := logger.GetLogger()

	// 3. Connect to Database
	db.ConnectDB(cfg.DatabaseURL, log)
	defer db.CloseDB()

	// 4. Initialize Validator
	utils.InitValidator()

	// Repositories
	userRepo := repository.NewUserRepository(db.GetPool())
	authRepo := repository.NewAuthRepository(db.GetPool())
	sessionRepo := repository.NewSessionRepository(db.GetPool())

	// Services
	authService := service.NewAuthService(userRepo, authRepo, sessionRepo, cfg, db.GetPool())

	// Handlers
	authHandler := handler.NewAuthHandler(authService)

	// 5. Setup Gin Router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ZapLogger(log))
	router.Use(middleware.NewCORS())

	// Swagger Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", authHandler.GoogleLogin)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}
	}

	// Health Check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// 6. Start HTTP Server with Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Info("Starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ListenAndServe failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}
