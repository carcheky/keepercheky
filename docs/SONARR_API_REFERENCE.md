# Sonarr API Reference

> **Documentación completa de la API REST de Sonarr v3 para KeeperCheky**

## 📚 Índice

- [Autenticación](#autenticación)
- [Series API](#series-api)
- [Episodes API](#episodes-api)
- [Episode Files API](#episode-files-api)
- [Campos de Series](#campos-de-series)
- [Campos de Episodios](#campos-de-episodios)
- [Campos de Temporadas](#campos-de-temporadas)
- [Queue API](#queue-api)
- [History API](#history-api)
- [Calendar API](#calendar-api)
- [Quality Profiles API](#quality-profiles-api)
- [Tags API](#tags-api)
- [System Info API](#system-info-api)
- [Ejemplos Prácticos](#ejemplos-prácticos)

---

## 🔐 Autenticación

### Headers Requeridos

Todas las peticiones a la API de Sonarr requieren autenticación mediante API key:

```http
X-Api-Key: {apiKey}
```

**Ejemplo:**
```bash
curl -H "X-Api-Key: abc123def456" http://localhost:8989/api/v3/series
```

---

## 📺 Series API

### GET `/api/v3/series`

Obtiene todas las series monitoreadas en Sonarr.

**Response:**
```json
[
  {
    "id": 1,
    "title": "Breaking Bad",
    "sortTitle": "breaking bad",
    "status": "ended",
    "ended": true,
    "overview": "A high school chemistry teacher turned methamphetamine producer...",
    "network": "AMC",
    "airTime": "21:00",
    "images": [
      {
        "coverType": "poster",
        "url": "/sonarr/MediaCover/1/poster.jpg",
        "remoteUrl": "https://artworks.thetvdb.com/..."
      },
      {
        "coverType": "banner",
        "url": "/sonarr/MediaCover/1/banner.jpg",
        "remoteUrl": "https://artworks.thetvdb.com/..."
      }
    ],
    "seasons": [
      {
        "seasonNumber": 1,
        "monitored": true,
        "statistics": {
          "previousAiring": "2008-03-09T02:00:00Z",
          "episodeFileCount": 7,
          "episodeCount": 7,
          "totalEpisodeCount": 7,
          "sizeOnDisk": 4294967296,
          "percentOfEpisodes": 100.0
        }
      }
    ],
    "year": 2008,
    "path": "/tv/Breaking Bad",
    "qualityProfileId": 1,
    "languageProfileId": 1,
    "seasonFolder": true,
    "monitored": true,
    "useSceneNumbering": false,
    "runtime": 47,
    "tvdbId": 81189,
    "tvRageId": 18164,
    "tvMazeId": 169,
    "imdbId": "tt0903747",
    "titleSlug": "breaking-bad",
    "certification": "TV-MA",
    "genres": ["Crime", "Drama", "Thriller"],
    "tags": [1, 2],
    "added": "2020-01-01T00:00:00Z",
    "ratings": {
      "votes": 15678,
      "value": 9.5
    },
    "statistics": {
      "seasonCount": 5,
      "episodeFileCount": 62,
      "episodeCount": 62,
      "totalEpisodeCount": 62,
      "sizeOnDisk": 53687091200,
      "percentOfEpisodes": 100.0
    }
  }
]
```

### GET `/api/v3/series/{id}`

Obtiene una serie específica por ID.

**Parámetros:**
- `id` (path): ID interno de Sonarr para la serie

**Response:** Objeto serie individual (ver estructura arriba)

---

## 📋 Campos de Series

### Campos Principales

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | int | ID interno de Sonarr |
| `title` | string | Título de la serie |
| `sortTitle` | string | Título optimizado para ordenamiento |
| `status` | string | Estado: `continuing`, `ended`, `upcoming` |
| `ended` | bool | Si la serie ha terminado |
| `overview` | string | Sinopsis/descripción de la serie |
| `network` | string | Red/canal original (HBO, Netflix, etc.) |
| `airTime` | string | Hora de emisión (formato HH:MM) |
| `year` | int | Año de estreno/inicio |
| `path` | string | Ruta del directorio de la serie |
| `runtime` | int | Duración promedio de episodio en minutos |
| `tvdbId` | int | **Requerido** - ID de TheTVDB |
| `tvRageId` | int | ID de TVRage (legacy) |
| `tvMazeId` | int | ID de TVMaze |
| `imdbId` | string | ID de IMDb (formato: ttXXXXXXX) |
| `titleSlug` | string | Slug URL-friendly del título |
| `certification` | string | Clasificación (TV-MA, TV-14, etc.) |
| `genres` | string[] | Lista de géneros |
| `tags` | int[] | IDs de tags asignados |
| `added` | datetime | Fecha cuando se añadió a Sonarr |
| `monitored` | bool | Si está monitoreada para descargas |
| `useSceneNumbering` | bool | Si usa numeración de scene |
| `seasonFolder` | bool | Si organiza en carpetas por temporada |
| `qualityProfileId` | int | ID del perfil de calidad |
| `languageProfileId` | int | ID del perfil de idioma |

### Campos de Imágenes

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `images` | array | Array de objetos de imágenes |
| `images[].coverType` | string | Tipo: `poster`, `banner`, `fanart` |
| `images[].url` | string | URL local (Sonarr) |
| `images[].remoteUrl` | string | URL remota (TheTVDB, etc.) |

### Campos de Ratings

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `ratings.votes` | int | Número de votos |
| `ratings.value` | float | Calificación (0-10) |

### Campos de Estadísticas (Series)

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `statistics.seasonCount` | int | Número de temporadas |
| `statistics.episodeFileCount` | int | Episodios descargados |
| `statistics.episodeCount` | int | Episodios monitoreados |
| `statistics.totalEpisodeCount` | int | Total de episodios existentes |
| `statistics.sizeOnDisk` | int64 | Tamaño total en bytes |
| `statistics.percentOfEpisodes` | float | Porcentaje descargado (0-100) |

---

## 🎬 Campos de Temporadas

Cada serie contiene un array `seasons` con las siguientes propiedades:

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `seasonNumber` | int | Número de temporada (0 = especiales) |
| `monitored` | bool | Si está monitoreada |
| `statistics.previousAiring` | datetime | Última emisión |
| `statistics.episodeFileCount` | int | Archivos de episodios |
| `statistics.episodeCount` | int | Episodios monitoreados |
| `statistics.totalEpisodeCount` | int | Total de episodios |
| `statistics.sizeOnDisk` | int64 | Tamaño en bytes |
| `statistics.percentOfEpisodes` | float | Porcentaje completo |

**Ejemplo:**
```json
{
  "seasonNumber": 1,
  "monitored": true,
  "statistics": {
    "previousAiring": "2008-03-09T02:00:00Z",
    "episodeFileCount": 7,
    "episodeCount": 7,
    "totalEpisodeCount": 7,
    "sizeOnDisk": 4294967296,
    "percentOfEpisodes": 100.0
  }
}
```

---

## 📺 Episodes API

### GET `/api/v3/episode?seriesId={seriesId}`

Obtiene todos los episodios de una serie específica.

**Parámetros Query:**
- `seriesId` (required): ID de la serie

**Response:**
```json
[
  {
    "id": 1,
    "seriesId": 1,
    "tvdbId": 349232,
    "episodeFileId": 1,
    "seasonNumber": 1,
    "episodeNumber": 1,
    "title": "Pilot",
    "airDate": "2008-01-20",
    "airDateUtc": "2008-01-20T02:00:00Z",
    "overview": "When an unassuming high school chemistry teacher...",
    "hasFile": true,
    "monitored": true,
    "absoluteEpisodeNumber": 1,
    "sceneAbsoluteEpisodeNumber": null,
    "sceneEpisodeNumber": null,
    "sceneSeasonNumber": null,
    "unverifiedSceneNumbering": false,
    "grabbed": false
  }
]
```

### GET `/api/v3/episode/{id}`

Obtiene un episodio específico.

**Response:** Objeto episodio individual (ver estructura arriba)

---

## 📋 Campos de Episodios

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | int | ID interno de Sonarr |
| `seriesId` | int | ID de la serie padre |
| `tvdbId` | int | ID de TheTVDB para el episodio |
| `episodeFileId` | int | ID del archivo (si existe) |
| `seasonNumber` | int | Número de temporada |
| `episodeNumber` | int | Número de episodio |
| `title` | string | Título del episodio |
| `airDate` | string | Fecha de emisión (YYYY-MM-DD) |
| `airDateUtc` | datetime | Fecha/hora UTC de emisión |
| `overview` | string | Sinopsis del episodio |
| `hasFile` | bool | Si tiene archivo descargado |
| `monitored` | bool | Si está monitoreado |
| `absoluteEpisodeNumber` | int | Número absoluto (anime) |
| `sceneEpisodeNumber` | int | Número según scene |
| `sceneSeasonNumber` | int | Temporada según scene |
| `unverifiedSceneNumbering` | bool | Numeración scene sin verificar |
| `grabbed` | bool | Si Sonarr lo está descargando |

---

## 📁 Episode Files API

### GET `/api/v3/episodefile?seriesId={seriesId}`

Obtiene información de archivos de episodios para una serie.

**Response:**
```json
[
  {
    "id": 1,
    "seriesId": 1,
    "seasonNumber": 1,
    "relativePath": "Season 01/Breaking Bad - S01E01 - Pilot.mkv",
    "path": "/tv/Breaking Bad/Season 01/Breaking Bad - S01E01 - Pilot.mkv",
    "size": 614572800,
    "dateAdded": "2020-01-01T00:00:00Z",
    "sceneName": "Breaking.Bad.S01E01.720p.BluRay.x264-GRP",
    "releaseGroup": "GRP",
    "quality": {
      "quality": {
        "id": 4,
        "name": "HDTV-720p",
        "source": "television",
        "resolution": 720
      },
      "revision": {
        "version": 1,
        "real": 0,
        "isRepack": false
      }
    },
    "mediaInfo": {
      "audioBitrate": 0,
      "audioChannels": 5.1,
      "audioCodec": "AC3",
      "audioLanguages": "eng",
      "audioStreamCount": 1,
      "videoBitDepth": 8,
      "videoBitrate": 0,
      "videoCodec": "h264",
      "videoFps": 23.976,
      "resolution": "1280x720",
      "runTime": "58:17",
      "scanType": "Progressive",
      "subtitles": "eng"
    },
    "qualityCutoffNotMet": false,
    "languageCutoffNotMet": false
  }
]
```

### Campos de Archivos de Episodios

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | int | ID del archivo |
| `seriesId` | int | ID de la serie |
| `seasonNumber` | int | Número de temporada |
| `relativePath` | string | Ruta relativa desde carpeta serie |
| `path` | string | Ruta completa del archivo |
| `size` | int64 | Tamaño del archivo en bytes |
| `dateAdded` | datetime | Fecha cuando se añadió |
| `sceneName` | string | Nombre del release original |
| `releaseGroup` | string | Grupo de release |
| `quality` | object | Objeto de calidad |
| `mediaInfo` | object | Información técnica del archivo |
| `qualityCutoffNotMet` | bool | Si no cumple calidad objetivo |
| `languageCutoffNotMet` | bool | Si no cumple idioma objetivo |

### Campos de MediaInfo

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `audioBitrate` | int | Bitrate de audio |
| `audioChannels` | float | Canales de audio (2.0, 5.1, 7.1) |
| `audioCodec` | string | Códec de audio (AC3, AAC, DTS) |
| `audioLanguages` | string | Idiomas de audio |
| `audioStreamCount` | int | Número de streams de audio |
| `videoBitDepth` | int | Profundidad de color (8, 10 bits) |
| `videoBitrate` | int | Bitrate de video |
| `videoCodec` | string | Códec de video (h264, h265, etc.) |
| `videoFps` | float | FPS del video |
| `resolution` | string | Resolución (WIDTHxHEIGHT) |
| `runTime` | string | Duración (HH:MM:SS) |
| `scanType` | string | Progressive o Interlaced |
| `subtitles` | string | Idiomas de subtítulos |

---

## 📥 Queue API

### GET `/api/v3/queue`

Obtiene la cola de descargas actual.

**Parámetros Query:**
- `pageSize` (optional): Número de resultados (default: 20)
- `includeUnknownSeriesItems` (optional): Incluir items desconocidos (default: true)

**Response:**
```json
{
  "page": 1,
  "pageSize": 20,
  "totalRecords": 3,
  "records": [
    {
      "id": 1,
      "seriesId": 1,
      "episodeId": 123,
      "title": "Breaking Bad - S01E01 - Pilot",
      "size": 614572800,
      "sizeleft": 102428800,
      "status": "downloading",
      "trackedDownloadStatus": "ok",
      "trackedDownloadState": "downloading",
      "statusMessages": [],
      "downloadId": "abc123",
      "protocol": "torrent",
      "downloadClient": "qBittorrent",
      "indexer": "RARBG",
      "outputPath": "/downloads/Breaking.Bad.S01E01.720p",
      "timedOut": false,
      "estimatedCompletionTime": "2025-11-11T01:30:00Z"
    }
  ]
}
```

---

## 📜 History API

### GET `/api/v3/history`

Obtiene el historial de eventos.

**Parámetros Query:**
- `pageSize` (optional): Resultados por página (default: 20, max: 100)
- `sortKey` (optional): Campo de ordenamiento (default: `date`)
- `sortDirection` (optional): `ascending` o `descending` (default: `descending`)
- `eventType` (optional): Tipo de evento: `grabbed`, `downloadFolderImported`, `downloadFailed`, etc.

**Response:**
```json
{
  "page": 1,
  "pageSize": 50,
  "totalRecords": 150,
  "records": [
    {
      "id": 1,
      "episodeId": 123,
      "seriesId": 1,
      "sourceTitle": "Breaking.Bad.S01E01.720p.BluRay.x264-GRP",
      "quality": {
        "quality": {
          "name": "HDTV-720p"
        }
      },
      "date": "2025-11-11T00:00:00Z",
      "eventType": "downloadFolderImported",
      "downloadId": "abc123"
    }
  ]
}
```

**Tipos de Eventos:**
- `grabbed`: Episodio agregado a cola de descarga
- `downloadFolderImported`: Descarga completada e importada
- `downloadFailed`: Descarga falló
- `episodeFileDeleted`: Archivo de episodio eliminado
- `episodeFileRenamed`: Archivo renombrado

---

## 📅 Calendar API

### GET `/api/v3/calendar`

Obtiene episodios que se emitirán en un rango de fechas.

**Parámetros Query:**
- `start` (required): Fecha inicio (YYYY-MM-DD)
- `end` (required): Fecha fin (YYYY-MM-DD)
- `unmonitored` (optional): Incluir no monitoreados (default: false)

**Response:**
```json
[
  {
    "id": 456,
    "seriesId": 2,
    "episodeFileId": 0,
    "seasonNumber": 2,
    "episodeNumber": 5,
    "title": "Episode Title",
    "airDate": "2025-11-15",
    "airDateUtc": "2025-11-15T01:00:00Z",
    "hasFile": false,
    "monitored": true,
    "series": {
      "title": "Series Name",
      "images": [...]
    }
  }
]
```

---

## 🎯 Quality Profiles API

### GET `/api/v3/qualityprofile`

Obtiene los perfiles de calidad configurados.

**Response:**
```json
[
  {
    "id": 1,
    "name": "HD-720p/1080p",
    "upgradeAllowed": true,
    "cutoff": 4,
    "items": [
      {
        "id": 1,
        "quality": {
          "id": 1,
          "name": "SDTV"
        },
        "allowed": false
      },
      {
        "id": 4,
        "quality": {
          "id": 4,
          "name": "HDTV-720p"
        },
        "allowed": true
      }
    ]
  }
]
```

---

## 🏷️ Tags API

### GET `/api/v3/tag`

Obtiene todos los tags definidos.

**Response:**
```json
[
  {
    "id": 1,
    "label": "anime"
  },
  {
    "id": 2,
    "label": "favorite"
  }
]
```

---

## ℹ️ System Info API

### GET `/api/v3/system/status`

Obtiene información del sistema Sonarr.

**Response:**
```json
{
  "version": "4.0.0.400",
  "buildTime": "2024-01-15T10:30:00Z",
  "isDebug": false,
  "isProduction": true,
  "isAdmin": false,
  "isUserInteractive": false,
  "startupPath": "/app/sonarr/bin",
  "appData": "/config",
  "osName": "ubuntu",
  "osVersion": "22.04",
  "isMonoRuntime": false,
  "isMono": false,
  "isLinux": true,
  "isOsx": false,
  "isWindows": false,
  "mode": "console",
  "branch": "main",
  "authentication": "forms",
  "sqliteVersion": "3.40.1",
  "urlBase": "",
  "runtimeVersion": "7.0.14",
  "runtimeName": ".NET"
}
```

---

## 💡 Ejemplos Prácticos

### Obtener Series con Episodios Faltantes

```bash
# 1. Obtener todas las series
curl -H "X-Api-Key: $API_KEY" http://localhost:8989/api/v3/series

# 2. Para cada serie, filtrar las que tienen episodios faltantes
# statistics.episodeFileCount < statistics.totalEpisodeCount
```

### Obtener Solo Series Activas (Continuing)

```bash
curl -H "X-Api-Key: $API_KEY" http://localhost:8989/api/v3/series \
  | jq '.[] | select(.status == "continuing")'
```

### Obtener Episodios de una Serie

```bash
curl -H "X-Api-Key: $API_KEY" \
  "http://localhost:8989/api/v3/episode?seriesId=1"
```

### Obtener Información de Archivos con MediaInfo

```bash
curl -H "X-Api-Key: $API_KEY" \
  "http://localhost:8989/api/v3/episodefile?seriesId=1" \
  | jq '.[] | {title: .relativePath, codec: .mediaInfo.videoCodec, resolution: .mediaInfo.resolution}'
```

### Obtener Cola de Descargas con Progreso

```bash
curl -H "X-Api-Key: $API_KEY" \
  "http://localhost:8989/api/v3/queue?pageSize=100" \
  | jq '.records[] | {title: .title, progress: ((.size - .sizeleft) / .size * 100)}'
```

### Obtener Próximos Estrenos (Próximos 7 Días)

```bash
START_DATE=$(date +%Y-%m-%d)
END_DATE=$(date -d "+7 days" +%Y-%m-%d)

curl -H "X-Api-Key: $API_KEY" \
  "http://localhost:8989/api/v3/calendar?start=$START_DATE&end=$END_DATE"
```

### Buscar Series por TVDb ID

```bash
curl -H "X-Api-Key: $API_KEY" \
  "http://localhost:8989/api/v3/series" \
  | jq '.[] | select(.tvdbId == 81189)'
```

---

## 🔍 Notas Importantes

### 1. **Identificadores Externos**

- **tvdbId**: Es el identificador principal y **requerido** en Sonarr
- **imdbId**: Útil para cross-referencing pero no se puede usar para añadir series
- **tmdbId**: No soportado directamente, debe mapearse a tvdbId

### 2. **Numeración de Episodios**

- `episodeNumber`: Numeración estándar de la serie
- `absoluteEpisodeNumber`: Numeración continua (usado en anime)
- `sceneEpisodeNumber`: Numeración usada por scene groups
- Usar `sceneEpisodeNumber` cuando `useSceneNumbering` es `true`

### 3. **Estadísticas de Series vs Temporadas**

- **Series statistics**: Suma de todas las temporadas
- **Season statistics**: Datos específicos de cada temporada
- `totalEpisodeCount` incluye episodios especiales (Season 0)

### 4. **Estados de Series**

| Estado | Descripción |
|--------|-------------|
| `continuing` | Serie en emisión activa |
| `ended` | Serie finalizada |
| `upcoming` | Serie anunciada pero no estrenada |
| `deleted` | Serie eliminada de Sonarr |

### 5. **Tamaños de Archivos**

- Todos los tamaños están en **bytes**
- `sizeOnDisk`: Espacio real ocupado en disco
- `size`: Tamaño del archivo/descarga

### 6. **Fechas y Zonas Horarias**

- `airDate`: Fecha local (YYYY-MM-DD)
- `airDateUtc`: Fecha/hora en UTC (ISO 8601)
- `airTime`: Hora local de emisión (HH:MM)

### 7. **Monitoreo**

- `monitored` en series: Buscar nuevos episodios
- `monitored` en temporadas: Buscar episodios de esta temporada
- `monitored` en episodios: Buscar este episodio específico

### 8. **Quality vs MediaInfo**

- `quality`: Perfil/configuración de calidad de Sonarr
- `mediaInfo`: Información técnica real del archivo
- Pueden no coincidir si el archivo tiene mejor calidad que el perfil

---

## 📚 Referencias

- **Documentación Oficial:** https://sonarr.tv/docs/api/
- **GitHub Sonarr:** https://github.com/Sonarr/Sonarr
- **pyarr (Python wrapper):** https://docs.totaldebug.uk/pyarr/
- **ArrAPI (wrapper):** https://arrapi.metamanager.wiki/
- **TheTVDB API:** https://thetvdb.com/api-information

---

## ⚠️ Diferencias con Radarr

| Característica | Sonarr | Radarr |
|----------------|--------|--------|
| Tipo de contenido | Series/Episodios | Películas |
| ID principal | tvdbId | tmdbId |
| Estructura | Series → Seasons → Episodes | Movies |
| Estadísticas | Por serie y temporada | Por película |
| Monitoreo | Serie, temporada y episodio | Solo película |
| Campos específicos | `airTime`, `network`, `seasons` | `studio`, `inCinemas` |

---

**Última Actualización:** 2025-11-11  
**Versión de Sonarr:** v4.0.0+  
**API Version:** v3  
**Mantenido por:** KeeperCheky Project
