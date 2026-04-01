# 🏗️ План миграции на Микросервисы + gRPC

> **Цель**: Переход от монолитной архитектуры к микросервисам с gRPC коммуникацией

---

## 📊 Текущая vs Целевая архитектура

### Текущая архитектура (Монолит)
```
┌─────────────────────────────────────┐
│         Frontend (React)            │
└────────────┬────────────────────────┘
             │ HTTP/REST
┌────────────▼────────────────────────┐
│      Backend Monolith (Go)          │
│  ┌──────────────────────────────┐   │
│  │ Handlers (HTTP)              │   │
│  ├──────────────────────────────┤   │
│  │ Services (Business Logic)    │   │
│  ├──────────────────────────────┤   │
│  │ Repositories (Data Access)   │   │
│  └──────────────────────────────┘   │
└────────┬────────────┬────────────────┘
         │            │
    ┌────▼───┐   ┌────▼─────┐
    │ Redis  │   │  Postgres│
    └────────┘   └──────────┘
```

### Целевая архитектура (Микросервисы + gRPC)
```
┌─────────────────────────────────────┐
│         Frontend (React)            │
└────────────┬────────────────────────┘
             │ HTTP/REST
┌────────────▼────────────────────────┐
│      API Gateway (Go)               │
│    - REST to gRPC translation       │
│    - Authentication                 │
│    - Rate limiting                  │
│    - Request routing                │
└───┬────┬────┬────┬────┬─────────────┘
    │    │    │    │    │ gRPC
    │    │    │    │    │
┌───▼──┐ │ ┌──▼──┐│ ┌──▼────────┐
│Auth  │ │ │Player││ │Crosshair  │
│Service│ │Service││ │Service    │
└──┬───┘ │└─┬───┘│ └─┬─────────┘
   │     │  │    │   │
┌──▼──┐  │┌─▼──┐ │┌──▼────┐
│Postgres││Redis││Postgres│
└─────┘  │└────┘ │└───────┘
         │       │
    ┌────▼──┐ ┌──▼────────┐
    │Match  │ │Analytics  │
    │Service│ │Service    │
    └───┬───┘ └─┬─────────┘
        │       │
    ┌───▼──┐ ┌──▼──────┐
    │Postgres││TimescaleDB│
    └──────┘ └─────────┘

┌─────────────────────────────────────┐
│   Message Broker (NATS/RabbitMQ)   │
│   - Event-driven communication      │
│   - Async processing                │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│   Service Mesh (Optional: Istio)    │
│   - Service discovery               │
│   - Load balancing                  │
│   - Observability                   │
└─────────────────────────────────────┘
```

---

## 🎯 Разбивка на микросервисы

### 1. **Auth Service** 👤
**Ответственность:**
- Steam OpenID authentication
- JWT token generation/validation
- User sessions management
- Permission checks

**API (gRPC):**
```protobuf
service AuthService {
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
  rpc GetUserPermissions(GetUserPermissionsRequest) returns (GetUserPermissionsResponse);
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
```

**База данных:**
- PostgreSQL (users, sessions, permissions)
- Redis (token blacklist, session cache)

---

### 2. **Player Service** 🎮
**Ответственность:**
- Player profiles
- Player statistics
- MMR tracking
- Personal records

**API (gRPC):**
```protobuf
service PlayerService {
  rpc GetPlayerProfile(GetPlayerProfileRequest) returns (PlayerProfileResponse);
  rpc GetPlayerStats(GetPlayerStatsRequest) returns (PlayerStatsResponse);
  rpc SearchPlayers(SearchPlayersRequest) returns (SearchPlayersResponse);
  rpc UpdatePlayerStats(UpdatePlayerStatsRequest) returns (UpdatePlayerStatsResponse);
  rpc GetMMRHistory(GetMMRHistoryRequest) returns (GetMMRHistoryResponse);
}

message GetPlayerProfileRequest {
  string steam_id = 1;
  bool include_stats = 2;
  bool include_mmr_history = 3;
}

message PlayerProfileResponse {
  PlayerProfile profile = 1;
  repeated HeroStat hero_stats = 2;
  repeated MMRPoint mmr_history = 3;
}
```

