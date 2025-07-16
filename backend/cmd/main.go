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
	"github.com/quenyu/deadlock-stats/internal/handlers"
	"github.com/quenyu/deadlock-stats/internal/repository"
	"github.com/quenyu/deadlock-stats/internal/services"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	db, err := connectDB(cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	userRepository := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepository, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1")
	authGroup := v1Group.Group("/auth")
	steamGroup := authGroup.Group("/steam")

	steamGroup.GET("/login", authHandler.LoginHandler)
	steamGroup.GET("/callback", authHandler.CallbackHandler)

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
