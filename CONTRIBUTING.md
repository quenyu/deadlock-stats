# 🤝 Contributing to Deadlock Stats

First off, thank you for considering contributing to Deadlock Stats! 🎉

We welcome contributions from everyone. This document will help you get started.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Style Guidelines](#style-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Community](#community)

## 📜 Code of Conduct

This project follows a Code of Conduct. By participating, you are expected to uphold this code.

**Quick version**: Be respectful, be kind, be professional.

## 🎯 How Can I Contribute?

### 🐛 Reporting Bugs

Before creating bug reports:
- Check the [existing issues](https://github.com/quenyu/deadlock-stats/issues)
- Search closed issues - it might have been fixed already

When creating a bug report, use the **Bug Report** template and include:
- Clear title and description
- Steps to reproduce
- Expected vs actual behavior
- Screenshots if applicable
- Environment details

### ✨ Suggesting Features

Feature suggestions are tracked as GitHub issues. Use the **Feature Request** template and include:
- Clear description of the feature
- Use case / problem it solves
- Possible implementation ideas
- Mockups or examples

### 💻 Contributing Code

1. **Find an issue to work on**
   - Look for issues labeled `good first issue` for beginners
   - Comment on the issue to let others know you're working on it

2. **Fork & Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/deadlock-stats.git
   cd deadlock-stats
   ```

3. **Create a branch**
   ```bash
   git checkout -b fix/issue-description
   ```

4. **Make your changes**
   - Write code
   - Add tests
   - Update documentation

5. **Test thoroughly**
   ```bash
   # Backend
   go test ./...
   
   # Frontend
   npm test
   ```

6. **Commit & Push**
   ```bash
   git add .
   git commit -m "fix(scope): description"
   git push origin fix/issue-description
   ```

7. **Create Pull Request**
   - Fill out the PR template
   - Link related issues
   - Wait for review

## 🚀 Getting Started

### Prerequisites
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/) & Docker Compose
- [Go 1.23+](https://go.dev/) (for backend development)
- [Node.js 18+](https://nodejs.org/) (for frontend development)

### Setup Development Environment

See [GETTING_STARTED.md](GETTING_STARTED.md) for detailed setup instructions.

**Quick start:**
```bash
# Clone repository
git clone https://github.com/quenyu/deadlock-stats.git
cd deadlock-stats

# Start all services
docker-compose up
```

## 🔄 Development Workflow

See [DEVELOPMENT_WORKFLOW.md](DEVELOPMENT_WORKFLOW.md) for detailed workflow.

### Branch Naming

```
fix/      - Bug fixes
feat/     - New features
refactor/ - Code refactoring
chore/    - Routine tasks
docs/     - Documentation
test/     - Tests
perf/     - Performance improvements
```

### Example Workflow

```bash
# 1. Update main branch
git checkout main
git pull origin main

# 2. Create feature branch
git checkout -b feat/new-feature

# 3. Make changes
# ... write code ...

# 4. Test
go test ./...        # Backend
npm test            # Frontend

# 5. Commit
git add .
git commit -m "feat(scope): add new feature"

# 6. Push
git push origin feat/new-feature

# 7. Create Pull Request on GitHub
```

## 🎨 Style Guidelines

### Go Code Style

Follow [Effective Go](https://go.dev/doc/effective_go) guidelines:

```go
// ✅ Good
func GetPlayerProfile(steamID string) (*PlayerProfile, error) {
    if steamID == "" {
        return nil, ErrInvalidSteamID
    }
    
    // Implementation...
    return profile, nil
}

// ❌ Bad
func getprofile(id string) *PlayerProfile {
    // No error handling
    // Poor naming
}
```

**Key points:**
- Use `gofmt` for formatting
- Comment exported functions
- Handle all errors
- Use meaningful variable names
- Keep functions small and focused

### TypeScript/React Style

Follow [React best practices](https://react.dev/learn):

```typescript
// ✅ Good
interface PlayerProfileProps {
  steamId: string
  onLoad?: (profile: PlayerProfile) => void
}

export const PlayerProfile: React.FC<PlayerProfileProps> = ({ 
  steamId, 
  onLoad 
}) => {
  const { data, loading, error } = usePlayerProfile(steamId)
  
  if (loading) return <Skeleton />
  if (error) return <ErrorMessage error={error} />
  
  return <ProfileCard profile={data} />
}

// ❌ Bad
function playerProfile(props: any) {
  // Using 'any'
  // Not handling loading/error states
  // Poor component structure
}
```

**Key points:**
- Use TypeScript strictly (no `any`)
- Functional components + hooks
- Proper props typing
- Handle loading/error states
- Extract reusable logic into hooks

### File Organization

**Backend:**
```
internal/
├── handlers/      # HTTP handlers (thin layer)
├── services/      # Business logic (thick layer)
├── repositories/  # Data access (thin layer)
├── domain/        # Domain models
├── dto/           # Data Transfer Objects
└── middleware/    # Custom middlewares
```

**Frontend (FSD):**
```
src/
├── app/          # App initialization
├── pages/        # Route pages
├── widgets/      # Complex UI blocks
├── features/     # User scenarios
├── entities/     # Business entities
└── shared/       # Shared utilities
```

## 📝 Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `fix:` - Bug fix
- `feat:` - New feature
- `refactor:` - Code refactoring
- `docs:` - Documentation
- `test:` - Tests
- `perf:` - Performance
- `chore:` - Routine tasks
- `style:` - Code formatting
- `ci:` - CI/CD changes

### Examples

```bash
# Simple fix
git commit -m "fix(auth): handle expired JWT tokens"

# Feature with body
git commit -m "feat(builds): add voting system

- Users can upvote/downvote builds
- Vote count shown on build cards
- Sorting by votes implemented

Closes #123"

# Breaking change
git commit -m "refactor(api): change response format

BREAKING CHANGE: API responses now use snake_case instead of camelCase

Migration guide:
- Update all API calls to use new format
- Run migration script: npm run migrate"
```

### Rules

- ✅ Use present tense ("add feature" not "added feature")
- ✅ Use imperative mood ("move cursor to..." not "moves cursor to...")
- ✅ First line max 72 characters
- ✅ Reference issues and PRs
- ✅ Explain **what** and **why**, not **how**

## 🔍 Pull Request Process

### Before Submitting

1. **Test your changes**
   ```bash
   # Backend
   go test ./...
   go test -cover ./...
   
   # Frontend
   npm test
   npm run lint
   ```

2. **Update documentation**
   - Update README if needed
   - Add JSDoc/comments for complex code
   - Update API docs if endpoints changed

3. **Clean up**
   - Remove console.log
   - Remove commented code
   - Fix linting errors

### PR Checklist

- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added to complex code
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] No new warnings
- [ ] Branch up to date with main

### PR Template

Use the provided PR template. Include:
- Description of changes
- Type of change
- Testing performed
- Screenshots (for UI changes)
- Breaking changes (if any)

### Review Process

1. **Automated checks** run (CI/CD)
2. **Code review** by maintainer
3. **Address feedback** if requested
4. **Approval & merge** when ready

### After Merge

1. Delete your branch
2. Update your local main
3. Celebrate! 🎉

## 🧪 Testing

### Backend Tests

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Specific package
go test ./internal/services/...
```

### Frontend Tests

```bash
# Run tests
npm test

# Watch mode
npm test -- --watch

# Coverage
npm test -- --coverage

# Specific test
npm test -- PlayerProfile
```

### Writing Tests

**Backend:**
```go
func TestGetPlayerProfile(t *testing.T) {
    // Arrange
    mockRepo := new(MockRepository)
    service := NewPlayerProfileService(mockRepo)
    
    // Act
    profile, err := service.GetPlayerProfile("12345")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, profile)
    assert.Equal(t, "TestPlayer", profile.Nickname)
}
```

**Frontend:**
```typescript
describe('PlayerProfile', () => {
  it('should display player info', async () => {
    render(<PlayerProfile steamId="12345" />)
    
    await waitFor(() => {
      expect(screen.getByText('TestPlayer')).toBeInTheDocument()
    })
  })
})
```

## 📚 Documentation

### Code Comments

**When to comment:**
- Complex algorithms
- Non-obvious decisions
- Public APIs
- Edge cases

**Good comments:**
```go
// CalculateMMR computes the player's Match Making Rating based on
// recent performance. Uses weighted average with exponential decay
// to prioritize recent matches.
func CalculateMMR(matches []Match) int {
    // ...
}
```

**Bad comments:**
```go
// Get player
func GetPlayer() {} // Don't state the obvious
```

### API Documentation

Update Swagger annotations when changing APIs:

```go
// @Summary Get player profile
// @Description Get detailed player profile by Steam ID
// @Tags players
// @Accept json
// @Produce json
// @Param steamId path string true "Steam ID"
// @Success 200 {object} PlayerProfile
// @Failure 404 {object} ErrorResponse
// @Router /players/{steamId} [get]
func (h *PlayerProfileHandler) GetPlayerProfile(c echo.Context) error {
    // ...
}
```

## 🎯 Areas That Need Help

Looking for contributions in these areas:

### High Priority 🔴
- [ ] Error handling improvements
- [ ] Security hardening
- [ ] Test coverage
- [ ] Performance optimization

### Features 🎮
- [ ] Hero builds system
- [ ] Crosshairs editor
- [ ] Leaderboards
- [ ] Meta analysis

### Infrastructure 🏗️
- [ ] CI/CD pipeline
- [ ] Monitoring & alerts
- [ ] Documentation
- [ ] Automated testing

See [ROADMAP.md](ROADMAP.md) for complete list.

## 💬 Community

### Getting Help

- **Documentation**: Read [GETTING_STARTED.md](GETTING_STARTED.md)
- **Issues**: Search [existing issues](https://github.com/quenyu/deadlock-stats/issues)
- **Discussions**: Ask in [GitHub Discussions](https://github.com/quenyu/deadlock-stats/discussions)
- **Discord**: Join our Discord (link TBD)

### Communication Guidelines

- Be respectful and constructive
- Search before posting
- Provide context and details
- Follow up on your posts
- Help others when you can

## 🏆 Recognition

Contributors will be:
- Listed in README
- Mentioned in release notes
- Invited to private Discord channel (when available)

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

## 🙏 Thank You!

Every contribution matters, no matter how small. Thank you for helping make Deadlock Stats better!

**Questions?** Feel free to ask in [Discussions](https://github.com/quenyu/deadlock-stats/discussions)

---

**Happy Contributing! 🚀**

