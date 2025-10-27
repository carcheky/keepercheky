# Ejemplo de Integración de Health Components en Files Page

Este documento muestra cómo integrar los componentes de salud en la página de archivos existente.

## Modificaciones Sugeridas para `files.html`

### 1. Añadir Health Statistics Cards al Header

Reemplazar o complementar las stats cards existentes:

```html
<!-- Antes del tab de filtros, añadir health stats -->
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
    <!-- Archivos OK -->
    <div x-data="healthCard('ok', files.filter(f => !f.healthIssues).length, 'ok', () => filterByHealth('ok'))"
        @click="filterByStatus()"
        :class="bgColorClasses"
        class="rounded-lg shadow-lg p-4 cursor-pointer transition-all border">
        <div class="flex items-center justify-between mb-2">
            <span class="text-3xl" x-text="icon"></span>
            <span class="text-2xl font-bold text-dark-text" x-text="count"></span>
        </div>
        <h3 class="text-sm font-semibold text-dark-text mb-1" x-text="title"></h3>
        <p class="text-xs text-dark-muted" x-text="description"></p>
    </div>
    
    <!-- Descargas Huérfanas -->
    <div x-data="healthCard('orphan_download', 
        files.filter(f => f.in_qbittorrent && !f.in_jellyfin).length, 
        'warning', 
        () => filterByHealth('orphan_download'))"
        @click="filterByStatus()"
        :class="bgColorClasses"
        class="rounded-lg shadow-lg p-4 cursor-pointer transition-all border">
        <div class="flex items-center justify-between mb-2">
            <span class="text-3xl" x-text="icon"></span>
            <span class="text-2xl font-bold text-dark-text" x-text="count"></span>
        </div>
        <h3 class="text-sm font-semibold text-dark-text mb-1" x-text="title"></h3>
        <p class="text-xs text-dark-muted" x-text="description"></p>
    </div>
    
    <!-- Solo Hardlinks -->
    <div x-data="healthCard('only_hardlink', 
        files.filter(f => f.is_hardlink && !f.in_qbittorrent).length, 
        'info', 
        () => filterByHealth('only_hardlink'))"
        @click="filterByStatus()"
        :class="bgColorClasses"
        class="rounded-lg shadow-lg p-4 cursor-pointer transition-all border">
        <div class="flex items-center justify-between mb-2">
            <span class="text-3xl" x-text="icon"></span>
            <span class="text-2xl font-bold text-dark-text" x-text="count"></span>
        </div>
        <h3 class="text-sm font-semibold text-dark-text mb-1" x-text="title"></h3>
        <p class="text-xs text-dark-muted" x-text="description"></p>
    </div>
    
    <!-- Duplicados -->
    <div x-data="healthCard('duplicate', 
        0, 
        'critical', 
        () => filterByHealth('duplicate'))"
        @click="filterByStatus()"
        :class="bgColorClasses"
        class="rounded-lg shadow-lg p-4 cursor-pointer transition-all border">
        <div class="flex items-center justify-between mb-2">
            <span class="text-3xl" x-text="icon"></span>
            <span class="text-2xl font-bold text-dark-text" x-text="count"></span>
        </div>
        <h3 class="text-sm font-semibold text-dark-text mb-1" x-text="title"></h3>
        <p class="text-xs text-dark-muted" x-text="description"></p>
    </div>
    
    <!-- Sin Metadatos -->
    <div x-data="healthCard('missing_metadata', 
        files.filter(f => !f.title || f.title === '').length, 
        'warning', 
        () => filterByHealth('missing_metadata'))"
        @click="filterByStatus()"
        :class="bgColorClasses"
        class="rounded-lg shadow-lg p-4 cursor-pointer transition-all border">
        <div class="flex items-center justify-between mb-2">
            <span class="text-3xl" x-text="icon"></span>
            <span class="text-2xl font-bold text-dark-text" x-text="count"></span>
        </div>
        <h3 class="text-sm font-semibold text-dark-text mb-1" x-text="title"></h3>
        <p class="text-xs text-dark-muted" x-text="description"></p>
    </div>
</div>
```

