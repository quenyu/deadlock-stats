# ✅ TODO Summary - Что нужно сделать

> **Быстрый обзор приоритетных задач для Deadlock Stats**

## 🎯 Как работать с этим списком

1. Задачи отсортированы по **приоритету** (🔴 Critical → 🟡 High → 🟢 Medium → 🔵 Low)
2. Для каждой задачи указана **ветка** (следуйте naming convention)
3. Начинайте с **Critical** задач, затем двигайтесь вниз
4. После выполнения задачи - отметьте в TODO списке IDE

---

## 🔴 CRITICAL PRIORITY - Делать немедленно!

### 1. Error Handling Backend
**Branch**: `fix/error-handling-backend`  
**Estimated**: 4-6 hours  
**Files**: `backend/internal/errors/`, `backend/internal/handlers/`

**Что сделать**:
```go
// 1. Создать backend/internal/errors/errors.go
var (
    ErrPlayerNotFound  = errors.New("player not found")
    ErrInvalidSteamID  = errors.New("invalid steam id")
    ErrRateLimited     = errors.New("rate limited")
    ErrAPIUnavailable  = errors.New("external api unavailable")
)

// 2. Обновить handleServiceError в handlers
func (h *PlayerProfileHandler) handleServiceError(c echo.Context, err error) error {
    switch {
    case errors.Is(err, ErrPlayerNotFound):
        return c.JSON(http.StatusNotFound, echo.Map{"error": "Player not found"})
    case errors.Is(err, ErrInvalidSteamID):
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Steam ID"})
    // ... остальные
    }
}
```

### 2. Goroutine Error Channel Fix
**Branch**: `fix/goroutine-error-channel`  
**Estimated**: 2-3 hours  
**File**: `backend/internal/services/player_profile_service.go:163-164`

**Проблема**: 
```go
errs := make(chan error, 5)  // Фиксированный размер - может блокировать!
```

**Решение**:
```go
// Вариант 1: Unbuffered channel с правильной обработкой
errs := make(chan error)

// Вариант 2: Динамический размер
errs := make(chan error, len(goroutines))

// Вариант 3: Не использовать канал, собирать в slice с mutex
```

### 3. Input Validation
**Branch**: `fix/input-validation`  
**Estimated**: 3-4 hours  
**Files**: `backend/internal/validators/steam_id.go` (новый)

**Что сделать**:
```go
// 1. Создать validators пакет
func ValidateSteamID(steamID string) error {
    if steamID == "" {
        return ErrInvalidSteamID
    }
    
    id, err := strconv.ParseInt(steamID, 10, 64)
    if err != nil || id <= 0 {
        return ErrInvalidSteamID
    }
    
    return nil
}

// 2. Использовать в handlers
steamID, err := h.validateSteamIDParam(c)
if err := validators.ValidateSteamID(steamID); err != nil {
    return c.JSON(400, echo.Map{"error": err.Error()})
}
```

### 4. Frontend Error Handling
**Branch**: `fix/frontend-error-handling`  
**Estimated**: 4-5 hours  
**Files**: `frontend/src/shared/lib/ErrorBoundary.tsx` (новый)

**Что сделать**:
```typescript
// 1. Создать Error Boundary
class ErrorBoundary extends React.Component {
    state = { hasError: false, error: null }
    
    static getDerivedStateFromError(error: Error) {
        return { hasError: true, error }
    }
    
    render() {
        if (this.state.hasError) {
            return <ErrorFallback error={this.state.error} />
        }
        return this.props.children
    }
}

// 2. Обернуть приложение
<ErrorBoundary>
    <App />
</ErrorBoundary>

// 3. Улучшить error states в stores
set({ 
    error: error.response?.data?.message || 'Failed to load data',
    loading: false 
})
```

### 5. Remove Console.log
**Branch**: `fix/remove-console-logs`  
**Estimated**: 1-2 hours  
**File**: `frontend/src/shared/lib/logger.ts` (новый)

**Что сделать**:
```typescript
// 1. Создать logger
const isDev = import.meta.env.DEV

export const logger = {
  log: (...args: any[]) => isDev && console.log(...args),
  error: (...args: any[]) => isDev && console.error(...args),
  warn: (...args: any[]) => isDev && console.warn(...args),
}

// 2. Найти все console.log
grep -r "console.log" frontend/src/

// 3. Заменить на logger.log
- console.log('data:', data)
+ logger.log('data:', data)
```

---

## 🟡 HIGH PRIORITY - Сделать в первую очередь

### 6. Rate Limiting
**Branch**: `fix/rate-limiting`  
**Estimated**: 3-4 hours

```go
import "golang.org/x/time/rate"

func RateLimitMiddleware() echo.MiddlewareFunc {
    limiter := rate.NewLimiter(rate.Every(time.Second), 10)
    
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if !limiter.Allow() {
                return c.JSON(429, echo.Map{"error": "Rate limit exceeded"})
            }
            return next(c)
        }
    }
}
```

