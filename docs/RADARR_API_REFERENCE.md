# Radarr API Reference

> **Documentación completa de la API REST de Radarr v3 para KeeperCheky**

## 📚 Índice

- [Autenticación](#autenticación)
- [Movies API](#movies-api)
- [Campos Disponibles](#campos-disponibles)
- [Queue API](#queue-api)
- [History API](#history-api)
- [Calendar API](#calendar-api)
- [Quality Profiles](#quality-profiles)
- [Tags API](#tags-api)
- [System Info](#system-info)
- [Ejemplos Prácticos](#ejemplos-prácticos)

---

## 🔐 Autenticación

### API Key Authentication

Radarr usa autenticación basada en API Key mediante un header HTTP.

**Header Requerido:**
```http
X-Api-Key: your-api-key-here
```

**Obtener API Key:**
1. Ir a Settings → General en la interfaz web de Radarr
2. Copiar el valor de "API Key"

**Ejemplo de Request:**
```bash
curl -H "X-Api-Key: your-api-key" http://localhost:7878/api/v3/movie
```

---

## 🎬 Movies API

### GET `/api/v3/movie`

Obtiene todas las películas de la biblioteca de Radarr con metadata completa.

**Parámetros Query:**

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `tmdbId` | number | Filtrar por ID de TMDb (opcional) |
| `excludeLocalCovers` | boolean | Excluir imágenes locales (opcional) |

**Response Completo:**
```json
[
  {
    "id": 1,
    "title": "The Matrix",
    "originalTitle": "The Matrix",
    "alternateTitles": [
      {
        "sourceType": "tmdb",
        "movieMetadataId": 1,
        "title": "Matrix",
        "id": 1
      }
    ],
    "secondaryYear": null,
    "secondaryYearSourceId": 0,
    "sortTitle": "matrix",
    "sizeOnDisk": 2147483648,
    "status": "released",
    "overview": "Set in the 22nd century, The Matrix tells the story of a computer hacker...",
    "inCinemas": "1999-03-24T00:00:00Z",
    "physicalRelease": "1999-09-21T00:00:00Z",
    "digitalRelease": "1999-09-21T00:00:00Z",
    "images": [
      {
        "coverType": "poster",
        "url": "https://image.tmdb.org/t/p/original/abc123.jpg",
        "remoteUrl": "https://image.tmdb.org/t/p/original/abc123.jpg"
      },
      {
        "coverType": "fanart",
        "url": "https://image.tmdb.org/t/p/original/xyz789.jpg",
        "remoteUrl": "https://image.tmdb.org/t/p/original/xyz789.jpg"
      }
    ],
    "website": "https://www.warnerbros.com/movies/matrix",
    "remotePoster": "https://image.tmdb.org/t/p/original/abc123.jpg",
    "year": 1999,
    "youTubeTrailerId": "vKQi3bBA1y8",
    "studio": "Warner Bros. Pictures",
    "path": "/movies/The Matrix (1999)",
    "qualityProfileId": 4,
    "hasFile": true,
    "movieFileId": 1,
    "monitored": true,
    "minimumAvailability": "released",
    "isAvailable": true,
    "folderName": "/movies/The Matrix (1999)",
    "runtime": 136,
    "cleanTitle": "thematrix",
    "imdbId": "tt0133093",
    "tmdbId": 603,
    "titleSlug": "the-matrix-1999",
    "certification": "R",
    "genres": ["Action", "Science Fiction"],
    "tags": [1, 2],
    "added": "2023-01-15T10:30:00Z",
    "ratings": {
      "imdb": {
        "votes": 1800000,
        "value": 8.7,
        "type": "user"
      },
      "tmdb": {
        "votes": 24000,
        "value": 8.2,
        "type": "user"
      },
      "metacritic": {
        "votes": 0,
        "value": 73,
        "type": "user"
      },
      "rottenTomatoes": {
        "votes": 0,
        "value": 88,
        "type": "user"
      }
    },
    "movieFile": {
      "movieId": 1,
      "relativePath": "The Matrix (1999).mkv",
      "path": "/movies/The Matrix (1999)/The Matrix (1999).mkv",
      "size": 2147483648,
      "dateAdded": "2023-01-15T10:45:00Z",
      "sceneName": "The.Matrix.1999.1080p.BluRay.x264-GROUP",
      "indexerFlags": 0,
      "quality": {
        "quality": {
          "id": 7,
          "name": "Bluray-1080p",
          "source": "bluray",
          "resolution": 1080,
          "modifier": "none"
        },
        "revision": {
          "version": 1,
          "real": 0,
          "isRepack": false
        }
      },
      "customFormatScore": 0,
      "customFormats": [],
      "mediaInfo": {
        "audioBitrate": 0,
        "audioChannels": 5.1,
        "audioCodec": "DTS",
        "audioLanguages": "eng",
        "audioStreamCount": 1,
        "videoBitDepth": 8,
        "videoBitrate": 0,
        "videoCodec": "h264",
        "videoFps": 23.976,
        "videoDynamicRange": "SDR",
        "videoDynamicRangeType": "",
        "resolution": "1920x1080",
        "runTime": "2:16:00",
        "scanType": "Progressive",
        "subtitles": "eng, spa"
      },
      "qualityCutoffNotMet": false,
      "languages": [
        {
          "id": 1,
          "name": "English"
        }
      ],
      "releaseGroup": "GROUP",
      "edition": "",
      "id": 1
    },
    "collection": {
      "title": "The Matrix Collection",
      "tmdbId": 2344,
      "monitored": true,
      "qualityProfileId": 4,
      "searchOnAdd": false,
      "minimumAvailability": "released",
      "images": [],
      "added": "2023-01-15T10:30:00Z",
      "id": 1
    },
    "popularity": 82.458
  }
]
```

### GET `/api/v3/movie/{id}`

Obtiene una película específica por su ID.

**Path Parameters:**
- `id` - ID de la película en Radarr

**Response:** Mismo formato que el array anterior, pero un solo objeto.

---

## 📋 Campos Disponibles

### Campos de Movie

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | number | ID único en Radarr |
| `title` | string | Título principal de la película |
| `originalTitle` | string | Título original (idioma nativo) |
| `alternateTitles` | array | Títulos alternativos de la película |
| `secondaryYear` | number | Año secundario (re-releases) |
| `sortTitle` | string | Título usado para ordenamiento |
| `sizeOnDisk` | number | Tamaño total en bytes |
| `status` | string | Estado: `tba`, `announced`, `inCinemas`, `released`, `deleted` |
| `overview` | string | Sinopsis de la película |
| `inCinemas` | datetime | Fecha de estreno en cines |
| `physicalRelease` | datetime | Fecha de lanzamiento físico (DVD/Blu-ray) |
| `digitalRelease` | datetime | Fecha de lanzamiento digital |
| `images` | array | URLs de imágenes (poster, fanart, etc.) |
| `website` | string | Sitio web oficial |
| `remotePoster` | string | URL del póster principal |
| `year` | number | Año de producción |
| `youTubeTrailerId` | string | ID del trailer de YouTube |
| `studio` | string | Estudio de producción |
| `path` | string | Ruta de la película en el filesystem |
| `qualityProfileId` | number | ID del perfil de calidad |
| `hasFile` | boolean | Indica si tiene archivo descargado |
| `movieFileId` | number | ID del archivo de película |
| `monitored` | boolean | Indica si está monitoreada |
| `minimumAvailability` | string | Disponibilidad mínima: `tba`, `announced`, `inCinemas`, `released`, `preDB` |
| `isAvailable` | boolean | Indica si está disponible según `minimumAvailability` |
| `folderName` | string | Nombre de la carpeta |
| `runtime` | number | Duración en minutos |
| `cleanTitle` | string | Título limpio (sin caracteres especiales) |
| `imdbId` | string | ID de IMDb (ej: tt0133093) |
| `tmdbId` | number | ID de The Movie Database |
| `titleSlug` | string | Slug para URLs |
| `certification` | string | Clasificación (G, PG, PG-13, R, NC-17) |
| `genres` | array[string] | Lista de géneros |
| `tags` | array[number] | IDs de tags asignados |
| `added` | datetime | Fecha en que se añadió a Radarr |
| `ratings` | object | Calificaciones de diferentes servicios |
| `movieFile` | object | Información del archivo (ver abajo) |
| `collection` | object | Colección a la que pertenece |
| `popularity` | number | Popularidad en TMDb |

### Campos de MovieFile

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `movieId` | number | ID de la película |
| `relativePath` | string | Ruta relativa del archivo |
| `path` | string | Ruta absoluta del archivo |
| `size` | number | Tamaño del archivo en bytes |
| `dateAdded` | datetime | Fecha de importación |
| `sceneName` | string | Nombre de la release scene |
| `indexerFlags` | number | Flags del indexer |
| `quality` | object | Información de calidad |
| `customFormatScore` | number | Puntuación de custom formats |
| `customFormats` | array | Custom formats aplicados |
| `mediaInfo` | object | Información técnica del archivo |
| `qualityCutoffNotMet` | boolean | Indica si no alcanza el cutoff de calidad |
| `languages` | array | Idiomas del archivo |
| `releaseGroup` | string | Grupo de release |
| `edition` | string | Edición (Extended, Director's Cut, etc.) |

### Campos de MediaInfo

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `audioBitrate` | number | Bitrate de audio |
| `audioChannels` | number | Número de canales de audio |
| `audioCodec` | string | Codec de audio (AAC, DTS, AC3, etc.) |
| `audioLanguages` | string | Idiomas de audio separados por coma |
| `audioStreamCount` | number | Número de streams de audio |
| `videoBitDepth` | number | Profundidad de color (8, 10 bits) |
| `videoBitrate` | number | Bitrate de video |
| `videoCodec` | string | Codec de video (h264, h265, AV1, etc.) |
| `videoFps` | number | Frames por segundo |
| `videoDynamicRange` | string | Rango dinámico (SDR, HDR) |
| `videoDynamicRangeType` | string | Tipo de HDR (HDR10, DolbyVision, etc.) |
| `resolution` | string | Resolución (1920x1080, 3840x2160, etc.) |
| `runTime` | string | Duración en formato HH:MM:SS |
| `scanType` | string | Tipo de escaneo (Progressive, Interlaced) |
| `subtitles` | string | Idiomas de subtítulos separados por coma |

### Campos de Ratings

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `imdb.votes` | number | Número de votos en IMDb |
| `imdb.value` | number | Calificación en IMDb (0-10) |
| `tmdb.votes` | number | Número de votos en TMDb |
| `tmdb.value` | number | Calificación en TMDb (0-10) |
| `metacritic.value` | number | Calificación en Metacritic (0-100) |
| `rottenTomatoes.value` | number | Calificación en Rotten Tomatoes (0-100) |

### Campos de Quality

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `quality.id` | number | ID de la calidad |
| `quality.name` | string | Nombre de la calidad |
| `quality.source` | string | Fuente (bluray, webdl, webrip, hdtv, etc.) |
| `quality.resolution` | number | Resolución (480, 720, 1080, 2160) |
| `quality.modifier` | string | Modificador (remux, brdisk, regional, etc.) |
| `revision.version` | number | Versión del archivo |
| `revision.real` | number | Versión "real" |
| `revision.isRepack` | boolean | Indica si es un repack |

---

## 📥 Queue API

### GET `/api/v3/queue`

Obtiene la cola de descargas actual con información de progreso.

**Parámetros Query:**

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `pageSize` | number | Número de items por página (default: 20) |
| `page` | number | Número de página (default: 1) |
| `sortKey` | string | Campo de ordenamiento |
| `sortDirection` | string | `ascending` o `descending` |
| `includeUnknownMovieItems` | boolean | Incluir items de películas desconocidas |

**Response:**
```json
{
  "page": 1,
  "pageSize": 20,
  "sortKey": "timeleft",
  "sortDirection": "ascending",
  "totalRecords": 2,
  "records": [
    {
      "id": 1,
      "movieId": 123,
      "title": "Test Movie 2024",
      "size": 1000000000,
      "sizeleft": 500000000,
      "status": "downloading",
      "trackedDownloadStatus": "ok",
      "trackedDownloadState": "downloading",
      "statusMessages": [],
      "downloadId": "abc123",
      "protocol": "torrent",
      "downloadClient": "qBittorrent",
      "indexer": "Example Indexer",
      "outputPath": "/downloads/Test Movie 2024",
      "timedOut": false,
      "estimatedCompletionTime": "2024-01-01T12:00:00Z",
      "quality": {
        "quality": {
          "id": 7,
          "name": "Bluray-1080p"
        }
      }
    }
  ]
}
```

---

## 📜 History API

### GET `/api/v3/history`

Obtiene el historial de eventos de Radarr.

**Parámetros Query:**

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `pageSize` | number | Items por página (1-100) |
| `page` | number | Número de página |
| `sortKey` | string | Campo de ordenamiento (`date`, `movie`, `quality`) |
| `sortDirection` | string | `ascending` o `descending` |
| `movieId` | number | Filtrar por película específica |
| `eventType` | number | Filtrar por tipo de evento (ver abajo) |

**Event Types:**
- `1` - grabbed (descarga iniciada)
- `3` - downloadFolderImported (archivo importado)
- `4` - downloadFailed (descarga fallida)
- `5` - movieFileDeleted (archivo eliminado)
- `6` - movieFileRenamed (archivo renombrado)
- `7` - downloadIgnored (descarga ignorada)

**Response:**
```json
{
  "page": 1,
  "pageSize": 50,
  "sortKey": "date",
  "sortDirection": "descending",
  "totalRecords": 100,
  "records": [
    {
      "id": 1,
      "movieId": 123,
      "sourceTitle": "Test.Movie.2024.1080p.BluRay.x264",
      "quality": {
        "quality": {
          "name": "Bluray-1080p"
        }
      },
      "date": "2024-01-01T10:00:00Z",
      "eventType": "grabbed",
      "downloadId": "abc123",
      "data": {
        "indexer": "Example Indexer",
        "releaseGroup": "GROUP",
        "age": "2",
        "downloadClient": "qBittorrent",
        "protocol": "torrent"
      }
    }
  ]
}
```

---

## 📅 Calendar API

### GET `/api/v3/calendar`

Obtiene películas próximas a estrenarse o lanzarse.

**Parámetros Query:**

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `start` | date | Fecha de inicio (YYYY-MM-DD) |
| `end` | date | Fecha de fin (YYYY-MM-DD) |
| `unmonitored` | boolean | Incluir películas no monitoreadas |

**Response:**
```json
[
  {
    "id": 1,
    "title": "Upcoming Movie",
    "inCinemas": "2024-01-15T00:00:00Z",
    "physicalRelease": "2024-02-01T00:00:00Z",
    "digitalRelease": "2024-01-20T00:00:00Z",
    "year": 2024,
    "hasFile": false,
    "monitored": true,
    "status": "announced",
    "overview": "Description...",
    "images": [...],
    "genres": ["Action"],
    "ratings": {...}
  }
]
```

---

## 🎨 Quality Profiles

### GET `/api/v3/qualityprofile`

Obtiene los perfiles de calidad configurados.

**Response:**
```json
[
  {
    "id": 1,
    "name": "Any",
    "upgradeAllowed": true,
    "cutoff": 20,
    "items": [
      {
        "id": 0,
        "quality": {
          "id": 0,
          "name": "Unknown",
          "source": "unknown",
          "resolution": 0
        },
        "items": [],
        "allowed": false
      }
    ],
    "minFormatScore": 0,
    "cutoffFormatScore": 0,
    "formatItems": []
  }
]
```

---

## 🏷️ Tags API

### GET `/api/v3/tag`

Obtiene todos los tags configurados.

**Response:**
```json
[
  {
    "id": 1,
    "label": "4k"
  },
  {
    "id": 2,
    "label": "classics"
  }
]
```

### POST `/api/v3/tag`

Crea un nuevo tag.

**Request Body:**
```json
{
  "label": "new-tag"
}
```

---

## 🖥️ System Info

### GET `/api/v3/system/status`

Obtiene información del sistema Radarr.

**Response:**
```json
{
  "version": "4.3.2.6858",
  "buildTime": "2024-01-01T00:00:00Z",
  "isDebug": false,
  "isProduction": true,
  "isAdmin": false,
  "isUserInteractive": false,
  "startupPath": "/app/radarr/bin",
  "appData": "/config",
  "osName": "ubuntu",
  "osVersion": "22.04",
  "isMonoRuntime": false,
  "isMono": false,
  "isLinux": true,
  "isOsx": false,
  "isWindows": false,
  "mode": "console",
  "branch": "master",
  "authentication": "forms",
  "sqliteVersion": "3.36.0",
  "urlBase": "",
  "runtimeVersion": "6.0.25",
  "runtimeName": ".NET Core",
  "packageVersion": "4.3.2.6858",
  "packageAuthor": "Team Radarr",
  "packageUpdateMechanism": "docker"
}
```

---

## 💡 Ejemplos Prácticos

### Obtener Todas las Películas con Archivos

```bash
curl -H "X-Api-Key: YOUR_KEY" http://localhost:7878/api/v3/movie | \
  jq '.[] | select(.hasFile == true)'
```

### Buscar Película por IMDb ID

```bash
curl -H "X-Api-Key: YOUR_KEY" http://localhost:7878/api/v3/movie | \
  jq '.[] | select(.imdbId == "tt0133093")'
```

### Obtener Películas con Calidad Específica

```bash
curl -H "X-Api-Key: YOUR_KEY" http://localhost:7878/api/v3/movie | \
  jq '.[] | select(.movieFile.quality.quality.name == "Bluray-1080p")'
```

### Obtener Películas Sin Archivo

```bash
curl -H "X-Api-Key: YOUR_KEY" http://localhost:7878/api/v3/movie | \
  jq '.[] | select(.hasFile == false and .monitored == true)'
```

### Filtrar por Género

```bash
curl -H "X-Api-Key: YOUR_KEY" http://localhost:7878/api/v3/movie | \
  jq '.[] | select(.genres | contains(["Action"]))'
```

### Obtener Películas Añadidas Recientemente

```bash
curl -H "X-Api-Key: YOUR_KEY" http://localhost:7878/api/v3/movie | \
  jq 'sort_by(.added) | reverse | .[0:10]'
```

### Calcular Tamaño Total de la Biblioteca

```bash
curl -H "X-Api-Key: YOUR_KEY" http://localhost:7878/api/v3/movie | \
  jq '[.[] | .sizeOnDisk] | add'
```

---

## 🔍 Búsqueda y Filtrado

### POST `/api/v3/movie/lookup`

Buscar películas en TMDb para añadir a Radarr.

**Query Parameters:**
- `term` - Término de búsqueda (puede ser título o IMDb ID con "imdb:")

**Ejemplo:**
```bash
curl -H "X-Api-Key: YOUR_KEY" \
  "http://localhost:7878/api/v3/movie/lookup?term=The%20Matrix"
```

---

## 🔄 Comandos y Acciones

### POST `/api/v3/command`

Ejecutar comandos en Radarr.

**Comandos Disponibles:**

#### RefreshMovie
```json
{
  "name": "RefreshMovie",
  "movieIds": [1, 2, 3]
}
```

#### RescanMovie
```json
{
  "name": "RescanMovie",
  "movieIds": [1]
}
```

#### MoviesSearch
```json
{
  "name": "MoviesSearch",
  "movieIds": [1, 2, 3]
}
```

#### RenameMovie
```json
{
  "name": "RenameMovie",
  "movieIds": [1],
  "files": [1, 2]
}
```

---

## 🚨 Códigos de Error

| Código | Descripción |
|--------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content (éxito sin respuesta) |
| 400 | Bad Request (parámetros inválidos) |
| 401 | Unauthorized (API key inválida) |
| 404 | Not Found (recurso no existe) |
| 405 | Method Not Allowed |
| 409 | Conflict (recurso ya existe) |
| 500 | Internal Server Error |

---

## 🔑 Notas Importantes

### 1. **Versionado de API**

Radarr usa versionado en la URL (`/api/v3/`). Siempre usar v3 para compatibilidad.

### 2. **Paginación**

Endpoints que devuelven listas grandes (queue, history) usan paginación:
- Default page size: 20
- Max page size: 100
- Siempre verificar `totalRecords` para paginación completa

### 3. **Fechas y Horas**

Todas las fechas están en formato ISO 8601 con timezone UTC:
```
2024-01-15T10:30:00Z
```

### 4. **IDs de Calidad**

Los IDs de calidad son consistentes entre instancias:
- 0: Unknown
- 1: SDTV
- 2: DVD
- 3: WEBDL-1080p
- 4: HDTV-720p
- 5: WEBDL-720p
- 6: Bluray-720p
- 7: Bluray-1080p
- 8: WEBDL-480p
- 9: HDTV-1080p
- 10: Raw-HD
- 16: HDTV-2160p
- 18: WEBDL-2160p
- 19: Bluray-2160p
- 20: Bluray-1080p Remux
- 21: Bluray-2160p Remux

### 5. **Custom Formats**

Los custom formats son específicos de cada instalación. Consultar:
```bash
GET /api/v3/customformat
```

### 6. **Rate Limiting**

Radarr no impone rate limiting estricto, pero se recomienda:
- Máximo 10 requests por segundo
- Usar caché cuando sea posible
- Implementar backoff exponencial en errores

---

## 📚 Referencias

- **Documentación Oficial:** https://radarr.video/docs/api/
- **Swagger UI:** `http://your-radarr-url/api/docs`
- **GitHub:** https://github.com/Radarr/Radarr
- **Wiki:** https://wiki.servarr.com/radarr

---

## 🆕 Changelog de la API

### v3 (Actual)
- Soporte completo para custom formats
- Mejoras en mediaInfo
- Soporte para colecciones
- Quality revision tracking

### v2 (Deprecated)
- No usar en nuevas implementaciones

### v1 (Obsoleto)
- Completamente removida

---

**Última Actualización:** 2025-11-11  
**Versión de Radarr:** 4.3.2+  
**Mantenido por:** KeeperCheky Project
