# 📊 Deadlock Stats - Project Overview

> **Comprehensive stats platform for Valve's Deadlock**

## 🎯 Project Vision

Создать **лучшую stats платформу** для Deadlock, объединяющую:
- 📈 **Детальную статистику** игроков и матчей
- 🏗️ **Систему билдов** с AI рекомендациями
- 🎯 **Кастомные прицелы** с визуальным редактором
- 📊 **Meta аналитику** и insights
- 👥 **Социальные функции** (friends, teams, tournaments)
- 💰 **Premium tier** с эксклюзивными фичами

---

## 🏗️ Tech Stack

### Backend
- **Language**: Go 1.23
- **Framework**: Echo (web framework)
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **ORM**: GORM
- **Logging**: Zap
- **Config**: Viper
- **Migrations**: golang-migrate

### Frontend
- **Framework**: React 19
- **Language**: TypeScript 5.8
- **Build Tool**: Vite 7
- **State Management**: Zustand (→ migrating to React Query)
- **UI Components**: Radix UI
- **Styling**: Tailwind CSS 4
- **Charts**: Recharts
- **Routing**: React Router v7
- **Forms**: React Hook Form (planned)
- **Validation**: Zod 4

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Reverse Proxy**: Nginx
- **Monitoring**: Prometheus + Grafana (planned)
- **Error Tracking**: Sentry (planned)
- **CI/CD**: GitHub Actions (planned)

### Architecture
- **Backend**: Clean Architecture (Handlers → Services → Repositories)
- **Frontend**: Feature-Sliced Design (FSD)
- **API**: RESTful (GraphQL planned for future)
- **Auth**: JWT + Steam OpenID

---

## 📁 Project Structure

```
deadlock-stats/
├── backend/
│   ├── cmd/
│   │   └── main.go                    # Entry point
│   ├── internal/
│   │   ├── handlers/                  # HTTP handlers (controllers)
│   │   ├── services/                  # Business logic
│   │   ├── repositories/              # Data access layer
│   │   ├── domain/                    # Domain models
│   │   ├── dto/                       # Data Transfer Objects
│   │   ├── clients/                   # External API clients
│   │   ├── middleware/                # Custom middlewares
│   │   └── config/                    # Configuration
│   ├── migrations/                    # Database migrations
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── app/                       # App initialization
│   │   │   ├── providers/             # Global providers (router, theme, query)
│   │   │   └── styles/                # Global styles
│   │   ├── pages/                     # Page components (routes)
│   │   ├── widgets/                   # Complex UI blocks
│   │   ├── features/                  # User scenarios (AuthBySteam, PlayerSearch)
│   │   ├── entities/                  # Business entities (player, user, deadlock)
│   │   │   ├── player/
│   │   │   │   ├── api/               # API calls
│   │   │   │   ├── model/             # State management
│   │   │   │   ├── types/             # TypeScript types
│   │   │   │   └── utils/             # Helper functions
│   │   │   └── ...
│   │   └── shared/                    # Shared utilities
│   │       ├── ui/                    # UI components (Button, Card, etc.)
│   │       ├── lib/                   # Utils, helpers
│   │       ├── api/                   # Axios instance
│   │       └── constants/             # Constants
│   ├── Dockerfile
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
│
├── docker-compose.yml                 # Docker orchestration
├── ROADMAP.md                         # Development roadmap
├── DEVELOPMENT_WORKFLOW.md            # Git workflow & conventions
├── GETTING_STARTED.md                 # Quick start guide
└── README.md                          # Main documentation
```

---

## 🚀 Current Status

### ✅ Implemented Features

**Authentication**
- ✅ Steam OpenID login
- ✅ JWT token management
- ✅ Protected routes

**Player Profiles**
- ✅ Extended player profiles with stats
- ✅ Match history
- ✅ Hero statistics
- ✅ MMR tracking & history
- ✅ Performance dynamics
- ✅ Personal records
- ✅ Featured heroes
- ✅ Peak rank tracking