**База данных:**
- PostgreSQL (player_stats, hero_stats, personal_records)
- Redis (profile cache, stats cache)

---

### 3. **Match Service** 🎲
**Ответственность:**
- Match history
- Match details
- Player performance in matches
- Match statistics aggregation

**API (gRPC):**
```protobuf
service MatchService {
  rpc GetMatchHistory(GetMatchHistoryRequest) returns (GetMatchHistoryResponse);
  rpc GetMatchDetails(GetMatchDetailsRequest) returns (GetMatchDetailsResponse);
  rpc GetPlayerMatchStats(GetPlayerMatchStatsRequest) returns (GetPlayerMatchStatsResponse);
  rpc StreamLiveMatches(StreamLiveMatchesRequest) returns (stream LiveMatchUpdate);
}

message GetMatchHistoryRequest {
  string steam_id = 1;
  int32 limit = 2;
  int32 offset = 3;
  MatchFilters filters = 4;
}
```

**База данных:**
- PostgreSQL (matches, player_match_stats)
- Redis (recent matches cache)

---

### 4. **Crosshair Service** 🎯
**Ответственность:**
- Crosshair CRUD
- Crosshair likes/votes
- Crosshair gallery
- User submissions

**API (gRPC):**
```protobuf
service CrosshairService {
  rpc CreateCrosshair(CreateCrosshairRequest) returns (CrosshairResponse);
  rpc GetCrosshair(GetCrosshairRequest) returns (CrosshairResponse);
  rpc UpdateCrosshair(UpdateCrosshairRequest) returns (CrosshairResponse);
  rpc DeleteCrosshair(DeleteCrosshairRequest) returns (DeleteCrosshairResponse);
  rpc ListCrosshairs(ListCrosshairsRequest) returns (ListCrosshairsResponse);
  rpc LikeCrosshair(LikeCrosshairRequest) returns (LikeCrosshairResponse);
}
```

**База данных:**
- PostgreSQL (crosshairs, crosshair_likes)

---

### 5. **Build Service** 🏗️
**Ответственность:**
- Hero builds CRUD
- Build voting/comments
- Build recommendations
- Build analytics

**API (gRPC):**
```protobuf
service BuildService {
  rpc CreateBuild(CreateBuildRequest) returns (BuildResponse);
  rpc GetBuild(GetBuildRequest) returns (BuildResponse);
  rpc ListBuilds(ListBuildsRequest) returns (ListBuildsResponse);
  rpc VoteBuild(VoteBuildRequest) returns (VoteBuildResponse);
  rpc CommentBuild(CommentBuildRequest) returns (CommentBuildResponse);
  rpc GetRecommendedBuilds(GetRecommendedBuildsRequest) returns (RecommendedBuildsResponse);
}
```

**База данных:**
- PostgreSQL (builds, votes, comments)
- Redis (popular builds cache)

---

### 6. **Analytics Service** 📊
**Ответственность:**
- Meta analysis
- Win rates calculation
- Tier lists
- Trend analysis
- Leaderboards

**API (gRPC):**
```protobuf
service AnalyticsService {
  rpc GetHeroWinRates(GetHeroWinRatesRequest) returns (HeroWinRatesResponse);
  rpc GetMetaAnalysis(GetMetaAnalysisRequest) returns (MetaAnalysisResponse);
  rpc GetLeaderboard(GetLeaderboardRequest) returns (LeaderboardResponse);
  rpc GetTrendAnalysis(GetTrendAnalysisRequest) returns (TrendAnalysisResponse);
  rpc StreamMetrics(StreamMetricsRequest) returns (stream MetricsUpdate);
}
```

**База данных:**
- TimescaleDB (time-series data для метрик)
- Redis (leaderboard cache)

---

### 7. **External API Service** 🌐
**Ответственность:**
- Интеграция с Deadlock API
- Data synchronization
- API rate limiting
- Data transformation

