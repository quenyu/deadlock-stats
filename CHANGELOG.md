# 📝 Changelog - Deadlock Stats

Все значимые изменения проекта документируются в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/ru/1.0.0/),
проект придерживается [Semantic Versioning](https://semver.org/lang/ru/).

---

## [Unreleased]

### 🎉 Added
- **Crosshairs System** - Полноценная система для кастомных прицелов
  - CRUD API для crosshairs (`backend/internal/handlers/crosshair_handler.go`)
  - Система лайков crosshairs
  - Галерея crosshairs на frontend
  - Builder для создания crosshairs
- **Error Handling** - Централизованная обработка ошибок
  - 20+ типизированных ошибок (`backend/internal/errors/errors.go`)
  - Централизованный ErrorHandler (`backend/internal/handlers/error_handler.go`)
  - Правильные HTTP статусы (404, 400, 401, 403, 429, 500, 503)
- **Player Search** - Расширенный поиск игроков
  - Advanced search с фильтрами
  - Search API (`backend/internal/handlers/player_search_handler.go`)
  - Search service (`backend/internal/services/player_search_service.go`)

### 📚 Documentation
- Полная документация проекта (11 файлов)
- README с features, tech stack, quick start
- ROADMAP на 9+ месяцев
- TODO_SUMMARY с приоритизированными задачами
- CONTRIBUTING guide
- GitHub Issue/PR templates

### 🗄️ Database
- Миграции для crosshairs (000016, 000017)
- Расширенная схема для likes
- Всего 17 миграций

### 🐛 Fixed
- Error handling теперь возвращает правильные статусы вместо generic 500
- Исправлены handlers для использования централизованного ErrorHandler

---

## [0.1.0-alpha] - 2025-01-XX

### Initial Release

#### Backend Features
- **Authentication** - Steam OpenID login + JWT
- **Player Profiles** - Extended profiles с детальной статистикой
- **Match History** - История матчей игрока
- **Hero Statistics** - Статистика по героям
- **MMR Tracking** - Отслеживание MMR и истории
- **Caching** - Redis кэширование с TTL
- **Database** - PostgreSQL с миграциями

#### Frontend Features
- **Modern UI** - React 19 + TypeScript + Tailwind CSS
- **Dark Theme** - Красивая тёмная тема
- **Responsive Design** - Работает на всех устройствах
- **Charts** - Recharts для визуализации статистики
- **State Management** - Zustand для управления состоянием

#### Infrastructure
- Docker + Docker Compose setup
- Nginx reverse proxy configuration
- Development environment ready

---

## Текущее состояние проекта

### ✅ Реализовано (13.3% от общего плана)
1. Error Handling Backend ✅
2. Crosshairs API ✅
3. Crosshairs UI ✅
4. Documentation ✅

### 🔨 В процессе
1. Goroutine Error Channel Fix (3 места)

### 📋 Следующие задачи (по приоритету)
1. fix/input-validation (CRITICAL)
2. fix/frontend-error-handling (CRITICAL)
3. fix/remove-console-logs (HIGH)
4. fix/rate-limiting (HIGH)
5. fix/db-connection-pool (HIGH)

---

## Статистика проекта

**Backend:**
- Files: 56+
- Lines of code: ~6000+
- Migrations: 17
- Error types: 20+
- Test coverage: ~45%

**Frontend:**
- Files: 120+
- Lines of code: ~10000+
- Components: 60+
- Pages: 6
- Test coverage: ~20%

**Documentation:**
- Files: 11
- Lines: ~3500+
- Coverage: 100% ✅

---

## Что дальше?

### Week 1-2: Critical Fixes
- [ ] Goroutine deadlock fix
- [ ] Input validation
- [ ] Frontend error handling
- [ ] Remove console.log

### Week 3-4: Infrastructure
- [ ] Rate limiting
- [ ] DB optimization
- [ ] Prometheus metrics
- [ ] CI/CD pipeline

### Week 5+: Features & Quality
- [ ] React Query migration
- [ ] Unit tests (60%+)
- [ ] Hero Builds system
- [ ] Performance optimization

---

**Last Updated**: 2025-10-07  
**Maintainer**: [@wqeqadas](https://github.com/wqeqadas)


