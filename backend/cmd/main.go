package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quenyu/deadlock-stats/internal/clients/deadlockapi"
	"github.com/quenyu/deadlock-stats/internal/config"
	"github.com/quenyu/deadlock-stats/internal/database/pool"
	"github.com/quenyu/deadlock-stats/internal/handlers"
	customMiddleware "github.com/quenyu/deadlock-stats/internal/middleware"
	"github.com/quenyu/deadlock-stats/internal/middleware/ratelimit"
	"github.com/quenyu/deadlock-stats/internal/middleware/security"
	"github.com/quenyu/deadlock-stats/internal/repositories"
	"github.com/quenyu/deadlock-stats/internal/services"
	"github.com/redis/go-redis/v9"
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

	poolConfig := &pool.Config{
		Host:                cfg.Database.Host,
		Port:                cfg.Database.Port,
		User:                cfg.Database.User,
		Password:            cfg.Database.Password,
		DBName:              cfg.Database.Name,
		SSLMode:             cfg.Database.SSLMode,
		MaxOpenConns:        cfg.Database.Pool.MaxOpenConns,
		MaxIdleConns:        cfg.Database.Pool.MaxIdleConns,
		ConnMaxLifetime:     cfg.Database.Pool.ConnMaxLifetime,
		ConnMaxIdleTime:     cfg.Database.Pool.ConnMaxIdleTime,
		HealthCheckInterval: cfg.Database.Pool.HealthCheckInterval,
		EnableMetrics:       cfg.Database.Pool.EnableMetrics,
	}

	var poolManager *pool.Manager
	for i := 0; i < 5; i++ {
		poolManager, err = pool.NewManager(poolConfig, logger)
		if err == nil {
			break
		}
		logger.Warn("failed to initialize database pool, retrying in 5 seconds...",
			zap.Error(err),
			zap.Int("attempt", i+1),
		)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		logger.Fatal("failed to initialize database pool after multiple retries", zap.Error(err))
	}

	if err := poolManager.WaitForHealthy(30 * time.Second); err != nil {
		logger.Fatal("database did not become healthy", zap.Error(err))
	}

	db := poolManager.DB()
	sqlDB := poolManager.SqlDB()

	defer func() {
		if err := poolManager.Close(); err != nil {
			logger.Error("error closing database pool", zap.Error(err))
		}
	}()

	rdb := connectRedis(cfg.Redis, logger)

	if err := runMigrations(sqlDB, logger); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	staticDataService := services.NewStaticDataService(logger)
	if err := staticDataService.LoadStaticData(); err != nil {
		logger.Fatal("failed to load static data", zap.Error(err))
	}

	userRepository := repositories.NewUserRepository(db)
	playerProfileRepository := repositories.NewPlayerProfilePostgresRepository(db)

	var deadlockAPIClient *deadlockapi.Client
	if cfg.API.EnableRetry {
		deadlockAPIClient = deadlockapi.NewClientWithCustomTimeout(cfg.API.Timeout)
	} else {
		deadlockAPIClient = deadlockapi.NewClient()
	}

	authService := services.NewAuthService(userRepository, cfg, logger)

	playerSearchService := services.NewPlayerSearchService(
		playerProfileRepository,
		userRepository,
		authService,
		deadlockAPIClient,
		rdb,
		cfg.Steam.APIKey,
		logger,
	)

	playerProfileService := services.NewPlayerProfileService(playerProfileRepository, userRepository, authService, deadlockAPIClient, staticDataService, rdb, logger)

	crosshairRepository := repositories.NewCrosshairRepository(db)
	crosshairService := services.NewCrosshairService(crosshairRepository)

	authHandler := handlers.NewAuthHandler(authService, cfg)
	playerSearchHandler := handlers.NewPlayerSearchHandler(playerSearchService, logger)
	playerProfileHandler := handlers.NewPlayerProfileHandler(playerProfileService)
	crosshairHandler := handlers.NewCrosshairHandler(crosshairService)
	healthHandler := handlers.NewHealthHandler(poolManager, logger)
	jwtMiddleware := customMiddleware.NewJWTMiddleware(cfg)

	e := echo.New()

	// Global middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Security middleware (modular: Headers, CSP, CORS, CSRF)
	securityManager := security.NewManager(buildSecurityConfig(cfg, logger))
	e.Use(securityManager.Middleware())

	if cfg.RateLimit.Enabled {
		rateLimitManager, err := ratelimit.NewManager(&ratelimit.ManagerConfig{
			Config: &ratelimit.Config{
				Enabled:           cfg.RateLimit.Enabled,
				Strategy:          ratelimit.Strategy(cfg.RateLimit.Strategy),
				RequestsPerSecond: cfg.RateLimit.RequestsPerSecond,
				Burst:             cfg.RateLimit.Burst,
				UseRedis:          cfg.RateLimit.UseRedis,
				RedisKeyTTL:       cfg.RateLimit.RedisKeyTTL,
				PerEndpoint:       cfg.RateLimit.PerEndpoint,
				Whitelist:         cfg.RateLimit.Whitelist,
				TrustedProxies:    cfg.RateLimit.TrustedProxies,
			},
			RedisClient: rdb,
			Logger:      logger,
		})
		if err != nil {
			logger.Fatal("failed to initialize rate limiter", zap.Error(err))
		}

		rateLimitManager.SetOnLimitReached(func(c echo.Context, key string) {
			logger.Warn("rate limit exceeded",
				zap.String("ip", c.RealIP()),
				zap.String("path", c.Path()),
				zap.String("key", key),
			)
		})

		e.Use(rateLimitManager.Middleware())

		defer func() {
			if err := rateLimitManager.Close(); err != nil {
				logger.Error("error closing rate limiter", zap.Error(err))
			}
		}()

		logger.Info("rate limiting enabled",
			zap.String("strategy", string(rateLimitManager.Config().Strategy)),
			zap.Int("requests_per_second", rateLimitManager.Config().RequestsPerSecond),
		)
	}

	// Unprotected routes
	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1")
	authGroup := v1Group.Group("/auth")
	steamGroup := authGroup.Group("/steam")

	steamGroup.GET("/login", authHandler.LoginHandler)
	steamGroup.GET("/callback", authHandler.CallbackHandler)

	v1Group.GET("/players/search", playerSearchHandler.SearchPlayers)
	v1Group.GET("/players/search/debug", playerSearchHandler.SearchPlayersDebug)
	v1Group.GET("/players/search/autocomplete", playerSearchHandler.SearchPlayersAutocomplete)
	v1Group.GET("/players/search/filters", playerSearchHandler.SearchPlayersWithFilters)
	v1Group.GET("/players/popular", playerSearchHandler.GetPopularPlayers)
	v1Group.GET("/players/recently-active", playerSearchHandler.GetRecentlyActivePlayers)

	v1Group.GET("/players/:steamId", playerProfileHandler.GetPlayerProfileV2)
	v1Group.GET("/players/:steamId/metrics", playerProfileHandler.GetPlayerProfileWithMetrics)
	v1Group.GET("/players/:steamId/matches", playerProfileHandler.GetRecentMatches)
	v1Group.GET("/ranks", staticDataService.GetRanksHandler)

	// Crosshair routes (public)
	v1Group.GET("/crosshairs", crosshairHandler.GetAll)
	v1Group.GET("/crosshairs/:id", crosshairHandler.GetByID)
	v1Group.GET("/authors/:author_id/crosshairs", crosshairHandler.GetByAuthorID)

	// Logout route
	v1Group.POST("/auth/logout", authHandler.LogoutHandler)

	// Protected routes
	protectedGroup := v1Group.Group("")
	protectedGroup.Use(jwtMiddleware.Authorization)
	protectedGroup.GET("/users/me", authHandler.GetUserMe)

	// Protected crosshair routes
	protectedGroup.POST("/crosshairs", crosshairHandler.Create)
	protectedGroup.POST("/crosshairs/:id/like", crosshairHandler.Like)
	protectedGroup.DELETE("/crosshairs/:id/like", crosshairHandler.Unlike)
	protectedGroup.DELETE("/crosshairs/:id", crosshairHandler.Delete)

	e.GET("/health", healthHandler.HealthCheck)
	e.GET("/health/detailed", healthHandler.HealthCheckDetailed)
	e.GET("/metrics/db", healthHandler.MetricsHandler)

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

