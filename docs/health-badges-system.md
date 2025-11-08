# Sistema de Badges de Salud - KeeperCheky

## Descripción General

Los **health badges** (badges de salud) son indicadores visuales que muestran el estado de cada archivo de media en el sistema. Aparecen en los headers de las tarjetas de archivos en la vista `/files`.

## Componentes

### 1. Backend: Datos del Archivo

**Endpoint**: `/api/files`  
**Handler**: `internal/handler/files.go` → `GetFilesAPI()`  
**Modelo**: `internal/models/models.go` → `Media` struct

#### Campos Clave para Health Status

```go
type Media struct {
    // Service flags
    InRadarr      bool    // En Radarr
    InSonarr      bool    // En Sonarr
    InJellyfin    bool    // En biblioteca Jellyfin
    InQBittorrent bool    // En cliente torrent
    
    // Torrent info
    TorrentState    string  // Estado: "error", "uploading", "seeding", etc.
    IsSeeding       bool    // Seedeando activamente
    SeedRatio       float64 // Ratio de seeds
    TorrentHash     string  // Hash del torrent
    
    // Filesystem
    IsHardlink      bool    // Tiene hardlinks
    HardlinkPaths   string  // Rutas de hardlinks
    
    // Viewing
    HasBeenWatched  bool    // Reproducido en Jellyfin
    TotalPlayCount  int     // Número de reproducciones
}
```

### 2. Frontend: Componente Alpine.js

**Archivo**: `web/static/js/file-health-components.js` (global)  
**También**: `web/templates/pages/files.html` (inline, línea 1533)

#### Función Principal

```javascript
function healthStatusBadge(status, severity, size = 'md') {
    return {
        icon,         // Emoji del estado
        label,        // Texto descriptivo
        colorClasses, // Clases CSS para colores
        sizeClasses   // Clases CSS para tamaño
    };
}
```

### 3. Lógica de Detección de Estado

**Archivo**: `web/templates/pages/files.html` (línea 2197)

```javascript
getFileHealthStatus(file) {
    // Prioridad de detección (de mayor a menor gravedad)
    if (file.torrent_state === 'error') return 'critical';
    if (file.in_qbittorrent && !file.in_jellyfin) return 'orphan_download';
    if (file.is_hardlink) return 'only_hardlink';
    if (!file.has_been_watched && file.in_jellyfin) return 'missing_metadata';
    return 'ok';
}

getFileSeverity(file) {
    if (file.torrent_state === 'error' || file.torrent_state === 'missingFiles') 
        return 'critical';
    if (file.in_qbittorrent && !file.in_jellyfin) 
        return 'warning';
    if (!file.has_been_watched && file.in_jellyfin) 
        return 'warning';
    return 'ok';
}
```

## Estados de Salud

### ✅ OK (Saludable)
- **Condición**: Archivo gestionado correctamente
- **Criterios**:
  - En Jellyfin ✅
  - En Radarr o Sonarr ✅
  - No hay errores ✅
- **Color**: Verde
- **Clases**: `bg-green-900/40 border-green-600/50 text-green-300`

### ⚠️ Orphan Download (Huérfano en Descargas)
- **Condición**: `in_qbittorrent = true` Y `in_jellyfin = false`
- **Significado**: Archivo descargado pero NO importado a biblioteca
- **Acción**: Importar a Radarr/Sonarr
- **Color**: Amarillo
- **Clases**: `bg-yellow-900/40 border-yellow-600/50 text-yellow-300`

### 🔗 Only Hardlink (Solo Hardlink)
- **Condición**: `is_hardlink = true`
- **Significado**: Torrent original eliminado, solo quedan hardlinks
- **Acción**: Limpiar hardlink de downloads sin perder archivo en biblioteca
- **Color**: Azul
- **Clases**: `bg-blue-900/40 border-blue-600/50 text-blue-300`

### 🔴 Critical (Crítico)
- **Condición**: `torrent_state = "error"` O `torrent_state = "missingFiles"`
- **Significado**: Torrent con errores graves
- **Acción**: Revisar en qBittorrent o eliminar
- **Color**: Rojo
- **Clases**: `bg-red-900/40 border-red-600/50 text-red-300`

