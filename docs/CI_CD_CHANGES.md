# 🔄 Cambios en CI/CD - Octubre 2025

## Resumen

Este documento describe los cambios realizados en los workflows de GitHub Actions para optimizar el proceso de CI/CD de KeeperCheky.

## Problema

Anteriormente, el workflow `docker-build.yml` se ejecutaba en **cada Pull Request** a las ramas `stable` y `develop`, lo cual:

- ❌ Consumía muchos recursos de GitHub Actions
- ❌ Tardaba ~10-15 minutos por PR (build multi-arquitectura)
- ❌ Era innecesario - solo necesitamos validar que el código compila, no generar imágenes Docker completas

## Solución

### 1. Separación de responsabilidades

Ahora tenemos **3 workflows independientes**:

```
┌─────────────────────────────────────────────────────────────┐
│                    FLUJO CI/CD OPTIMIZADO                   │
└─────────────────────────────────────────────────────────────┘

┌──────────────┐
│ Pull Request │
└──────┬───────┘
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│ ci.yml - Validaciones rápidas (~2-5 min)                    │
├─────────────────────────────────────────────────────────────┤
│ ✓ Lint (go fmt, go vet)                                     │
│ ✓ Test (go test -race -cover)                               │
│ ✓ Build (compilación binario)                               │
│ ✓ Docker Build Check (validar Dockerfile, sin push)         │
└─────────────────────────────────────────────────────────────┘
       │
       │ merge
       ▼
┌──────────────────┐
│ Push a develop/  │
│     stable       │
└────────┬─────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│ release.yml - Semantic Release (~5-10 min)                  │
├─────────────────────────────────────────────────────────────┤
│ 1. Analizar commits convencionales                          │
│ 2. ¿Hay nueva versión?                                      │
│    └─ No → Fin                                              │
│    └─ Sí → 3. Generar CHANGELOG                             │
│            4. Crear tag Git                                 │
│            5. Build Docker multi-arch                       │
│            6. Push a GHCR                                   │
└─────────────────────────────────────────────────────────────┘
         │
         │ tag creado (v1.0.0)
         ▼
┌─────────────────────────────────────────────────────────────┐
│ docker-build.yml - Build desde tag (manual, opcional)       │
├─────────────────────────────────────────────────────────────┤
│ • Trigger: push de tag v*                                   │
│ • Uso: Reconstruir versión específica manualmente           │
│ • Normalmente no se usa (release.yml ya construye)          │
└─────────────────────────────────────────────────────────────┘
```

### 2. Detalles de cada workflow

#### `ci.yml` - Nuevo workflow de CI
**Trigger:** PRs y pushes a `develop`/`stable`

**Jobs (en paralelo):**
1. **Lint** - Formato y análisis estático
2. **Test** - Tests unitarios con coverage
3. **Build** - Compilación del binario
4. **Docker Build Check** - Validar que Dockerfile compila (linux/amd64 only, sin push)

**Tiempo:** ~2-5 minutos
**Recursos:** Mínimos (solo compilación Go y build Docker de prueba)

#### `release.yml` - Sin cambios
**Trigger:** Push a `develop`/`stable`

**Jobs (secuenciales):**
1. **semantic-release** - Analiza commits, genera versión
2. **build-and-push** - Si hay nueva versión, construye y publica Docker
3. **notify** - Resumen del proceso

**Tiempo:** ~5-10 minutos (solo si hay release)

#### `docker-build.yml` - Simplificado
**Trigger:** Push de tags `v*`

**Cambios:**
- ❌ Removido trigger de `pull_request`
- ✅ Solo ejecuta en tags
- ✅ Un solo job `build-and-push` (antes eran 2)

**Uso típico:** Reconstruir manualmente una versión. Generalmente no se usa porque `release.yml` ya construye las imágenes.

## Comparación: Antes vs Después

### Antes
```
PR → docker-build.yml
     ├─ Build multi-arch (amd64, arm64)  ~10 min
     └─ No push
     
Total: ~10-15 minutos por PR ❌
```

### Después
```
PR → ci.yml
     ├─ Lint                              ~30 seg
     ├─ Test                              ~1 min
     ├─ Build                             ~1 min
     └─ Docker Build Check (amd64)        ~2 min
     
Total: ~2-5 minutos por PR ✅
```

## Beneficios

### 🚀 Velocidad
- **Antes:** ~10-15 min por PR
- **Ahora:** ~2-5 min por PR
- **Mejora:** 2-3x más rápido

### 💰 Recursos
- **Antes:** Build multi-arch completo en cada PR
- **Ahora:** Solo validación de compilación
- **Ahorro:** ~70-80% de minutos de GitHub Actions

### ✅ Feedback
- Validación de código más rápida
- Desarrolladores reciben feedback en minutos, no en ~15 min
- Cancel-in-progress: builds antiguos se cancelan automáticamente

### 🎯 Enfoque
- CI workflow enfocado en validación de código
- Docker builds solo cuando realmente se necesitan (releases)
- Separación clara de responsabilidades

## Migraciones necesarias

### Para desarrolladores

**Nada cambia en el flujo de trabajo:**
1. Crear PR como siempre
2. Ahora verás 4 checks en lugar de 1:
   - ✅ Lint
   - ✅ Test
   - ✅ Build
   - ✅ Docker Build Check
3. Merge cuando todos pasen

### Para releases

**Flujo sin cambios:**
1. Hacer commits convencionales a `develop`
2. `release.yml` se ejecuta automáticamente
3. Si hay `feat`/`fix`/`perf`, se crea release + imagen Docker
4. Todo automático

## Validación

✅ Sintaxis YAML validada con `yamllint`
✅ Documentación actualizada en `docs/RELEASE_WORKFLOW.md`
✅ Workflows probados localmente con act (opcional)

## Referencias

- [GitHub Actions: Workflow syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

**Fecha:** 2025-11-01  
**PR:** #[número se asignará al crear]
