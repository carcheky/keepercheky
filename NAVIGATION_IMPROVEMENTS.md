# 🧭 Files List Navigation Improvements

## Summary of Changes

This document describes the improvements made to the Files list navigation in response to issue #92.

## ✅ Improvements Implemented

### 1. Enhanced Pagination Controls

**Before:**
- Basic "Anterior/Siguiente" buttons
- Page number displayed without context
- No way to jump to specific pages
- Items per page buried in a small dropdown

**After:**
- ✅ Full pagination info: "Mostrando 1-25 de 77 archivos"
- ✅ First/Last page buttons (⏮️ ⏭️)
- ✅ Numbered page buttons with ellipsis for large datasets
- ✅ Quick jump to specific page (when >10 pages)
- ✅ Prominent items per page selector (10, 25, 50, 100)
- ✅ Better disabled states for navigation buttons

**Code Example:**
```html
<!-- New pagination controls -->
<div class="bg-dark-surface border border-dark-border rounded-lg p-4">
    <!-- Info and items per page -->
    <div class="flex items-center justify-between mb-4">
        <div>Mostrando 1 - 25 de 77 archivo(s)</div>
        <select x-model.number="itemsPerPage">
            <option value="10">10</option>
            <option value="25">25</option>
            <option value="50">50</option>
            <option value="100">100</option>
        </select>
    </div>
    
    <!-- Navigation: First | Previous | 1 2 3 ... 8 | Next | Last -->
    <div class="flex items-center justify-center gap-1">
        <!-- Numbered pages with smart ellipsis -->
    </div>
    
    <!-- Quick jump (for large datasets) -->
    <div x-show="totalPages > 10">
        <input type="number" placeholder="Ir a página...">
        <button>Ir</button>
    </div>
</div>
```

### 2. Always-Visible Search and Filters

**Before:**
- Filters collapsed by default behind a toggle button
- Search hidden in collapsed section
- Had to click to reveal filters

**After:**
- ✅ Search box always visible at the top
- ✅ All filters expanded and immediately accessible
- ✅ Type filter dropdown (Todos / Solo películas / Solo series)
- ✅ Service filter dropdown (Todos / qBittorrent / Radarr / Sonarr / Jellyfin / Huérfanos)
- ✅ Active filter count indicator
- ✅ Clear all filters button

**Code Example:**
```html
<!-- Always visible search -->
<input type="text" 
    x-model="searchQuery" 
    @input.debounce.300ms="handleSearchChange()" 
    placeholder="Buscar por título, ruta o hash de torrent (mínimo 3 caracteres)...">

<!-- Type filter -->
<select x-model="filterType" @change="handleFilterChange()">
    <option value="">Todos (películas y series)</option>
    <option value="movie">Solo películas</option>
    <option value="series">Solo series</option>
</select>

<!-- Service filter -->
<select x-model="filterService" @change="handleFilterChange()">
    <option value="">Todos los servicios</option>
    <option value="qbittorrent">Solo en qBittorrent</option>
    <option value="radarr">Solo en Radarr</option>
    <option value="sonarr">Solo en Sonarr</option>
    <option value="jellyfin">Solo en Jellyfin</option>
    <option value="orphan">Huérfanos (sin gestionar)</option>
</select>
```

### 3. Backend Filter Support

**New Query Parameters:**
- `type`: Filter by content type (`movie`, `series`)
- `service`: Filter by service presence (`qbittorrent`, `radarr`, `sonarr`, `jellyfin`, `orphan`)
- `search`: Text search (minimum 3 characters)

**Code Example:**
```go
// Type filter
if typeFilter != "" {
    switch typeFilter {
    case "movie":
        query = query.Where("type = ?", "movie")
    case "series":
        query = query.Where("type IN (?)", []string{"series", "episode"})
    }
}

// Service filter
if serviceFilter != "" {
    switch serviceFilter {
    case "qbittorrent":
        query = query.Where("in_q_bittorrent = ?", true)
    case "orphan":
        query = query.Where("in_q_bittorrent = ? AND in_radarr = ? AND in_sonarr = ? AND in_jellyfin = ?",
            true, false, false, false)
    }
}
```

### 4. Consistent Filtering Across Views

**Both list view and organized view now support:**
- ✅ Search by title, path, or hash
- ✅ Type filtering (movies, series)
- ✅ Service filtering
- ✅ Tab filtering (Healthy, Attention, Critical, etc.)
- ✅ Sorting options

