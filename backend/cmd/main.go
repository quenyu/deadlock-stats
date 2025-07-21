package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"github.com/quenyu/deadlock-stats/internal/config"
	"github.com/quenyu/deadlock-stats/internal/handlers"
	customMiddleware "github.com/quenyu/deadlock-stats/internal/middleware"
	"github.com/quenyu/deadlock-stats/internal/repositories"
	"github.com/quenyu/deadlock-stats/internal/services"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger, _ := zap.NewDevelopment()
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

	rdb := connectRedis(cfg.Redis, logger)

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get underlying sql.DB", zap.Error(err))
	}

	if err := runMigrations(sqlDB, logger); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	staticDataService := services.NewStaticDataService(logger)
	if err := staticDataService.LoadStaticData(); err != nil {
		logger.Fatal("failed to load static data", zap.Error(err))
	}

	userRepository := repositories.NewUserRepository(db)
	playerProfileRepository := repositories.NewPlayerProfilePostgresRepository(db)
	deadlockAPIClient := deadlockapi.NewClient()

	authService := services.NewAuthService(userRepository, cfg, logger)
	playerProfileService := services.NewPlayerProfileService(playerProfileRepository, deadlockAPIClient, staticDataService, rdb, logger)

	authHandler := handlers.NewAuthHandler(authService, cfg)
	playerProfileHandler := handlers.NewPlayerProfileHandler(playerProfileService)
	jwtMiddleware := customMiddleware.NewJWTMiddleware(cfg)

	e := echo.New()

	// Global middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Unprotected routes
	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1")
	authGroup := v1Group.Group("/auth")
	steamGroup := authGroup.Group("/steam")

	steamGroup.GET("/login", authHandler.LoginHandler)
	steamGroup.GET("/callback", authHandler.CallbackHandler)

	v1Group.GET("/players/search", playerProfileHandler.SearchPlayers)
	v1Group.GET("/players/:steamId", playerProfileHandler.GetPlayerProfileV2)
	v1Group.GET("/players/:steamId/matches", playerProfileHandler.GetRecentMatches)

	// Logout route
	authGroup.GET("/logout", authHandler.LogoutHandler)

	// Protected routes
	protectedGroup := v1Group.Group("")
	protectedGroup.Use(jwtMiddleware.Authorization)
	protectedGroup.GET("/users/me", authHandler.GetUserMe)

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

func connectRedis(cfg config.RedisConfig, logger *zap.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}

	return rdb
}

func connectDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(gorm_postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB, logger *zap.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	logger.Info("database migrations applied successfully")
	return nil
}
