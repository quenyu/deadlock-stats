<!--
SYNC IMPACT REPORT
==================
Version Change: Initial → 1.0.0
Constitution established: 2025-10-08

New Principles Added:
- I. Clean Architecture (Backend)
- II. Feature-Sliced Design (Frontend)
- III. Test-Driven Quality (NON-NEGOTIABLE)
- IV. Performance & Scalability
- V. Security First

Sections Added:
- Core Principles (5 principles)
- Technology Standards
- Development Workflow
- Governance

Templates Status:
- ✅ plan-template.md: Constitution Check section aligns with new principles
- ✅ spec-template.md: User stories and requirements align with quality standards
- ✅ tasks-template.md: Task organization supports independent testing (aligns with Principle III)
- ✅ All templates: Generic guidance maintained, no agent-specific references

Follow-up TODOs: None - all placeholders filled

Next Review: 2025-11-08
-->

# Deadlock Stats Constitution

## Core Principles

### I. Clean Architecture (Backend)

Every backend feature MUST follow Clean Architecture layers with clear separation of concerns:

- **Handlers** (HTTP layer) → handle requests/responses only, no business logic
- **Services** (Business logic) → orchestrate domain operations, transaction boundaries
- **Repositories** (Data access) → database operations only, return domain models
- **Domain Models** → pure data structures with business meaning, no infrastructure dependencies

**Rationale**: Separation ensures testability, maintainability, and enables independent evolution of each layer. Business logic remains framework-agnostic and can be tested without HTTP or database infrastructure.

**Rules**:
- Handlers MUST NOT contain business logic
- Services MUST NOT depend on HTTP-specific types
- Repositories MUST return domain models, not ORM entities
- Dependencies flow inward: Handlers → Services → Repositories → Domain

### II. Feature-Sliced Design (Frontend)

Frontend architecture MUST follow FSD methodology with strict layer hierarchy:

- **app/** → application initialization (providers, routes, global styles)
- **pages/** → route components, composition only
- **widgets/** → complex UI blocks combining features
- **features/** → user scenarios (AuthBySteam, PlayerSearch, BuildVoting)
- **entities/** → business entities (player, match, build, crosshair)
- **shared/** → reusable UI components, utilities, API client

**Rationale**: FSD provides predictable structure, prevents cyclic dependencies, and enables independent feature development. Each layer has clear responsibility and import restrictions.

**Rules**:
- Layers MUST NOT import from layers above them (e.g., entities cannot import features)
- Features MUST be self-contained and independently deletable
- Shared layer MUST NOT depend on business logic
- Pages MUST only compose, not implement logic

### III. Test-Driven Quality (NON-NEGOTIABLE)

Testing discipline MUST be maintained at all times:

- **Backend**: 60%+ code coverage REQUIRED before merging
- **Frontend**: Tests for critical user flows REQUIRED
- **Integration Tests**: REQUIRED for API contracts and external service interactions
- **Test-First**: Write tests → Verify failure → Implement → Verify pass

**Rationale**: Quality is not negotiable. Tests prevent regressions, document behavior, and enable confident refactoring. Test-first ensures testability is built in, not bolted on.

**Rules**:
- All PRs MUST include tests for new functionality
- Critical bugs MUST have regression tests before fixing
- Integration tests REQUIRED for: new API endpoints, external API changes, database schema changes
- Tests MUST fail before implementation (Red-Green-Refactor)
- Existing tests MUST NOT be deleted to make PRs pass

### IV. Performance & Scalability

System MUST meet performance targets and scale gracefully:

- **API Response Time**: <100ms (p95) for all endpoints
- **Page Load Time**: <2s for initial load
- **Time to Interactive**: <3s
- **Database Queries**: Properly indexed, N+1 queries prohibited
- **Caching Strategy**: Redis caching for expensive operations, TTL-based invalidation

**Rationale**: User experience depends on speed. Performance degradation is a bug. Scalability must be designed in from the start, not retrofitted later.

**Rules**:
- All database tables MUST have appropriate indexes
- Expensive operations MUST be cached
- API endpoints exceeding 100ms (p95) MUST be investigated and optimized
- Frontend bundles MUST use code splitting for routes
- Images MUST be lazy-loaded and optimized

### V. Security First

Security MUST be treated as a first-class concern:

- **Input Validation**: All user inputs MUST be validated and sanitized
- **Authentication**: JWT tokens with proper expiration and refresh
- **Authorization**: Role-based access control for protected resources
- **SQL Injection**: GORM parameterized queries ONLY, no string concatenation
- **Rate Limiting**: All public endpoints MUST have rate limits
- **Error Handling**: NEVER expose internal errors to clients

**Rationale**: Security vulnerabilities put users at risk and damage trust. Security must be built in, not added later. Every feature is a potential attack vector.

**Rules**:
- All Steam IDs MUST be validated before database queries
- Passwords MUST NEVER be logged or exposed
- Error messages MUST be generic (no stack traces to clients)
- Rate limiting REQUIRED: 100 req/min per IP for public endpoints
- Security headers MUST be configured (CORS, CSRF, XSS protection)
- Dependencies MUST be updated monthly for security patches

## Technology Standards

**Backend**:
- Language: Go 1.23+
- Framework: Echo (web framework)
- Database: PostgreSQL 16+ with GORM
- Cache: Redis 7+
- Logging: Zap (structured logging)
- Migrations: golang-migrate

**Frontend**:
- Language: TypeScript 5.8+ (strict mode, NO `any` types)
- Framework: React 19
- Build Tool: Vite 7
- State Management: Zustand (global UI state) + React Query (server state, planned)
- UI Library: Radix UI (accessible components)
- Styling: Tailwind CSS 4
- Validation: Zod (planned)

**Infrastructure**:
- Containerization: Docker + Docker Compose
- Monitoring: Prometheus + Grafana (planned)
- Error Tracking: Sentry (planned)
- CI/CD: GitHub Actions (planned)

## Development Workflow

**Branch Naming**: `<type>/<scope>-<description>`
- Types: `fix/`, `feat/`, `refactor/`, `chore/`, `docs/`, `test/`, `perf/`
- Example: `fix/error-handling-backend`, `feat/builds-api`

**Commit Messages**: Conventional Commits format
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Code Review Requirements**:
- All PRs MUST pass linter (golangci-lint, ESLint)
- All PRs MUST pass existing tests
- All PRs MUST include tests for new functionality
- Breaking changes MUST be documented in PR description
- Security-sensitive changes REQUIRE additional review

**Deployment Gates**:
- Lint passes
- Tests pass (60%+ backend coverage maintained)
- No console.log in production frontend
- No TODO comments in merged code (move to issues)

## Governance

**Constitution Authority**:
This constitution supersedes all other development practices and guidelines. When conflicts arise, constitution principles take precedence.

**Compliance Verification**:
- All PRs MUST be verified for compliance with Core Principles
- Architecture violations MUST be justified in `Complexity Tracking` section of plan.md
- Unjustified complexity MUST be rejected

**Amendment Process**:
1. Propose amendment with justification (GitHub issue)
2. Document impact on existing codebase
3. Create migration plan if breaking changes required
4. Require approval from project maintainer
5. Update version (see Versioning Policy)
6. Propagate changes to all dependent templates

**Versioning Policy** (Semantic Versioning):
- **MAJOR**: Backward incompatible principle removals or redefinitions (e.g., removing Clean Architecture requirement)
- **MINOR**: New principles added or material expansion of guidance (e.g., adding new principle VI)
- **PATCH**: Clarifications, wording fixes, non-semantic refinements (e.g., fixing typos, clarifying existing rules)

**Runtime Guidance**:
For day-to-day development guidance not in the constitution, refer to:
- `DEVELOPMENT_WORKFLOW.md` - Git workflow, commands, troubleshooting
- `PROJECT_OVERVIEW.md` - Architecture details, tech stack rationale
- `ROADMAP.md` - Feature priorities and timeline

**Enforcement**:
- Maintainers MUST enforce constitution compliance in code reviews
- Violations blocking merge: test coverage below 60%, missing input validation, architecture violations
- Warnings (fix in follow-up): missing indexes, console.log in non-critical paths, minor performance issues

**Version**: 1.0.0 | **Ratified**: 2025-10-08 | **Last Amended**: 2025-10-08
