# KeeperCheky - Índice de Documentación

## 📚 Guía de Navegación

Bienvenido a la documentación del proyecto **KeeperCheky**, un gestor moderno de limpieza de bibliotecas multimedia.

---

## 🎯 Documentación Activa

### 📺 Features Implementadas

- **[MOVIES_TAB.md](MOVIES_TAB.md)** ⭐ - Pestaña de películas (implementación actual)
  - Arquitectura backend y frontend
  - Problemas conocidos y soluciones
  - Próximos pasos y roadmap
  - **NOTA:** Pestaña de Series pendiente (similar a Movies)

- **[ORGANIZED_VIEW.md](ORGANIZED_VIEW.md)** - Vista organizada de archivos por series/temporadas

- **[FILE_ACTIONS_API.md](FILE_ACTIONS_API.md)** - API de acciones sobre archivos

### 🔧 Integraciones y APIs

- **[JELLYFIN_INTEGRATION.md](JELLYFIN_INTEGRATION.md)** - Integración con Jellyfin
- **[RADARR_API.md](RADARR_API.md)** - Referencia de API de Radarr
- **[RADARR_IMPLEMENTATION_SUMMARY.md](RADARR_IMPLEMENTATION_SUMMARY.md)** - Detalles de implementación Radarr

### ⚙️ Configuración y Optimización

- **[SETTINGS_REFACTORING.md](SETTINGS_REFACTORING.md)** - Refactorización de settings
- **[PERFORMANCE_OPTIMIZATION_FILES_TAB.md](PERFORMANCE_OPTIMIZATION_FILES_TAB.md)** - Optimización de rendimiento
- **[TESTING_SETUP.md](TESTING_SETUP.md)** - Configuración de entorno de testing

### 🚀 CI/CD y Release

- **[CI_CD.md](CI_CD.md)** - Documentación de pipeline CI/CD
- **[CI_CD_CHANGES.md](CI_CD_CHANGES.md)** - Cambios recientes en CI/CD
- **[RELEASE_WORKFLOW.md](RELEASE_WORKFLOW.md)** - Workflow de releases y semantic versioning
- **[health-badges-system.md](health-badges-system.md)** - Sistema de health badges

### 📊 Planificación y Progreso

- **[PROGRESS.md](PROGRESS.md)** - Estado actual del proyecto y roadmap
- **[AGENTS_MD_ANALYSIS.md](AGENTS_MD_ANALYSIS.md)** - Análisis de AGENTS.md
- **[AGENTS_MD_INTEGRATION_SUMMARY.md](AGENTS_MD_INTEGRATION_SUMMARY.md)** - Resumen de integración AGENTS.md

---

## 📐 Templates y Especificaciones

### Diseño de UI
- **[templates/card.md](templates/card.md)** - Especificación de diseño de tarjetas para películas/series

---

## 📦 Documentación Obsoleta

La documentación histórica, propuestas antiguas y análisis previos se han movido a:
- **`obsolete/`** - Archivos antiguos y no actualizados
- **`obsolete/propuestas/`** - Propuestas de stack iniciales (stack final elegido: Go + Alpine.js)

---

## 🚀 Inicio Rápido

### Para Desarrolladores
1. Lee **[../.github/copilot-instructions.md](../.github/copilot-instructions.md)** - Guía principal de desarrollo
2. Lee **[../AGENTS.md](../AGENTS.md)** - Quick reference para agentes AI
3. Revisa **[MOVIES_TAB.md](MOVIES_TAB.md)** - Ejemplo de feature actual

### Para Colaboradores
1. Revisa **[PROGRESS.md](PROGRESS.md)** - Estado del proyecto
2. Lee **[CI_CD.md](CI_CD.md)** - Pipeline de integración continua
3. Consulta **[TESTING_SETUP.md](TESTING_SETUP.md)** - Setup de testing

---

## 📝 Convenciones de Documentación

- **Idioma:** Documentación técnica en **Inglés**, comunicación con usuarios en **Español**
- **Formato:** Markdown con estructura clara y navegable
- **Ubicación:**
  - Features y guías → `docs/`
  - Instrucciones de desarrollo → `.github/`
  - Specs de diseño → `docs/templates/`
  - Obsoleto → `docs/obsolete/`

---

## 🔄 Mantenimiento de Docs

**IMPORTANTE:** Al implementar nuevas features:
1. ✅ Crear/actualizar documentación en `docs/`
2. ✅ Mover docs obsoletos a `obsolete/`
3. ✅ Actualizar este índice (README.md)
4. ✅ Referenciar en AGENTS.md si es relevante

---

**Última actualización:** 2025-11-11  
**Stack actual:** Go 1.22+ • Fiber v2 • Alpine.js 3.x • Tailwind CSS • GORM v2
