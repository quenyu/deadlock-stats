# 🚀 Development Workflow Guide

## 📋 Conventional Commits & Branch Naming

### Типы коммитов и веток

```
fix/      - Исправление багов
feat/     - Новый функционал
refactor/ - Рефакторинг без изменения функциональности
chore/    - Рутинные задачи (deps, CI/CD, configs)
docs/     - Документация
test/     - Тесты
perf/     - Оптимизация производительности
style/    - Форматирование кода (не CSS!)
build/    - Изменения в build системе
ci/       - Изменения в CI/CD
```

### Naming Convention для веток

```bash
# Формат: <type>/<scope>-<short-description>
# Примеры:

fix/error-handling-backend
feat/builds-api
refactor/react-query-integration
chore/prometheus-metrics
docs/api-swagger
test/backend-unit-tests
perf/code-splitting
```

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Примеры:**

```bash
# Fix
git commit -m "fix(handlers): add proper error status codes

- Added typed errors in internal/errors
- Updated handleServiceError to return 404, 400, 429, 500
- Closes #123"

# Feature
git commit -m "feat(builds): implement builds CRUD API

- Added builds endpoints (GET, POST, PUT, DELETE)
- Implemented vote system
- Added filtering and sorting
- Refs #45"

# Refactor
git commit -m "refactor(frontend): migrate from Zustand to React Query

- Setup QueryClient provider
- Created usePlayerProfile hook
- Removed duplicate Zustand stores
- BREAKING CHANGE: API structure changed"

# Chore
git commit -m "chore(ci): add GitHub Actions workflow

- Added lint, test, build stages
- Configured Docker build & push
- Added pre-commit hooks"
```

## 🔄 Git Workflow

### 1. Создание новой ветки

```bash
# Всегда создавайте ветку от актуального main
git checkout main
git pull origin main

# Создайте feature branch
git checkout -b fix/error-handling-backend
```

### 2. Работа над задачей

```bash
# Делайте атомарные коммиты
git add .
git commit -m "fix(handlers): add typed error definitions"

git add .
git commit -m "fix(handlers): implement error status mapping"

git add .
git commit -m "test(handlers): add error handling tests"
```

### 3. Подготовка к PR

```bash
# Обновите свою ветку с main
git fetch origin
git rebase origin/main

# Если есть конфликты, разрешите их
git rebase --continue

# Push в remote
git push origin fix/error-handling-backend

# Если делали rebase после push
git push --force-with-lease origin fix/error-handling-backend
```

### 4. Pull Request

**Шаблон PR:**

```markdown
## Description
Кратко опишите изменения

## Type of change
- [ ] Bug fix (fix/)
- [ ] New feature (feat/)
- [ ] Refactoring (refactor/)
- [ ] Documentation (docs/)
- [ ] Tests (test/)
- [ ] Performance (perf/)

## Changes
- Добавлено X
- Изменено Y
- Удалено Z

## Testing
Как тестировали изменения

## Checklist
- [ ] Код проходит линтер
- [ ] Добавлены тесты
- [ ] Обновлена документация
- [ ] Проверено в dev окружении
- [ ] Нет breaking changes (или они документированы)

## Related Issues
Closes #123
Refs #456
```

## 🎯 Приоритеты задач

### 🔴 CRITICAL (Делать немедленно)
Проблемы безопасности и критические баги

### 🟡 HIGH (Делать в первую очередь)
Важные фичи и улучшения стабильности

### 🟢 MEDIUM (Делать после HIGH)
Улучшения качества кода и DX

### 🔵 LOW (Nice to have)
Оптимизации и дополнительные фичи

## 📊 Workflow пример

### Week 1-2: Critical Fixes 🔴

```bash
# Day 1-2
git checkout -b fix/error-handling-backend
# Работа над задачей
git push && create PR

# Day 3-4
git checkout -b fix/goroutine-error-channel
# Работа над задачей
git push && create PR

# Day 5
git checkout -b fix/input-validation
# Работа над задачей
git push && create PR

# Day 6-7
git checkout -b fix/frontend-error-handling
# Работа над задачей
git push && create PR
```

### Week 3-4: High Priority 🟡

```bash
# Week 3
fix/rate-limiting
fix/db-connection-pool
fix/add-db-indexes

# Week 4
chore/prometheus-metrics
test/backend-unit-tests
```

### Week 5-6: Medium Priority 🟢

