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
	"plan-balance-service/internal/logger"
	"plan-balance-service/internal/middleware"
	"plan-balance-service/internal/validator"
)

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
	validator.InitValidator()

	// 5. Setup Gin Router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ZapLogger(log))
	router.Use(middleware.NewCORS())

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
