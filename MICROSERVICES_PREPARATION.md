# 🔧 Подготовка кода к микросервисам

> **Что делать СЕЙЧАС**, чтобы потом легко перейти на микросервисы

---

## 🎯 Immediate Actions (Можно делать прямо сейчас)

### 1. **Разделение на модули** (1-2 недели)

#### Текущая структура:
```
backend/internal/
├── handlers/
├── services/
└── repositories/
```

#### Целевая структура:
```
backend/
├── modules/
│   ├── auth/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   └── domain/
│   ├── player/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   └── domain/
│   ├── match/
│   ├── crosshair/
│   └── build/
└── pkg/              # Shared code
    ├── errors/
    ├── middleware/
    └── database/
```

#### Как мигрировать:
```bash
# 1. Создать структуру модулей
mkdir -p backend/modules/{auth,player,match,crosshair,build}

# 2. Переместить код по модулям
mv backend/internal/handlers/auth_handler.go backend/modules/auth/handler/
mv backend/internal/services/auth_service.go backend/modules/auth/service/
mv backend/internal/repositories/user_repository.go backend/modules/auth/repository/

# 3. Обновить импорты
# Было:
# import "github.com/quenyu/deadlock-stats/internal/services"
# Стало:
# import "github.com/quenyu/deadlock-stats/modules/auth/service"
```

---

### 2. **Создать proto definitions** (1 неделя)

Даже если пока используем REST, proto файлы помогут:
- Документировать API
- Подготовиться к gRPC
- Использовать для validation

#### Установка:
```bash
# Install protoc
# https://grpc.io/docs/protoc-installation/

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install buf (optional but recommended)
# https://docs.buf.build/installation
```

#### Создать proto файлы:
```bash
mkdir -p api/{auth,player,match,common}/v1
```

```protobuf
// api/auth/v1/auth.proto
syntax = "proto3";
package auth.v1;

option go_package = "github.com/quenyu/deadlock-stats/api/auth/v1;authv1";

service AuthService {
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}

message LoginRequest {
  string steam_id = 1;
  string steam_token = 2;
}

message LoginResponse {
  string access_token = 1;
  string refresh_token = 2;
  User user = 3;
}

message User {
  string id = 1;
  string steam_id = 2;
  string nickname = 3;
  string avatar_url = 4;
}

message ValidateTokenRequest {
  string token = 1;
}

message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
}
```

#### Makefile для генерации:
```makefile
# Makefile
.PHONY: proto-gen
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/**/*.proto

.PHONY: proto-lint
proto-lint:
	buf lint

.PHONY: proto-breaking
proto-breaking:
	buf breaking --against '.git#branch=main'
```

---

### 3. **Добавить interface для services** (2-3 дня)

Это позволит легко заменить реализацию на gRPC клиент

#### Сейчас:
```go
// backend/modules/auth/service/auth_service.go
type AuthService struct {
    userRepo *repository.UserRepository
}

func (s *AuthService) Login(ctx context.Context, steamID, token string) (*domain.User, error) {
    // ...
}
```

#### Добавить интерфейс:
```go
// backend/modules/auth/service/interface.go
package service

import (
    "context"
    "github.com/quenyu/deadlock-stats/modules/auth/domain"
)

// AuthService определяет контракт для auth сервиса
// Этот интерфейс используется для:
// 1. Dependency injection
// 2. Mocking в тестах
// 3. Замены на gRPC client в будущем
type AuthService interface {
    Login(ctx context.Context, steamID, token string) (*domain.User, string, string, error)
    ValidateToken(ctx context.Context, token string) (*domain.Claims, error)
    RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

// authServiceImpl - имплементация интерфейса
type authServiceImpl struct {
    userRepo repository.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
    return &authServiceImpl{
        userRepo: userRepo,
        jwtSecret: jwtSecret,
    }
}

func (s *authServiceImpl) Login(ctx context.Context, steamID, token string) (*domain.User, string, string, error) {
    // Implementation
}
```

#### В handler используем интерфейс:
```go
// backend/modules/auth/handler/auth_handler.go
type AuthHandler struct {
    authService service.AuthService  // Интерфейс, не конкретная реализация!
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}
```

