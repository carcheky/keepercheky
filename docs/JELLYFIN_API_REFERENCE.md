# Jellyfin API Reference

> **Documentación completa de la API REST de Jellyfin para KeeperCheky**

## 📚 Índice

- [Autenticación](#autenticación)
- [Items API](#items-api)
- [Tipos de Items](#tipos-de-items)
- [Campos Disponibles](#campos-disponibles)
- [Ordenamiento](#ordenamiento)
- [Imágenes](#imágenes)
- [Playback State](#playback-state)
- [Ejemplos Prácticos](#ejemplos-prácticos)

---

## 🔐 Autenticación

### POST `/Users/Authenticate`

Autentica usuarios con usuario y contraseña.

**Headers Requeridos:**
```http
X-Emby-Client: My Client Application
X-Emby-Device-Name: Device Name
X-Emby-Device-Id: unique-device-id
X-Emby-Version: 1.0.0
```

**Request Body:**
```json
{
  "Username": "demo",
  "Password": ""
}
```

**Response (200 OK):**
```json
{
  "User": {
    "Name": "demo",
    "Id": "some-user-id"
  },
  "SessionInfo": {},
  "AccessToken": "a-valid-access-token",
  "ServerId": "some-server-id"
}
```

---

## 📁 Items API

### GET `/Items`

Obtiene items de la biblioteca con filtros y paginación.

**Parámetros Query:**

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `userId` | string | **Requerido** - ID del usuario autenticado |
| `parentId` | string | ID de la biblioteca padre (opcional) |
| `includeItemTypes` | string[] | **Filtro por tipo**: `Movie`, `Series`, `Episode`, `Season` |
| `fields` | string[] | Campos adicionales a incluir (ver [Campos Disponibles](#campos-disponibles)) |
| `sortBy` | string[] | Campo de ordenamiento (ver [Ordenamiento](#ordenamiento)) |
| `sortOrder` | string[] | `Ascending` o `Descending` |
| `startIndex` | number | Índice inicial para paginación (default: 0) |
| `limit` | number | Número máximo de resultados |
| `recursive` | boolean | Buscar recursivamente en subdirectorios |
| `enableImages` | boolean | Incluir URLs de imágenes |
| `enableUserData` | boolean | Incluir estado de reproducción del usuario |

**Ejemplo de Request:**
```http
GET /Items?userId=abc123&includeItemTypes=Movie&fields=Overview,Genres,MediaStreams&sortBy=PremiereDate&sortOrder=Descending&startIndex=0&limit=20&recursive=true
```

**Response:**
```json
{
  "Items": [
    {
      "Name": "The Matrix",
      "Id": "movie-id-123",
      "Type": "Movie",
      "ProductionYear": 1999,
      "Genres": ["Action", "Sci-Fi"],
      "OfficialRating": "R",
      "Overview": "A hacker discovers the true nature of reality...",
      "MediaSources": [
        {
          "Id": "source-123",
          "Path": "/movies/The Matrix (1999)/The Matrix (1999).mkv",
          "Size": 2147483648
        }
      ],
      "ImageTags": {
        "Primary": "abc123def456"
      },
      "UserData": {
        "Played": false,
        "PlayCount": 0,
        "IsFavorite": false
      }
    }
  ],
  "TotalRecordCount": 150
}
```

---

## 🎬 Tipos de Items

Jellyfin clasifica el contenido multimedia en varios tipos:

| Tipo | Descripción | Uso en KeeperCheky |
|------|-------------|-------------------|
| `Movie` | Película | ✅ Pestaña Películas |
| `Series` | Serie de TV | ⏳ Futuro: Pestaña Series |
| `Season` | Temporada de serie | ❌ Filtrar |
| `Episode` | Episodio de serie | ⏳ Futuro: Vista de episodios |
| `Audio` | Audio/Música | ❌ No usado |
| `MusicAlbum` | Álbum de música | ❌ No usado |
| `MusicArtist` | Artista musical | ❌ No usado |
| `Book` | Libro | ❌ No usado |
| `AudioBook` | Audiolibro | ❌ No usado |

### ⚠️ Importante: Filtrado de Tipos

**SIEMPRE** filtrar por `includeItemTypes` para evitar mezclar contenido:

```go
// ❌ MAL - No filtrar devuelve Movies, Series, Episodes, Seasons
jellyfinMedia, err := jellyfinClient.GetLibrary(ctx)

// ✅ BIEN - Filtrar solo Movies
params := map[string]interface{}{
    "includeItemTypes": []string{"Movie"},
    "recursive": true,
}
```

---

## 📋 Campos Disponibles

Campos que se pueden incluir con el parámetro `fields`:

| Campo | Descripción |
|-------|-------------|
| `Overview` | Sinopsis/descripción |
| `Genres` | Géneros |
| `MediaStreams` | Streams de video/audio/subtítulos |
| `MediaSources` | Fuentes de medios (archivos) |
| `Path` | Ruta del archivo |
| `ProviderIds` | IDs externos (IMDb, TMDb, etc.) |
| `Studios` | Estudios de producción |
| `People` | Cast y crew |
| `ProductionYear` | Año de producción |
| `CommunityRating` | Calificación de la comunidad |
| `CriticRating` | Calificación de críticos |
| `OfficialRating` | Clasificación oficial (PG, R, etc.) |
| `DateCreated` | Fecha de creación en Jellyfin |
| `DateModified` | Fecha de modificación |
| `PlayAccess` | Acceso de reproducción |
| `RemoteTrailers` | Trailers remotos |

---

## 🔢 Ordenamiento

Valores válidos para el parámetro `sortBy`:

| Valor | Descripción |
|-------|-------------|
| `PremiereDate` | Fecha de estreno |
| `ProductionYear` | Año de producción |
| `Name` | Nombre/título |
| `SortName` | Nombre ordenado |
| `Random` | Aleatorio |
| `CommunityRating` | Calificación |
| `CriticRating` | Calificación de críticos |
| `DateCreated` | Fecha añadida a Jellyfin |
| `DatePlayed` | Última reproducción |
| `PlayCount` | Número de reproducciones |
| `Runtime` | Duración |

---

## 🖼️ Imágenes

### Tipos de Imágenes

| Tipo | Descripción |
|------|-------------|
| `Primary` | Póster principal |
| `Backdrop` | Fondos/fanart |
| `Logo` | Logo del título |
| `Thumb` | Miniatura |
| `Art` | Arte adicional |
| `Banner` | Banner horizontal |

### URLs de Imágenes

**Formato:**
```
{baseURL}/Items/{itemId}/Images/{imageType}?tag={imageTag}&maxWidth={width}&maxHeight={height}&quality={quality}
```

**Ejemplo:**
```
https://demo.jellyfin.org/stable/Items/abc123/Images/Primary?tag=xyz789&maxWidth=400&maxHeight=600&quality=90
```

**Parámetros:**
- `tag`: Tag de la imagen (del campo `ImageTags`)
- `maxWidth`: Ancho máximo en píxeles
- `maxHeight`: Alto máximo en píxeles
- `quality`: Calidad JPEG (1-100)

---

## ▶️ Playback State

### POST `/Sessions/Playing`

Reportar inicio de reproducción.

### POST `/Sessions/Playing/Progress`

Reportar progreso de reproducción (llamar periódicamente).

### POST `/Sessions/Playing/Stopped`

Reportar fin de reproducción.

### POST `/PlayedItems/{itemId}`

Marcar item como reproducido.

### DELETE `/UserPlayedItems/{itemId}`

Marcar item como no reproducido.

---

## 💡 Ejemplos Prácticos

### Obtener Solo Películas

```http
GET /Items?userId={userId}&includeItemTypes=Movie&recursive=true&fields=Overview,Genres,MediaSources&sortBy=Name&sortOrder=Ascending
```

### Obtener Series con Episodios

```http
GET /Items?userId={userId}&includeItemTypes=Series&recursive=true&fields=RecursiveItemCount
```

### Buscar por Año

```http
GET /Items?userId={userId}&includeItemTypes=Movie&years=2023&recursive=true
```

### Buscar por Género

```http
GET /Items?userId={userId}&includeItemTypes=Movie&genres=Action&recursive=true
```

### Solo Favoritos

```http
GET /Items?userId={userId}&includeItemTypes=Movie&filters=IsFavorite&recursive=true
```

### Solo No Vistos

```http
GET /Items?userId={userId}&includeItemTypes=Movie&filters=IsUnplayed&recursive=true
```

---

## 🔍 Búsqueda

### GET `/Search/Hints`

Búsqueda general en toda la biblioteca.

**Parámetros:**
- `searchTerm`: Término de búsqueda
- `userId`: ID del usuario
- `includeItemTypes`: Filtrar por tipos
- `limit`: Número máximo de resultados

**Ejemplo:**
```http
GET /Search/Hints?searchTerm=inception&userId={userId}&includeItemTypes=Movie&limit=10
```

---

## 📊 Información del Sistema

### GET `/System/Info/Public`

Obtener información pública del servidor (no requiere autenticación).

**Response:**
```json
{
  "ServerName": "My Jellyfin Server",
  "Version": "10.10.0",
  "ProductName": "Jellyfin Server",
  "OperatingSystem": "Linux",
  "Id": "server-id-123"
}
```

---

## 🔑 Notas Importantes

### 1. **Headers de Autenticación**

Todas las peticiones autenticadas deben incluir:
```http
X-Emby-Token: {accessToken}
```

O alternativamente:
```http
Authorization: MediaBrowser Token="{accessToken}"
```

### 2. **Formato de IDs**

Los IDs en Jellyfin son strings hexadecimales, por ejemplo:
```
abc123def456789
```

### 3. **Ticks**

Los tiempos en Jellyfin se miden en "ticks" (1 tick = 100 nanosegundos):
- 1 segundo = 10,000,000 ticks
- 1 minuto = 600,000,000 ticks
- 1 hora = 36,000,000,000 ticks

### 4. **Códigos de Error Comunes**

| Código | Descripción |
|--------|-------------|
| 401 | No autenticado o token inválido |
| 404 | Item o recurso no encontrado |
| 500 | Error interno del servidor |

---

## 📚 Referencias

- **Documentación Oficial:** https://api.jellyfin.org
- **SDK TypeScript:** https://github.com/jellyfin/jellyfin-sdk-typescript
- **Jellyfin Server:** https://github.com/jellyfin/jellyfin
- **OpenAPI Spec:** https://api.jellyfin.org/openapi/jellyfin-openapi-stable.json

---

**Última Actualización:** 2025-11-11  
**Versión de Jellyfin:** 10.10.0+  
**Mantenido por:** KeeperCheky Project