### 2. Añadir Filtros Avanzados

Añadir panel de filtros colapsable después de las tabs:

```html
<!-- Advanced Filters Panel -->
<div x-data="{ filtersOpen: false, filters: healthFilters() }" 
    @filters-changed="applyHealthFilters($event.detail)"
    class="mb-6">
    
    <!-- Toggle Button -->
    <button @click="filtersOpen = !filtersOpen"
        class="flex items-center gap-2 px-4 py-2 bg-dark-surface border border-dark-border rounded-lg hover:bg-slate-700/50 transition-colors">
        <svg class="w-4 h-4 transition-transform" :class="{ 'rotate-90': filtersOpen }" 
            fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
        </svg>
        <span class="text-dark-text font-medium">Filtros Avanzados</span>
        <span x-show="filters.hasActiveFilters" 
            class="px-2 py-0.5 bg-blue-900/40 border border-blue-600/50 text-blue-300 rounded text-xs">
            <span x-text="filters.activeFilterCount"></span> activo<span x-text="filters.activeFilterCount > 1 ? 's' : ''"></span>
        </span>
    </button>
    
    <!-- Filter Panel -->
    <div x-show="filtersOpen" 
        x-transition
        class="mt-2 bg-dark-surface border border-dark-border rounded-lg p-6">
        
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
            <!-- Problem Filter -->
            <div>
                <label class="block text-sm font-medium text-dark-text mb-2">Tipo de Problema</label>
                <select x-model="filters.selectedProblem" @change="filters.applyFilters()"
                    class="w-full px-3 py-2 bg-dark-bg border border-dark-border rounded-lg text-dark-text">
                    <template x-for="option in filters.problemOptions">
                        <option :value="option.value" x-text="option.label"></option>
                    </template>
                </select>
            </div>
            
            <!-- Service Filter -->
            <div>
                <label class="block text-sm font-medium text-dark-text mb-2">Servicio</label>
                <select x-model="filters.selectedService" @change="filters.applyFilters()"
                    class="w-full px-3 py-2 bg-dark-bg border border-dark-border rounded-lg text-dark-text">
                    <template x-for="option in filters.serviceOptions">
                        <option :value="option.value" x-text="option.label"></option>
                    </template>
                </select>
            </div>
            
            <!-- Size Filter -->
            <div>
                <label class="block text-sm font-medium text-dark-text mb-2">Tamaño</label>
                <select x-model="filters.selectedSize" @change="filters.applyFilters()"
                    class="w-full px-3 py-2 bg-dark-bg border border-dark-border rounded-lg text-dark-text">
                    <template x-for="option in filters.sizeOptions">
                        <option :value="option.value" x-text="option.label"></option>
                    </template>
                </select>
            </div>
            
            <!-- Age Filter -->
            <div>
                <label class="block text-sm font-medium text-dark-text mb-2">Antigüedad</label>
                <select x-model="filters.selectedAge" @change="filters.applyFilters()"
                    class="w-full px-3 py-2 bg-dark-bg border border-dark-border rounded-lg text-dark-text">
                    <template x-for="option in filters.ageOptions">
                        <option :value="option.value" x-text="option.label"></option>
                    </template>
                </select>
            </div>
        </div>
        
        <!-- Search and Actions -->
        <div class="flex items-center gap-4">
            <input type="text" x-model="filters.searchQuery" @input="filters.applyFilters()" 
                placeholder="Buscar por título o ruta..."
                class="flex-1 px-4 py-2 bg-dark-bg border border-dark-border rounded-lg text-dark-text">
            
            <button @click="filters.clearFilters()"
                :class="filters.hasActiveFilters ? 'bg-gray-600 hover:bg-gray-700' : 'bg-gray-800 cursor-not-allowed opacity-50'"
                class="px-4 py-2 text-white rounded-lg transition-colors">
                Limpiar Filtros
            </button>
        </div>
    </div>
</div>
```

### 3. Añadir Bulk Actions Bar

Añadir barra de acciones antes de la tabla/lista de archivos:

