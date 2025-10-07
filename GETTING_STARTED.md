# 🚀 Getting Started - Первые шаги

> **Быстрый старт для новых контрибьюторов**

## 📋 Начало работы с проектом

### 1️⃣ Setup Development Environment

#### Требования
- **Go** 1.23+
- **Node.js** 18+
- **Docker** & Docker Compose
- **Git**

#### Клонирование репозитория
```bash
git clone https://github.com/quenyu/deadlock-stats.git
cd deadlock-stats
```

#### Запуск проекта
```bash
# Запустить все сервисы (backend, frontend, postgres, redis)
docker-compose up

# Или с пересборкой
docker-compose up --build

# В фоне
docker-compose up -d
```

**Доступ к приложению**:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Health Check: http://localhost:8080/health

#### Локальная разработка без Docker

**Backend**:
```bash
cd backend

# Установить зависимости
go mod download

# Настроить конфиг
cp internal/config/config.example.yaml internal/config/config.yaml
# Отредактировать config.yaml с вашими настройками

# Запустить БД локально
docker-compose up postgres redis

# Запустить backend
go run cmd/main.go

# Или с hot reload (установить air)
go install github.com/cosmtrek/air@latest
air
```

**Frontend**:
```bash
cd frontend

# Установить зависимости
npm install

# Настроить env
echo "VITE_API_URL=http://localhost:8080/api/v1" > .env

# Запустить dev server
npm run dev
```

---

### 2️⃣ Выбор первой задачи

#### Для новичков 🟢

**Начните с простых задач**:

1. **fix/remove-console-logs** ⭐ EASY
   - Файл: `frontend/src/entities/deadlock/model/store.ts`
   - Что делать: Создать logger helper и заменить console.log
   - Время: 30-60 минут

2. **docs/update-readme** ⭐ EASY
   - Файл: `README.md`
   - Что делать: Добавить описание проекта, screenshots, setup инструкции
   - Время: 1-2 часа

3. **refactor/skeleton-loaders** ⭐⭐ MEDIUM
   - Файлы: `frontend/src/shared/ui/skeleton.tsx` + использование в pages
   - Что делать: Создать skeleton компоненты, заменить "Loading..." текст
   - Время: 2-3 часа

#### Для опытных разработчиков 🔴

1. **fix/error-handling-backend** ⭐⭐⭐ HARD
   - Файлы: `backend/internal/errors/`, `backend/internal/handlers/`
   - Что делать: Типизированные ошибки, правильные HTTP статусы
   - Время: 4-6 часов

2. **feat/builds-api** ⭐⭐⭐ HARD
   - Файлы: новый модуль `backend/internal/*/builds.go`
   - Что делать: CRUD endpoints для билдов
   - Время: 8-12 часов

---

### 3️⃣ Создание ветки и начало работы

#### Пример: fix/remove-console-logs

```bash
# 1. Обновить main
git checkout main
git pull origin main

# 2. Создать ветку
git checkout -b fix/remove-console-logs

# 3. Посмотреть TODO
# Открыть проект в редакторе, найти задачу в TODO списке

# 4. Начать работу
cd frontend/src/shared/lib
```

**Создать файл `logger.ts`**:
```typescript
// shared/lib/logger.ts
const isDev = import.meta.env.DEV

export const logger = {
  log: (...args: any[]) => {
    if (isDev) {
      console.log(...args)
    }
  },
  
  error: (...args: any[]) => {
    if (isDev) {
      console.error(...args)
    }
    // В production - отправить в Sentry (когда будет настроен)
  },
  
  warn: (...args: any[]) => {
    if (isDev) {
      console.warn(...args)
    }
  },
  
  info: (...args: any[]) => {
    if (isDev) {
      console.info(...args)
    }
  }
}
```

**Найти и заменить все console.log**:
```bash
# Найти все использования console.log
grep -r "console.log" src/

# Заменить в файлах
# В store.ts:
- console.log('Full API Response for Extended Player Profile:', data)
+ logger.log('Full API Response for Extended Player Profile:', data)

# Не забыть добавить import:
+ import { logger } from '@/shared/lib/logger'
```

**Коммит изменений**:
```bash
# Проверить изменения
git status
git diff

# Добавить файлы
git add .

# Коммит с правильным форматом
git commit -m "fix(logger): replace console.log with conditional logger

- Created shared/lib/logger.ts with dev-only logging
- Replaced all console.log calls with logger.log
- Prevents debug logs in production build
- Closes #XXX"

# Push в remote
git push origin fix/remove-console-logs
```