**API (gRPC):**
```protobuf
service ExternalAPIService {
  rpc FetchPlayerData(FetchPlayerDataRequest) returns (FetchPlayerDataResponse);
  rpc FetchMatchData(FetchMatchDataRequest) returns (FetchMatchDataResponse);
  rpc SyncPlayerStats(SyncPlayerStatsRequest) returns (SyncPlayerStatsResponse);
  rpc GetAPIStatus(GetAPIStatusRequest) returns (GetAPIStatusResponse);
}
```

**База данных:**
- Redis (API cache, rate limit counters)

---

## 🔄 План поэтапной миграции

### **Phase 1: Подготовка (2-3 недели)**

#### Week 1: Инфраструктура
```bash
# 1. Установить инструменты
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 2. Создать структуру проекта
deadlock-stats/
├── api/                    # Proto definitions
│   ├── auth/
│   │   └── v1/
│   │       └── auth.proto
│   ├── player/
│   │   └── v1/
│   │       └── player.proto
│   └── common/
│       └── v1/
│           └── common.proto
├── services/              # Микросервисы
│   ├── auth/
│   ├── player/
│   ├── match/
│   └── ...
├── gateway/               # API Gateway
├── pkg/                   # Shared packages
│   ├── grpcutil/
│   ├── middleware/
│   └── errors/
└── docker-compose.yml

# 3. Настроить buf для управления proto
# buf.yaml
version: v1
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
```

#### Week 2-3: Proto Definitions
```protobuf
// api/common/v1/common.proto
syntax = "proto3";
package common.v1;

option go_package = "github.com/quenyu/deadlock-stats/api/common/v1;commonv1";

message Empty {}

message Error {
  string code = 1;
  string message = 2;
  map<string, string> details = 3;
}

message PaginationRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message PaginationResponse {
  int32 page = 1;
  int32 page_size = 2;
  int32 total_items = 3;
  int32 total_pages = 4;
}

// api/player/v1/player.proto
syntax = "proto3";
package player.v1;

import "common/v1/common.proto";
import "google/protobuf/timestamp.proto";

option go_package = "github.com/quenyu/deadlock-stats/api/player/v1;playerv1";

service PlayerService {
  rpc GetPlayerProfile(GetPlayerProfileRequest) returns (GetPlayerProfileResponse);
  rpc SearchPlayers(SearchPlayersRequest) returns (SearchPlayersResponse);
  rpc UpdatePlayerStats(UpdatePlayerStatsRequest) returns (UpdatePlayerStatsResponse);
}

message PlayerProfile {
  string steam_id = 1;
  string nickname = 2;
  string avatar_url = 3;
  int32 current_mmr = 4;
  int32 peak_mmr = 5;
  google.protobuf.Timestamp last_updated = 6;
}
```

---

### **Phase 2: Первый микросервис - Auth Service (3-4 недели)**

#### Структура Auth Service
```
services/auth/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── server/
│   │   └── server.go        # gRPC server implementation
│   ├── service/
│   │   └── auth_service.go  # Business logic
│   ├── repository/
│   │   └── user_repository.go
│   └── config/
│       └── config.go
├── proto/
│   └── auth/
│       └── v1/
│           └── auth.proto
├── Dockerfile
└── go.mod
```

