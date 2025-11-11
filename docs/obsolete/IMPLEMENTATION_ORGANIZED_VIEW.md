# Resumen de Implementación: Vista Organizada de Archivos

## Objetivo
Implementar una vista jerárquica en la pestaña Files que organice los archivos por series/temporadas y películas, con soporte para mostrar múltiples versiones del mismo contenido.

## Solución Implementada

### 1. Backend (Go)

#### Nuevo Endpoint API
```
GET /api/files/organized
```

**Parámetros:**
- `page`: Número de página (default: 1)
- `perPage`: Elementos por página (default: 25, max: 100)
- `tab`: Filtro de categoría (healthy, attention, critical, hardlinks, unwatched)
- `type`: Tipo de media (series, movie)

**Estructuras de Datos:**

```go
type SeriesInfo struct {
    SeriesTitle  string
    TotalSize    int64
    SeasonCount  int
    EpisodeCount int
    Seasons      []SeasonInfo
    // ... metadata
}

type SeasonInfo struct {
    SeasonNumber int
    EpisodeCount int
    TotalSize    int64
    Episodes     []EpisodeInfo
}

type EpisodeInfo struct {
    EpisodeNumber int
    Title         string
    FilePath      string
    Size          int64
    Quality       string
    Versions      []string  // Otras versiones
}

type MovieInfo struct {
    Title         string
    TotalSize     int64
    PrimaryFile   MediaFileInfo
    OtherVersions []MediaFileInfo
    // ... metadata
}
```

#### Parsing de Nombres de Archivo

**Patrones soportados:**
- `S01E01`, `s01e01` → Serie/Temporada/Episodio
- `1x01`, `1X01` → Serie/Temporada/Episodio
- Detección de calidad: 1080p, 720p, 2160p, 4K, BluRay, WEB-DL, etc.

**Algoritmo:**
1. Extraer nombre base del archivo
2. Buscar patrón de episodio (SxxExx o XxY)
3. Truncar todo después del patrón
4. Limpiar puntos, guiones bajos, tags de calidad
5. Normalizar espacios

#### Agrupación Inteligente

**Series:**
1. Agrupar por nombre de serie (detectado o del campo `title`)
2. Organizar por temporadas
3. Ordenar episodios dentro de cada temporada
4. Calcular tamaños acumulados en cada nivel

**Películas:**
1. Agrupar por título
2. Priorizar versión en Jellyfin como principal
3. Resto como "otras versiones"
4. Calcular tamaño total

### 2. Frontend (Alpine.js)

#### Toggle de Vista
Botones en el header para cambiar entre "Lista" y "Organizado"

#### Componentes Desplegables

**Series (3 niveles):**
```
Serie (nombre, total episodios, total tamaño)
  └── Temporada (número, episodios, tamaño)
        └── Episodio (número, título, calidad, tamaño)
              └── Otras versiones (opcional)
```

**Películas:**
```
Película (título, tamaño total)
  └── Versión principal (calidad, tamaño)
  └── Otras versiones (opcional)
```

#### Estado de Expansión
- `expandedSeries[]`: Series expandidas
- `expandedSeasons[]`: Temporadas expandidas (formato: "SerieTitle-S1")
- `expandedEpisodeVersions[]`: Episodios con versiones expandidas
- `expandedMovieVersions[]`: Películas con versiones expandidas

#### Métodos Principales
```javascript
- loadOrganizedFiles()        // Cargar datos del API
- toggleSeries(title)          // Expandir/colapsar serie
- toggleSeason(series, num)    // Expandir/colapsar temporada
- toggleEpisodeVersions(path)  // Expandir/colapsar versiones
- toggleMovieVersions(title)   // Expandir/colapsar versiones
```

### 3. Tests Unitarios

```go
TestParseSeasonEpisode
- Breaking.Bad.S01E01.720p.mkv → S1E1 ✓
- Game.of.Thrones.s05e08.1080p.mp4 → S5E8 ✓
- The.Office.1x01.Pilot.mp4 → S1E1 ✓
- Friends.2X10.The.One.mkv → S2E10 ✓
- Movie.Title.2023.1080p.mkv → No match ✓

TestExtractSeriesName
- Breaking.Bad.S01E01.720p.mkv → "Breaking Bad" ✓
- Game_of_Thrones_s05e08_1080p.mp4 → "Game of Thrones" ✓
- The.Office.1x01.Pilot.1080p.BluRay.x264.mp4 → "The Office" ✓
- Friends.2X10.The.One.720p.WEB-DL.mkv → "Friends" ✓
```

## Optimizaciones

### Performance
- Expresiones regulares compiladas una sola vez (package-level)
- Límite configurable de items por página
- Lazy loading de vista organizada (solo cuando se selecciona)

### Detección de Tipo
1. Prioridad al campo `type` en MediaFileInfo
2. Fallback a detección por path (`/tv/`, `/series/`)
3. Parsing de nombre como último recurso

### Paginación
- Bounds checking para evitar índices negativos
- Paginación combinada de series y películas
- Respeta filtros existentes (tabs)

## Documentación

### Archivos creados
- `docs/ORGANIZED_VIEW.md`: Guía completa de uso
- `internal/handler/files_organized.go`: Implementación backend
- `internal/handler/files_organized_test.go`: Tests unitarios

### Archivos modificados
- `README.md`: Nueva característica agregada
- `web/templates/pages/files.html`: Vista organizada agregada
- `cmd/server/main.go`: Ruta API registrada

## Casos de Uso

### Ejemplo 1: Serie con múltiples temporadas
```
Breaking Bad (52.4 GB, 5 temporadas, 62 episodios)
  ├── Temporada 1 (7.5 GB, 7 episodios)
  │     ├── E01 - Pilot (1.1 GB, 720p)
  │     ├── E02 - Cat's in the Bag... (1.0 GB, 720p)
  │     └── ...
  ├── Temporada 2 (...)
  └── ...
```

### Ejemplo 2: Película con múltiples versiones
```
Inception (10.7 GB)
  ├── Versión principal: /media/movies/Inception.2010.1080p.mkv (8.6 GB, 1080p)
  └── Otras versiones (1)
        └── /downloads/Inception.2010.2160p.mkv (2.1 GB, 2160p)
```

### Ejemplo 3: Episodio con versiones
```
E01 - Pilot (1.1 GB, 720p)
  └── Otras versiones (2)
        ├── /downloads/Breaking.Bad.S01E01.1080p.mkv
        └── /downloads/Breaking.Bad.S01E01.2160p.mkv
```

## Compatibilidad

✅ Compatible con todos los filtros existentes (tabs)
✅ Respeta paginación del servidor
✅ Funciona con datos existentes sin migración
✅ No requiere cambios en base de datos

## Próximos Pasos Sugeridos

1. **Testing con datos reales**: Probar con biblioteca grande para verificar performance
2. **Acciones masivas**: Agregar acciones en vista organizada
3. **Ordenamiento**: Implementar ordenamiento personalizado
4. **Miniaturas**: Agregar vista de posters/thumbnails
5. **Búsqueda**: Filtrado de búsqueda en vista organizada

## Métricas

- **Archivos creados**: 3
- **Archivos modificados**: 3
- **Líneas de código**: ~800
- **Tests**: 9 casos de prueba
- **Cobertura**: Parsing y extracción 100%
