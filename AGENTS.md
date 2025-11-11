# AGENTS - KeeperCheky Quick Reference

> **Media Library Cleanup Manager** - Go + Fiber + Alpine.js | [Full Guide](.github/copilot-instructions.md)

**Stack:** Go 1.22+ • Fiber v2 • GORM v2 • Alpine.js 3.x • Tailwind CSS • Docker

---

**✅ You CAN:**

- Read logs: `cat/tail/grep logs/keepercheky-dev.log`
- Debug containers: `docker exec -it keepercheky-app sh`
- Build/test: `go build`, `go test ./...`
- Edit code and config files

> **🔄 Need to test? Ask user to restart service**

---

## 🔗 CRITICAL: PR and Issue Linking

**⚠️ WHEN CREATING A PR FROM AN ASSIGNED ISSUE:**

**ALWAYS start PR description with:**
```markdown
Closes #[ISSUE_NUMBER]
```

**Why:** This automatically links and closes the issue when PR is merged.

**Supported keywords:** `Closes`, `Fixes`, `Resolves` (with or without `:`)

**Example:**
```markdown
Closes #42

## Summary
Implemented feature X...
```

**⚠️ NEVER create a PR without linking it to its issue!**

---

## 📂 Project Structure

```txt
keepercheky/
├── cmd/server/main.go              # Entry point
├── internal/                       # Private code (NOT importable)
│   ├── config/, models/, repository/
│   ├── service/                    # Business logic
│   │   ├── clients/                # Radarr, Sonarr, etc.
│   │   ├── cleanup/                # Cleanup strategies
│   │   └── scheduler/              # Cron jobs
│   ├── handler/                    # HTTP handlers (Fiber)
│   └── middleware/
├── web/
│   ├── templates/                  # Go templates
│   └── static/                     # CSS, JS, images
├── pkg/                            # Public packages
└── docs/                           # Documentation
```

**Patterns:** Repository → Service → Handler → Client Interface

---

## ⚡ Quick Commands

### ✅ Validation (ALWAYS run before commit)

```bash
make validate          # 🔍 Full validation (MANDATORY before commit)
make validate-quick    # ⚡ Quick check (format, vet, test)
make check-and-fix     # 🔧 Auto-fix + validate
make lint-check        # 📝 Check format only
make lint-fix          # 🛠️ Fix format automatically
```

### Build & Test

```bash
# Build
go build -o bin/keepercheky ./cmd/server

# Test
go test ./...              # All tests
go test -cover ./...       # With coverage

# Code quality
gofmt -w .                 # Format
go vet ./...               # Vet
golangci-lint run          # Lint (if installed)
```

---

---

## 🐛 Debugging

### Logs (ALWAYS read after changes)

```bash
cat logs/keepercheky-dev.log           # Full log
tail -n 100 logs/keepercheky-dev.log   # Last 100 lines
tail -f logs/keepercheky-dev.log       # Follow real-time
grep -i error logs/keepercheky-dev.log # Search errors
```

### Containers

```bash
docker ps                              # List containers
docker logs keepercheky-app            # View logs
docker exec -it keepercheky-app sh     # Enter container
```

### Service Health

```bash
curl http://localhost:8000/health      # Health check
curl http://localhost:8000/api/media   # API test
```

---

---

## 🎨 Frontend (Alpine.js)

### Component Example

```html
<div x-data="mediaCard({{ .Media.ID }})" class="card">
    <h4 x-text="media.title"></h4>
    <button @click="exclude()">Exclude</button>
</div>

<script>
function mediaCard(mediaId) {
    return {
        media: null,
        async init() { await this.fetchMedia(); },
        async exclude() { /* ... */ }
    }
}
</script>
```

### API Calls Pattern

```javascript
{
    data: null,
    loading: false,
    error: null,
    
    async fetch() {
        this.loading = true;
        this.error = null;
        try {
            const res = await fetch('/api/endpoint');
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            this.data = await res.json();
        } catch (err) {
            this.error = err.message;
        } finally {
            this.loading = false;
        }
    }
}
```

---

---

## 📝 Code Guidelines

### Error Handling

```go
// ❌ BAD
result, _ := someFunction()

// ✅ GOOD
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed: %w", err)
}
```

### Logging

```go
logger.Info("Processing",
    "count", len(items),
    "dry_run", cfg.DryRun,
)
```

### Database Transactions

```go
return db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Save(item).Error; err != nil {
        return err
    }
    return nil
})
```

---

---

## 🔄 Git Workflow

### Commit Types

**⚠️ Only `feat`, `fix`, `perf` trigger releases!**

**Trigger builds:**

- `feat` - New user-facing feature
- `fix` - Runtime bug fix
- `perf` - Performance improvement

**No builds:**