#### Реализация
```go
// services/auth/internal/server/server.go
package server

import (
    "context"
    authv1 "github.com/quenyu/deadlock-stats/api/auth/v1"
    "google.golang.org/grpc"
)

type AuthServer struct {
    authv1.UnimplementedAuthServiceServer
    authService *service.AuthService
}

func NewAuthServer(authService *service.AuthService) *AuthServer {
    return &AuthServer{
        authService: authService,
    }
}

func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
    user, accessToken, refreshToken, err := s.authService.Login(ctx, req.SteamId, req.SteamToken)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
    }

    return &authv1.LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User: &authv1.User{
            Id:       user.ID,
            SteamId:  user.SteamID,
            Nickname: user.Nickname,
        },
    }, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
    claims, err := s.authService.ValidateToken(ctx, req.Token)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }

    return &authv1.ValidateTokenResponse{
        Valid:  true,
        UserId: claims.UserID,
    }, nil
}

// services/auth/cmd/server/main.go
func main() {
    cfg := config.Load()
    
    // Setup dependencies
    db := setupDatabase(cfg)
    repo := repository.NewUserRepository(db)
    authService := service.NewAuthService(repo, cfg.JWT)
    
    // Create gRPC server
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
            grpc_recovery.UnaryServerInterceptor(),
            grpc_prometheus.UnaryServerInterceptor,
            grpc_zap.UnaryServerInterceptor(logger),
        )),
    )
    
    // Register service
    authServer := server.NewAuthServer(authService)
    authv1.RegisterAuthServiceServer(grpcServer, authServer)
    
    // Start server
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    log.Printf("Auth service listening on :%d", cfg.Port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

---

### **Phase 3: API Gateway (2-3 недели)**

#### Задачи Gateway
1. **REST to gRPC translation**
2. **Authentication & Authorization**
3. **Rate limiting**
4. **Request aggregation**
5. **Response caching**

#### Реализация
```go
// gateway/internal/handler/player_handler.go
package handler

import (
    "context"
    "github.com/labstack/echo/v4"
    playerv1 "github.com/quenyu/deadlock-stats/api/player/v1"
    "net/http"
)

type PlayerHandler struct {
    playerClient playerv1.PlayerServiceClient
    authClient   authv1.AuthServiceClient
}

func (h *PlayerHandler) GetPlayerProfile(c echo.Context) error {
    // 1. Validate JWT token via Auth Service
    token := extractToken(c)
    _, err := h.authClient.ValidateToken(c.Request().Context(), &authv1.ValidateTokenRequest{
        Token: token,
    })
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
    }

    // 2. Call Player Service
    steamID := c.Param("steamId")
    resp, err := h.playerClient.GetPlayerProfile(c.Request().Context(), &playerv1.GetPlayerProfileRequest{
        SteamId:           steamID,
        IncludeStats:      true,
        IncludeMmrHistory: true,
    })
    if err != nil {
        return handleGRPCError(c, err)
    }

    // 3. Convert protobuf to JSON
    return c.JSON(http.StatusOK, convertToJSON(resp))
}

// gateway/internal/handler/aggregation_handler.go
// Агрегация данных из нескольких сервисов
func (h *AggregationHandler) GetExtendedPlayerProfile(c echo.Context) error {
    ctx := c.Request().Context()
    steamID := c.Param("steamId")

    // Параллельные запросы к разным сервисам
    var (
        profileResp *playerv1.GetPlayerProfileResponse
        matchResp   *matchv1.GetMatchHistoryResponse
        buildsResp  *buildv1.ListBuildsResponse
    )

    g, ctx := errgroup.WithContext(ctx)

    g.Go(func() error {
        var err error
        profileResp, err = h.playerClient.GetPlayerProfile(ctx, &playerv1.GetPlayerProfileRequest{
            SteamId: steamID,
        })
        return err
    })

    g.Go(func() error {
        var err error
        matchResp, err = h.matchClient.GetMatchHistory(ctx, &matchv1.GetMatchHistoryRequest{
            SteamId: steamID,
            Limit:   5,
        })
        return err
    })

    g.Go(func() error {
        var err error
        buildsResp, err = h.buildClient.ListBuilds(ctx, &buildv1.ListBuildsRequest{
            CreatedBy: steamID,
            Limit:     10,
        })
        return err
    })

    if err := g.Wait(); err != nil {
        return handleGRPCError(c, err)
    }

    // Агрегируем результаты
    response := map[string]interface{}{
        "profile":       profileResp.Profile,
        "recent_matches": matchResp.Matches,
        "builds":        buildsResp.Builds,
    }

    return c.JSON(http.StatusOK, response)
}
```

---

### **Phase 4: Остальные сервисы (4-6 недель)**

#### Порядок миграции
1. **Player Service** (Week 1-2)
2. **Match Service** (Week 2-3)
3. **Crosshair Service** (Week 3-4)
4. **Build Service** (Week 4-5)
5. **Analytics Service** (Week 5-6)

#### Шаблон для каждого сервиса
```bash
# 1. Создать proto definitions
# 2. Сгенерировать Go code
make proto-gen

