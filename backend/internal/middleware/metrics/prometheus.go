package metrics

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP метрики
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Database метрики
	dbConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	dbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	// Cache метрики
	cacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	cacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	// Business метрики
	playerProfilesViewed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "player_profiles_viewed_total",
			Help: "Total number of player profiles viewed",
		},
	)

	crosshairsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "crosshairs_created_total",
			Help: "Total number of crosshairs created",
		},
	)

	apiCallsExternal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_calls_external_total",
			Help: "Total number of external API calls",
		},
		[]string{"api_name", "status"},
	)
)

// PrometheusMiddleware создает middleware для сбора метрик
func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Выполняем запрос
			err := next(c)

			// Собираем метрики
			duration := time.Since(start).Seconds()
			status := strconv.Itoa(c.Response().Status)
			method := c.Request().Method
			endpoint := c.Path()

			// Упрощаем endpoint для группировки
			if endpoint == "" {
				endpoint = "unknown"
			}

			httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
			httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)

			return err
		}
	}
}

// RecordPlayerProfileView записывает просмотр профиля игрока
func RecordPlayerProfileView() {
	playerProfilesViewed.Inc()
}

// RecordCrosshairCreated записывает создание прицела
func RecordCrosshairCreated() {
	crosshairsCreated.Inc()
}

// RecordCacheHit записывает попадание в кэш
func RecordCacheHit(cacheType string) {
	cacheHits.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss записывает промах кэша
func RecordCacheMiss(cacheType string) {
	cacheMisses.WithLabelValues(cacheType).Inc()
}

// RecordExternalAPICall записывает вызов внешнего API
func RecordExternalAPICall(apiName, status string) {
	apiCallsExternal.WithLabelValues(apiName, status).Inc()
}

// UpdateDBConnections обновляет метрики подключений к БД
func UpdateDBConnections(active, idle int) {
	dbConnectionsActive.Set(float64(active))
	dbConnectionsIdle.Set(float64(idle))
}