// connectDB is deprecated - use database.NewPoolManager instead
// Kept for backwards compatibility
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

func buildSecurityConfig(cfg *config.Config, logger *zap.Logger) *security.ManagerConfig {
	// Convert SameSite string to http.SameSite
	var sameSite http.SameSite
	switch cfg.Security.CSRFCookieSameSite {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "lax":
		sameSite = http.SameSiteLaxMode
	case "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteStrictMode
	}

	return &security.ManagerConfig{
		Headers: &security.HeadersConfig{
			HSTSMaxAge:            cfg.Security.HSTSMaxAge,
			HSTSIncludeSubdomains: cfg.Security.HSTSIncludeSubdomains,
			HSTSPreload:           cfg.Security.HSTSPreload,
			XSSProtection:         cfg.Security.XSSProtection,
			XFrameOptions:         cfg.Security.XFrameOptions,
			ContentTypeNoSniff:    cfg.Security.ContentTypeNoSniff,
			ReferrerPolicy:        cfg.Security.ReferrerPolicy,
			PermissionsPolicy:     cfg.Security.PermissionsPolicy,
			XContentTypeOptions:   "nosniff",
			XDNSPrefetchControl:   "off",
			XDownloadOptions:      "noopen",
			XPermittedCrossDomain: "none",
			RemoveHeaders:         []string{"Server", "X-Powered-By"},
			Logger:                logger,
		},
		CSP: &security.CSPConfig{
			Enabled:    cfg.Security.CSPEnabled,
			ReportOnly: cfg.Security.CSPReportOnly,
			Directives: security.DefaultCSPDirectives(),
			Logger:     logger,
		},
		CORS: &security.CORSConfig{
			Enabled:          true,
			AllowOrigins:     []string{cfg.App.ClientURL},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token", "X-Request-ID"},
			ExposeHeaders:    []string{"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
			AllowCredentials: true,
			MaxAge:           86400,
			Logger:           logger,
		},
		CSRF: &security.CSRFConfig{
			Enabled:        cfg.Security.CSRFEnabled,
			CookieSecure:   cfg.Security.CSRFCookieSecure,
			CookieSameSite: sameSite,
			TokenLookup:    "header:X-CSRF-Token",
			SkipPaths:      []string{"/health", "/metrics"},
			Logger:         logger,
		},
		Logger: logger,
	}
}