# 3. Перенести business logic из монолита
# 4. Добавить gRPC server
# 5. Написать тесты
# 6. Создать Docker image
# 7. Обновить docker-compose.yml
```

---

### **Phase 5: Event-Driven Architecture (опционально, 2-3 недели)**

#### Добавление Message Broker (NATS/RabbitMQ)
```go
// Пример: При создании билда отправляем событие
func (s *BuildService) CreateBuild(ctx context.Context, req *buildv1.CreateBuildRequest) (*buildv1.BuildResponse, error) {
    // 1. Создаем билд
    build, err := s.repository.Create(ctx, req)
    if err != nil {
        return nil, err
    }

    // 2. Публикуем событие
    event := &events.BuildCreated{
        BuildId:   build.Id,
        CreatedBy: build.CreatedBy,
        HeroId:    build.HeroId,
        CreatedAt: timestamppb.Now(),
    }
    s.eventPublisher.Publish("builds.created", event)

    return &buildv1.BuildResponse{Build: build}, nil
}

// Analytics Service подписывается на события
func (s *AnalyticsService) SubscribeToBuildEvents() {
    s.eventSubscriber.Subscribe("builds.created", func(event *events.BuildCreated) {
        // Обновляем статистику
        s.updateBuildStatistics(event)
    })
}
```

---

## 🛠️ Технологический стек

### Core Technologies
- **gRPC**: google.golang.org/grpc
- **Protocol Buffers**: google.golang.org/protobuf
- **Buf**: buf.build (управление proto файлами)
- **grpc-gateway**: генерация REST API из proto (опционально)

### Middleware & Tools
```go
// grpc middleware
import (
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
    grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
    grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
    grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
)

server := grpc.NewServer(
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        grpc_recovery.UnaryServerInterceptor(),
        grpc_prometheus.UnaryServerInterceptor,
        grpc_zap.UnaryServerInterceptor(logger),
        grpc_auth.UnaryServerInterceptor(authFunc),
    )),
)
```

### Service Discovery
- **Consul**: Service registry
- **etcd**: Конфигурация
- **Kubernetes DNS**: Для K8s deployment

### Observability
- **Prometheus**: Metrics
- **Grafana**: Dashboards
- **Jaeger/Tempo**: Distributed tracing
- **Loki**: Log aggregation

---

## 📦 Docker Compose для разработки

```yaml
# docker-compose.microservices.yml
version: '3.8'

services:
  # Databases
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: deadlock_stats
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  # Message Broker
  nats:
    image: nats:latest
    ports:
      - "4222:4222"
      - "8222:8222"

  # Service Discovery
  consul:
    image: consul:latest
    ports:
      - "8500:8500"
    command: agent -server -ui -bootstrap-expect=1 -client=0.0.0.0

  # Microservices
  auth-service:
    build: ./services/auth
    ports:
      - "50051:50051"
    depends_on:
      - postgres
      - redis
      - consul
    environment:
      - DATABASE_URL=postgres://postgres:postgres@postgres:5432/deadlock_stats
      - REDIS_URL=redis://redis:6379
      - CONSUL_URL=consul:8500

  player-service:
    build: ./services/player
    ports:
      - "50052:50052"
    depends_on:
      - postgres
      - redis
      - consul

  match-service:
    build: ./services/match
    ports:
      - "50053:50053"
    depends_on:
      - postgres
      - redis

  crosshair-service:
    build: ./services/crosshair
    ports:
      - "50054:50054"
    depends_on:
      - postgres

  build-service:
    build: ./services/build
    ports:
      - "50055:50055"
    depends_on:
      - postgres
      - nats

  analytics-service:
    build: ./services/analytics
    ports:
      - "50056:50056"
    depends_on:
      - postgres
      - redis

  # API Gateway
  gateway:
    build: ./gateway
    ports:
      - "8080:8080"
    depends_on:
      - auth-service
      - player-service
      - match-service
      - crosshair-service
      - build-service
      - analytics-service
    environment:
      - AUTH_SERVICE_URL=auth-service:50051
      - PLAYER_SERVICE_URL=player-service:50052
      - MATCH_SERVICE_URL=match-service:50053

  # Frontend
  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - gateway

  # Monitoring
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3001:3000"
    depends_on:
      - prometheus

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
      - "14268:14268"