**Search**
- ✅ Player search by Steam ID
- ✅ Player search by nickname
- ✅ Fuzzy search with external API

**Caching**
- ✅ Redis caching with TTL
- ✅ Partial cache fallback
- ✅ Smart cache invalidation

**Database**
- ✅ 15 migrations (users, stats, matches, builds, crosshairs, votes, comments, tags)
- ✅ Normalized schema
- ✅ Proper relationships

**UI/UX**
- ✅ Dark theme
- ✅ Responsive design
- ✅ Modern UI with Radix components
- ✅ Charts for stats visualization

### 🚧 In Progress

- 🔨 Error handling improvements
- 🔨 Security hardening
- 🔨 Performance optimization
- 🔨 Testing infrastructure

### 📅 Planned Features

**Phase 1**: Stabilization (4-6 weeks)
- Error handling & validation
- Rate limiting
- Security headers
- Database optimization
- Monitoring & metrics
- CI/CD pipeline
- Unit tests (60%+ coverage)

**Phase 2**: Code Quality (4-6 weeks)
- React Query integration
- Zod validation
- Skeleton loaders
- Code splitting
- Documentation
- Frontend tests

**Phase 3**: Hero Builds (6-8 weeks)
- CRUD API for builds
- Vote & comment system
- Tagging & filtering
- Build creator UI
- AI recommendations

**Phase 4**: Crosshairs (4-6 weeks)
- Visual editor
- Gallery
- Pro configs
- Import/Export

**Phase 5**: Analytics (6-8 weeks)
- Global leaderboards
- Hero leaderboards
- Meta analysis
- Counter picks
- Item analytics

**Phase 6**: Social (4-6 weeks)
- Friends system
- Profile comparison
- Teams/Clans
- Match sharing

**Phase 7**: Premium (4-6 weeks)
- Subscription system
- Advanced analytics
- Profile customization
- Ad-free experience

**Phase 8**: Mobile (6-8 weeks)
- PWA
- React Native app (optional)

---

## 📊 Architecture Decisions

### Backend Patterns

**Clean Architecture** - разделение на слои:
```
HTTP Request → Handler → Service → Repository → Database
                  ↓
            Domain Models
```

**Dependency Injection** - через constructor injection:
```go
service := NewPlayerProfileService(
    repository,
    apiClient,
    cache,
    logger,
)
```

**Repository Pattern** - абстракция над БД:
```go
type PlayerProfileRepository interface {
    FindBySteamID(ctx context.Context, steamID string) (*PlayerProfile, error)
    Update(ctx context.Context, profile *PlayerProfile) error
}
```

### Frontend Patterns

**Feature-Sliced Design** - модульная архитектура:
```
entities/ - бизнес-сущности (player, user)
features/ - пользовательские сценарии (auth, search)
widgets/  - сложные UI блоки (navbar, profile card)
pages/    - страницы приложения
shared/   - переиспользуемые компоненты
```

**Custom Hooks** - переиспользование логики:
```typescript
const { profile, loading, error } = usePlayerProfile(steamId)
```

**State Management**:
- **Zustand** - для UI state (user, theme)
- **React Query** (planned) - для server state (API data)

---

## 🎯 Key Metrics & Goals

### Performance Targets
- ⚡ API Response Time: <100ms (p95)
- ⚡ Page Load Time: <2s
- ⚡ Time to Interactive: <3s
- ⚡ Lighthouse Score: 90+

### Quality Targets
- ✅ Test Coverage: 80%+
- ✅ Uptime: 99.9%
- ✅ Error Rate: <0.1%

### User Metrics (6 months)
- 👥 10,000+ registered users
- 👥 1,000+ daily active users
- 👥 100,000+ profiles viewed
- 👥 10,000+ builds created
- 👥 100+ premium subscribers

---

## 🔒 Security Considerations

### Current
- ✅ HTTPS enforcement
- ✅ JWT authentication
- ✅ GORM SQL injection protection
- ✅ CORS configuration

