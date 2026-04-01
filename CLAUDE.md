# Deadlock Stats Development Guidelines

**Project**: Deadlock Stats  
**Last updated**: 2025-10-08

## Active Technologies

- Go 1.23 + Echo (backend)
- TypeScript 5.8 + React 19 + Vite 7 (frontend)
- PostgreSQL 16 + Redis 7 (storage)
- Docker + Docker Compose (infrastructure)

## Architecture Overview

### Backend (Go)
```
backend/
├── cmd/main.go              # Entry point
├── internal/
│   ├── handlers/            # HTTP handlers (no business logic)
│   ├── services/            # Business logic layer
│   ├── repositories/        # Data access layer
│   ├── domain/              # Domain models
│   ├── dto/                 # Data Transfer Objects
│   ├── clients/             # External API clients
│   ├── middleware/          # Custom middlewares
│   └── config/              # Configuration
└── migrations/              # Database migrations
```

**Clean Architecture**: Handlers → Services → Repositories → Domain

### Frontend (TypeScript + React)
```
frontend/src/
├── app/                     # App initialization
│   ├── providers/           # Global providers
│   └── styles/              # Global styles
├── pages/                   # Route components
├── widgets/                 # Complex UI blocks
├── features/                # User scenarios
├── entities/                # Business entities
│   └── [entity]/
│       ├── api/             # API calls
│       ├── model/           # State management
│       ├── types/           # TypeScript types
│       └── utils/           # Helpers
└── shared/                  # Reusable components
    ├── ui/                  # UI components
    ├── lib/                 # Utilities
    ├── api/                 # Axios instance
    └── constants/           # Constants
```

**Feature-Sliced Design**: app → pages → widgets → features → entities → shared

## Core Principles (from Constitution v1.0.0)

### 1. Clean Architecture (Backend) 🏗️
- Handlers MUST NOT contain business logic
- Services MUST NOT depend on HTTP-specific types
- Repositories MUST return domain models
- Dependencies flow inward: Handlers → Services → Repositories → Domain

### 2. Feature-Sliced Design (Frontend) 📐
- Layers MUST NOT import from layers above them
- Features MUST be self-contained and independently deletable
- Shared layer MUST NOT depend on business logic
- Pages MUST only compose, not implement logic

### 3. Test-Driven Quality (NON-NEGOTIABLE) ✅
- Backend: 60%+ code coverage REQUIRED
- Frontend: Tests for critical user flows REQUIRED
- Integration tests REQUIRED for API contracts
- Test-First: Tests → Verify failure → Implement → Verify pass

### 4. Performance & Scalability ⚡
- API Response: <100ms (p95)
- Page Load: <2s
- All database tables MUST have indexes
- Expensive operations MUST be cached (Redis)
- Frontend MUST use code splitting

### 5. Security First 🔒
- All inputs MUST be validated
- Rate limiting: 100 req/min per IP
- No SQL string concatenation (use GORM parameterized queries)
- Error messages MUST be generic (no stack traces)
- Dependencies MUST be updated monthly

## Commands

### Backend (Go)
```bash
# Development
cd backend
go mod download
air                          # Hot reload

# Testing
go test ./...
go test -cover ./...
go test -coverprofile=coverage.out ./...

# Linting
golangci-lint run
go fmt ./...
goimports -w .
```

### Frontend (TypeScript + React)
```bash
# Development
cd frontend
npm install
npm run dev

# Testing
npm test
npm test -- --coverage
npm test -- --watch

# Linting
npm run lint
npm run lint -- --fix
```

### Docker
```bash
# Start all services
docker-compose up --build

# Logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Database access
docker-compose exec postgres psql -U postgres -d deadlock_stats
```

## Code Style

### Go
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Comment all exported functions
- Handle all errors explicitly
- No naked returns in functions >5 lines

### TypeScript/React
- Strict TypeScript mode (NO `any` types)
- Functional components + hooks only
- Proper types for everything
- Use Zod for runtime validation (planned)
- ESLint + Prettier compliance required

## Development Workflow

### Branch Naming
```
fix/<scope>-<description>       # Bug fixes
feat/<scope>-<description>      # New features
refactor/<scope>-<description>  # Refactoring
chore/<scope>-<description>     # Maintenance
docs/<scope>-<description>      # Documentation
test/<scope>-<description>      # Tests
perf/<scope>-<description>      # Performance
```

### Commit Messages (Conventional Commits)
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Examples**:
- `fix(handlers): add proper error status codes`
- `feat(builds): implement builds CRUD API`
- `refactor(frontend): migrate to React Query`

### PR Requirements
- ✅ Linter passes
- ✅ Tests pass (60%+ backend coverage maintained)
- ✅ Tests included for new functionality
- ✅ No console.log in production code
- ✅ Breaking changes documented

## Recent Changes

- Phase 1: Stabilization - Error handling, security, testing infrastructure
- Constitution v1.0.0: Established 5 core principles
- Current focus: Critical security fixes (rate limiting, input validation)

## References

- **Constitution**: `.specify/memory/constitution.md` - Core principles (v1.0.0)
- **Roadmap**: `ROADMAP.md` - Feature timeline and priorities
- **Workflow**: `DEVELOPMENT_WORKFLOW.md` - Git workflow, commands
- **Architecture**: `PROJECT_OVERVIEW.md` - Detailed architecture

## Current Priorities (Phase 1)

🔴 **Week 1-2: Critical Fixes**
- [ ] Error handling & typed errors
- [ ] Goroutine error channel fix
- [ ] Input validation
- [ ] Frontend error handling

🟡 **Week 3-4: Security & Infrastructure**
- [ ] Rate limiting
- [ ] Database optimization (indexes)
- [ ] Security headers
- [ ] Prometheus metrics

🟢 **Week 5-6: Testing & CI/CD**
- [ ] Backend unit tests (60%+ coverage)
- [ ] CI/CD pipeline
- [ ] API documentation

---

**See also**: `DEVELOPMENT_WORKFLOW.md` for detailed commands and troubleshooting