```

---

## 🎯 Преимущества микросервисной архитектуры

### ✅ Плюсы
1. **Независимое развёртывание**
   - Каждый сервис можно деплоить отдельно
   - Нет downtime всего приложения

2. **Масштабируемость**
   - Можно масштабировать только нужные сервисы
   - Player Service нагружен? → Запускаем ещё 5 инстансов

3. **Технологическая гибкость**
   - Analytics Service можно написать на Python
   - Player Service на Go
   - Build Service на Rust (если нужно)

4. **Изоляция отказов**
   - Если упал Crosshair Service, остальные работают

5. **Команды работают независимо**
   - Team A работает над Player Service
   - Team B над Analytics
   - Нет конфликтов в коде

6. **Легче тестировать**
   - Unit tests для каждого сервиса
   - Integration tests через gRPC

### ⚠️ Минусы
1. **Сложность инфраструктуры**
   - Нужен service discovery, load balancer
   - Monitoring сложнее

2. **Distributed transactions**
   - Нет ACID гарантий между сервисами
   - Нужны Saga patterns, compensating transactions

3. **Network overhead**
   - Больше сетевых вызовов
   - Latency может вырасти

4. **Debugging сложнее**
   - Нужен distributed tracing
   - Логи разбросаны по сервисам

5. **Data consistency**
   - Eventual consistency вместо strong consistency
   - Нужны стратегии синхронизации

---

## 📊 Когда переходить на микросервисы?

### ✅ Стоит переходить если:
- **>10 разработчиков** в команде
- **>100k пользователей**
- **Явные bounded contexts** в приложении
- **Разные требования к scaling** для разных модулей
- **Нужна fault isolation**

### ❌ НЕ стоит если:
- Команда <5 человек
- <10k пользователей
- Монолит справляется с нагрузкой
- Нет проблем с deployment

### 🎯 Для Deadlock Stats
**Рекомендация**: Начать с **модульного монолита**, переходить на микросервисы когда:
1. **100k+ MAU** (monthly active users)
2. **10+ разработчиков**
3. **Явные проблемы с масштабированием**

---

## 🚀 Альтернатива: Модульный монолит

Вместо полноценных микросервисов, можно сделать **модульный монолит**:

```
backend/
├── modules/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── domain.go
│   ├── player/
│   ├── match/
│   ├── crosshair/
│   └── build/
├── pkg/
│   ├── grpc/      # gRPC инфраструктура (готово к разделению)
│   ├── events/    # Event bus
│   └── database/
└── cmd/
    └── server/
        └── main.go
```

**Преимущества**:
- ✅ Простое развёртывание (один binary)
- ✅ Нет network overhead
- ✅ Легко рефакторить
- ✅ Готов к split на микросервисы

---

## 📝 Заключение

### Рекомендуемый план для Deadlock Stats:

#### **Сейчас (0-50k users):**
1. ✅ Продолжаем монолит
2. ✅ Улучшаем архитектуру (чистая архитектура)
3. ✅ Пишем proto definitions (готовимся к gRPC)

#### **Позже (50-100k users):**
1. Переходим на **модульный монолит** с gRPC ready кодом
2. Добавляем event-driven коммуникацию внутри монолита

#### **Будущее (100k+ users):**
1. Выделяем первый сервис (Auth или Analytics)
2. Постепенно переводим остальные модули
3. Полноценная микросервисная архитектура

### Следующие шаги:
1. ✅ Завершить все CRITICAL fixes
2. ✅ Написать proto definitions (даже для монолита)
3. ✅ Внедрить event bus (NATS) внутри монолита
4. ✅ Добавить service layer абстракции

Хотите я создам детальный план по подготовке текущего кода к будущей миграции?