```html
<!-- Bulk Actions Component -->
<div x-data="bulkActions()" 
    @bulk-action-complete="loadFiles()"
    class="mb-4">
    
    <!-- Bulk Actions Bar (shown when items selected) -->
    <div x-show="hasSelection" 
        x-transition
        class="p-4 bg-blue-900/20 border border-blue-600/50 rounded-lg flex items-center justify-between mb-4">
        
        <div class="flex items-center gap-4">
            <span class="text-blue-300 font-semibold">
                <span x-text="selectedCount"></span> archivo<span x-text="selectedCount > 1 ? 's' : ''"></span> seleccionado<span x-text="selectedCount > 1 ? 's' : ''"></span>
            </span>
            <button @click="clearSelection()"
                class="px-3 py-1 text-sm bg-gray-600 hover:bg-gray-700 text-white rounded transition-colors">
                Deseleccionar Todos
            </button>
        </div>
        
        <div class="flex gap-2">
            <button @click="executeBulkAction('import')"
                :disabled="actionInProgress"
                class="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded transition-colors">
                <span x-show="!actionInProgress">📚 Importar a Biblioteca</span>
                <span x-show="actionInProgress">Procesando...</span>
            </button>
            
            <button @click="executeBulkAction('delete')"
                :disabled="actionInProgress"
                class="px-4 py-2 bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white rounded transition-colors">
                <span x-show="!actionInProgress">🗑️ Eliminar</span>
                <span x-show="actionInProgress">Procesando...</span>
            </button>
        </div>
    </div>
    
    <!-- Add checkbox to table header -->
    <!-- In the table header, add: -->
    <th scope="col" class="px-6 py-3">
        <input type="checkbox" 
            :checked="selectAll"
            @change="toggleSelectAll(filteredFiles)"
            class="w-5 h-5">
    </th>
    
    <!-- Add checkbox to each row -->
    <!-- In each table row, add: -->
    <td class="px-6 py-4">
        <input type="checkbox" 
            :checked="isSelected(file.id)"
            @change="toggleSelect(file.id)"
            class="w-5 h-5">
    </td>
</div>
```

### 4. Usar Service Status Indicators en la Tabla

Reemplazar los badges actuales con los componentes:

```html
<!-- En lugar de los badges existentes, usar: -->
<div class="flex gap-2 text-xs mt-2 flex-wrap">
    <!-- Service indicators -->
    <template x-if="file.in_radarr">
        <span x-data="serviceStatusIndicator('radarr', true, {})" 
            :class="colorClasses"
            class="px-2 py-0.5 rounded border inline-flex items-center gap-1">
            <span x-text="icon"></span>
            <span x-text="label"></span>
        </span>
    </template>
    
    <template x-if="file.in_sonarr">
        <span x-data="serviceStatusIndicator('sonarr', true, {})" 
            :class="colorClasses"
            class="px-2 py-0.5 rounded border inline-flex items-center gap-1">
            <span x-text="icon"></span>
            <span x-text="label"></span>
        </span>
    </template>
    
    <template x-if="file.in_jellyfin">
        <span x-data="serviceStatusIndicator('jellyfin', true, {})" 
            :class="colorClasses"
            class="px-2 py-0.5 rounded border inline-flex items-center gap-1">
            <span x-text="icon"></span>
            <span x-text="label"></span>
        </span>
    </template>
    
    <template x-if="file.in_qbittorrent">
        <div x-data="serviceStatusIndicator('qbittorrent', true, {
            state: file.torrent_state,
            ratio: file.seed_ratio,
            category: file.torrent_category,
            tags: file.torrent_tags,
            seeding: file.is_seeding
        })" 
        @mouseenter="showTooltip = true"
        @mouseleave="showTooltip = false"
        class="relative inline-block">
            <span :class="colorClasses"
                class="px-2 py-0.5 rounded border inline-flex items-center gap-1 cursor-help">
                <span x-text="icon"></span>
                <span x-text="label"></span>
            </span>
            
            <!-- Tooltip -->
            <div x-show="showTooltip" 
                class="absolute z-50 bg-slate-900 border border-indigo-500/50 text-slate-200 text-xs rounded-lg p-3 shadow-xl w-64 bottom-full mb-2 left-0"
                style="display: none;">
                <div class="space-y-1">
                    <div x-show="details.state">
                        <span class="text-indigo-400">Estado:</span> <span x-text="details.state"></span>
                    </div>
                    <div x-show="details.category">
                        <span class="text-indigo-400">Categoría:</span> <span x-text="details.category"></span>
                    </div>
                    <div x-show="details.ratio !== undefined">
                        <span class="text-indigo-400">Ratio:</span> <span x-text="details.ratio?.toFixed(2)"></span>
                    </div>
                    <div x-show="details.tags">
                        <span class="text-indigo-400">Tags:</span> <span x-text="details.tags"></span>
                    </div>
                    <div x-show="details.seeding">
                        <span class="text-green-400">⬆️ Seedeando</span>
                    </div>
                </div>
            </div>
        </div>
    </template>
</div>
```