#### Потом легко заменить на gRPC:
```go
// backend/modules/auth/client/grpc_client.go
package client

import (
    authv1 "github.com/quenyu/deadlock-stats/api/auth/v1"
    "github.com/quenyu/deadlock-stats/modules/auth/service"
)

// gRPC client реализует тот же интерфейс!
type authGRPCClient struct {
    client authv1.AuthServiceClient
}

func NewAuthGRPCClient(conn *grpc.ClientConn) service.AuthService {
    return &authGRPCClient{
        client: authv1.NewAuthServiceClient(conn),
    }
}

func (c *authGRPCClient) Login(ctx context.Context, steamID, token string) (*domain.User, string, string, error) {
    resp, err := c.client.Login(ctx, &authv1.LoginRequest{
        SteamId: steamID,
        SteamToken: token,
    })
    // Convert proto to domain
}
```

---

### 4. **Event-driven communication** (1-2 недели)

Добавить event bus для асинхронной коммуникации между модулями

#### Установка NATS:
```bash
# docker-compose.yml
nats:
  image: nats:latest
  ports:
    - "4222:4222"
    - "8222:8222"
```

#### Event Bus интерфейс:
```go
// backend/pkg/events/bus.go
package events

import "context"

type EventBus interface {
    Publish(ctx context.Context, topic string, event interface{}) error
    Subscribe(topic string, handler EventHandler) error
}

type EventHandler func(ctx context.Context, event interface{}) error

// NATS implementation
type natsEventBus struct {
    conn *nats.Conn
}

func NewNATSEventBus(url string) (EventBus, error) {
    conn, err := nats.Connect(url)
    if err != nil {
        return nil, err
    }
    return &natsEventBus{conn: conn}, nil
}
```

#### Использование:
```go
// backend/modules/build/service/build_service.go
type BuildService struct {
    repository repository.BuildRepository
    eventBus   events.EventBus
}

func (s *BuildService) CreateBuild(ctx context.Context, req *domain.CreateBuildRequest) (*domain.Build, error) {
    build, err := s.repository.Create(ctx, req)
    if err != nil {
        return nil, err
    }

    // Публикуем событие
    event := &events.BuildCreated{
        BuildID:   build.ID,
        CreatedBy: build.CreatedBy,
        HeroID:    build.HeroID,
        CreatedAt: time.Now(),
    }
    s.eventBus.Publish(ctx, "builds.created", event)

    return build, nil
}

// backend/modules/analytics/service/analytics_service.go
func (s *AnalyticsService) Start(ctx context.Context) error {
    // Подписываемся на события
    s.eventBus.Subscribe("builds.created", s.handleBuildCreated)
    s.eventBus.Subscribe("match.finished", s.handleMatchFinished)
    return nil
}

func (s *AnalyticsService) handleBuildCreated(ctx context.Context, event interface{}) error {
    buildEvent := event.(*events.BuildCreated)
    // Обновляем статистику
    return s.updateBuildStats(ctx, buildEvent)
}
```

---

### 5. **Добавить health checks** (1 день)

```go
// backend/pkg/health/health.go
package health

import (
    "context"
    "database/sql"
    "github.com/redis/go-redis/v9"
)

type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

type HealthStatus struct {
    Status   string            `json:"status"`
    Services map[string]string `json:"services"`
}

func (h *HealthChecker) Check(ctx context.Context) HealthStatus {
    status := HealthStatus{
        Status:   "healthy",
        Services: make(map[string]string),
    }

    // Check database
    if err := h.db.PingContext(ctx); err != nil {
        status.Services["database"] = "unhealthy"
        status.Status = "degraded"
    } else {
        status.Services["database"] = "healthy"
    }

    // Check redis
    if err := h.redis.Ping(ctx).Err(); err != nil {
        status.Services["redis"] = "unhealthy"
        status.Status = "degraded"
    } else {
        status.Services["redis"] = "healthy"
    }

    return status
}

// backend/cmd/main.go
e.GET("/health", func(c echo.Context) error {
    status := healthChecker.Check(c.Request().Context())
    if status.Status == "healthy" {
        return c.JSON(200, status)
    }
    return c.JSON(503, status)
})
```

---

### 6. **Structured logging с context** (1-2 дня)