```bash
# Week 5
refactor/react-query-integration
refactor/skeleton-loaders
chore/ci-cd-pipeline

# Week 6
test/frontend-unit-tests
docs/api-swagger
```

### Week 7+: Features & Optimization 🔵

```bash
# Features
feat/builds-api
feat/builds-ui
feat/leaderboard

# Performance
perf/code-splitting
perf/image-optimization
```

## 🧪 Testing Strategy

### Backend Tests

```bash
# Запуск всех тестов
go test ./...

# С покрытием
go test -cover ./...

# Детальное покрытие
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Конкретный пакет
go test ./internal/services/...

# С verbose
go test -v ./internal/handlers/...
```

### Frontend Tests

```bash
# Запуск тестов
npm test

# С покрытием
npm test -- --coverage

# Watch mode
npm test -- --watch

# Конкретный файл
npm test -- PlayerProfilePage
```

## 🔍 Code Review Checklist

### Автор PR
- [ ] Код соответствует style guide
- [ ] Добавлены тесты
- [ ] Все тесты проходят
- [ ] Нет console.log
- [ ] Обновлена документация
- [ ] PR описание заполнено
- [ ] Связан с issue

### Reviewer
- [ ] Код понятен и читаем
- [ ] Нет дублирования кода
- [ ] Обработаны все edge cases
- [ ] Тесты покрывают функционал
- [ ] Нет проблем с безопасностью
- [ ] Нет performance issues
- [ ] Breaking changes документированы

## 📚 Полезные команды

### Git

```bash
# Посмотреть статус
git status

# История коммитов
git log --oneline --graph --decorate --all

# Посмотреть изменения
git diff

# Отменить последний коммит (сохранив изменения)
git reset --soft HEAD~1

# Изменить последний коммит
git commit --amend

# Переключиться на main и обновить
git checkout main && git pull

# Удалить локальную ветку
git branch -d fix/old-branch

# Удалить remote ветку
git push origin --delete fix/old-branch

# Stash изменения
git stash
git stash pop
```

### Docker

```bash
# Поднять весь проект
docker-compose up

# С пересборкой
docker-compose up --build

# В фоне
docker-compose up -d

# Остановить
docker-compose down

# Удалить все контейнеры и volumes
docker-compose down -v

# Логи
docker-compose logs -f backend
docker-compose logs -f frontend

# Войти в контейнер
docker-compose exec backend sh
docker-compose exec postgres psql -U postgres -d deadlock_stats
```

### Database

```bash
# Подключиться к Postgres
docker-compose exec postgres psql -U postgres -d deadlock_stats

# Запустить миграции
docker-compose exec backend ./migrate -path ./migrations -database "postgres://..." up

# Откатить миграцию
docker-compose exec backend ./migrate -path ./migrations -database "postgres://..." down 1

# Создать новую миграцию
migrate create -ext sql -dir migrations -seq add_new_feature
```

## 🏗️ Development Environment Setup

### Backend

```bash
# Установка зависимостей
cd backend
go mod download

# Запуск в dev режиме с hot reload (установить air)
go install github.com/cosmtrek/air@latest
air

# Линтинг
golangci-lint run

# Форматирование
go fmt ./...
goimports -w .
```

### Frontend

```bash
# Установка зависимостей
cd frontend
npm install

# Запуск dev server
npm run dev

# Линтинг
npm run lint

# Форматирование (если настроен prettier)
npm run format

# Build для продакшена
npm run build

# Preview production build
npm run preview
```

## 🚨 Troubleshooting

### Проблемы с Docker

```bash
# Пересоздать контейнеры
docker-compose down -v
docker-compose up --build

# Очистить Docker
docker system prune -a --volumes
```

### Проблемы с миграциями

```bash
# Проверить версию миграций
SELECT * FROM schema_migrations;

# Откатить все
migrate -path ./migrations -database "postgres://..." down

# Применить заново
migrate -path ./migrations -database "postgres://..." up
```

### Проблемы с Git

```bash
# Откатить все изменения
git reset --hard HEAD

# Синхронизироваться с remote
git fetch origin
git reset --hard origin/main
```

## 📖 Дополнительные ресурсы

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Go Best Practices](https://github.com/golang-standards/project-layout)
- [React Best Practices](https://react.dev/learn)
- [Feature-Sliced Design](https://feature-sliced.design/)
- [Testing Best Practices](https://testingjavascript.com/)

---

**Happy Coding! 🚀**

