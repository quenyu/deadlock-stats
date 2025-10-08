# Speckit - Development Workflow System

Эта директория содержит систему Speckit для управления разработкой фич в проекте Deadlock Stats.

## 📁 Структура

```
.specify/
├── memory/
│   └── constitution.md          # Конституция проекта (v1.0.0)
├── scripts/
│   └── powershell/
│       ├── create-new-feature.ps1       # Создание новой фичи
│       ├── setup-plan.ps1               # Подготовка плана
│       ├── update-agent-context.ps1     # Обновление контекста AI
│       └── common.ps1                   # Общие функции
├── templates/
│   ├── spec-template.md         # Шаблон спецификации
│   ├── plan-template.md         # Шаблон implementation plan
│   ├── tasks-template.md        # Шаблон списка задач
│   ├── checklist-template.md    # Шаблон чеклиста
│   └── agent-file-template.md   # Шаблон для AI контекста
└── README.md                    # Этот файл
```

## 🎯 Доступные команды

### 1. `/speckit.constitution` ✅
**Назначение**: Создание/обновление конституции проекта

**Статус**: ✅ Выполнено (v1.0.0, 2025-10-08)

**Результат**: `.specify/memory/constitution.md`

**Принципы**:
- I. Clean Architecture (Backend)
- II. Feature-Sliced Design (Frontend)
- III. Test-Driven Quality (60%+ coverage)
- IV. Performance & Scalability (<100ms API)
- V. Security First (validation, rate limiting)

---

### 2. `/speckit.specify <описание фичи>`
**Назначение**: Создание спецификации фичи

**Пример**:
```
/speckit.specify Система rate limiting для защиты API от DDoS атак
```

**Что создаёт**:
- Новую ветку (например: `001-rate-limiting`)
- Директорию `specs/001-rate-limiting/`
- Файл `specs/001-rate-limiting/spec.md`

**Содержит**:
- User scenarios & testing
- Functional requirements
- Success criteria
- Key entities

---

### 3. `/speckit.plan`
**Назначение**: Создание implementation plan на основе спецификации

**Предусловия**: Должен существовать `spec.md`

**Что создаёт**:
- `specs/[feature]/plan.md` - технический план
- `specs/[feature]/research.md` - результаты исследования
- `specs/[feature]/data-model.md` - модель данных
- `specs/[feature]/quickstart.md` - quick start guide
- `specs/[feature]/contracts/` - API контракты

**Фазы планирования**:
- Phase 0: Research (изучение существующего кода)
- Phase 1: Design (проектирование решения)
- Phase 2: Implementation (готовность к разработке)

---

### 4. `/speckit.tasks`
**Назначение**: Создание списка задач из implementation plan

**Предусловия**: Должны существовать `spec.md` и `plan.md`

**Что создаёт**:
- `specs/[feature]/tasks.md` - список конкретных задач

**Организация задач**:
- Phase 1: Setup (инициализация)
- Phase 2: Foundational (фундамент)
- Phase 3+: User Stories (по приоритетам P1, P2, P3)
- Final Phase: Polish (доработка)

**Формат задачи**: `[ID] [P?] [Story] Description`
- `[P]` = можно выполнять параллельно
- `[Story]` = принадлежность к user story (US1, US2, ...)

---

### 5. `/speckit.checklist`
**Назначение**: Создание чеклистов для проверки качества

**Типы чеклистов**:
- Requirements checklist - проверка полноты требований
- Implementation checklist - проверка реализации
- Security checklist - проверка безопасности
- Performance checklist - проверка производительности

**Что создаёт**:
- `specs/[feature]/checklists/[type].md`

---

### 6. `/speckit.clarify`
**Назначение**: Уточнение неясных требований в спецификации

**Использование**:
- Автоматически находит `[NEEDS CLARIFICATION]` маркеры
- Предлагает варианты ответов
- Обновляет спецификацию после получения ответов

---

## 🔄 Типичный workflow

### Вариант 1: Полный цикл (новая фича)

```bash
# 1. Создать спецификацию
/speckit.specify Система голосования за билды с рейтингом

# 2. Уточнить требования (если нужно)
/speckit.clarify

# 3. Создать implementation plan
/speckit.plan

# 4. Создать список задач
/speckit.tasks

# 5. Начать разработку по задачам из tasks.md
```

### Вариант 2: Быстрый старт (известные требования)

```bash
# Если требования уже ясны, можно сразу создать план
/speckit.specify Добавить кэширование Redis для API профилей
/speckit.plan
/speckit.tasks

# И начать разработку
```

### Вариант 3: Обновление конституции

```bash
# При добавлении новых принципов или изменении существующих
/speckit.constitution

# Система автоматически:
# - Увеличит версию (MAJOR/MINOR/PATCH)
# - Проверит зависимые шаблоны
# - Создаст Sync Impact Report
```

---

## 📊 Структура фичи (создается автоматически)