```go
// backend/pkg/logger/logger.go
package logger

import (
    "context"
    "go.uber.org/zap"
)

type contextKey string

const loggerKey contextKey = "logger"

// Добавить logger в context
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

// Получить logger из context
func FromContext(ctx context.Context) *zap.Logger {
    if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
        return logger
    }
    return zap.L() // fallback to global logger
}

// Middleware для добавления request_id
func LoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            requestID := c.Request().Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = uuid.New().String()
            }

            // Создаем logger с request_id
            reqLogger := logger.With(
                zap.String("request_id", requestID),
                zap.String("method", c.Request().Method),
                zap.String("path", c.Request().URL.Path),
            )

            // Добавляем в context
            ctx := WithLogger(c.Request().Context(), reqLogger)
            c.SetRequest(c.Request().WithContext(ctx))

            return next(c)
        }
    }
}

// Использование в service
func (s *PlayerService) GetPlayerProfile(ctx context.Context, steamID string) (*domain.PlayerProfile, error) {
    logger := logger.FromContext(ctx)
    logger.Info("fetching player profile", zap.String("steam_id", steamID))
    
    // ...
}
```

---

## 📊 Checklist для подготовки

### Архитектура
- [ ] Код разбит на модули (auth, player, match, etc.)
- [ ] Каждый модуль независим (нет circular dependencies)
- [ ] Shared код вынесен в `pkg/`
- [ ] Proto definitions созданы для всех API

### Code Quality
- [ ] Все services имеют интерфейсы
- [ ] Dependency injection через конструкторы
- [ ] Context везде передается первым параметром
- [ ] Structured logging с request_id

### Infrastructure
- [ ] Health checks для всех зависимостей
- [ ] Event bus настроен (NATS/RabbitMQ)
- [ ] Metrics экспортируются (Prometheus)
- [ ] Distributed tracing готов (OpenTelemetry)

### Testing
- [ ] Unit tests для каждого модуля
- [ ] Integration tests с моками
- [ ] Contract tests для proto definitions

### Deployment
- [ ] Docker image для каждого модуля
- [ ] docker-compose.yml с раздельными сервисами
- [ ] Environment-specific configs

---

## 🎯 Пример: Подготовка Auth модуля

### Шаг 1: Создать структуру
```bash
mkdir -p backend/modules/auth/{handler,service,repository,domain}
mkdir -p backend/api/auth/v1
```

### Шаг 2: Proto definition
```protobuf
// api/auth/v1/auth.proto
syntax = "proto3";
package auth.v1;

service AuthService {
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}
```

### Шаг 3: Domain models
```go
// backend/modules/auth/domain/user.go
package domain

import "time"

type User struct {
    ID        string
    SteamID   string
    Nickname  string
    AvatarURL string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Claims struct {
    UserID  string
    SteamID string
}
```

### Шаг 4: Repository interface
```go
// backend/modules/auth/repository/interface.go
package repository

type UserRepository interface {
    FindBySteamID(ctx context.Context, steamID string) (*domain.User, error)
    Create(ctx context.Context, user *domain.User) error
    Update(ctx context.Context, user *domain.User) error
}
```

### Шаг 5: Service interface
```go
// backend/modules/auth/service/interface.go
package service

type AuthService interface {
    Login(ctx context.Context, steamID, token string) (*domain.User, string, string, error)
    ValidateToken(ctx context.Context, token string) (*domain.Claims, error)
}
```

### Шаг 6: Handler
```go
// backend/modules/auth/handler/auth_handler.go
package handler

type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}
```

---

## 🚀 Когда начинать миграцию?

### Метрики для принятия решения:

#### ✅ Готовы к микросервисам когда:
- **100,000+ MAU** (monthly active users)
- **50+ requests/second** на один эндпоинт
- **10+ разработчиков** в команде
- **Database становится bottleneck**
- **Разные требования к scaling** для разных модулей

#### ⏸️ Продолжаем монолит если:
- <50,000 MAU
- Команда <5 человек
- Монолит справляется с нагрузкой
- Нет проблем с deployment

---

## 📚 Полезные ресурсы

- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [Protocol Buffers Tutorial](https://developers.google.com/protocol-buffers/docs/gotutorial)
- [Buf Documentation](https://docs.buf.build/)
- [Microservices Patterns (книга)](https://microservices.io/patterns/index.html)
- [NATS Go Client](https://github.com/nats-io/nats.go)

---

**Следующий шаг**: Начните с рефакторинга в модульный монолит, потом добавьте proto definitions!

