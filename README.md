# 🎮 Deadlock Stats

> **Comprehensive statistics and analytics platform for Valve's Deadlock**

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.8-3178C6?logo=typescript)](https://www.typescriptlang.org/)

## 📋 Table of Contents

## ✨ Features

### 🎯 Current Features

- **Player Profiles** - Детальная статистика игроков
  - MMR tracking & history
  - Hero statistics with win rates
  - Match history with performance metrics
  - Personal records (max kills, best KDA, etc.)
  - Peak rank tracking
  - Performance dynamics charts

- **Search System** - Поиск игроков
  - Search by Steam ID
  - Search by nickname
  - Fuzzy search integration

- **Authentication** - Steam OpenID
  - Secure JWT tokens
  - Protected routes

- **Modern UI/UX**
  - Dark theme
  - Responsive design
  - Real-time data visualization

### 🚧 Coming Soon

- **Hero Builds** - Build creator with voting system
- **Crosshairs** - Visual crosshair editor & gallery
- **Leaderboards** - Global & hero-specific rankings
- **Meta Analysis** - Win rates, tier lists, trends
- **Social Features** - Friends, teams, profile comparison
- **Premium Tier** - Advanced analytics & customization

_See [ROADMAP.md](ROADMAP.md) for detailed timeline_

## 🏗️ Tech Stack

### Backend
- **Go 1.23** - Modern, performant language
- **Echo** - High-performance web framework
- **PostgreSQL 16** - Reliable relational database
- **Redis 7** - Fast caching layer
- **GORM** - Developer-friendly ORM
- **Zap** - Structured logging

### Frontend
- **React 19** - Latest React features
- **TypeScript 5.8** - Type-safe development
- **Vite 7** - Lightning-fast build tool
- **Tailwind CSS 4** - Utility-first styling
- **Radix UI** - Accessible components
- **Zustand** - Lightweight state management
- **React Query** _(planned)_ - Data fetching & caching

### Infrastructure
- **Docker** - Containerization
- **Nginx** - Reverse proxy
- **Prometheus + Grafana** _(planned)_ - Monitoring
- **GitHub Actions** _(planned)_ - CI/CD

## 🚀 Quick Start

### Prerequisites

- [Docker](https://www.docker.com/) & Docker Compose
- [Go 1.23+](https://go.dev/) _(for local development)_
- [Node.js 18+](https://nodejs.org/) _(for local development)_

### Installation

```bash
# Clone repository
git clone https://github.com/quenyu/deadlock-stats.git
cd deadlock-stats

# Start all services (backend, frontend, postgres, redis)
docker-compose up

# Or with rebuild
docker-compose up --build
```

### Access

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### Configuration

```bash
# Backend configuration
cp backend/internal/config/config.example.yaml backend/internal/config/config.yaml
# Edit config.yaml with your settings

# Frontend environment
echo "VITE_API_URL=http://localhost:8080/api/v1" > frontend/.env
```

## 📚 Documentation

- **[PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md)** - Project vision, architecture, status
- **[GETTING_STARTED.md](GETTING_STARTED.md)** - Step-by-step guide for contributors
- **[DEVELOPMENT_WORKFLOW.md](DEVELOPMENT_WORKFLOW.md)** - Git workflow, conventions, best practices
- **[ROADMAP.md](ROADMAP.md)** - Development roadmap & timeline
- **[API Documentation](docs/api.md)** _(coming soon)_ - API endpoints reference

## 📁 Project Structure

```
deadlock-stats/
├── backend/                # Go backend
│   ├── cmd/               # Entry point
│   ├── internal/          # Private application code
│   │   ├── handlers/      # HTTP handlers
│   │   ├── services/      # Business logic
│   │   ├── repositories/  # Data access
│   │   ├── domain/        # Domain models
│   │   └── ...
│   └── migrations/        # Database migrations
│
├── frontend/              # React frontend
│   └── src/
│       ├── app/          # App initialization
│       ├── pages/        # Route pages
│       ├── widgets/      # Complex UI blocks
│       ├── features/     # User scenarios
│       ├── entities/     # Business entities
│       └── shared/       # Shared utilities
│
└── docker-compose.yml    # Docker orchestration
```

_See [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md) for detailed architecture_

## 💻 Development

### Local Development (without Docker)

**Backend**:
```bash
cd backend

# Install dependencies
go mod download

# Setup database (via Docker)
docker-compose up postgres redis

# Run backend with hot reload
go install github.com/cosmtrek/air@latest
air
```

**Frontend**:
```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev
```

### Testing

**Backend**:
```bash
cd backend

# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Frontend**:
```bash
cd frontend

# Run tests
npm test

# With coverage
npm test -- --coverage

# Watch mode
npm test -- --watch
```

### Linting & Formatting

**Backend**:
```bash
# Format code
go fmt ./...

# Run linter (install golangci-lint first)
golangci-lint run
```

**Frontend**:
```bash
# Lint
npm run lint

# Fix auto-fixable issues
npm run lint -- --fix
```

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Read documentation**
   - [GETTING_STARTED.md](GETTING_STARTED.md) - Setup guide
   - [DEVELOPMENT_WORKFLOW.md](DEVELOPMENT_WORKFLOW.md) - Git workflow

2. **Find or create an issue**
   - Check [open issues](https://github.com/quenyu/deadlock-stats/issues)
   - Create new issue for bugs/features

3. **Create a branch**
   ```bash
   git checkout -b fix/error-handling-backend
   # or
   git checkout -b feat/builds-api
   ```

4. **Make your changes**
   - Follow code style guidelines
   - Add tests for new features
   - Update documentation

5. **Submit a Pull Request**
   - Fill out PR template
   - Reference related issues
   - Wait for review

### Branch Naming Convention

```
fix/      - Bug fixes
feat/     - New features
refactor/ - Code refactoring
chore/    - Routine tasks (deps, CI/CD)
docs/     - Documentation
test/     - Tests
perf/     - Performance improvements
```

### Commit Message Format

```
<type>(<scope>): <subject>

Examples:
fix(auth): handle expired JWT tokens
feat(builds): implement voting system
refactor(ui): extract skeleton components
docs(readme): add setup instructions
```

## 🗺️ Roadmap

### Phase 1: Stabilization (4-6 weeks) 🔴
- ✅ Error handling & validation
- ✅ Rate limiting & security
- ✅ Database optimization
- ✅ Monitoring & CI/CD
- ✅ Unit tests (60%+ coverage)

### Phase 2: Code Quality (4-6 weeks) 🟢
- React Query integration
- Zod validation
- Performance optimization
- Enhanced documentation

### Phase 3: Hero Builds (6-8 weeks) 🎮
- CRUD API for builds
- Vote & comment system
- Build creator UI
- AI recommendations

### Phase 4-9: Advanced Features
- Crosshairs system
- Leaderboards & analytics
- Social features
- Premium tier
- Mobile app (PWA)

_See [ROADMAP.md](ROADMAP.md) for complete timeline_

## 📊 Current Status

**Version**: 0.1.0-alpha
**Status**: Active Development 🚧

### Metrics
- **Code Coverage**: Backend ~40%, Frontend ~20%
- **Database**: 15 migrations
- **API Endpoints**: 11 endpoints
- **UI Pages**: 5 pages

### Known Issues
- ⚠️ No proper error handling (generic 500 errors)
- ⚠️ No rate limiting (vulnerable to abuse)
- ⚠️ Missing database indexes (slow queries)
- ⚠️ Console.log in production frontend

_See TODO list in IDE for full details_

## 📸 Screenshots

_Coming soon..._

## 🎯 Performance Targets

- ⚡ API Response: <100ms (p95)
- ⚡ Page Load: <2s
- ⚡ Lighthouse Score: 90+
- ⚡ Uptime: 99.9%

## 🔒 Security

- ✅ HTTPS enforcement
- ✅ JWT authentication
- ✅ SQL injection protection (GORM)
- ✅ CORS configuration
- 🚧 Rate limiting _(planned)_
- 🚧 Input validation _(planned)_
- 🚧 CSRF protection _(planned)_

Report security vulnerabilities to: [security@example.com]

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

### Technologies
- [Go](https://go.dev/) - The Go Authors
- [Echo](https://echo.labstack.com/) - LabStack
- [React](https://react.dev/) - Meta
- [Radix UI](https://www.radix-ui.com/) - WorkOS
- [Valve](https://www.valvesoftware.com/) - For Deadlock!

### Data Sources
- [Deadlock Wiki](https://deadlock.wiki/) - Game data & documentation
- Deadlock API - Player stats & match data

### Contributors
Thanks to all contributors! 🎉

_See [Contributors](https://github.com/quenyu/deadlock-stats/graphs/contributors)_

## 📞 Contact

- **GitHub**: [@quenyu](https://github.com/wqeqadas)
- **Issues**: [Create an Issue](https://github.com/quenyu/deadlock-stats/issues/new)
- **Discussions**: [Join Discussion](https://github.com/quenyu/deadlock-stats/discussions)

## 🌟 Show Your Support

If you like this project, please consider:
- ⭐ Starring the repository
- 🐛 Reporting bugs
- 💡 Suggesting new features
- 🤝 Contributing code
- 📢 Sharing with the community

---

**Made with ❤️ for the Deadlock community**

_Last updated: 2025-10-07_

