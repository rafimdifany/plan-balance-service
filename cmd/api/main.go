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
	categoryRepo := repository.NewCategoryRepository(db.GetPool())
	assetRepo := repository.NewAssetRepository(db.GetPool())
	transactionRepo := repository.NewTransactionRepository(db.GetPool())
	todoRepo := repository.NewTodoRepository(db.GetPool())
	goalRepo := repository.NewGoalRepository(db.GetPool())

	// Services
	categoryService := service.NewCategoryService(categoryRepo)
	assetService := service.NewAssetService(assetRepo)
	transactionService := service.NewTransactionService(transactionRepo, assetRepo, categoryRepo, db.GetPool())
	todoService := service.NewTodoService(todoRepo, categoryRepo)
	goalService := service.NewGoalService(goalRepo, assetRepo, transactionRepo)
	dashboardService := service.NewDashboardService(assetRepo, transactionRepo, goalRepo, todoRepo, goalService, transactionService)
	authService := service.NewAuthService(userRepo, authRepo, sessionRepo, categoryService, cfg, db.GetPool())

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	assetHandler := handler.NewAssetHandler(assetService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	todoHandler := handler.NewTodoHandler(todoService)
	goalHandler := handler.NewGoalHandler(goalService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

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

		// Protected Routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			categories := protected.Group("/categories")
			{
				categories.POST("", categoryHandler.Create)
				categories.GET("", categoryHandler.GetAll)
				categories.GET("/:id", categoryHandler.GetByID)
				categories.PUT("/:id", categoryHandler.Update)
				categories.DELETE("/:id", categoryHandler.Delete)
			}

			assets := protected.Group("/assets")
			{
				assets.POST("", assetHandler.Create)
				assets.GET("", assetHandler.GetAll)
				assets.GET("/:id", assetHandler.GetByID)
				assets.PUT("/:id", assetHandler.Update)
				assets.DELETE("/:id", assetHandler.Delete)
			}

			transactions := protected.Group("/transactions")
			{
				transactions.POST("", transactionHandler.Create)
				transactions.GET("", transactionHandler.List)
				transactions.GET("/summary", transactionHandler.GetSummary)
				transactions.GET("/:id", transactionHandler.GetByID)
				transactions.PUT("/:id", transactionHandler.Update)
				transactions.DELETE("/:id", transactionHandler.Delete)
			}

			todos := protected.Group("/todos")
			{
				todos.POST("", todoHandler.Create)
				todos.GET("", todoHandler.List)
				todos.GET("/:id", todoHandler.GetByID)
				todos.PUT("/:id", todoHandler.Update)
				todos.PATCH("/:id/status", todoHandler.PatchStatus)
				todos.DELETE("/:id", todoHandler.Delete)
			}

			goals := protected.Group("/goals")
			{
				goals.POST("", goalHandler.Create)
				goals.GET("", goalHandler.List)
				goals.GET("/:id", goalHandler.GetByID)
				goals.PUT("/:id", goalHandler.Update)
				goals.DELETE("/:id", goalHandler.Delete)
			}

			v1.GET("/dashboard/summary", dashboardHandler.GetSummary)
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