## 🎯 User Experience Improvements

### Scenario 1: Finding Orphan Downloads
**Before:**
1. Click "Atención" tab → Shows 0 (bug)
2. Try to search → Must open filters first
3. Can't filter by service
4. Manual review of multiple pages

**After:**
1. Click "Atención" tab → Shows correct count
2. Search box already visible → Type to filter
3. Select "Huérfanos" from service filter
4. Navigate with numbered pages or quick jump

### Scenario 2: Reviewing Large Libraries
**Before:**
- See "Página 1" with no context
- Click "Siguiente" multiple times
- No idea how many pages remain
- Can't jump to end

**After:**
- See "Mostrando 1-25 de 250 archivos"
- See page numbers: 1 2 3 ... 10
- Can jump to last page
- Can increase items per page to 100

### Scenario 3: Mixed Content Management
**Before:**
- All content mixed together
- No way to separate movies and series
- Hard to review specific content types

**After:**
- Filter by "Solo películas" or "Solo series"
- Combine with service filters
- Active filter count shows what's applied
- Easy reset with "Limpiar filtros"

## 📊 Technical Details

### Frontend Changes
- **File:** `web/templates/pages/files.html`
- **Lines Changed:** ~200 additions/modifications
- **New Features:**
  - `filterType` and `filterService` reactive properties
  - `hasActiveFilters` computed property
  - `activeFilterCount` computed property
  - `goToPage()` method for quick navigation
  - `handleFilterChange()` for filter updates
  - Improved pagination UI with numbered pages

### Backend Changes
- **Files:** 
  - `internal/handler/files.go` - List view API
  - `internal/handler/files_organized.go` - Organized view API
- **New Query Parameters:**
  - `type` - Content type filter
  - `service` - Service presence filter
- **Improved:**
  - Search validation (min 3 characters)
  - Consistent filtering across both views
  - Better query performance with indexed fields

## ✅ Testing

All existing tests pass:
```bash
$ go test -v ./internal/handler/...
PASS
ok      github.com/carcheky/keepercheky/internal/handler    0.092s
```

## 🚀 Future Enhancements (Not in Scope)

- Client-side caching for better performance
- Saved filter presets
- Export filtered results to CSV
- Virtual scrolling for very large lists
- Keyboard shortcuts (j/k navigation, / for search)

## 📝 API Examples

### List View with Filters
```
GET /api/files?page=1&perPage=25&type=movie&service=qbittorrent&search=inception
```

### Organized View with Filters
```
GET /api/files/organized?page=1&perPage=50&type=series&service=jellyfin&tab=unwatched
```

## 🎨 UI Mockup

```
┌─────────────────────────────────────────────────────────────┐
│ 🏥 Salud del Almacenamiento                    [Sync] [View]│
├─────────────────────────────────────────────────────────────┤
│ [✅ 45] [⚠️ 12] [🔗 8] [💀 3] [👁️ 22]                      │
├─────────────────────────────────────────────────────────────┤
│ 🔍 Buscar: [________________________] [Limpiar]             │
│                                                              │
│ 🎬 Tipo: [Todos ▼]  🔧 Servicio: [Todos ▼]                 │
│ 📊 Ordenar: [Ruta ▼]  🔄 Orden: [Asc ▼]                    │
│                                                              │
│ 🔵 2 filtros activos           [🔄 Limpiar filtros]        │
├─────────────────────────────────────────────────────────────┤
│ [OK (45)] [Atención (12)] [Críticos (3)] [Hardlinks (8)]   │
├─────────────────────────────────────────────────────────────┤
│ Mostrando 1-25 de 77 archivos      Mostrar: [25 ▼] por pág│
├─────────────────────────────────────────────────────────────┤
│ 📄 Movie Title 1                              ✅ OK   45GB │
│ 📄 Series S01E01                              ⚠️ Atención  │
│ ...                                                          │
├─────────────────────────────────────────────────────────────┤
│         [⏮️] [◀️ Anterior] 1 [2] 3 ... 8 [Siguiente ▶️] [⏭️]│
│                                                              │
│         Ir a página: [__] [Ir]                              │
└─────────────────────────────────────────────────────────────┘
```

## 📚 Documentation References

- Issue: #92 - 🧭 [P2] Navegación vista lista confusa
- PR: [TBD]
- Related: #87, #88 (Tab filtering bugs)
