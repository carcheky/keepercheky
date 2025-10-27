# Alpine.js Health Components - Guía de Uso

Este documento describe los componentes Alpine.js reutilizables creados para la gestión de salud de archivos en KeeperCheky.

## 📋 Tabla de Contenidos

1. [Instalación](#instalación)
2. [Componentes Disponibles](#componentes-disponibles)
3. [Helpers Compartidos](#helpers-compartidos)
4. [Ejemplos de Uso](#ejemplos-de-uso)
5. [Integración con Backend](#integración-con-backend)

## 🔧 Instalación

Los componentes están automáticamente disponibles en todas las páginas ya que se incluyen en el layout principal (`layouts/main.html`).

Para usarlos, simplemente utiliza las funciones en tu template HTML con Alpine.js:

```html
<div x-data="healthCard('ok', 45, 'ok')">
  <!-- Tu contenido -->
</div>
```

## 🧩 Componentes Disponibles

### 1. `healthCard(status, count, severity, onFilter)`

Card de estadística que muestra el estado de salud con acción de filtrado.

**Parámetros:**
- `status` (string): Tipo de estado de salud (ej: 'ok', 'orphan_download', 'duplicate')
- `count` (number): Número de archivos en esta categoría
- `severity` (string): Nivel de severidad ('ok', 'warning', 'critical', 'info')
- `onFilter` (function, opcional): Callback para filtrar por este estado

**Propiedades Computadas:**
- `icon`: Emoji del icono según el status
- `colorClasses`: Clases CSS de colores según severity
- `bgColorClasses`: Clases CSS de fondo según severity
- `title`: Título legible del estado
- `description`: Descripción del estado

**Métodos:**
- `filterByStatus()`: Ejecuta el callback onFilter si está definido

**Ejemplo:**
```html
<div x-data="healthCard('orphan_download', 12, 'warning', (status) => filterFiles(status))"
    @click="filterByStatus()"
    :class="bgColorClasses"
    class="rounded-lg shadow-lg p-6 cursor-pointer transition-all border">
    <div class="flex items-center justify-between mb-2">
        <span class="text-4xl" x-text="icon"></span>
        <span class="text-3xl font-bold" x-text="count"></span>
    </div>
    <h3 class="text-lg font-semibold mb-1" x-text="title"></h3>
    <p class="text-sm" x-text="description"></p>
</div>
```

### 2. `healthStatusBadge(status, severity, size)`

Badge reutilizable para mostrar estados de salud en diferentes tamaños.

**Parámetros:**
- `status` (string): Tipo de estado
- `severity` (string): Nivel de severidad ('ok', 'warning', 'critical', 'info')
- `size` (string): Tamaño del badge ('sm', 'md', 'lg')

**Propiedades Computadas:**
- `icon`: Emoji del icono
- `colorClasses`: Clases CSS de colores
- `sizeClasses`: Clases CSS de tamaño
- `label`: Etiqueta corta del estado

**Ejemplo:**
```html
<span x-data="healthStatusBadge('orphan_download', 'warning', 'md')" 
    :class="[colorClasses, sizeClasses]"
    class="inline-flex items-center gap-1 rounded border">
    <span x-text="icon"></span>
    <span x-text="label"></span>
</span>
```

### 3. `serviceStatusIndicator(service, isActive, details)`

Indicador de estado para servicios con tooltips informativos.

**Parámetros:**
- `service` (string): Nombre del servicio ('radarr', 'sonarr', 'jellyfin', 'qbittorrent', etc.)
- `isActive` (boolean): Si el archivo está en este servicio
- `details` (object): Detalles adicionales del servicio (ratio, state, tags, etc.)

**Propiedades:**
- `showTooltip` (boolean): Control de visibilidad del tooltip

**Propiedades Computadas:**
- `icon`: Emoji del servicio
- `label`: Nombre legible del servicio
- `colorClasses`: Clases CSS de colores según estado
- `hasDetails`: Si tiene detalles para mostrar
- `tooltipContent`: Contenido formateado del tooltip

**Ejemplo:**
```html
<div x-data="serviceStatusIndicator('qbittorrent', true, { state: 'seeding', ratio: 2.5, seeding: true })" 
    @mouseenter="showTooltip = true"
    @mouseleave="showTooltip = false"
    class="relative">
    <span :class="colorClasses" 
        class="inline-flex items-center gap-1 px-2 py-1 rounded border text-sm">
        <span x-text="icon"></span>
        <span x-text="label"></span>
    </span>
    <div x-show="showTooltip" 
        x-transition
        class="absolute z-10 bg-slate-900 border text-xs rounded-lg p-3 shadow-xl w-64">
        <div x-html="tooltipContent.split('\n').join('<br>')"></div>
    </div>
</div>
```

### 4. `healthFilters()`

Sistema de filtros avanzados para archivos.

**Estado:**
- `selectedProblem`: Filtro por tipo de problema
- `selectedService`: Filtro por servicio
- `selectedSize`: Filtro por rango de tamaño
- `selectedAge`: Filtro por antigüedad
- `searchQuery`: Búsqueda por texto

**Propiedades Computadas:**
- `activeFilterCount`: Número de filtros activos
- `hasActiveFilters`: Si hay algún filtro activo
- `problemOptions`: Opciones de problemas
- `serviceOptions`: Opciones de servicios
- `sizeOptions`: Opciones de tamaños
- `ageOptions`: Opciones de antigüedad

**Métodos:**
- `applyFilters()`: Aplica los filtros (dispara evento 'filters-changed')
- `clearFilters()`: Limpia todos los filtros
- `getFilteredFiles(files)`: Filtra archivos del lado del cliente

**Eventos:**
- `filters-changed`: Emitido cuando cambian los filtros, con datos del filtro

**Ejemplo:**
```html
<div x-data="healthFilters()" @filters-changed="handleFilterChange($event.detail)">
    <select x-model="selectedProblem" @change="applyFilters()">
        <template x-for="option in problemOptions">
            <option :value="option.value" x-text="option.label"></option>
        </template>
    </select>
    <button @click="clearFilters()">Limpiar Filtros</button>
</div>
```

### 5. `bulkActions()`

Sistema de selección múltiple y acciones en masa.

**Estado:**
- `selectedIds`: Set de IDs seleccionados
- `selectAll`: Estado del checkbox "seleccionar todos"
- `actionInProgress`: Si hay una acción en progreso

**Propiedades Computadas:**
- `selectedCount`: Número de elementos seleccionados
- `hasSelection`: Si hay elementos seleccionados

**Métodos:**
- `toggleSelect(fileId)`: Alterna selección de un archivo
- `isSelected(fileId)`: Verifica si un archivo está seleccionado
- `toggleSelectAll(files)`: Alterna selección de todos
- `clearSelection()`: Limpia todas las selecciones
- `executeBulkAction(action, options)`: Ejecuta acción en masa
- `confirmBulkAction(action)`: Muestra confirmación

**Eventos:**
- `bulk-action-complete`: Emitido cuando termina una acción bulk

**Ejemplo:**
```html
<div x-data="bulkActions()" @bulk-action-complete="reloadFiles()">
    <!-- Select all -->
    <input type="checkbox" 
        :checked="selectAll"
        @change="toggleSelectAll(files)">
    
    <!-- File items -->
    <template x-for="file in files">
        <div :class="{ 'bg-blue-900/20': isSelected(file.id) }">
            <input type="checkbox" 
                :checked="isSelected(file.id)"
                @change="toggleSelect(file.id)">
            <span x-text="file.title"></span>
        </div>
    </template>
    
    <!-- Actions -->
    <button x-show="hasSelection" 
        @click="executeBulkAction('delete')"
        :disabled="actionInProgress">
        Eliminar (<span x-text="selectedCount"></span>)
    </button>
</div>
```

### 6. `fileHealthCard(file, healthReport)`

Card individual de archivo con información detallada de salud.

**Parámetros:**
- `file` (object): Objeto con datos del archivo
- `healthReport` (object, opcional): Reporte de salud con status, issues, suggestions, actions

**Estado:**
- `detailsExpanded`: Si los detalles están expandidos
- `actionInProgress`: Si hay una acción en progreso

**Propiedades Computadas:**
- `healthStatus`: Estado de salud del archivo
- `severity`: Nivel de severidad
- `issues`: Array de problemas detectados
- `suggestions`: Array de sugerencias
- `actions`: Array de acciones disponibles
- `hasIssues`: Si tiene problemas
- `hasSuggestions`: Si tiene sugerencias

**Métodos:**
- `toggleDetails()`: Alterna expansión de detalles
- `executeAction(action)`: Ejecuta una acción en el archivo
- `getActionLabel(action)`: Obtiene etiqueta de acción
- `getActionColor(action)`: Obtiene color de acción

**Eventos:**
- `file-action-complete`: Emitido cuando termina una acción en archivo

**Ejemplo:**
```html
<div x-data="fileHealthCard(
    { id: 1, title: 'Movie.mkv', file_path: '/media/movie.mkv', size: 8589934592 },
    { 
        status: 'orphan_download', 
        severity: 'warning',
        issues: ['No encontrado en biblioteca'],
        suggestions: ['Importar a Radarr'],
        actions: ['import', 'delete']
    }
)" @file-action-complete="reloadFile()">
    <h3 x-text="file.title"></h3>
    
    <div x-show="hasIssues">
        <template x-for="issue in issues">
            <li x-text="issue"></li>
        </template>
    </div>
    
    <template x-for="action in actions">
        <button @click="executeAction(action)"
            :disabled="actionInProgress"
            :class="getActionColor(action)">
            <span x-text="getActionLabel(action)"></span>
        </button>
    </template>
</div>
```

## 🛠️ Helpers Compartidos

### `formatBytes(bytes)`

Formatea bytes a formato legible.

```javascript
formatBytes(8589934592) // "8 GB"
formatBytes(524288000)  // "500 MB"
```

### `formatDate(timestamp)`

Formatea fecha a formato relativo en español.

```javascript
formatDate(new Date(Date.now() - 3*24*60*60*1000)) // "hace 3 días"
formatDate(new Date(Date.now() - 2*60*60*1000))    // "hace 2 horas"
```

### `getStatusIcon(status)`

Retorna emoji de icono según el estado.

```javascript
getStatusIcon('ok')               // "✅"
getStatusIcon('orphan_download')  // "⚠️"
getStatusIcon('duplicate')        // "📋"
```

### `getSeverityColor(severity)`

Retorna clases CSS de colores según severidad.

```javascript
getSeverityColor('ok')       // "bg-green-900/60 border-green-700 text-green-300"
getSeverityColor('warning')  // "bg-yellow-900/60 border-yellow-700 text-yellow-300"
getSeverityColor('critical') // "bg-red-900/60 border-red-700 text-red-300"
```

### `getSeverityBgColor(severity)`

Retorna clases CSS de fondo según severidad.

```javascript
getSeverityBgColor('warning') // "bg-yellow-900/40 hover:bg-yellow-900/60"
```

## 📚 Ejemplos de Uso

### Página con Health Cards y Filtros

```html
<div x-data="{
    files: [],
    filters: healthFilters(),
    bulk: bulkActions(),
    
    async init() {
        await this.loadFiles();
    },
    
    async loadFiles() {
        const response = await fetch('/api/files/health');
        this.files = await response.json();
    },
    
    handleFilterChange(filterData) {
        // Aplicar filtros
        this.files = this.filters.getFilteredFiles(this.originalFiles);
    }
}">
    <!-- Health Statistics -->
    <div class="grid grid-cols-4 gap-4 mb-6">
        <div x-data="healthCard('ok', 45, 'ok', () => filters.selectedProblem = 'ok')"
            @click="filterByStatus()">
            <!-- Card content -->
        </div>
    </div>
    
    <!-- Filters -->
    <div x-data="filters">
        <!-- Filter controls -->
    </div>
    
    <!-- Bulk Actions Bar -->
    <div x-show="bulk.hasSelection" x-data="bulk">
        <span x-text="`${selectedCount} seleccionados`"></span>
        <button @click="executeBulkAction('delete')">Eliminar</button>
    </div>
    
    <!-- File List -->
    <template x-for="file in files">
        <div x-data="fileHealthCard(file, file.healthReport)">
            <!-- File card content -->
        </div>
    </template>
</div>
```

## 🔌 Integración con Backend

### Endpoints Necesarios

Los componentes esperan los siguientes endpoints:

1. **GET /api/files/health**
   - Retorna archivos con información de salud
   - Response: `{ files: [...], stats: {...} }`

2. **POST /api/files/bulk-action**
   - Ejecuta acción en masa
   - Body: `{ action: string, ids: number[], options: {} }`
   - Response: `{ success_count: number, failure_count: number, results: {} }`

3. **POST /api/files/:id/:action**
   - Ejecuta acción en archivo individual
   - Response: `{ success: boolean, message: string }`

### Estructura de Datos

**FileHealthReport:**
```typescript
{
  status: 'ok' | 'orphan_download' | 'only_hardlink' | 'duplicate' | ...,
  severity: 'ok' | 'warning' | 'critical' | 'info',
  issues: string[],
  suggestions: string[],
  actions: string[]  // 'import', 'delete', 'ignore', 'rescan', 'fix'
}
```

## 📝 Notas

- Todos los componentes están diseñados para trabajar con Alpine.js 3.x
- Los estilos usan Tailwind CSS con el tema dark personalizado de KeeperCheky
- Los componentes son reactivos y actualizan la UI automáticamente
- Se incluyen confirmaciones para acciones destructivas
- Los errores se manejan con mensajes claros al usuario

## 🎨 Demo

Visita `/health-demo` para ver una demostración completa de todos los componentes en acción.