**Создать Pull Request**:
1. Перейти на GitHub
2. Создать PR из вашей ветки в main
3. Заполнить шаблон PR
4. Запросить review

---

### 4️⃣ Пример более сложной задачи

#### Пример: refactor/skeleton-loaders

**Шаг 1: Создать skeleton компонент**

```bash
git checkout -b refactor/skeleton-loaders
cd frontend/src/shared/ui
```

**Создать `skeleton.tsx`**:
```typescript
// shared/ui/skeleton.tsx
import { cn } from '@/shared/lib/utils'

interface SkeletonProps {
  className?: string
}

export function Skeleton({ className }: SkeletonProps) {
  return (
    <div
      className={cn(
        'animate-pulse rounded-md bg-muted',
        className
      )}
    />
  )
}

// Специфичные skeleton компоненты
export function PlayerProfileSkeleton() {
  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8">
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-8">
          {/* Header skeleton */}
          <div className="flex items-center gap-4">
            <Skeleton className="h-24 w-24 rounded-full" />
            <div className="space-y-2">
              <Skeleton className="h-8 w-48" />
              <Skeleton className="h-4 w-32" />
            </div>
          </div>
          
          {/* Stats skeleton */}
          <div className="grid grid-cols-3 gap-4">
            <Skeleton className="h-32" />
            <Skeleton className="h-32" />
            <Skeleton className="h-32" />
          </div>
          
          {/* Chart skeleton */}
          <Skeleton className="h-64" />
        </div>
        
        <div className="space-y-8">
          <Skeleton className="h-48" />
          <Skeleton className="h-96" />
        </div>
      </div>
    </div>
  )
}

export function MatchCardSkeleton() {
  return (
    <div className="flex items-center gap-4 p-4 border rounded-lg">
      <Skeleton className="h-16 w-16 rounded-full" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-32" />
        <Skeleton className="h-4 w-48" />
      </div>
      <Skeleton className="h-8 w-16" />
    </div>
  )
}
```

**Шаг 2: Использовать skeleton в pages**

```typescript
// pages/player-profile/PlayerProfilePage.tsx
import { PlayerProfileSkeleton } from '@/shared/ui/skeleton'

export const PlayerProfilePage = () => {
  const { steamId } = useParams<{ steamId: string }>()
  const { profile, loading, error, fetchProfile } = useExtendedProfileStore()

  // ...

  if (loading) {
-   return <div className="text-center py-10">Loading profile...</div>
+   return <PlayerProfileSkeleton />
  }

  // ...
}
```

**Шаг 3: Тестирование**

```bash
# Запустить dev server
cd frontend
npm run dev

# Проверить в браузере:
# 1. Перейти на любой профиль игрока
# 2. Должен показаться skeleton вместо текста
# 3. Проверить что плавно переходит от skeleton к контенту
```

**Шаг 4: Коммит**

```bash
git add .
git commit -m "refactor(ui): add skeleton loaders for better UX

- Created reusable Skeleton component
- Added PlayerProfileSkeleton component
- Added MatchCardSkeleton component
- Replaced loading text with skeleton UI
- Improves perceived performance"

git push origin refactor/skeleton-loaders
```

---

### 5️⃣ Code Review Process

#### Checklist перед созданием PR

```markdown
## Pre-PR Checklist ✅

- [ ] Код работает локально
- [ ] Нет console.log (кроме logger)
- [ ] Код отформатирован (eslint, prettier)
- [ ] Нет type errors
- [ ] Добавлены комментарии к сложной логике
- [ ] Обновлена документация (если нужно)
- [ ] Тесты добавлены (если нужно)
- [ ] Коммиты соответствуют Conventional Commits
- [ ] PR описание заполнено
```

#### После создания PR

1. **Automatic checks** - дождитесь прохождения CI
2. **Request review** - запросите review у maintainers
3. **Address feedback** - исправьте замечания
4. **Merge** - после approval ветка будет смержена

---

### 6️⃣ Полезные команды

#### Git

```bash
# Посмотреть статус
git status

# Посмотреть изменения
git diff

# Добавить все файлы
git add .

# Добавить конкретный файл
git add src/shared/lib/logger.ts

# Коммит
git commit -m "fix(logger): add conditional logger"

# Изменить последний коммит
git commit --amend

# Откатить изменения (не закоммиченные)
git checkout -- .

# Переключиться на main
git checkout main

# Обновить main
git pull origin main

# Удалить ветку
git branch -d fix/old-branch
```

#### Docker