### 5. Modificaciones en el Script de filesPage()

```javascript
function filesPage() {
    return {
        files: [],
        lastSync: '',
        selectedTab: 'all',
        syncing: false,
        syncMessage: '',
        syncMessageType: 'info',
        syncMessageTimeout: null,
        
        // Nuevo: Health filter
        healthFilter: null,
        
        async init() {
            await this.loadFiles();
        },
        
        // Nuevo: Filtrar por health status
        filterByHealth(status) {
            this.healthFilter = status;
            // Aplicar filtro a la vista actual
            // Esto se puede combinar con los filtros existentes
        },
        
        // Nuevo: Aplicar filtros de health
        applyHealthFilters(filterData) {
            // Implementar lógica de filtrado
            console.log('Applying health filters:', filterData);
            // Puedes filtrar client-side o hacer una nueva llamada al API
        },
        
        // Resto de métodos existentes...
        
        // Modificar filteredFiles para incluir health filter
        get filteredFiles() {
            let filtered = this.files;
            
            // Aplicar health filter si está activo
            if (this.healthFilter) {
                filtered = filtered.filter(f => {
                    switch (this.healthFilter) {
                        case 'ok':
                            return f.in_jellyfin && (f.in_radarr || f.in_sonarr);
                        case 'orphan_download':
                            return f.in_qbittorrent && !f.in_jellyfin;
                        case 'only_hardlink':
                            return f.is_hardlink && !f.in_qbittorrent;
                        case 'missing_metadata':
                            return !f.title || f.title === '';
                        default:
                            return true;
                    }
                });
            }
            
            // Aplicar filtro de tab existente
            if (this.selectedTab === 'all') {
                return filtered;
            } else if (this.selectedTab === 'movies') {
                return filtered.filter(f => f.type === 'movie');
            } else if (this.selectedTab === 'series') {
                return filtered.filter(f => f.type === 'series' || f.type === 'episode');
            } else if (this.selectedTab === 'orphan') {
                return filtered.filter(f => !f.type || (f.type !== 'movie' && f.type !== 'series' && f.type !== 'episode'));
            }
            return filtered;
        },
        
        // Resto de métodos...
    }
}
```

## Notas de Implementación

1. **Los componentes ya están disponibles**: No necesitas importarlos, están cargados globalmente desde `layouts/main.html`

2. **Reutiliza helpers existentes**: La función `formatBytes` ya existe en la página, los componentes la hacen disponible globalmente

3. **Integración gradual**: Puedes implementar los componentes de uno en uno sin romper la funcionalidad existente

4. **Backend pendiente**: Los endpoints de `/api/files/bulk-action` y `/api/files/{id}/{action}` necesitan ser implementados en el backend

5. **Cálculo de health status**: Puedes calcular el health status client-side basándote en las propiedades existentes del archivo, o añadirlo en el backend

## Próximos Pasos

1. Integrar health cards en el header de files.html
2. Añadir panel de filtros avanzados
3. Implementar bulk actions con checkboxes
4. Reemplazar badges con service indicators
5. Implementar endpoints de backend para acciones
6. Añadir cálculo de health status en backend o frontend
7. Testear con datos reales