### 7. Database Connection Pool
**Branch**: `fix/db-connection-pool`  
**Estimated**: 1 hour  
**File**: `backend/cmd/main.go:55`

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
sqlDB.SetConnMaxIdleTime(5 * time.Minute)
```

### 8. Database Indexes
**Branch**: `fix/add-db-indexes`  
**Estimated**: 2 hours

```sql
-- migrations/000016_add_indexes.up.sql
CREATE INDEX IF NOT EXISTS idx_pms_user_id ON player_match_stats(user_id);
CREATE INDEX IF NOT EXISTS idx_pms_match_id ON player_match_stats(match_id);
CREATE INDEX IF NOT EXISTS idx_matches_time ON matches(match_time DESC);
CREATE INDEX IF NOT EXISTS idx_users_nickname ON users USING gin(nickname gin_trgm_ops);

-- migrations/000016_add_indexes.down.sql
DROP INDEX IF EXISTS idx_pms_user_id;
DROP INDEX IF EXISTS idx_pms_match_id;
DROP INDEX IF EXISTS idx_matches_time;
DROP INDEX IF EXISTS idx_users_nickname;
```

---

## 🟢 MEDIUM PRIORITY - После HIGH

### 9-12. Refactoring & Quality

- **refactor/zod-validation** - Zod schemas для API responses (4 hours)
- **refactor/react-query-integration** - Миграция на React Query (8 hours)
- **refactor/skeleton-loaders** - Skeleton UI components (3 hours)
- **fix/security-headers** - CORS, CSRF, XSS protection (2 hours)

### 13-15. Infrastructure

- **chore/prometheus-metrics** - Prometheus метрики (6 hours)
- **chore/improve-logging** - Structured logging (4 hours)
- **chore/ci-cd-pipeline** - GitHub Actions (8 hours)

### 16-18. Testing & Docs

- **test/backend-unit-tests** - Backend тесты 60%+ (12+ hours)
- **test/frontend-unit-tests** - Frontend тесты (8 hours)
- **docs/api-swagger** - Swagger документация (4 hours)

---

## 🎮 FEATURES - После стабилизации

### Phase 3: Hero Builds (6-8 weeks)
- **feat/builds-api** - CRUD для билдов (16 hours)
- **feat/builds-ui** - UI для билдов (20 hours)

### Phase 4: Crosshairs (4-6 weeks)
- **feat/crosshairs-api** - CRUD для прицелов (12 hours)
- **feat/crosshairs-ui** - Visual editor (16 hours)

### Phase 5: Analytics (6-8 weeks)
- **feat/leaderboard** - Leaderboards (12 hours)
- **feat/meta-analysis** - Meta dashboard (16 hours)
- **feat/compare-profiles** - Сравнение игроков (8 hours)

---

## ⚡ PERFORMANCE - Low priority

### 27-29. Optimization
- **perf/code-splitting** - Lazy loading (4 hours)
- **perf/image-optimization** - Image optimization (3 hours)
- **perf/virtual-scrolling** - Virtual scrolling (4 hours)

---

## 🚀 Рекомендуемый план на 4 недели

### Week 1: Critical Fixes 🔴
```
Day 1-2:   fix/error-handling-backend
Day 3:     fix/goroutine-error-channel
Day 4:     fix/input-validation
Day 5-6:   fix/frontend-error-handling
Day 7:     fix/remove-console-logs
```

### Week 2: High Priority 🟡
```
Day 8:     fix/rate-limiting
Day 9:     fix/db-connection-pool + fix/add-db-indexes
Day 10-11: chore/prometheus-metrics
Day 12-13: chore/improve-logging
Day 14:    fix/security-headers
```

### Week 3: Refactoring 🟢
```
Day 15-16: refactor/react-query-integration
Day 17:    refactor/skeleton-loaders
Day 18-19: refactor/zod-validation
Day 20-21: chore/ci-cd-pipeline
```

### Week 4: Testing 🧪
```
Day 22-24: test/backend-unit-tests
Day 25-27: test/frontend-unit-tests
Day 28:    docs/api-swagger
```

После этого - готовы к features! 🎉

---

## 📊 Progress Tracking

**Completed**: 0/30 (0%)

### Критические задачи
- [ ] fix/error-handling-backend
- [ ] fix/goroutine-error-channel
- [ ] fix/input-validation
- [ ] fix/frontend-error-handling
- [ ] fix/remove-console-logs

### Важные задачи
- [ ] fix/rate-limiting
- [ ] fix/db-connection-pool
- [ ] fix/add-db-indexes
- [ ] chore/prometheus-metrics

### Тесты
- [ ] test/backend-unit-tests (60%+)
- [ ] test/frontend-unit-tests

### CI/CD
- [ ] chore/ci-cd-pipeline

**После выполнения всех выше - проект готов к production! 🚀**

---

## 🎯 Quick Start Guide

1. **Выберите задачу** из списка выше
2. **Создайте ветку**: `git checkout -b fix/error-handling-backend`
3. **Сделайте изменения** согласно описанию
4. **Тестируйте**: `go test ./...` или `npm test`
5. **Коммит**: `git commit -m "fix(handlers): add proper error handling"`
6. **PR**: Создайте Pull Request с описанием
7. **Повторите** для следующей задачи! 🔄

---

## 📚 Полезные ссылки

- [GETTING_STARTED.md](GETTING_STARTED.md) - Как начать работу
- [DEVELOPMENT_WORKFLOW.md](DEVELOPMENT_WORKFLOW.md) - Git workflow
- [ROADMAP.md](ROADMAP.md) - Полный roadmap
- [CONTRIBUTING.md](CONTRIBUTING.md) - Как контрибьютить

---

**Удачи в разработке! 🚀**

