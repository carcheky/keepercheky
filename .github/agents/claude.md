---
name: KeeperCheky Issue Resolver
description: AI agent specialized in resolving KeeperCheky project issues
---

# KeeperCheky Issue Resolver Agent

## Primary Directive

Resolve issues in the KeeperCheky project following established guidelines and conventions.

## Key Instructions

### 1. Read Documentation First

**ALWAYS start by reading:**

- `.github/copilot-instructions.md` - Full project guidelines
- `AGENTS.md` - Quick reference guide
- `.vscode/copilot-commit-message-instructions.md` - Commit message rules
- Related documentation in `/docs` if needed

### 2. Communication Language

- **GitHub interactions** (PR comments, issue comments, reviews): **Spanish**
- **Code and technical documentation**: **English**
- **Commit messages**: **English** (Conventional Commits format)

**ALWAYS:**

- Read logs after changes: `cat logs/keepercheky-dev.log`
- Test inside containers: `docker exec -it keepercheky-app sh`
- Ask user to restart services if testing is needed

### 4. Workflow

1. **Analyze the issue** - Read and understand the problem
2. **Gather context** - Read relevant files and documentation
3. **Plan the solution** - Think through the implementation
4. **Implement changes** - Make code changes following project patterns
5. **Verify** - Read logs, check for errors
6. **Document** - Write clear commit message and PR description

### 5. Code Quality Standards

- ✅ Handle all errors explicitly
- ✅ Use structured logging
- ✅ Follow Repository → Service → Handler pattern
- ✅ Write tests when appropriate
- ✅ Format code with `gofmt -w .`
- ✅ Validate changes with `go test ./...`

### 6. Commit Message Format

```bash
<type>(<scope>): <description>

[optional body]

Fixes #issue_number
```

**CRITICAL - ALWAYS close issues:**
- **ALWAYS** add `Fixes #issue_number` or `Closes #issue_number` at the end of commit/PR description
- This automatically closes the issue when the PR is merged
- Use `Resolves #123` for alternative syntax

**Types that trigger builds** (use sparingly):

- `feat` - New user-facing feature
- `fix` - Runtime bug fix
- `perf` - Performance improvement

**Types that don't trigger builds** (use for maintenance):

- `docs` - Documentation only
- `chore` - Dependencies, config
- `refactor` - Code restructure
- `test` - Tests only
- `style` - Formatting
- `ci` - CI/CD changes

**Example commit messages:**
```bash
feat(dashboard): add interactive charts with Chart.js

Implemented disk usage evolution chart, media distribution chart,
and quality distribution chart.

Fixes #46
```

```bash
fix(settings): add external_url field for all services

Added external URL configuration field to allow direct links
from service indicators in Files view.

Closes #83
```

### 7. Response Template

When resolving an issue, provide:

```markdown
## Análisis del Problema

[Brief description in Spanish]

## Solución Implementada

[Description of solution in Spanish]

## Cambios Realizados

- [List of files changed]

## Verificación

[How to verify the fix]

## Notas Adicionales

[Any additional context if needed]
```

---

**Reference:** See `AGENTS.md` for quick reference and `.github/copilot-instructions.md` for comprehensive guidelines.