После выполнения полного цикла получаем:

```
specs/001-feature-name/
├── spec.md              # Спецификация (ЧТО нужно)
├── plan.md              # Технический план (КАК делать)
├── research.md          # Результаты исследования
├── data-model.md        # Модель данных
├── quickstart.md        # Quick start guide
├── tasks.md             # Список задач
├── contracts/           # API контракты
│   ├── endpoint-1.md
│   └── endpoint-2.md
└── checklists/          # Чеклисты
    ├── requirements.md
    ├── implementation.md
    └── security.md
```

---

## 🎓 Лучшие практики

### При создании спецификации
- ✅ Описывайте **ЧТО** нужно пользователям, а не **КАК** реализовать
- ✅ Избегайте технических деталей (фреймворки, БД, API)
- ✅ Пишите для бизнес-стейкхолдеров, не для разработчиков
- ✅ Каждое требование должно быть тестируемым
- ❌ Максимум 3 `[NEEDS CLARIFICATION]` маркера

### При создании плана
- ✅ Следуйте Constitution Check
- ✅ Обосновывайте сложность в "Complexity Tracking"
- ✅ Указывайте конкретные пути к файлам
- ❌ Не нарушайте принципы конституции без обоснования

### При создании задач
- ✅ Группируйте по User Stories
- ✅ Помечайте параллельные задачи `[P]`
- ✅ Указывайте конкретные пути к файлам
- ✅ Начинайте с тестов (Test-First)
- ❌ Не смешивайте разные User Stories в одной задаче

---

## 🔍 Constitution Check

Каждый план проверяется на соответствие принципам конституции:

**Backend (Clean Architecture)**:
- [ ] Handlers только для HTTP
- [ ] Services содержат бизнес-логику
- [ ] Repositories только для БД
- [ ] Domain модели чистые от инфраструктуры

**Frontend (Feature-Sliced Design)**:
- [ ] Нет импортов вверх по слоям
- [ ] Фичи независимы
- [ ] Shared не зависит от бизнес-логики
- [ ] Pages только композиция

**Test-Driven Quality**:
- [ ] Тесты написаны до реализации
- [ ] Backend coverage 60%+
- [ ] Есть integration тесты для API

**Performance**:
- [ ] API <100ms (p95)
- [ ] Есть индексы БД
- [ ] Кэширование для дорогих операций

**Security**:
- [ ] Валидация всех входных данных
- [ ] Rate limiting на публичных endpoints
- [ ] Параметризованные запросы (GORM)
- [ ] Обобщенные ошибки для клиентов

---

## 🛠️ Скрипты PowerShell

### create-new-feature.ps1
Создает новую фичу с веткой и директорией.

```powershell
.\.specify\scripts\powershell\create-new-feature.ps1 -Json "Описание фичи"
```

**Возвращает JSON**:
```json
{
  "BRANCH_NAME": "001-feature-name",
  "SPEC_FILE": "D:/path/specs/001-feature-name/spec.md",
  "FEATURE_DIR": "D:/path/specs/001-feature-name"
}
```

### update-agent-context.ps1
Обновляет контекстные файлы для AI агентов.

```powershell
# Обновить для Cursor
.\.specify\scripts\powershell\update-agent-context.ps1 -AgentType cursor

# Обновить для всех существующих
.\.specify\scripts\powershell\update-agent-context.ps1
```

**Поддерживаемые агенты**:
- `claude` → `CLAUDE.md`
- `cursor` → `.cursor/rules/specify-rules.mdc`
- `copilot` → `.github/copilot-instructions.md`
- `gemini`, `qwen`, `windsurf`, `kilocode`, etc.

---

## 📖 Дополнительная документация

- **Constitution**: `.specify/memory/constitution.md` - Принципы проекта
- **Agent Context**: `CLAUDE.md` / `.cursor/rules/specify-rules.mdc` - Контекст для AI
- **Project Docs**: `PROJECT_OVERVIEW.md`, `ROADMAP.md`, `DEVELOPMENT_WORKFLOW.md`

---

## ❓ FAQ

**Q: Когда использовать `/speckit.specify`?**  
A: Когда начинаете работу над новой фичей. Это создаст ветку и структуру.

**Q: Обязательно ли следовать всему workflow?**  
A: Для простых изменений (fix багов) можно пропустить. Для фич - настоятельно рекомендуется.

**Q: Что делать, если Constitution Check не проходит?**  
A: Либо упростите решение, либо обоснуйте сложность в секции "Complexity Tracking".

**Q: Можно ли изменить конституцию?**  
A: Да, через `/speckit.constitution`. Система увеличит версию и проверит зависимости.

**Q: Где хранятся созданные фичи?**  
A: В директории `specs/[###-feature-name]/`

---

**Версия Speckit**: 1.0  
**Последнее обновление**: 2025-10-08  
**Версия Constitution**: 1.0.0