- `docs` - Documentation only
- `chore` - Dependencies, config
- `refactor` - Code restructure (no behavior change)
- `test` - Tests only
- `style` - Formatting (gofmt)
- `ci` - CI/CD changes

### Examples

```bash
# TRIGGERS BUILD
feat(api): add bulk deletion endpoint
fix(sync): correct hash matching

# NO BUILD
docs(config): update .env.example
chore(deps): update dependencies
```

> See [commit instructions](.vscode/copilot-commit-message-instructions.md) for details

---

---

## 🚀 Common Development Tasks

| Task | Command | Notes |
|------|---------|-------|
| **Validate before commit** | `make validate` | ⚠️ **MANDATORY** before every commit |
| Quick validation | `make validate-quick` | Format, vet, tests |
| Auto-fix + validate | `make check-and-fix` | Fix format, then validate |
| Format code | `make lint-fix` | Auto-fix gofmt issues |
| Build binary | `go build -o bin/keepercheky ./cmd/server` | Development build |
| Run tests | `go test ./...` | All tests |
| Vet code | `go vet ./...` | Static analysis |
| Read logs | `cat logs/keepercheky-dev.log` | After changes |
| Tail logs | `tail -f logs/keepercheky-dev.log` | Real-time |
| Check health | `curl http://localhost:8000/health` | Service status |
| List containers | `docker ps` | Running services |
| Exec in container | `docker exec -it keepercheky-app sh` | Debug inside |

---

## 📊 Makefile Commands Reference

**⚠️ DO NOT RUN THESE - Only for reference:**

| Command | Purpose | **DO NOT USE** |
|---------|---------|----------------|
| `make dev` | Start development environment | ❌ User only |
| `make stop` | Stop services | ❌ User only |
| `make restart` | Restart services | ❌ User only |
| `make logs` | View logs | ❌ Use `cat` instead |
| `make build` | Build Docker image | ❌ User only |

**What you CAN use:**

```bash
# Read files
cat logs/keepercheky-dev.log
cat config/config.yaml

# Search files
grep -r "pattern" internal/

# Format and test
gofmt -w .
go test ./...
```

---

## 🔐 Security Best Practices

- ✅ Validate all user input
- ✅ Sanitize file paths to prevent directory traversal
- ✅ Use parameterized queries (GORM handles this)
- ✅ Never log sensitive data (API keys, passwords)
- ✅ Validate file types and sizes for uploads

```go
// Example: Path validation
func validatePath(path string) error {
    cleanPath := filepath.Clean(path)
    if strings.Contains(cleanPath, "..") {
        return fmt.Errorf("invalid path: contains '..'")
    }
    return nil
}
```

---

## 📚 Key Dependencies

```go
// Primary dependencies
github.com/gofiber/fiber/v2          // Web framework
github.com/gofiber/template/html/v2  // Template engine
gorm.io/gorm                         // ORM
gorm.io/driver/sqlite                // SQLite driver (dev)
gorm.io/driver/postgres              // PostgreSQL driver (prod)
github.com/go-resty/resty/v2         // HTTP client
github.com/robfig/cron/v3            // Job scheduler
github.com/spf13/viper               // Configuration
go.uber.org/zap                      // Structured logging
```

---

## 🗣️ Communication Guidelines

**REMEMBER:**

- 📢 Always respond to users in **Spanish**
- 📝 Technical documentation and code comments in **English**
- 🐛 Issue titles and descriptions in **Spanish**
- 💬 Pull request descriptions in **Spanish**
- ✅ **ALL GitHub interactions MUST be in Spanish** (PR comments, code reviews, etc.)

---

## 🎯 Before Committing - Checklist

- [ ] **RUN `make validate`** ⚠️ MANDATORY
- [ ] All validation checks pass (format, vet, tests, build)
- [ ] No sensitive data in code
- [ ] Error messages are descriptive
- [ ] Logs use structured logging
- [ ] Commit message follows Conventional Commits format
- [ ] Commit message is in **English**

### Quick Pre-Commit Workflow

```bash
make check-and-fix    # Auto-fix + validate
git add .
git commit -m "feat: my changes"
git push
```

---

## 📖 Additional Documentation

For more detailed guidelines, see:

- `.github/copilot-instructions.md` - Full project guidelines and philosophy
- `.vscode/copilot-commit-message-instructions.md` - Detailed commit message rules
- `docs/` - Additional documentation and analysis
- `docs/JELLYFIN_API_REFERENCE.md` - Complete Jellyfin API documentation
- `docs/JELLYSEERR_API_REFERENCE.md` - Complete Jellyseerr API documentation

---

**Last Updated:** 2025-11-02  
**Format:** AGENTS.md v1.0 (OpenAI standard)  
**Project:** KeeperCheky - Media Library Cleanup Manager
