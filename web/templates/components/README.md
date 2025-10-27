# Health Components

Componentes Alpine.js reutilizables para la gestión de salud de archivos.

## 📦 Componentes Incluidos

Este directorio contiene componentes Alpine.js modulares y reutilizables:

### `health-components.html`

Librería principal que incluye:

1. **healthCard** - Cards de estadísticas con acción de filtrado
2. **healthStatusBadge** - Badges de estado en múltiples tamaños
3. **serviceStatusIndicator** - Indicadores de estado de servicio
4. **healthFilters** - Sistema de filtros avanzados
5. **bulkActions** - Gestión de selección múltiple y acciones en masa
6. **fileHealthCard** - Card individual de archivo con detalles de salud

### Helpers Globales

- `formatBytes()` - Formato de tamaño de archivo
- `formatDate()` - Formato de fecha relativa
- `getStatusIcon()` - Iconos según estado
- `getSeverityColor()` - Colores según severidad
- `getSeverityBgColor()` - Fondos según severidad

## 🚀 Uso Rápido

Los componentes están disponibles globalmente en todas las páginas:

```html
<!-- Usar un componente -->
<div x-data="healthCard('orphan_download', 12, 'warning')">
    <span x-text="icon"></span>
    <span x-text="count"></span>
    <h3 x-text="title"></h3>
</div>

<!-- Usar helper -->
<span x-text="formatBytes(8589934592)"></span>
```

## 📖 Documentación Completa

Ver [HEALTH_COMPONENTS.md](../../docs/HEALTH_COMPONENTS.md) para:
- Descripción detallada de cada componente
- Parámetros y propiedades
- Ejemplos de uso completos
- Integración con backend
- Estructura de datos

## 🎨 Demo

Visita `/health-demo` en la aplicación para ver todos los componentes en acción.

## 🔧 Desarrollo

Para modificar los componentes:

1. Edita `health-components.html`
2. Los cambios se recargan automáticamente en desarrollo
3. Prueba en `/health-demo` antes de usar en producción

## ✨ Características

- ✅ Componentes modulares y reutilizables
- ✅ State management con Alpine.js
- ✅ Estilos con Tailwind CSS
- ✅ Tema dark integrado
- ✅ Tooltips informativos
- ✅ Confirmaciones para acciones destructivas
- ✅ Error handling robusto
- ✅ Loading states automáticos
- ✅ Eventos personalizados para comunicación entre componentes
