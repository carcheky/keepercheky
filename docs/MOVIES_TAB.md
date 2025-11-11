# Pestaña de Películas

## Descripción General

La pestaña de películas muestra todas las películas disponibles en las bibliotecas de Jellyfin, con una interfaz moderna tipo tarjetas (cards).

## Arquitectura Backend

### Handler: `internal/handler/movies.go`

**Responsabilidades:**
- Obtener películas desde Jellyfin via `GetLibrary()`
- Filtrar solo items de tipo "movie" (excluir series, episodes, seasons)
- Convertir datos de Jellyfin a formato MovieFile para el frontend
- Proporcionar API JSON para consumo del frontend

**Endpoints:**
- `GET /movies` - Renderiza la página HTML
- `GET /api/movies` - Devuelve JSON con las películas

**Estructuras de Datos:**

```go
type MovieFile struct {
    Path      string                 `json:"path"`
    Title     string                 `json:"title"`
    Year      int                    `json:"year,omitempty"`
    Size      int64                  `json:"size"`
    Extension string                 `json:"extension"`
    Directory string                 `json:"directory"`
    PosterURL string                 `json:"poster_url,omitempty"`
    Jellyfin  *JellyfinMovieMetadata `json:"jellyfin,omitempty"`
}

type JellyfinMovieMetadata struct {
    ID              string  `json:"id"`
    Overview        string  `json:"overview,omitempty"`
    CommunityRating float64 `json:"community_rating,omitempty"`
    PlayCount       int     `json:"play_count"`
    LastPlayed      string  `json:"last_played,omitempty"`
    DateAdded       string  `json:"date_added,omitempty"`
}
```

### Cliente Jellyfin: `internal/service/clients/jellyfin.go`

**Filtrado de Tipos:**
El cliente de Jellyfin devuelve diferentes tipos de items:
- `"Movie"` → convertido a type "movie" ✅
- `"Series"` → convertido a type "series" (filtrado)
- `"Episode"` → convertido a type "episode" (filtrado)
- `"Season"` → **PROBLEMA ACTUAL**: se convierte a "movie" por defecto

**Issue Conocido:**
Las temporadas (Seasons) de series están siendo devueltas como "movies" porque el código asume que todo lo que no es "Series" ni "Episode" es una película. Esto necesita ser corregido.

**Solución Pendiente:**
Agregar filtro explícito para `item.Type == "Season"` en `convertToMedia()`.

## Arquitectura Frontend

### Template: `web/templates/pages/movies.html`

**Componente Alpine.js:** `moviesPage()`

**Estado:**
```javascript
{
    movies: [],              // Array de películas
    libraryPaths: [],        // Rutas de bibliotecas de Jellyfin
    loading: true,           // Estado de carga
    error: null,             // Errores
    searchQuery: '',         // Búsqueda
    sortBy: 'title',         // Ordenamiento
    selectedMovies: []       // Películas seleccionadas
}
```

**Funcionalidades:**
- ✅ Búsqueda por título
- ✅ Ordenamiento (título, tamaño)
- ✅ Selección múltiple con checkboxes
- ✅ Vista de carátulas de Jellyfin
- ✅ Popup con metadatos de Jellyfin
- ✅ Diseño responsive con grid

### Diseño de Tarjetas

**Especificación (según `docs/templates/card.md`):**
```
----------------------------------
[ ] Título - Año - Tamaño - ICONOS - Carátula
----------------------------------
```

**Implementación Actual:**
- Checkbox en esquina superior izquierda de la carátula
- Carátula: `poster_url` de Jellyfin (aspect-ratio 2:3)
- Título + Año (si disponible)
- Tamaño en formato legible (KB, MB, GB)
- Icono Jellyfin (🎞️) con popup de metadatos

**Iconos Futuros Planificados:**
- 🔽 qBittorrent (estado de torrent)
- 📡 Radarr (información de Radarr)
- 📺 Sonarr (para series)

## Flujo de Datos

```
1. Usuario accede a /movies
2. Frontend carga y hace fetch a /api/movies
3. Backend llama a jellyfinClient.GetLibrary()
4. Jellyfin devuelve todos los items (movies, series, episodes, seasons)
5. convertJellyfinToMovies() filtra solo type="movie"
6. Backend devuelve JSON con películas filtradas
7. Frontend renderiza las tarjetas con Alpine.js
```

## Problemas Conocidos

### 1. Temporadas (Seasons) Incluidas como Películas

**Estado:** 🔴 PENDIENTE

**Descripción:** Las temporadas de series están siendo incluidas en el conteo de películas porque `convertToMedia()` asigna type="movie" por defecto a todo lo que no sea "Series" o "Episode".

**Ejemplo:**
- API devuelve: 19 items (12 movies + 7 seasons)
- Esperado: 12 items (solo movies)

**Solución:**
Agregar en `jellyfin.go`:
```go
func (c *JellyfinClient) convertToMedia(item *jellyfinItem) *models.Media {
    mediaType := "movie"
    if item.Type == "Series" {
        mediaType = "series"
    } else if item.Type == "Episode" {
        mediaType = "episode"
    } else if item.Type == "Season" {
        mediaType = "season"  // Nueva línea
    }
    // ...
}
```

Y en `movies.go`:
```go
if media.Type != "movie" {
    // Filtrar todo excepto movies
    continue
}
```

### 2. Metadatos Limitados de Jellyfin

**Estado:** 📝 TODO

**Descripción:** Actualmente solo mostramos metadatos básicos (ID, fechas, play_count) porque están en `models.Media`. Los metadatos enriquecidos (Overview, CommunityRating) requieren llamadas adicionales a la API de Jellyfin.

**Solución Futura:**
- Enriquecer `models.Media` con más campos de Jellyfin
- O hacer llamadas adicionales por película (impacto en performance)
- O cachear metadatos enriquecidos

### 3. Extracción de Año

**Estado:** ⚠️ MEJORABLE

**Descripción:** El año se extrae del título buscando patrón `(YYYY)`. Esto es frágil y puede fallar con títulos no estándar.

**Alternativa:** Jellyfin tiene campo `ProductionYear` que sería más confiable.

## Próximos Pasos

### Corto Plazo
- [ ] Corregir filtrado de Seasons
- [ ] Agregar campo `year` a `models.Media` desde Jellyfin
- [ ] Probar con biblioteca grande (100+ películas)
- [ ] Agregar paginación si es necesario

### Mediano Plazo
- [ ] Crear pestaña de Series (similar a Movies)
- [ ] Agregar integración con qBittorrent (icono + estado)
- [ ] Agregar integración con Radarr (icono + acciones)
- [ ] Enriquecer metadatos con llamadas adicionales a Jellyfin

### Largo Plazo
- [ ] Acciones bulk (eliminar, mover, etc.)
- [ ] Filtros avanzados (género, calidad, fecha)
- [ ] Estadísticas de visualización
- [ ] Integración con Jellyseerr (solicitudes)

## Notas para Series Tab

**RECORDATORIO:** Más tarde crearemos la pestaña de Series con arquitectura similar:
- `SeriesHandler` en `internal/handler/series.go`
- Filtrar `media.Type == "series"`
- Mostrar temporadas y episodios
- Diseño de tarjetas adaptado para series (con temporadas/episodios)

## Referencias

- Especificación de tarjetas: `docs/templates/card.md`
- Cliente Jellyfin: `internal/service/clients/jellyfin.go`
- Modelo Media: `internal/models/models.go`
- Instrucciones generales: `.github/copilot-instructions.md`

---

**Última actualización:** 2025-11-11  
**Estado:** En desarrollo activo  
**Prioridad:** Alta
