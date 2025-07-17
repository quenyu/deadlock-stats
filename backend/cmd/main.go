package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quenyu/deadlock-stats/internal/config"
	"github.com/quenyu/deadlock-stats/internal/domain"
	"github.com/quenyu/deadlock-stats/internal/handlers"
	customMiddleware "github.com/quenyu/deadlock-stats/internal/middleware"
	"github.com/quenyu/deadlock-stats/internal/repositories"
	"github.com/quenyu/deadlock-stats/internal/services"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadConfig("./internal/config/config.yaml")
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	var db *gorm.DB
	for i := 0; i < 5; i++ {
		db, err = connectDB(cfg.Database)
		if err == nil {
			break
		}
		logger.Warn("failed to connect to database, retrying in 5 seconds...", zap.Error(err))
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		logger.Fatal("failed to connect to database after multiple retries", zap.Error(err))
	}

	// Run database migrations
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		logger.Fatal("failed to migrate database", zap.Error(err))
	}

	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository, cfg, logger)
	authHandler := handlers.NewAuthHandler(authService, cfg)
	jwtMiddleware := customMiddleware.NewJWTMiddleware(cfg)

	e := echo.New()

	// Global middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Unprotected routes
	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1")
	authGroup := v1Group.Group("/auth")
	steamGroup := authGroup.Group("/steam")

	steamGroup.GET("/login", authHandler.LoginHandler)
	steamGroup.GET("/callback", authHandler.CallbackHandler)

	// Protected routes
	profileGroup := v1Group.Group("/profile")
	profileGroup.Use(jwtMiddleware.Authorization)
	profileGroup.GET("/me", authHandler.GetMyProfileHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "ok",
			"version": viper.GetString("app.version"),
		})
	})

	go func() {
		port := viper.GetString("server.port")
		if port == "" {
			port = "8080"
		}

		if err := e.Start(":" + port); err != nil {
			logger.Fatal("shutting down the server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("error during server shutdown", zap.Error(err))
	}
}

func connectDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
