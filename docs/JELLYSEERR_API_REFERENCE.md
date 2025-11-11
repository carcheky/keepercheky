# Jellyseerr API Reference

> **Documentación completa de la API REST de Jellyseerr para KeeperCheky**

## 📚 Índice

- [Autenticación](#autenticación)
- [Request API](#request-api)
- [Media Object](#media-object)
- [User Object](#user-object)
- [Estados de Solicitudes](#estados-de-solicitudes)
- [Paginación y Filtros](#paginación-y-filtros)
- [Integración con Servicios](#integración-con-servicios)
- [Ejemplos Prácticos](#ejemplos-prácticos)
- [Mejores Prácticas](#mejores-prácticas)

---

## 🔐 Autenticación

### Headers Requeridos

**Para API Key Authentication:**
```http
X-Api-Key: your-api-key-here
```

**Para Cookie Authentication:**
```http
Cookie: connect.sid=your-session-id
```

### Obtener API Key

1. Login en Jellyseerr/Overseerr
2. Ir a Settings → General → API Key
3. Copiar o regenerar API key

---

## 📋 Request API

### GET `/api/v1/request`

Obtiene todas las solicitudes de medios con paginación y filtros.

**Parámetros Query:**

| Parámetro | Tipo | Descripción | Default |
|-----------|------|-------------|---------|
| `take` | number | Resultados por página | 20 |
| `skip` | number | Número de resultados a saltar | 0 |
| `filter` | string | Filtrar por tipo (`all`, `available`, `approved`, `pending`, `unavailable`) | `all` |
| `sort` | string | Ordenar por campo (`added`, `modified`) | `added` |
| `requestedBy` | number | Filtrar por ID de usuario solicitante | - |

**Ejemplo Request:**
```http
GET /api/v1/request?take=100&skip=0&filter=all&sort=added
```

**Response:**
```json
{
  "pageInfo": {
    "pages": 5,
    "pageSize": 100,
    "results": 47,
    "page": 1
  },
  "results": [
    {
      "id": 123,
      "status": 2,
      "createdAt": "2024-11-10T10:30:00.000Z",
      "updatedAt": "2024-11-10T15:45:00.000Z",
      "type": "movie",
      "is4k": false,
      "serverId": 1,
      "profileId": 4,
      "rootFolder": "/movies",
      "languageProfileId": null,
      "tags": [],
      "isAutoRequest": false,
      "media": {
        "id": 456,
        "mediaType": "movie",
        "tmdbId": 550,
        "tvdbId": null,
        "status": 5,
        "status4k": 1,
        "externalServiceId": 123,
        "externalServiceId4k": null,
        "externalServiceSlug": "fight-club-1999",
        "ratingKey": "12345",
        "serviceId": 1,
        "serviceId4k": null
      },
      "requestedBy": {
        "id": 1,
        "email": "user@example.com",
        "username": "username",
        "plexToken": null,
        "plexUsername": "plexuser",
        "userType": 1,
        "permissions": 2,
        "avatar": "/avatar.jpg",
        "createdAt": "2024-01-01T00:00:00.000Z",
        "updatedAt": "2024-11-10T10:00:00.000Z",
        "requestCount": 5,
        "displayName": "Display Name"
      },
      "modifiedBy": {
        "id": 2,
        "displayName": "Admin User"
      },
      "seasons": []
    }
  ]
}
```

### GET `/api/v1/request/{requestId}`

Obtiene una solicitud específica por su ID.

**Response:**
```json
{
  "id": 123,
  "status": 2,
  "createdAt": "2024-11-10T10:30:00.000Z",
  "updatedAt": "2024-11-10T15:45:00.000Z",
  "type": "movie",
  "media": { /* Media object completo */ },
  "requestedBy": { /* User object completo */ },
  "modifiedBy": { /* User object completo */ },
  "seasons": [ /* Para series */ ]
}
```

### DELETE `/api/v1/request/{requestId}`

Elimina una solicitud específica.

**Response:** `204 No Content` (éxito) o `200 OK`

---

## 🎬 Media Object

El objeto Media contiene información completa sobre el contenido solicitado.

### Campos Principales

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | number | ID interno de Jellyseerr |
| `mediaType` | string | `"movie"` o `"tv"` |
| `tmdbId` | number | ID de TMDB (The Movie Database) |
| `tvdbId` | number | ID de TVDB (para series) |
| `imdbId` | string | ID de IMDB (opcional) |
| `status` | number | Estado del media (ver [Estados](#estados-de-solicitudes)) |
| `status4k` | number | Estado para versión 4K |
| `externalServiceId` | number | ID en Radarr/Sonarr |
| `externalServiceId4k` | number | ID en Radarr/Sonarr para 4K |
| `externalServiceSlug` | string | Slug del título en servicio externo |
| `ratingKey` | string | Key de rating (Plex) |
| `serviceId` | number | ID de instancia de servicio configurada |
| `serviceId4k` | number | ID de instancia 4K |
| `createdAt` | string | Fecha de creación (ISO 8601) |
| `updatedAt` | string | Fecha de última actualización |

### Media Type

- **movie**: Película individual
- **tv**: Serie de televisión (con seasons)

---

## 👤 User Object

Información del usuario que solicitó o modificó la request.

### Campos del Usuario

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | number | ID del usuario |
| `email` | string | Email del usuario |
| `username` | string | Username interno |
| `plexUsername` | string | Username de Plex (si aplica) |
| `plexToken` | string | Token de Plex (no exponer) |
| `jellyfinUsername` | string | Username de Jellyfin |
| `jellyfinUserId` | string | User ID de Jellyfin |
| `userType` | number | Tipo: `1` = Plex, `2` = Local, `3` = Jellyfin |
| `permissions` | number | Máscara de bits de permisos |
| `avatar` | string | URL del avatar |
| `createdAt` | string | Fecha de creación |
| `updatedAt` | string | Fecha de última actualización |
| `requestCount` | number | Número total de solicitudes |
| `displayName` | string | Nombre a mostrar |

---

## 📊 Estados de Solicitudes

### Request Status (Campo `status`)

| Valor | Constante | Descripción |
|-------|-----------|-------------|
| `1` | `PENDING` | Pendiente de aprobación |
| `2` | `APPROVED` | Aprobada, enviada a Radarr/Sonarr |
| `3` | `DECLINED` | Rechazada por administrador |
| `4` | `AVAILABLE` | Disponible en biblioteca |

### Media Status (Campo `media.status`)

| Valor | Constante | Descripción |
|-------|-----------|-------------|
| `1` | `UNKNOWN` | Estado desconocido |
| `2` | `PENDING` | Pendiente |
| `3` | `PROCESSING` | Procesando/descargando |
| `4` | `PARTIALLY_AVAILABLE` | Parcialmente disponible (series) |
| `5` | `AVAILABLE` | Completamente disponible |

---

## 🔄 Seasons (Para Series)

Cuando `type = "tv"`, el campo `seasons` contiene las temporadas solicitadas.

### Season Request Object

```json
{
  "id": 1,
  "seasonNumber": 1,
  "status": 2,
  "createdAt": "2024-11-10T10:30:00.000Z",
  "updatedAt": "2024-11-10T15:45:00.000Z"
}
```

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | number | ID de la season request |
| `seasonNumber` | number | Número de temporada (0 = Specials) |
| `status` | number | Estado de la temporada |
| `createdAt` | string | Fecha de creación |
| `updatedAt` | string | Fecha de actualización |

---

## 📄 Paginación y Filtros

### PageInfo Object

```json
{
  "pages": 5,
  "pageSize": 100,
  "results": 47,
  "page": 1
}
```

| Campo | Descripción |
|-------|-------------|
| `pages` | Total de páginas disponibles |
| `pageSize` | Tamaño de página solicitado |
| `results` | Total de resultados encontrados |
| `page` | Página actual |

### Filtros Disponibles

**Por Estado:**
- `all`: Todas las solicitudes
- `pending`: Solo pendientes
- `approved`: Solo aprobadas
- `available`: Solo disponibles
- `unavailable`: No disponibles aún

**Ordenamiento:**
- `added`: Por fecha de creación (más recientes primero)
- `modified`: Por fecha de modificación

---

## 🔗 Integración con Servicios

### Radarr (Películas)

Jellyseerr envía solicitudes a Radarr usando TMDB ID:

- **profileId**: ID del perfil de calidad en Radarr
- **rootFolder**: Carpeta raíz donde se guardan películas
- **tags**: Tags de Radarr para organización
- **externalServiceId**: ID de película en Radarr

### Sonarr (Series)

Jellyseerr envía solicitudes a Sonarr usando TVDB ID:

- **profileId**: ID del perfil de calidad en Sonarr
- **rootFolder**: Carpeta raíz donde se guardan series
- **languageProfileId**: Perfil de idioma (opcional)
- **tags**: Tags de Sonarr
- **externalServiceId**: ID de serie en Sonarr

**⚠️ Nota sobre TMDB vs TVDB:**
- Jellyseerr usa TMDB como fuente primaria de metadata
- Sonarr usa TVDB para series
- Pueden existir discrepancias en numeración de temporadas/episodios

### Service ID

El campo `serviceId` indica qué instancia de Radarr/Sonarr está configurada:

- Multiple instancias pueden configurarse (ej: Radarr HD, Radarr 4K)
- Cada instancia tiene un ID único
- `serviceId` vincula la request con la instancia correcta

---

## 💡 Ejemplos Prácticos

### Obtener Todas las Solicitudes Pendientes

```http
GET /api/v1/request?filter=pending&take=100
X-Api-Key: your-api-key
```

### Obtener Solicitudes de un Usuario

```http
GET /api/v1/request?requestedBy=5&take=50
X-Api-Key: your-api-key
```

### Obtener una Solicitud Específica

```http
GET /api/v1/request/123
X-Api-Key: your-api-key
```

### Eliminar una Solicitud

```http
DELETE /api/v1/request/123
X-Api-Key: your-api-key
```

---

## 📊 System Info

### GET `/api/v1/status`

Obtiene información del sistema Jellyseerr.

**Response:**
```json
{
  "version": "1.9.2",
  "commitTag": "v1.9.2",
  "updateAvailable": false,
  "commitsBehind": 0
}
```

| Campo | Descripción |
|-------|-------------|
| `version` | Versión actual de Jellyseerr |
| `commitTag` | Tag del commit actual |
| `updateAvailable` | Si hay actualización disponible |
| `commitsBehind` | Commits detrás de la última versión |

---

## 🔑 Mejores Prácticas

### 1. **Paginación**

- Usar `take` y `skip` para paginar grandes cantidades de requests
- Máximo recomendado: 100 items por página
- Siempre revisar `pageInfo` para navegación

### 2. **Filtrado**

- Filtrar por estado para reducir payload
- Usar `requestedBy` para requests de usuario específico
- Combinar filtros con ordenamiento

### 3. **Performance**

- No solicitar todos los requests de una vez
- Cachear datos que no cambien frecuentemente
- Usar `updatedAt` para detectar cambios

### 4. **Seguridad**

- **NUNCA** exponer API keys en código cliente
- Usar headers seguros para autenticación
- Validar permisos antes de eliminar requests

### 5. **Manejo de Estados**

```javascript
// Mapping de estados para UI
const REQUEST_STATUS = {
  1: 'pending',
  2: 'approved', 
  3: 'declined',
  4: 'available'
};

const MEDIA_STATUS = {
  1: 'unknown',
  2: 'pending',
  3: 'processing',
  4: 'partially_available',
  5: 'available'
};
```

### 6. **Polling vs WebSockets**

- Para updates en tiempo real, considerar polling cada 30-60 segundos
- No hacer polling agresivo (< 10 segundos)
- Implementar exponential backoff en errores

---

## 🔍 Request Complete Object Example

```json
{
  "id": 123,
  "status": 2,
  "createdAt": "2024-11-10T10:30:00.000Z",
  "updatedAt": "2024-11-10T15:45:00.000Z",
  "type": "tv",
  "is4k": false,
  "serverId": 1,
  "profileId": 4,
  "rootFolder": "/tv",
  "languageProfileId": 1,
  "tags": [5, 9],
  "isAutoRequest": true,
  "media": {
    "id": 789,
    "mediaType": "tv",
    "tmdbId": 1399,
    "tvdbId": 121361,
    "imdbId": "tt0944947",
    "status": 5,
    "status4k": 1,
    "createdAt": "2024-01-01T00:00:00.000Z",
    "updatedAt": "2024-11-10T15:45:00.000Z",
    "lastSeasonChange": "2024-11-10T15:00:00.000Z",
    "mediaAddedAt": "2024-11-10T16:00:00.000Z",
    "serviceId": 1,
    "serviceId4k": null,
    "externalServiceId": 45,
    "externalServiceId4k": null,
    "externalServiceSlug": "game-of-thrones",
    "ratingKey": "67890",
    "ratingKey4k": null
  },
  "requestedBy": {
    "id": 1,
    "email": "user@example.com",
    "username": "user01",
    "plexToken": null,
    "plexUsername": "plexuser01",
    "jellyfinUsername": null,
    "jellyfinUserId": null,
    "userType": 1,
    "permissions": 2,
    "avatar": "/avatar/1.jpg",
    "createdAt": "2024-01-01T00:00:00.000Z",
    "updatedAt": "2024-11-10T10:00:00.000Z",
    "requestCount": 12,
    "displayName": "User Name"
  },
  "modifiedBy": {
    "id": 2,
    "email": "admin@example.com",
    "username": "admin",
    "displayName": "Admin User",
    "permissions": 2
  },
  "seasons": [
    {
      "id": 1,
      "seasonNumber": 1,
      "status": 5,
      "createdAt": "2024-11-10T10:30:00.000Z",
      "updatedAt": "2024-11-10T16:00:00.000Z"
    },
    {
      "id": 2,
      "seasonNumber": 2,
      "status": 5,
      "createdAt": "2024-11-10T10:30:00.000Z",
      "updatedAt": "2024-11-10T16:00:00.000Z"
    }
  ]
}
```

---

## 🚨 Códigos de Error

| Código | Descripción | Solución |
|--------|-------------|----------|
| `401` | No autenticado | Verificar API key o cookie de sesión |
| `403` | Sin permisos | Usuario no tiene permisos suficientes |
| `404` | Request no encontrado | Verificar que el ID existe |
| `500` | Error del servidor | Revisar logs de Jellyseerr |

---

## 📚 Referencias

- **Documentación Oficial:** https://docs.seerr.dev/
- **GitHub Jellyseerr:** https://github.com/Fallenbagel/jellyseerr
- **GitHub Overseerr:** https://github.com/sct/overseerr
- **API Docs (Overseerr):** https://api-docs.overseerr.dev/
- **OpenAPI Spec:** Disponible en la ruta `/api-docs` de tu instancia

---

## 🔄 Changelog del Documento

| Fecha | Versión | Cambios |
|-------|---------|---------|
| 2025-11-11 | 1.0.0 | Creación inicial del documento |

---

**Última Actualización:** 2025-11-11  
**Versión de Jellyseerr:** 1.9.2+  
**Mantenido por:** KeeperCheky Project
