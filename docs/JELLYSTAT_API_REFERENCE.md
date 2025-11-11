# Jellystat API Reference

> **Documentación completa de la API REST de Jellystat para KeeperCheky**

## 📚 Índice

- [Visión General](#visión-general)
- [Autenticación](#autenticación)
- [Endpoints de Configuración](#endpoints-de-configuración)
- [Endpoints de Estadísticas](#endpoints-de-estadísticas)
- [Endpoints de Bibliotecas](#endpoints-de-bibliotecas)
- [Endpoints de Usuarios](#endpoints-de-usuarios)
- [Endpoints de Actividad](#endpoints-de-actividad)
- [Endpoints de Items](#endpoints-de-items)
- [Filtros Temporales](#filtros-temporales)
- [Ejemplos Prácticos](#ejemplos-prácticos)
- [Estructuras de Datos](#estructuras-de-datos)
- [Mejores Prácticas](#mejores-prácticas)

---

## 🔍 Visión General

Jellystat es una aplicación de estadísticas y monitoreo para Jellyfin que proporciona:
- Seguimiento de actividad de usuarios
- Estadísticas de reproducción
- Métricas de bibliotecas
- Historial de reproducción
- Análisis de contenido más popular
- Monitoreo de sesiones activas

**Backend:** Node.js con Express  
**Base de Datos:** PostgreSQL  
**Repositorio:** https://github.com/CyferShepard/Jellystat

---

## 🔐 Autenticación

Todos los endpoints de Jellystat requieren autenticación mediante API Key.

### Header Requerido

```http
x-api-token: YOUR_API_KEY_HERE
```

### Obtener API Key

La API Key se genera desde la interfaz web de Jellystat:
1. Navega a Settings → API Keys
2. Genera una nueva API Key
3. Copia y guárdala de forma segura

---

## ⚙️ Endpoints de Configuración

### GET `/api/getconfig`

Obtiene la configuración pública del servidor Jellystat.

**Autenticación:** No requerida (endpoint público)

**Response:**
```json
{
  "JF_HOST": "http://jellyfin:8096",
  "APP_USER": "admin",
  "REQUIRE_LOGIN": false,
  "IS_JELLYFIN": true
}
```

**Campos:**
- `JF_HOST`: URL del servidor Jellyfin conectado
- `APP_USER`: Usuario administrador de Jellystat
- `REQUIRE_LOGIN`: Si requiere login para acceder
- `IS_JELLYFIN`: true para Jellyfin, false para Emby

---

## 📊 Endpoints de Estadísticas

### GET `/api/statistics`

Obtiene estadísticas generales de contenido agregado en un período de tiempo.

**Parámetros Query:**
- `days` (int, required): Número de días para el período de estadísticas

**Ejemplo Request:**
```http
GET /api/statistics?days=30
x-api-token: your-api-key
```

**Response:**
```json
{
  "days": 30,
  "movies": 150,
  "episodes": 500,
  "songs": 1000
}
```

**Campos:**
- `days`: Número de días analizados
- `movies`: Total de películas en el período
- `episodes`: Total de episodios en el período
- `songs`: Total de canciones en el período

**Caso de Uso:**  
Monitorear crecimiento de biblioteca, generar reportes de contenido agregado.

---

### GET `/api/stats/getViewsByLibraryType`

Obtiene el número de reproducciones agregadas por tipo de biblioteca.

**Parámetros Query:**
- `days` (int, required): Número de días para el análisis

**Ejemplo Request:**
```http
GET /api/stats/getViewsByLibraryType?days=7
x-api-token: your-api-key
```

**Response:**
```json
{
  "music": 10,
  "movie": 50,
  "episode": 100,
  "book": 5
}
```

**Campos:**
- `music`: Reproducciones de contenido musical
- `movie`: Reproducciones de películas
- `episode`: Reproducciones de episodios de TV
- `book`: Reproducciones de audiolibros

**Caso de Uso:**  
Identificar qué tipo de contenido es más consumido, optimizar almacenamiento.

---

## 📚 Endpoints de Bibliotecas

### GET `/api/stats/getLibraryStats`

Obtiene estadísticas detalladas de todas las bibliotecas.

**Parámetros Query:**
- `days` (int, required): Número de días para el análisis

**Ejemplo Request:**
```http
GET /api/stats/getLibraryStats?days=30
x-api-token: your-api-key
```

**Response:**
```json
[
  {
    "library_id": "lib1",
    "library_name": "Movies",
    "total_items": 100,
    "total_plays": 500,
    "total_minutes": 30000
  },
  {
    "library_id": "lib2",
    "library_name": "TV Shows",
    "total_items": 50,
    "total_plays": 800,
    "total_minutes": 45000
  }
]
```

**Campos:**
- `library_id`: Identificador único de la biblioteca
- `library_name`: Nombre de la biblioteca
- `total_items`: Total de items en la biblioteca
- `total_plays`: Total de reproducciones
- `total_minutes`: Total de minutos reproducidos

**Caso de Uso:**  
Comparar popularidad de bibliotecas, identificar contenido poco usado.

---

## 👥 Endpoints de Usuarios

### GET `/api/stats/getUserActivity`

Obtiene estadísticas de actividad de todos los usuarios.

**Parámetros Query:**
- `days` (int, required): Número de días para el análisis

**Ejemplo Request:**
```http
GET /api/stats/getUserActivity?days=30
x-api-token: your-api-key
```

**Response:**
```json
[
  {
    "user_id": "user1",
    "user_name": "John Doe",
    "total_plays": 50,
    "total_minutes": 3000
  },
  {
    "user_id": "user2",
    "user_name": "Jane Smith",
    "total_plays": 30,
    "total_minutes": 1500
  }
]
```

**Campos:**
- `user_id`: Identificador único del usuario (UUID de Jellyfin)
- `user_name`: Nombre del usuario
- `total_plays`: Total de reproducciones del usuario
- `total_minutes`: Total de minutos reproducidos

**Caso de Uso:**  
Identificar usuarios más activos, generar reportes de uso, detectar cuentas inactivas.

**⚠️ Nota:** Este endpoint puede devolver 404 en versiones antiguas de Jellystat.

---

## 📺 Endpoints de Actividad

### GET `/api/sessions`

Obtiene información sobre sesiones activas de reproducción.

**Response:**
```json
[
  {
    "session_id": "session123",
    "user_id": "user1",
    "user_name": "John Doe",
    "device": "Chrome",
    "now_playing_item": {
      "item_id": "movie123",
      "item_name": "The Matrix",
      "item_type": "Movie"
    },
    "playback_position": 1234567890,
    "play_method": "DirectPlay",
    "is_paused": false
  }
]
```

**Campos:**
- `session_id`: ID de la sesión
- `user_id`: ID del usuario
- `user_name`: Nombre del usuario
- `device`: Dispositivo o cliente usado
- `now_playing_item`: Item que se está reproduciendo actualmente
- `playback_position`: Posición de reproducción en ticks
- `play_method`: Método de reproducción (DirectPlay, Transcode, etc.)
- `is_paused`: Si la reproducción está pausada

**Caso de Uso:**  
Monitoreo en tiempo real de reproducciones activas, detección de transcodificación.

---

### GET `/api/history`

Obtiene el historial completo de reproducciones.

**Parámetros Query (opcionales):**
- `user_id` (string): Filtrar por usuario específico
- `item_id` (string): Filtrar por item específico
- `days` (int): Límite de días hacia atrás
- `limit` (int): Número máximo de resultados
- `offset` (int): Offset para paginación

**Ejemplo Request:**
```http
GET /api/history?days=7&limit=50
x-api-token: your-api-key
```

**Response:**
```json
[
  {
    "playback_id": "play123",
    "user_id": "user1",
    "user_name": "John Doe",
    "item_id": "movie123",
    "item_name": "The Matrix",
    "item_type": "Movie",
    "played_at": "2025-11-10T14:30:00Z",
    "play_duration": 5400,
    "play_method": "DirectPlay",
    "device": "Chrome"
  }
]
```

**Campos:**
- `playback_id`: ID único de la reproducción
- `user_id`: ID del usuario
- `user_name`: Nombre del usuario
- `item_id`: ID del item reproducido
- `item_name`: Nombre del item
- `item_type`: Tipo de item (Movie, Episode, etc.)
- `played_at`: Timestamp ISO 8601 de la reproducción
- `play_duration`: Duración de reproducción en segundos
- `play_method`: Método de reproducción
- `device`: Dispositivo usado

**Caso de Uso:**  
Análisis de patrones de visualización, auditoría de uso, reportes históricos.

---

## 🎬 Endpoints de Items

### GET `/api/items/most-watched`

Obtiene los items más reproducidos.

**Parámetros Query:**
- `days` (int, required): Número de días para el análisis
- `limit` (int, optional, default: 10): Número de items a retornar
- `type` (string, optional): Filtrar por tipo (Movie, Episode, etc.)

**Ejemplo Request:**
```http
GET /api/items/most-watched?days=30&limit=10&type=Movie
x-api-token: your-api-key
```

**Response:**
```json
[
  {
    "item_id": "movie123",
    "item_name": "The Matrix",
    "item_type": "Movie",
    "library_name": "Movies",
    "play_count": 45,
    "unique_users": 12,
    "total_minutes": 4500,
    "average_completion": 92.5
  }
]
```

**Campos:**
- `item_id`: ID del item en Jellyfin
- `item_name`: Nombre del item
- `item_type`: Tipo de contenido
- `library_name`: Biblioteca que contiene el item
- `play_count`: Total de reproducciones
- `unique_users`: Número de usuarios únicos que lo vieron
- `total_minutes`: Total de minutos reproducidos
- `average_completion`: Porcentaje promedio de finalización

**Caso de Uso:**  
Identificar contenido popular, decisiones de retención de contenido.

---

### GET `/api/items/recently-added`

Obtiene items recientemente agregados a Jellyfin.

**Parámetros Query:**
- `limit` (int, optional, default: 20): Número de items a retornar
- `type` (string, optional): Filtrar por tipo

**Response:**
```json
[
  {
    "item_id": "movie456",
    "item_name": "New Movie 2024",
    "item_type": "Movie",
    "library_name": "Movies",
    "date_added": "2025-11-10T10:00:00Z",
    "date_first_played": "2025-11-10T14:30:00Z",
    "play_count": 3
  }
]
```

**Caso de Uso:**  
Monitorear adopción de contenido nuevo, verificar sincronización con Jellyfin.

---

## ⏱️ Filtros Temporales

La mayoría de endpoints de estadísticas aceptan el parámetro `days` para filtrar por período:

| Valor | Descripción |
|-------|-------------|
| `1` | Últimas 24 horas |
| `7` | Última semana |
| `30` | Último mes |
| `90` | Últimos 3 meses |
| `365` | Último año |
| `0` | Todo el tiempo (sin filtro) |

**Ejemplo:**
```http
GET /api/stats/getUserActivity?days=7
```

---

## 💡 Ejemplos Prácticos

### Dashboard de Estadísticas Generales

```bash
# Obtener estadísticas de los últimos 30 días
curl -H "x-api-token: YOUR_API_KEY" \
  "http://jellystat:3000/api/statistics?days=30"

# Obtener vistas por tipo de biblioteca
curl -H "x-api-token: YOUR_API_KEY" \
  "http://jellystat:3000/api/stats/getViewsByLibraryType?days=30"

# Obtener actividad de usuarios
curl -H "x-api-token: YOUR_API_KEY" \
  "http://jellystat:3000/api/stats/getUserActivity?days=30"
```

### Monitoreo de Popularidad de Contenido

```bash
# Top 10 películas más vistas del último mes
curl -H "x-api-token: YOUR_API_KEY" \
  "http://jellystat:3000/api/items/most-watched?days=30&limit=10&type=Movie"

# Top usuarios más activos
curl -H "x-api-token: YOUR_API_KEY" \
  "http://jellystat:3000/api/stats/getUserActivity?days=7" \
  | jq 'sort_by(-.total_plays) | .[0:5]'
```

### Verificar Sesiones Activas

```bash
# Ver quién está reproduciendo ahora
curl -H "x-api-token: YOUR_API_KEY" \
  "http://jellystat:3000/api/sessions" \
  | jq '.[] | select(.now_playing_item != null)'
```

---

## 📋 Estructuras de Datos

### JellystatSystemInfo

```go
type JellystatSystemInfo struct {
    Version string `json:"version"`
    Status  string `json:"status"`
}
```

### JellystatStatistics

```go
type JellystatStatistics struct {
    Days     int `json:"days"`
    Movies   int `json:"movies"`
    Episodes int `json:"episodes"`
    Songs    int `json:"songs"`
    Total    int `json:"total"`
}
```

### ViewsByLibraryType

```go
type ViewsByLibraryType struct {
    Music   int `json:"music"`
    Movie   int `json:"movie"`
    Episode int `json:"episode"`
    Book    int `json:"book"`
}
```

### UserActivity

```go
type UserActivity struct {
    UserID       string `json:"user_id"`
    UserName     string `json:"user_name"`
    TotalPlays   int    `json:"total_plays"`
    TotalMinutes int    `json:"total_minutes"`
}
```

### JellystatLibraryStats

```go
type JellystatLibraryStats struct {
    LibraryID    string `json:"library_id"`
    LibraryName  string `json:"library_name"`
    TotalItems   int    `json:"total_items"`
    TotalPlays   int    `json:"total_plays"`
    TotalMinutes int    `json:"total_minutes"`
}
```

### ItemPlaybackStats

```go
type ItemPlaybackStats struct {
    ItemID            string    `json:"item_id"`
    ItemName          string    `json:"item_name"`
    ItemType          string    `json:"item_type"`
    LibraryName       string    `json:"library_name"`
    PlayCount         int       `json:"play_count"`
    UniqueUsers       int       `json:"unique_users"`
    TotalMinutes      int       `json:"total_minutes"`
    AverageCompletion float64   `json:"average_completion"`
    DateAdded         time.Time `json:"date_added,omitempty"`
    DateFirstPlayed   time.Time `json:"date_first_played,omitempty"`
}
```

### PlaybackHistory

```go
type PlaybackHistory struct {
    PlaybackID   string    `json:"playback_id"`
    UserID       string    `json:"user_id"`
    UserName     string    `json:"user_name"`
    ItemID       string    `json:"item_id"`
    ItemName     string    `json:"item_name"`
    ItemType     string    `json:"item_type"`
    PlayedAt     time.Time `json:"played_at"`
    PlayDuration int       `json:"play_duration"`
    PlayMethod   string    `json:"play_method"`
    Device       string    `json:"device"`
}
```

### ActiveSession

```go
type ActiveSession struct {
    SessionID    string `json:"session_id"`
    UserID       string `json:"user_id"`
    UserName     string `json:"user_name"`
    Device       string `json:"device"`
    NowPlayingItem *struct {
        ItemID   string `json:"item_id"`
        ItemName string `json:"item_name"`
        ItemType string `json:"item_type"`
    } `json:"now_playing_item"`
    PlaybackPosition int64  `json:"playback_position"`
    PlayMethod       string `json:"play_method"`
    IsPaused         bool   `json:"is_paused"`
}
```

---

## 🎯 Mejores Prácticas

### 1. **Manejo de Errores**

```go
func (c *JellystatClient) GetStatistics(ctx context.Context, days int) (*JellystatStatistics, error) {
    // Validar parámetros
    if days < 0 {
        return nil, fmt.Errorf("days must be non-negative")
    }
    
    // Llamar con retry
    err := c.callWithRetry(ctx, func() error {
        resp, err := c.client.R().
            SetContext(ctx).
            SetResult(&stats).
            SetQueryParam("days", fmt.Sprintf("%d", days)).
            Get("/api/statistics")
            
        if err != nil {
            return fmt.Errorf("failed to get statistics: %w", err)
        }
        
        // Manejar 404 (endpoint no disponible)
        if resp.StatusCode() == 404 {
            c.logger.Warn("Statistics endpoint not available")
            return nil
        }
        
        if resp.StatusCode() != 200 {
            return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
        }
        
        return nil
    })
    
    return &stats, err
}
```

### 2. **Cacheo de Datos**

Las estadísticas no cambian frecuentemente, considera cachear resultados:

```go
type CachedStats struct {
    data      *JellystatStatistics
    fetchedAt time.Time
    ttl       time.Duration
}

func (c *CachedStats) IsValid() bool {
    return time.Since(c.fetchedAt) < c.ttl
}
```

### 3. **Logging Estructurado**

```go
c.logger.Info("Retrieved Jellystat statistics",
    zap.Int("days", days),
    zap.Int("movies", stats.Movies),
    zap.Int("episodes", stats.Episodes),
    zap.Int("total", stats.Total),
)
```

### 4. **Timeouts Apropiados**

```go
client := resty.New()
client.SetTimeout(30 * time.Second)
client.SetRetryCount(3)
client.SetRetryWaitTime(2 * time.Second)
```

### 5. **Paginación para Historial**

```go
// Iterar sobre historial en lotes
offset := 0
limit := 100

for {
    history, err := client.GetHistory(ctx, HistoryParams{
        Days:   30,
        Limit:  limit,
        Offset: offset,
    })
    
    if err != nil || len(history) == 0 {
        break
    }
    
    // Procesar lote
    processHistory(history)
    
    offset += limit
}
```

### 6. **Validación de Conexión**

Siempre verificar conexión antes de usar el cliente:

```go
if err := jellystatClient.TestConnection(ctx); err != nil {
    return fmt.Errorf("Jellystat not available: %w", err)
}
```

---

## 🔔 Limitaciones Conocidas

1. **API No Completamente Documentada**  
   Jellystat está en desarrollo activo. Algunos endpoints pueden no estar documentados oficialmente.

2. **Versiones**  
   Endpoints como `getUserActivity` y `getLibraryStats` pueden devolver 404 en versiones antiguas.

3. **Rate Limiting**  
   No hay rate limiting oficial documentado, pero se recomienda no exceder 10 req/s.

4. **Datos Históricos**  
   Si usas el plugin Jellyfin Playback Reporting, los datos solo están disponibles después de su instalación.

5. **Tiempo de Procesamiento**  
   Estadísticas con períodos largos (`days=365`) pueden tardar varios segundos.

---

## 📚 Referencias

- **Repositorio Oficial:** https://github.com/CyferShepard/Jellystat
- **Docker Hub:** https://hub.docker.com/r/cyfershepard/jellystat
- **Documentación Jellyfin:** https://jellyfin.org/docs/
- **Playback Reporting Plugin:** https://github.com/jellyfin/jellyfin-plugin-playbackreporting

---

## 🔄 Cambios Recientes

### v1.1.6+
- Refactorización de endpoints `/statistics` a GET (antes POST)
- Nuevo endpoint `/stats/getViewsByLibraryType`
- Mejoras en performance de consultas de historial

---

**Última Actualización:** 2025-11-11  
**Versión de Jellystat:** 1.1.6+  
**Mantenido por:** KeeperCheky Project