### 👁️ Missing Metadata (Sin Reproducir)
- **Condición**: `in_jellyfin = true` Y `has_been_watched = false`
- **Significado**: Nunca reproducido
- **Acción**: Considerar eliminar si no es de interés
- **Color**: Amarillo
- **Severidad**: Warning

## Renderizado en HTML

```html
<span x-data="healthStatusBadge(
        getFileHealthStatus(file), 
        getFileSeverity(file), 
        'md')" 
    :class="[colorClasses, sizeClasses]"
    class="inline-flex items-center gap-1 rounded border">
    <span x-text="icon"></span>
    <span x-text="label"></span>
</span>
```

## Indicadores Adicionales de Servicios

Además del health badge, cada archivo muestra **service status indicators** para:

- 🧲 **qBittorrent**: Estado torrent, ratio, seeding
- 🎬 **Radarr**: Gestión de películas
- 📺 **Sonarr**: Gestión de series
- 🍿 **Jellyfin**: En biblioteca, reproducido
- 🎫 **Jellyseerr**: Peticiones
- 📊 **Jellystat**: Estadísticas

**Tooltips** muestran información detallada al hacer hover.

## Troubleshooting

### Badges Aparecen Vacíos o Incorrectos

1. **Verificar sincronización**:
   - Click en botón "Sincronizar" en `/files`
   - Esto llama `/api/sync/files` que enriquece datos

2. **Verificar datos en backend**:
   ```bash
   # Inspeccionar respuesta de API
   curl http://localhost:8000/api/files?page=1&perPage=10
   ```

3. **Verificar campos en base de datos**:
   ```sql
   SELECT id, title, in_qbittorrent, in_jellyfin, 
          torrent_state, is_seeding, seed_ratio
   FROM media LIMIT 10;
   ```

4. **Verificar console del navegador**:
   - Abrir DevTools (F12)
   - Tab Console
   - Buscar errores JavaScript

### Dependencias

Los badges dependen de datos correctos de:
- ✅ qBittorrent enrichment (Issue #88)
- ✅ qBittorrent endpoint registrado (Issue #87)
- ✅ Sync ejecutado recientemente

## Arquitectura de Componentes

```
┌─────────────────────────────────────────┐
│  /files Page (Alpine.js)                │
│  ┌─────────────────────────────────┐    │
│  │ filesPage() component           │    │
│  │ - loadFiles() → /api/files      │    │
│  │ - files array                   │    │
│  │                                 │    │
│  │ ┌─────────────────────────┐    │    │
│  │ │ For each file:          │    │    │
│  │ │ ┌─────────────────────┐ │    │    │
│  │ │ │ healthStatusBadge() │ │    │    │
│  │ │ │ - icon              │ │    │    │
│  │ │ │ - label             │ │    │    │
│  │ │ │ - colorClasses ✅   │ │    │    │
│  │ │ └─────────────────────┘ │    │    │
│  │ │                         │    │    │
│  │ │ ┌─────────────────────┐ │    │    │
│  │ │ │ serviceStatusIndicator() │ │    │
│  │ │ │ - qBittorrent       │ │    │    │
│  │ │ │ - Radarr/Sonarr     │ │    │    │
│  │ │ │ - Jellyfin          │ │    │    │
│  │ │ └─────────────────────┘ │    │    │
│  │ └─────────────────────────┘    │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
            ↓ API Call
┌─────────────────────────────────────────┐
│  Backend                                │
│  /api/files (Handler)                   │
│  → Query DB                             │
│  → Return MediaFileInfo[]               │
│    - torrent_state                      │
│    - in_qbittorrent                     │
│    - in_jellyfin                        │
│    - is_seeding                         │
│    - etc.                               │
└─────────────────────────────────────────┘
```

## Historial de Cambios

### 2025-11-08: Fix colorClasses Missing
- **Issue**: #89 - Badges sin información
- **Fix**: Agregado getter `colorClasses` a `file-health-components.js`
- **Archivo**: `web/static/js/file-health-components.js`
- **Commit**: `fix(frontend): add colorClasses getter to healthStatusBadge component`

## Referencias

- **Issue Original**: carcheky/keepercheky#89
- **Dependencias**: Issues #87, #88 (qBittorrent)
- **Backend Handler**: `internal/handler/files.go`
- **Frontend Component**: `web/static/js/file-health-components.js`
- **Template**: `web/templates/pages/files.html`