### Planned
- 🔨 Rate limiting
- 🔨 Input validation
- 🔨 CSRF protection
- 🔨 Security headers
- 🔨 Secure cookies
- 🔨 SQL injection audits
- 🔨 Dependency scanning

---

## 📈 Scalability Strategy

### Current Setup
- Single server deployment
- PostgreSQL with connection pooling
- Redis caching
- Docker containerization

### Future (when needed)
- **Horizontal Scaling**: Multiple backend instances behind load balancer
- **Database**: Read replicas for read-heavy operations
- **Caching**: Redis Cluster for distributed cache
- **CDN**: Cloudflare/CloudFront for static assets
- **Message Queue**: RabbitMQ/Redis for async tasks
- **Monitoring**: Prometheus + Grafana + Loki

---

## 🤝 Contributing

### How to Contribute

1. **Pick a task** from TODO list or Issues
2. **Create a branch** following naming convention
3. **Make changes** following code style
4. **Test thoroughly** before pushing
5. **Create PR** with detailed description
6. **Address review** comments
7. **Merge** after approval

### Code Style

**Go**:
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Comment exported functions
- Handle all errors

**TypeScript/React**:
- Strict TypeScript (no `any`)
- Functional components + hooks
- Proper types for everything
- ESLint + Prettier

### Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):
```
<type>(<scope>): <subject>

fix(auth): handle expired tokens
feat(builds): add voting system
refactor(ui): extract common components
docs(readme): add setup instructions
test(services): add player service tests
```

---

## 📚 Learning Resources

### Deadlock
- [Deadlock Wiki](https://deadlock.wiki/) - Game data & stats
- [Deadlock API Docs](https://docs.deadlock-api.com/) - API reference

### Technologies
- [Go Documentation](https://go.dev/doc/)
- [Echo Framework](https://echo.labstack.com/)
- [React Documentation](https://react.dev/)
- [Feature-Sliced Design](https://feature-sliced.design/)
- [Radix UI](https://www.radix-ui.com/)

### Best Practices
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [12 Factor App](https://12factor.net/)
- [API Design Best Practices](https://swagger.io/resources/articles/best-practices-in-api-design/)

---

## 🐛 Known Issues

### Critical 🔴
- No proper error handling on backend (returns generic 500)
- Potential goroutine deadlock in `fetchAllData`
- No input validation for Steam IDs

### High Priority 🟡
- No rate limiting (vulnerable to abuse)
- Console.log in production frontend
- Missing database indexes (slow queries)

### Medium Priority 🟢
- No React Query caching on frontend
- No Zod validation for API responses
- Loading states show plain text instead of skeletons

### Low Priority 🔵
- Missing code splitting (large bundle)
- No PWA support
- No error tracking (Sentry)

_See TODO list for full details and action items_

---

## 📞 Contact & Support

### Maintainers
- **GitHub**: [@quenyu](https://github.com/quenyu)

### Links
- **Repository**: https://github.com/quenyu/deadlock-stats
- **Issues**: https://github.com/quenyu/deadlock-stats/issues
- **Discussions**: https://github.com/quenyu/deadlock-stats/discussions

### Getting Help
1. Check documentation files (README, GETTING_STARTED, etc.)
2. Search existing issues
3. Create new issue with detailed description
4. Join Discord (if available)

---

## 🎉 Acknowledgments

### Technologies
Спасибо авторам и мейнтейнерам:
- Go Team
- Echo Framework
- React Team
- Radix UI Team
- Valve (за Deadlock!)

### Contributors
Спасибо всем контрибьюторам проекта! 🙏

---

## 📄 License

_TBD - Нужно добавить LICENSE файл_

Рекомендуется: **MIT License** (для open-source проектов)

---

**Last Updated**: 2025-10-07
**Version**: 0.1.0-alpha
**Status**: Active Development 🚧

---

> 💡 **Tip**: Этот документ регулярно обновляется. Для актуальной информации смотрите также ROADMAP.md и TODO list.

