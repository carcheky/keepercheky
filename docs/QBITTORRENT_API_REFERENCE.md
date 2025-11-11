# qBittorrent Web API Reference

> **Documentación completa de la API REST de qBittorrent para KeeperCheky**

## 📚 Índice

- [Autenticación](#autenticación)
- [Endpoints Principales](#endpoints-principales)
- [Torrents Info](#torrents-info)
- [Torrent Properties](#torrent-properties)
- [Estados de Torrents](#estados-de-torrents)
- [Files API](#files-api)
- [Trackers API](#trackers-api)
- [Peers API](#peers-api)
- [Transfer Info](#transfer-info)
- [Server State](#server-state)
- [Categories & Tags](#categories--tags)
- [Ejemplos Prácticos](#ejemplos-prácticos)

---

## 🔐 Autenticación

### POST `/api/v2/auth/login`

Autentica usuarios con usuario y contraseña. La autenticación se realiza mediante cookies (SID).

**Request Body (form-data):**
```
username=admin
password=adminpass
```

**Response (200 OK):**
```
Ok.
```

**Cookie Devuelta:**
```
SID=your-session-id
```

**Uso posterior:**
Todas las peticiones subsecuentes deben incluir la cookie `SID` en el header:
```http
Cookie: SID=your-session-id
```

### POST `/api/v2/auth/logout`

Cierra la sesión actual.

**Response:**
```
Ok.
```

---

## 📁 Endpoints Principales

### GET `/api/v2/app/version`

Obtiene la versión de qBittorrent.

**Response:**
```
v4.5.2
```

### GET `/api/v2/app/buildInfo`

Obtiene información de compilación.

**Response:**
```json
{
  "qt": "6.4.2",
  "libtorrent": "2.0.8.0",
  "boost": "1.81.0",
  "openssl": "3.0.7",
  "bitness": 64
}
```

### GET `/api/v2/app/preferences`

Obtiene las preferencias de la aplicación, incluyendo rutas de guardado.

**Response (extracto):**
```json
{
  "save_path": "/downloads",
  "temp_path": "/downloads/incomplete",
  "temp_path_enabled": true,
  "export_dir": "/downloads/.torrents",
  "export_dir_fin": "/downloads/.torrents/finished",
  "scan_dirs": {
    "/watch": 0
  },
  "max_ratio": 2.0,
  "max_seeding_time": -1,
  "download_limit": -1,
  "upload_limit": -1
}
```

---

## 🗂️ Torrents Info

### GET `/api/v2/torrents/info`

Obtiene lista de todos los torrents con su información básica.

**Parámetros Query (opcionales):**

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `filter` | string | Filtrar por estado: `all`, `downloading`, `seeding`, `completed`, `paused`, `active`, `inactive`, `resumed`, `stalled`, `stalled_uploading`, `stalled_downloading` |
| `category` | string | Filtrar por categoría |
| `tag` | string | Filtrar por tag (separar múltiples con coma) |
| `sort` | string | Campo de ordenamiento |
| `reverse` | boolean | Ordenar en reversa |
| `limit` | number | Límite de resultados |
| `offset` | number | Offset para paginación |
| `hashes` | string | Lista de hashes separados por `|` (pipe) |

**Ejemplo de Request:**
```http
GET /api/v2/torrents/info?filter=downloading&sort=added_on&reverse=true
```

**Response (array de torrents):**
```json
[
  {
    "added_on": 1699564800,
    "amount_left": 0,
    "auto_tmm": false,
    "availability": 1.0,
    "category": "movies",
    "completed": 1699565400,
    "completion_on": 1699565400,
    "content_path": "/downloads/Movie.2023.1080p",
    "dl_limit": -1,
    "dlspeed": 0,
    "download_path": "",
    "downloaded": 2147483648,
    "downloaded_session": 2147483648,
    "eta": 8640000,
    "f_l_piece_prio": false,
    "force_start": false,
    "hash": "abc123def456789",
    "infohash_v1": "abc123def456789",
    "infohash_v2": "",
    "last_activity": 1699566000,
    "magnet_uri": "magnet:?xt=urn:btih:abc123...",
    "max_ratio": -1,
    "max_seeding_time": -1,
    "name": "Movie.2023.1080p.BluRay.x264",
    "num_complete": 50,
    "num_incomplete": 10,
    "num_leechs": 2,
    "num_seeds": 5,
    "priority": 0,
    "progress": 1.0,
    "ratio": 1.5,
    "ratio_limit": -2,
    "save_path": "/downloads/",
    "seeding_time": 3600,
    "seeding_time_limit": -2,
    "seen_complete": 1699565400,
    "seq_dl": false,
    "size": 2147483648,
    "state": "uploading",
    "super_seeding": false,
    "tags": "hdr,bluray",
    "time_active": 7200,
    "total_size": 2147483648,
    "tracker": "https://tracker.example.com:443/announce",
    "trackers_count": 3,
    "up_limit": -1,
    "uploaded": 3221225472,
    "uploaded_session": 3221225472,
    "upspeed": 524288
  }
]
```

### Campos Completos del Torrent

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `hash` | string | Hash del torrent (identificador único) |
| `infohash_v1` | string | InfoHash v1 (BitTorrent v1) |
| `infohash_v2` | string | InfoHash v2 (BitTorrent v2, puede estar vacío) |
| `name` | string | Nombre del torrent |
| `size` | integer | Tamaño total en bytes |
| `total_size` | integer | Tamaño total incluyendo archivos no seleccionados |
| `progress` | float | Progreso de descarga (0.0 a 1.0) |
| `dlspeed` | integer | Velocidad de descarga (bytes/s) |
| `upspeed` | integer | Velocidad de subida (bytes/s) |
| `downloaded` | integer | Total descargado en bytes |
| `uploaded` | integer | Total subido en bytes |
| `downloaded_session` | integer | Descargado en esta sesión |
| `uploaded_session` | integer | Subido en esta sesión |
| `ratio` | float | Ratio de compartición (uploaded/downloaded) |
| `eta` | integer | Tiempo estimado para completar (segundos) |
| `state` | string | Estado del torrent (ver [Estados](#estados-de-torrents)) |
| `num_seeds` | integer | Número de seeds conectados |
| `num_complete` | integer | Total de seeds en el swarm |
| `num_leechs` | integer | Número de leechers conectados |
| `num_incomplete` | integer | Total de leechers en el swarm |
| `amount_left` | integer | Bytes restantes por descargar |
| `save_path` | string | Ruta de guardado del torrent |
| `content_path` | string | Ruta completa del contenido |
| `download_path` | string | Ruta temporal de descarga (si está habilitada) |
| `category` | string | Categoría del torrent |
| `tags` | string | Tags separados por comas |
| `added_on` | integer | Timestamp Unix de cuando se añadió |
| `completion_on` | integer | Timestamp Unix de cuando se completó |
| `completed` | integer | Total de bytes completados |
| `last_activity` | integer | Timestamp Unix de última actividad |
| `seeding_time` | integer | Tiempo total en seed (segundos) |
| `time_active` | integer | Tiempo total activo (segundos) |
| `tracker` | string | URL del tracker principal |
| `trackers_count` | integer | Número total de trackers |
| `magnet_uri` | string | URI magnet del torrent |
| `priority` | integer | Prioridad del torrent |
| `seq_dl` | boolean | Descarga secuencial habilitada |
| `f_l_piece_prio` | boolean | Prioridad first/last piece |
| `force_start` | boolean | Forzar inicio |
| `super_seeding` | boolean | Super seeding habilitado |
| `auto_tmm` | boolean | Automatic Torrent Management habilitado |
| `availability` | float | Disponibilidad de piezas (0.0 a 1.0+) |
| `dl_limit` | integer | Límite de descarga (-1 = sin límite) |
| `up_limit` | integer | Límite de subida (-1 = sin límite) |
| `max_ratio` | float | Ratio máximo (-1 = global, -2 = sin límite) |
| `max_seeding_time` | integer | Tiempo máximo de seed (-1 = global, -2 = sin límite) |
| `ratio_limit` | float | Límite de ratio aplicado |
| `seeding_time_limit` | integer | Límite de tiempo de seed aplicado |
| `seen_complete` | integer | Timestamp de última vez visto completo |

---

## 🔍 Torrent Properties

### GET `/api/v2/torrents/properties`

Obtiene propiedades detalladas de un torrent específico.

**Parámetros Query:**
- `hash` (requerido): Hash del torrent

**Ejemplo:**
```http
GET /api/v2/torrents/properties?hash=abc123def456789
```

**Response:**
```json
{
  "addition_date": 1699564800,
  "comment": "Uploaded by user123",
  "completion_date": 1699565400,
  "created_by": "uTorrent/3.5.5",
  "creation_date": 1699564700,
  "dl_limit": -1,
  "dl_speed": 0,
  "dl_speed_avg": 524288,
  "download_path": "",
  "eta": 8640000,
  "last_seen": 1699566000,
  "nb_connections": 10,
  "nb_connections_limit": 100,
  "peers": 7,
  "peers_total": 60,
  "piece_size": 4194304,
  "pieces_have": 512,
  "pieces_num": 512,
  "reannounce": 1200,
  "save_path": "/downloads/movies",
  "seeding_time": 3600,
  "seeds": 5,
  "seeds_total": 50,
  "share_ratio": 1.5,
  "time_elapsed": 7200,
  "total_downloaded": 2147483648,
  "total_downloaded_session": 2147483648,
  "total_size": 2147483648,
  "total_uploaded": 3221225472,
  "total_uploaded_session": 3221225472,
  "total_wasted": 1048576,
  "up_limit": -1,
  "up_speed": 524288,
  "up_speed_avg": 262144
}
```

### Campos de Propiedades Detalladas

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `addition_date` | integer | Timestamp Unix de cuando se añadió |
| `comment` | string | Comentario del torrent |
| `completion_date` | integer | Timestamp Unix de completación (-1 si no completado) |
| `created_by` | string | Programa que creó el torrent |
| `creation_date` | integer | Timestamp Unix de creación del torrent |
| `save_path` | string | Ruta de guardado |
| `download_path` | string | Ruta temporal de descarga |
| `piece_size` | integer | Tamaño de cada pieza en bytes |
| `pieces_num` | integer | Número total de piezas |
| `pieces_have` | integer | Número de piezas descargadas |
| `total_size` | integer | Tamaño total en bytes |
| `total_downloaded` | integer | Total descargado (incluye datos descartados) |
| `total_uploaded` | integer | Total subido |
| `total_wasted` | integer | Total de datos descartados (bytes) |
| `dl_speed` | integer | Velocidad de descarga actual (bytes/s) |
| `dl_speed_avg` | integer | Velocidad promedio de descarga |
| `up_speed` | integer | Velocidad de subida actual (bytes/s) |
| `up_speed_avg` | integer | Velocidad promedio de subida |
| `eta` | integer | ETA en segundos (8640000 = infinito) |
| `seeding_time` | integer | Tiempo total en seed (segundos) |
| `time_elapsed` | integer | Tiempo total transcurrido (segundos) |
| `share_ratio` | float | Ratio de compartición |
| `seeds` | integer | Seeds conectados |
| `seeds_total` | integer | Seeds totales en el swarm |
| `peers` | integer | Peers conectados |
| `peers_total` | integer | Peers totales en el swarm |
| `nb_connections` | integer | Número de conexiones |
| `nb_connections_limit` | integer | Límite de conexiones |
| `last_seen` | integer | Timestamp de última actividad |
| `reannounce` | integer | Segundos hasta próximo announce |
| `dl_limit` | integer | Límite de descarga |
| `up_limit` | integer | Límite de subida |

---

## 🚦 Estados de Torrents

Los torrents en qBittorrent pueden estar en diferentes estados:

| Estado | Descripción |
|--------|-------------|
| `error` | Ha ocurrido un error (ej: falta archivo) |
| `missingFiles` | Faltan archivos del torrent |
| `uploading` | Subiendo datos (seeding) |
| `pausedUP` | Pausado y descarga completa |
| `queuedUP` | En cola para subir |
| `stalledUP` | Subida estancada (esperando peers) |
| `checkingUP` | Verificando archivos (antes de subir) |
| `forcedUP` | Subida forzada |
| `allocating` | Asignando espacio en disco |
| `downloading` | Descargando datos |
| `metaDL` | Descargando metadata (magnet links) |
| `pausedDL` | Pausado y descarga incompleta |
| `queuedDL` | En cola para descargar |
| `stalledDL` | Descarga estancada (esperando seeds) |
| `checkingDL` | Verificando archivos (antes de descargar) |
| `forcedDL` | Descarga forzada |
| `checkingResumeData` | Verificando datos de reanudación |
| `moving` | Moviendo archivos del torrent |
| `unknown` | Estado desconocido |

### Estados de Seeding

Los siguientes estados indican que el torrent está en modo seeding (compartiendo):
- `uploading`
- `stalledUP`
- `checkingUP`
- `queuedUP`
- `forcedUP`

### Estados de Descarga

Los siguientes estados indican descarga activa o pendiente:
- `downloading`
- `metaDL`
- `stalledDL`
- `checkingDL`
- `queuedDL`
- `forcedDL`

---

## 📄 Files API

### GET `/api/v2/torrents/files`

Obtiene la lista de archivos de un torrent.

**Parámetros Query:**
- `hash` (requerido): Hash del torrent
- `indexes` (opcional): Lista de índices separados por `|`

**Ejemplo:**
```http
GET /api/v2/torrents/files?hash=abc123def456789
```

**Response:**
```json
[
  {
    "index": 0,
    "name": "Movie.2023.1080p.mkv",
    "size": 2147483648,
    "progress": 1.0,
    "priority": 7,
    "is_seed": true,
    "piece_range": [0, 511],
    "availability": 1.0
  },
  {
    "index": 1,
    "name": "Subtitles/Movie.2023.en.srt",
    "size": 102400,
    "progress": 1.0,
    "priority": 7,
    "is_seed": true,
    "piece_range": [512, 512],
    "availability": 1.0
  }
]
```

### Campos de Files

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `index` | integer | Índice del archivo en el torrent |
| `name` | string | Nombre/ruta del archivo (relativa al torrent) |
| `size` | integer | Tamaño del archivo en bytes |
| `progress` | float | Progreso de descarga (0.0 a 1.0) |
| `priority` | integer | Prioridad (0=No descargar, 1=Normal, 6=Alta, 7=Máxima) |
| `is_seed` | boolean | `true` si el archivo está completamente descargado |
| `piece_range` | array | Rango de piezas [primera, última] |
| `availability` | float | Disponibilidad del archivo en el swarm |

---

## 🌐 Trackers API

### GET `/api/v2/torrents/trackers`

Obtiene la lista de trackers de un torrent.

**Parámetros Query:**
- `hash` (requerido): Hash del torrent

**Ejemplo:**
```http
GET /api/v2/torrents/trackers?hash=abc123def456789
```

**Response:**
```json
[
  {
    "url": "https://tracker.example.com:443/announce",
    "status": 2,
    "tier": 0,
    "num_peers": 50,
    "num_seeds": 30,
    "num_leeches": 20,
    "num_downloaded": 500,
    "msg": "Working",
    "last_message": "",
    "next_announce": 1800,
    "min_announce": 60
  },
  {
    "url": "udp://tracker2.example.com:6969/announce",
    "status": 1,
    "tier": 1,
    "num_peers": 0,
    "num_seeds": 0,
    "num_leeches": 0,
    "num_downloaded": 0,
    "msg": "Not contacted yet",
    "last_message": "",
    "next_announce": 0,
    "min_announce": 0
  }
]
```

### Campos de Trackers

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `url` | string | URL del tracker |
| `status` | integer | Estado: 0=Disabled, 1=Not contacted, 2=Working, 3=Updating, 4=Not working |
| `tier` | integer | Nivel del tracker (trackers del mismo tier se usan simultáneamente) |
| `num_peers` | integer | Número de peers reportados |
| `num_seeds` | integer | Número de seeds reportados |
| `num_leeches` | integer | Número de leechers reportados |
| `num_downloaded` | integer | Número de descargas completas reportadas |
| `msg` | string | Mensaje de estado |
| `last_message` | string | Último mensaje del tracker |
| `next_announce` | integer | Segundos hasta el próximo announce |
| `min_announce` | integer | Intervalo mínimo entre announces |

### Estados de Trackers

| Valor | Descripción |
|-------|-------------|
| 0 | Deshabilitado |
| 1 | No contactado aún |
| 2 | Funcionando correctamente |
| 3 | Actualizando |
| 4 | No funcionando |

---

## 👥 Peers API

### GET `/api/v2/torrents/peers`

Obtiene la lista de peers conectados a un torrent.

**Parámetros Query:**
- `hash` (requerido): Hash del torrent
- `rid` (opcional): Request ID para obtener solo cambios desde última petición

**Ejemplo:**
```http
GET /api/v2/torrents/peers?hash=abc123def456789
```

**Response:**
```json
{
  "full_update": true,
  "peers": {
    "192.168.1.100:51234": {
      "ip": "192.168.1.100",
      "port": 51234,
      "client": "uTorrent 3.5.5",
      "connection": "uTP",
      "flags": "d E X I P",
      "flags_desc": "d = Downloading, E = Encrypted, X = uTorrent extension, I = Incoming, P = uTP",
      "country": "US",
      "country_code": "US",
      "dl_speed": 524288,
      "up_speed": 262144,
      "downloaded": 104857600,
      "uploaded": 52428800,
      "progress": 0.45,
      "relevance": 1.0,
      "files": "",
      "seed": false,
      "peer_id_client": "UT355"
    }
  }
}
```

### Campos de Peers

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `ip` | string | Dirección IP del peer |
| `port` | integer | Puerto del peer |
| `client` | string | Cliente BitTorrent del peer |
| `connection` | string | Tipo de conexión (BT, uTP, etc.) |
| `flags` | string | Flags del peer (ver descripción abajo) |
| `flags_desc` | string | Descripción de los flags |
| `country` | string | Nombre del país |
| `country_code` | string | Código ISO del país |
| `dl_speed` | integer | Velocidad de descarga desde este peer (bytes/s) |
| `up_speed` | integer | Velocidad de subida a este peer (bytes/s) |
| `downloaded` | integer | Bytes descargados desde este peer |
| `uploaded` | integer | Bytes subidos a este peer |
| `progress` | float | Progreso del peer (0.0 a 1.0) |
| `relevance` | float | Relevancia del peer |
| `files` | string | Archivos que tiene el peer |
| `seed` | boolean | `true` si el peer es un seed |
| `peer_id_client` | string | ID del cliente del peer |

### Peer Flags

| Flag | Descripción |
|------|-------------|
| `d` | Downloading (interested, not choked) |
| `D` | Your client is interested but peer is choking |
| `u` | Uploading (peer interested, you're not choking) |
| `U` | Peer interested but you're choking |
| `O` | Optimistic unchoke |
| `S` | Peer snubbing (not uploading despite unchoke) |
| `I` | Incoming connection |
| `K` | Peer unchoking you but you're not interested |
| `?` | Peer is choked but interested |
| `X` | Peer supports extensions protocol |
| `H` | Peer connected via DHT |
| `E` | Peer is encrypted |
| `L` | Peer is local (same subnet) |
| `P` | Peer is using uTP |

---

## 📊 Transfer Info

### GET `/api/v2/transfer/info`

Obtiene información global de transferencia.

**Response:**
```json
{
  "dl_info_speed": 524288,
  "dl_info_data": 1073741824,
  "up_info_speed": 262144,
  "up_info_data": 2147483648,
  "dl_rate_limit": -1,
  "up_rate_limit": -1,
  "dht_nodes": 450,
  "connection_status": "connected",
  "queueing": true,
  "use_alt_speed_limits": false,
  "refresh_interval": 1500
}
```

### Campos de Transfer Info

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `dl_info_speed` | integer | Velocidad global de descarga (bytes/s) |
| `dl_info_data` | integer | Datos descargados en esta sesión (bytes) |
| `up_info_speed` | integer | Velocidad global de subida (bytes/s) |
| `up_info_data` | integer | Datos subidos en esta sesión (bytes) |
| `dl_rate_limit` | integer | Límite global de descarga (-1 = sin límite) |
| `up_rate_limit` | integer | Límite global de subida (-1 = sin límite) |
| `dht_nodes` | integer | Nodos DHT conectados |
| `connection_status` | string | Estado de conexión: `connected`, `firewalled`, `disconnected` |
| `queueing` | boolean | Sistema de colas habilitado |
| `use_alt_speed_limits` | boolean | Límites de velocidad alternativos activos |
| `refresh_interval` | integer | Intervalo de refresco (milisegundos) |

---

## 🖥️ Server State

### GET `/api/v2/sync/maindata`

Obtiene el estado del servidor y cambios incrementales.

**Parámetros Query:**
- `rid` (opcional): Request ID para cambios incrementales

**Response:**
```json
{
  "rid": 42,
  "full_update": true,
  "server_state": {
    "dl_info_speed": 524288,
    "dl_info_data": 1073741824,
    "up_info_speed": 262144,
    "up_info_data": 2147483648,
    "dht_nodes": 450,
    "connection_status": "connected",
    "queueing": true,
    "use_alt_speed_limits": false,
    "free_space_on_disk": 107374182400,
    "global_ratio": "1.5",
    "alltime_dl": 10737418240,
    "alltime_ul": 21474836480,
    "total_peer_connections": 150,
    "read_cache_hits": "95.5%",
    "write_cache_overload": "0%",
    "total_buffers_size": 67108864,
    "refresh_interval": 1500
  },
  "torrents": {
    "abc123def456789": {
      "name": "Movie.2023.1080p",
      "state": "uploading",
      "progress": 1.0
    }
  },
  "categories": {
    "movies": {
      "name": "movies",
      "savePath": "/downloads/movies"
    }
  },
  "tags": ["hdr", "bluray", "4k"],
  "trackers": {}
}
```

---

## 📁 Categories & Tags

### GET `/api/v2/torrents/categories`

Obtiene todas las categorías y sus rutas.

**Response:**
```json
{
  "movies": {
    "name": "movies",
    "savePath": "/downloads/movies"
  },
  "tv": {
    "name": "tv",
    "savePath": "/downloads/tv"
  }
}
```

### GET `/api/v2/torrents/tags`

Obtiene todos los tags disponibles.

**Response:**
```json
["hdr", "bluray", "4k", "remux", "hevc"]
```

---

## 💡 Ejemplos Prácticos

### Obtener Todos los Torrents en Seeding

```http
GET /api/v2/torrents/info?filter=seeding
```

### Obtener Torrents de una Categoría

```http
GET /api/v2/torrents/info?category=movies
```

### Obtener Torrents con Tag Específico

```http
GET /api/v2/torrents/info?tag=hdr
```

### Obtener Solo Torrents Activos

```http
GET /api/v2/torrents/info?filter=active
```

### Obtener Información Completa de un Torrent

```bash
# 1. Info básica
GET /api/v2/torrents/info?hashes=abc123def456789

# 2. Propiedades detalladas
GET /api/v2/torrents/properties?hash=abc123def456789

# 3. Lista de archivos
GET /api/v2/torrents/files?hash=abc123def456789

# 4. Trackers
GET /api/v2/torrents/trackers?hash=abc123def456789

# 5. Peers conectados
GET /api/v2/torrents/peers?hash=abc123def456789
```

### Pausar/Reanudar Torrents

```http
POST /api/v2/torrents/pause
Body: hashes=abc123def456789|def456abc123789

POST /api/v2/torrents/resume
Body: hashes=abc123def456789
```

### Eliminar Torrent

```http
POST /api/v2/torrents/delete
Body: hashes=abc123def456789&deleteFiles=false
```

### Verificar Torrent

```http
POST /api/v2/torrents/recheck
Body: hashes=abc123def456789
```

---

## 🔍 Filtros Disponibles

Los filtros para el endpoint `/api/v2/torrents/info`:

| Filtro | Descripción |
|--------|-------------|
| `all` | Todos los torrents |
| `downloading` | Descargando |
| `seeding` | En seed (subiendo) |
| `completed` | Completados |
| `paused` | Pausados |
| `active` | Activos (con actividad de red) |
| `inactive` | Inactivos (sin actividad de red) |
| `resumed` | No pausados |
| `stalled` | Estancados (downloading o uploading) |
| `stalled_uploading` | Estancados en subida |
| `stalled_downloading` | Estancados en descarga |
| `errored` | Con errores |

---

## 🔑 Notas Importantes

### 1. **Autenticación con Cookies**

Todas las peticiones autenticadas requieren la cookie `SID`:
```http
Cookie: SID=your-session-id
```

### 2. **Formato de Hashes**

Los hashes deben estar en minúsculas y ser strings hexadecimales de 40 caracteres:
```
abc123def456789012345678901234567890abcd
```

Para múltiples hashes, separar con `|` (pipe):
```
hash1|hash2|hash3
```

### 3. **Timestamps Unix**

Todos los timestamps están en formato Unix (segundos desde epoch):
```python
# Python
import time
timestamp = 1699564800  # 9 Nov 2023 16:00:00 GMT
date = time.strftime('%Y-%m-%d %H:%M:%S', time.gmtime(timestamp))
```

### 4. **Valores Especiales**

- `-1`: Sin límite / valor global
- `-2`: Sin límite específico
- `8640000`: ETA infinito (equivalente a "∞")

### 5. **Bytes vs Bits**

Todas las velocidades y tamaños están en **bytes**, no en bits:
- 1 MB/s = 8 Mbps
- Para convertir: `Mbps = (bytes_per_sec * 8) / 1,000,000`

### 6. **Estados Compuestos**

Un torrent puede tener información en múltiples campos para determinar su estado real:
- Usar `state` para el estado principal
- Verificar `progress` para saber si está completo
- Verificar `dlspeed` y `upspeed` para actividad real

### 7. **Ratio y Tiempo de Seed**

Para lógica de cleanup, considerar:
- `ratio`: Ratio de compartición (uploaded/downloaded)
- `seeding_time`: Tiempo total en seed en segundos
- `max_ratio` y `max_seeding_time`: Límites configurados

---

## 📚 Referencias

- **Documentación Oficial v5.0+:** https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-5.0)
- **Python Client API:** https://qbittorrent-api.readthedocs.io/
- **OpenAPI Demo:** https://www.qbittorrent.org/openapi-demo/
- **GitHub Repository:** https://github.com/qbittorrent/qBittorrent

---

**Última Actualización:** 2025-11-11  
**Versión de qBittorrent:** 4.5.0+  
**Mantenido por:** KeeperCheky Project
