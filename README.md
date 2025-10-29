# KeeperCheky

> Modern web-based media library cleanup manager - A complete rewrite of Janitorr with a beautiful UI

[![stable](https://img.shields.io/github/actions/workflow/status/carcheky/keepercheky/release.yml?branch=stable&label=stable&logo=github)](https://github.com/carcheky/keepercheky/actions/workflows/release.yml)
[![stable version](https://img.shields.io/github/v/release/carcheky/keepercheky?label=stable)](https://github.com/carcheky/keepercheky/releases)
[![develop](https://img.shields.io/github/actions/workflow/status/carcheky/keepercheky/release.yml?branch=develop&label=develop&logo=github)](https://github.com/carcheky/keepercheky/actions/workflows/release.yml)
[![develop version](https://img.shields.io/github/v/release/carcheky/keepercheky?include_prereleases&label=develop&filter=*-dev*)](https://github.com/carcheky/keepercheky/releases)
[![Docker Image](https://img.shields.io/badge/docker-ghcr.io-blue?logo=docker)](https://github.com/carcheky/keepercheky/pkgs/container/keepercheky)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/carcheky/keepercheky)](go.mod)

## 📖 What is KeeperCheky?

KeeperCheky is a complete rewrite of [Janitorr](https://github.com/Schaka/janitorr) with a modern web interface. It automatically manages and cleans up your media library by removing old or unwatched content based on configurable rules.

### Key Features

- 🎨 **Modern Web UI** - Beautiful dashboard accessible from any browser
- 🧹 **Automated Cleanup** - Smart deletion based on age and disk space
- 🏷️ **Tag-based Rules** - Custom cleanup schedules using Radarr/Sonarr tags
- 📺 **Series Management** - Special handling for weekly/daily shows
- ⏰ **Leaving Soon** - Preview content before deletion in Jellyfin/Emby
- 🔗 **Full Integration** - Works with Radarr, Sonarr, Jellyfin, Emby, Jellyseerr
- 🐳 **Docker Ready** - Easy deployment with Docker Compose
- 🔒 **Safe Mode** - Dry-run by default to prevent accidents

## 🎯 Project Status

**Current Phase**: 🚀 Active Development - MVP Implementation

We have completed the planning phase and chosen **Go + Alpine.js** (Proposal 3) for optimal performance and minimal resource usage. The project structure is now complete and ready for development.

### Documentation

📚 **[Start Here: Documentation Index](docs/README.md)**

Quick links:
- [Executive Summary](docs/RESUMEN_EJECUTIVO.md) - Project overview and analysis
- [Maintainerr Analysis](docs/ANALISIS_MAINTAINERR.md) - Evaluation of similar projects 🆕
- [Quick Comparison](docs/RESUMEN_COMPARATIVO.md) - Janitorr vs Maintainerr vs KeeperCheky 🆕
- [Detailed Comparison](docs/COMPARACION_Y_RECOMENDACIONES.md) - Stack analysis and recommendations
- [Proposals](docs/propuestas/) - 4 complete technical proposals

## 🏗️ Proposed Architecture

We've developed **4 complete proposals** with different technology stacks:

### 1. TypeScript Full-Stack (Next.js + NestJS)
- **Best for**: Modern UX, scalability
- **Resources**: 512MB-1GB RAM
- **Development**: 4-6 weeks

### 2. Python/HTMX (FastAPI + HTMX)
- **Best for**: Simplicity, rapid development
- **Resources**: 50-150MB RAM
- **Development**: 3-4 weeks

### 3. Go/Alpine.js (Fiber + Alpine) ⭐ **RECOMMENDED**
- **Best for**: Performance, minimal resources
- **Resources**: 20-50MB RAM
- **Development**: 3-4 weeks
- **Docker image**: 15-25MB

### 4. Rust/Leptos (Axum + WASM)
- **Best for**: Maximum security, long-term projects
- **Resources**: 30-80MB RAM
- **Development**: 5-7 weeks

## 🚀 Recommended Implementation

### Why Go + Alpine.js?

1. ✅ **Optimal balance** - Performance meets resources
2. ✅ **Minimal footprint** - 20-50MB RAM, tiny Docker image
3. ✅ **Easy deployment** - Single binary, no dependencies
4. ✅ **Fast development** - 3-4 weeks to MVP
5. ✅ **Scalable** - Goroutines for concurrency

See [Comparison Document](docs/COMPARACION_Y_RECOMENDACIONES.md) for detailed analysis.

## 📊 Comparison with Similar Projects

| Feature              | Janitorr        | Maintainerr                   | KeeperCheky (Goal) |
| -------------------- | --------------- | ----------------------------- | ------------------ |
| **Web Interface**    | ❌               | ✅                             | ✅                  |
| **Stack**            | Kotlin + Spring | TypeScript + NestJS + Next.js | **Go + Alpine.js** |
| **Docker Image**     | ~300MB          | ~500MB                        | **15-25MB** ✅      |
| **RAM Usage**        | ~256MB          | ~400-600MB                    | **20-50MB** ✅      |
| **Startup Time**     | ~10-15s         | ~15-25s                       | **<1s** ✅          |
| **Jellyfin Support** | ✅               | ❌                             | ✅                  |
| **Plex Support**     | ❌               | ✅                             | ✅ (future)         |
| **Rule Builder**     | Config file     | ✅ Visual                      | ✅ Visual           |
| **Leaving Soon**     | ❌               | ✅ Collections                 | ✅ Symlinks         |

**See [detailed comparison](docs/RESUMEN_COMPARATIVO.md) for more information.**

## 🎨 Planned Features

### Dashboard
- Real-time disk usage statistics
- Service health monitoring
- Upcoming deletions preview
- Recent activity timeline

### Media Management
- Visual library browser with posters
- Advanced filtering and search
- One-click exclude/delete
- Detailed media information

### Cleanup Schedules
- Create/edit/delete schedules from UI
- Preview what will be deleted
- Manual execution with confirmation
- Enable/disable schedules

### Settings
- Service configuration forms
- Real-time connection testing
- Path validation
- No more manual YAML editing

### Logs & History
- Live log streaming
- Filter by level (INFO, ERROR, etc.)
- Search and export
- Complete action history

## 🛠️ Technology Stack (Recommended)

**Backend:**
- Go 1.22+
- Fiber v2 (web framework)
- GORM v2 (ORM)
- SQLite/PostgreSQL
- go-cron (scheduler)

**Frontend:**
- Alpine.js 3.x (15kb reactive framework)
- Tailwind CSS
- Chart.js (minimal charts)

**DevOps:**
- Docker (multi-stage builds)
- Single static binary
- ~15-25MB final image

## 📦 Installation (Coming Soon)

```bash
# Docker Compose (recommended)
docker-compose up -d

# Docker
docker run -d \
  --name keepercheky \
  -p 8000:8000 \
  -v ./config:/config \
  -v ./data:/data \
  -v /path/to/media:/media \
  keepercheky/keepercheky:latest

# Binary (standalone)
./keepercheky
```

## 🔧 Configuration Example (Coming Soon)

```yaml
# config.yml
general:
  dry_run: true
  leaving_soon_days: 14
  
clients:
  radarr:
    enabled: true
    url: "http://radarr:7878"
    api_key: "your-api-key"
    
  sonarr:
    enabled: true
    url: "http://sonarr:8989"
    api_key: "your-api-key"
    
  jellyfin:
    enabled: true
    url: "http://jellyfin:8096"
    api_key: "your-api-key"
    username: "janitor"
    password: "password"

schedules:
  media_cleanup:
    enabled: true
    expiration:
      5: 15d   # At 5% free space, delete 15+ day old media
      10: 30d
      15: 60d
      20: 90d
```

## 🗺️ Roadmap

### Phase 1: MVP (Weeks 1-4) - ✅ IN PROGRESS
- [x] Setup project structure
- [x] Core backend + database models
- [ ] Service clients (Radarr, Sonarr, Jellyfin) - in progress
- [ ] Cleanup logic implementation
- [x] Basic UI with Alpine.js
- [x] Dashboard + Media Management pages (templates ready)

### Phase 2: Features (Weeks 5-6)
- [ ] All service integrations
- [ ] Schedules management
- [ ] Settings page
- [ ] Logs viewer
- [ ] Docker optimization

### Phase 3: Polish (Week 7)
- [ ] Testing
- [ ] Documentation
- [ ] Docker Hub publish
- [ ] Release 1.0.0

## � AI Agent Development

This project includes comprehensive instructions for AI coding agents like GitHub Copilot:

- **[AGENTS.md](AGENTS.md)** - Practical guide for Copilot Coding Agent with commands and workflows
- **[.github/copilot-instructions.md](.github/copilot-instructions.md)** - Detailed project guidelines and architecture patterns

These files help AI agents understand the project structure, development workflows, and critical safety rules.

## �🤝 Contributing

This project is currently in the planning phase. Contributions will be welcome once we start implementation.

### Interested in helping?

1. Review the [technical proposals](docs/propuestas/)
2. Share your thoughts on the [recommended stack](docs/COMPARACION_Y_RECOMENDACIONES.md)
3. Check out [AGENTS.md](AGENTS.md) for development guidelines
4. Star this repo to stay updated

## 📝 License

MIT License - See [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

- **[Janitorr](https://github.com/Schaka/janitorr)** - Original project that inspired this rewrite
- **[Maintainerr](https://github.com/jorenn92/Maintainerr)** - Reference for UI/UX and advanced features
- All the *arr projects (Radarr, Sonarr, etc.)
- Jellyfin and Emby communities

## 📞 Links

- **Documentation**: [docs/README.md](docs/README.md)
- **Proposals**: [docs/propuestas/](docs/propuestas/)
- **Comparisons**: [docs/RESUMEN_COMPARATIVO.md](docs/RESUMEN_COMPARATIVO.md)
- **Janitorr Original**: [github.com/Schaka/janitorr](https://github.com/Schaka/janitorr)
- **Maintainerr**: [github.com/jorenn92/Maintainerr](https://github.com/jorenn92/Maintainerr)

---

**Note**: This is a rewrite/reimplementation project combining the best features of Janitorr (cleanup logic) and Maintainerr (beautiful UI), optimized for minimal resource usage with Go + Alpine.js.

**Status**: 🚀 MVP Implementation in progress - [Development Guide](DEVELOPMENT.md)