```bash
# Запустить всё
docker-compose up

# Остановить
docker-compose down

# Пересоздать
docker-compose up --build

# Логи
docker-compose logs -f backend
docker-compose logs -f frontend

# Войти в контейнер
docker-compose exec backend sh
```

#### Backend (Go)

```bash
# Запустить
go run cmd/main.go

# С hot reload
air

# Тесты
go test ./...

# Тесты с покрытием
go test -cover ./...

# Форматирование
go fmt ./...

# Линтинг (нужно установить)
golangci-lint run
```

#### Frontend (Node)

```bash
# Запустить dev
npm run dev

# Build
npm run build

# Preview production build
npm run preview

# Тесты
npm test

# Линтинг
npm run lint

# Fix lint errors
npm run lint -- --fix
```

#### Database

```bash
# Подключиться к Postgres
docker-compose exec postgres psql -U postgres -d deadlock_stats

# Запросы
\dt                    # список таблиц
\d users              # структура таблицы
SELECT * FROM users LIMIT 10;

# Выйти
\q
```

---

### 7️⃣ Troubleshooting

#### Проблема: Docker не запускается

```bash
# Удалить все контейнеры
docker-compose down -v

# Пересоздать
docker-compose up --build
```

#### Проблема: Frontend не подключается к Backend

```bash
# Проверить что backend запущен
curl http://localhost:8080/health

# Проверить .env файл
cat frontend/.env
# Должно быть: VITE_API_URL=http://localhost:8080/api/v1
```

#### Проблема: Ошибки миграций БД

```bash
# Посмотреть текущую версию
docker-compose exec postgres psql -U postgres -d deadlock_stats -c "SELECT * FROM schema_migrations;"

# Откатить все миграции
# (ВНИМАНИЕ: удалит все данные!)
docker-compose down -v
docker-compose up
```

#### Проблема: Git конфликты

```bash
# Обновить вашу ветку с main
git fetch origin
git rebase origin/main

# Если есть конфликты:
# 1. Открыть файлы с конфликтами
# 2. Разрешить конфликты (удалить <<<<<<, ======, >>>>>>)
# 3. Добавить файлы
git add .
git rebase --continue

# Если запутались - отменить rebase
git rebase --abort
```

---

### 8️⃣ Best Practices

#### Code Style

**Go**:
- Используйте `gofmt` для форматирования
- Комментарии для exported функций
- Error handling везде
- Используйте context для cancellation

**TypeScript/React**:
- Используйте TypeScript строго (no `any`)
- Functional components + hooks
- Props destructuring
- Правильные типы для всего

#### Git Commits

```bash
# ✅ Хорошо
git commit -m "fix(auth): handle expired tokens correctly"
git commit -m "feat(builds): add vote system"
git commit -m "refactor(ui): extract common components"

# ❌ Плохо
git commit -m "fix"
git commit -m "update code"
git commit -m "changes"
```

#### Testing

```bash
# Backend - всегда запускать тесты перед PR
go test ./...

# Frontend - тестировать критичные части
npm test

# Проверить что ничего не сломалось
# 1. Запустить приложение локально
# 2. Кликнуть по основным страницам
# 3. Проверить что всё работает
```

---

### 9️⃣ Где получить помощь?

#### Документация
- `README.md` - основная информация
- `DEVELOPMENT_WORKFLOW.md` - процесс разработки
- `ROADMAP.md` - план развития
- `GETTING_STARTED.md` - этот файл

#### Issues
- [GitHub Issues](https://github.com/quenyu/deadlock-stats/issues) - создать issue с вопросом

#### Community
- Discord (если есть) - задать вопрос в реальном времени
- GitHub Discussions - обсуждения

---

### 🎯 Следующие шаги

1. ✅ Setup development environment
2. ✅ Выбрать задачу из TODO
3. ✅ Создать ветку
4. ✅ Сделать изменения
5. ✅ Протестировать локально
6. ✅ Создать PR
7. ✅ Получить review
8. ✅ Merge и celebrate! 🎉

---

### 🌟 Tips для новых контрибьюторов

1. **Начните с малого** - не беритесь сразу за большие задачи
2. **Читайте код** - лучший способ понять проект
3. **Задавайте вопросы** - не стесняйтесь спрашивать
4. **Тестируйте** - всегда проверяйте свои изменения
5. **Будьте терпеливы** - review может занять время
6. **Учитесь** - каждый PR - это опыт

---

**Happy Coding! 🚀**

Если что-то непонятно - создайте issue или напишите в Discord!

