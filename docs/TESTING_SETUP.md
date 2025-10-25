# 🧪 Configuración del Entorno de Pruebas

Este documento describe cómo configurar y usar el entorno de desarrollo con servicios de prueba.

## 📦 Servicios Incluidos

El `docker-compose.dev.yml` incluye los siguientes servicios de prueba:

### KeeperCheky (Puerto 8000)
- **URL**: http://localhost:8000
- Dashboard de gestión
- Configuración pre-cargada con todos los servicios

### Radarr (Puerto 7878)
- **URL**: http://localhost:7878
- **Gestión**: Películas
- **API Key**: `test-radarr-api-key-12345` (configurar manualmente)

### Sonarr (Puerto 8989)
- **URL**: http://localhost:8989
- **Gestión**: Series de TV
- **API Key**: `test-sonarr-api-key-12345` (configurar manualmente)

### Jellyfin (Puerto 8096)
- **URL**: http://localhost:8096
- **Servidor**: Media streaming
- **API Key**: `test-jellyfin-api-key-12345` (configurar manualmente)

### Jellyseerr (Puerto 5055)
- **URL**: http://localhost:5055
- **Gestión**: Peticiones de media
- **API Key**: `test-jellyseerr-api-key-12345` (configurar manualmente)

## 🚀 Inicio Rápido

### 1. Levantar todos los servicios

```bash
docker compose -f docker-compose.dev.yml up -d
```

### 2. Configuración Inicial de Servicios

#### Radarr (http://localhost:7878)
1. Accede a la UI web
2. Completa el wizard de configuración inicial
3. Ve a **Settings → General**
4. Copia la **API Key** generada
5. **IMPORTANTE**: Reemplaza el API key en tu configuración de KeeperCheky

#### Sonarr (http://localhost:8989)
1. Accede a la UI web
2. Completa el wizard de configuración inicial
3. Ve a **Settings → General**
4. Copia la **API Key** generada
5. **IMPORTANTE**: Reemplaza el API key en tu configuración de KeeperCheky

#### Jellyfin (http://localhost:8096)
1. Accede a la UI web
2. Completa el wizard de configuración inicial
3. Crea un usuario administrador
4. Ve a **Dashboard → API Keys**
5. Crea una nueva API Key para "KeeperCheky"
6. **IMPORTANTE**: Reemplaza el API key en tu configuración de KeeperCheky

#### Jellyseerr (http://localhost:5055)
1. Accede a la UI web
2. Configura la conexión con Jellyfin
3. Ve a **Settings → General → API Key**
4. Copia la **API Key** generada
5. **IMPORTANTE**: Reemplaza el API key en tu configuración de KeeperCheky

### 3. Actualizar API Keys en KeeperCheky

Hay dos formas de configurar las API keys reales:

#### Opción A: Variables de Entorno (Recomendado)

Edita `docker-compose.dev.yml` y actualiza las variables:

```yaml
environment:
  - KEEPERCHEKY_SERVICES_RADARR_APIKEY=tu-api-key-real-de-radarr
  - KEEPERCHEKY_SERVICES_SONARR_APIKEY=tu-api-key-real-de-sonarr
  - KEEPERCHEKY_SERVICES_JELLYFIN_APIKEY=tu-api-key-real-de-jellyfin
  - KEEPERCHEKY_SERVICES_JELLYSEERR_APIKEY=tu-api-key-real-de-jellyseerr
```

Luego reinicia el contenedor:

```bash
docker compose -f docker-compose.dev.yml restart app
```

#### Opción B: Configuración Web

1. Ve a http://localhost:8000/settings
2. Ingresa las API keys reales en cada servicio
3. Haz clic en "Test Connection" para verificar
4. Guarda la configuración

### 4. Probar las Conexiones

```bash
# Probar Radarr
curl http://localhost:8000/api/config/test?service=radarr

# Probar Sonarr
curl http://localhost:8000/api/config/test?service=sonarr

# Probar Jellyfin
curl http://localhost:8000/api/config/test?service=jellyfin

# Probar Jellyseerr
curl http://localhost:8000/api/config/test?service=jellyseerr
```

### 5. Sincronizar Media

Una vez configuradas las API keys:

1. Ve a http://localhost:8000/
2. Haz clic en **"Sync Now"**
3. Espera a que se complete la sincronización
4. Verás las estadísticas actualizadas en el dashboard

## 📊 Añadir Media de Prueba

### Radarr - Añadir Películas

1. Ve a http://localhost:7878
2. **Add Movies → Add New Movie**
3. Busca una película (ej: "The Matrix")
4. Configura:
   - Root Folder: `/media/movies`
   - Quality Profile: Any
   - Monitor: Yes
5. Añade la película

### Sonarr - Añadir Series

1. Ve a http://localhost:8989
2. **Series → Add New**
3. Busca una serie (ej: "Breaking Bad")
4. Configura:
   - Root Folder: `/media/tv`
   - Quality Profile: Any
   - Monitor: All Episodes
5. Añade la serie

### Jellyfin - Escanear Biblioteca

1. Ve a http://localhost:8096
2. **Dashboard → Libraries**
3. Añade nueva biblioteca:
   - Type: Movies o Shows
   - Folder: `/media/movies` o `/media/tv`
4. Escanea la biblioteca

## 🔍 Verificar Datos

### Ver Media en KeeperCheky

```bash
# Ver todas las películas/series
curl http://localhost:8000/api/media | jq .

# Ver estadísticas
curl http://localhost:8000/api/stats | jq .
```

## 🛠️ Comandos Útiles

### Ver logs de todos los servicios

```bash
docker compose -f docker-compose.dev.yml logs -f
```

### Ver logs de un servicio específico

```bash
docker compose -f docker-compose.dev.yml logs -f app
docker compose -f docker-compose.dev.yml logs -f radarr
docker compose -f docker-compose.dev.yml logs -f sonarr
```

### Reiniciar un servicio

```bash
docker compose -f docker-compose.dev.yml restart app
```

### Detener todos los servicios

```bash
docker compose -f docker-compose.dev.yml down
```

### Detener y eliminar volúmenes (reset completo)

```bash
docker compose -f docker-compose.dev.yml down -v
```

## 🐛 Troubleshooting

### Los servicios no se conectan

1. Verifica que todos los contenedores estén corriendo:
   ```bash
   docker compose -f docker-compose.dev.yml ps
   ```

2. Verifica las API keys en la configuración web

3. Verifica los logs:
   ```bash
   docker compose -f docker-compose.dev.yml logs app
   ```

### No aparecen medias en KeeperCheky

1. Asegúrate de haber configurado las API keys correctas
2. Añade al menos una película en Radarr o serie en Sonarr
3. Ejecuta la sincronización manual desde el dashboard
4. Verifica los logs para errores de API

### Error "connection refused"

1. Verifica que todos los servicios estén en la misma red (`keepercheky-net`)
2. Usa los nombres de contenedor (no `localhost`) en las URLs internas
3. Ejemplo: `http://radarr:7878` (no `http://localhost:7878`)

## 📝 Notas

- Los API keys de prueba (`test-*-api-key-12345`) son placeholders
- **DEBES** reemplazarlos con los API keys reales de cada servicio
- Los datos persisten en volúmenes Docker entre reinicios
- Para un reset completo, usa `docker compose down -v`

## 🔗 Enlaces Rápidos

- **KeeperCheky**: http://localhost:8000
- **Radarr**: http://localhost:7878
- **Sonarr**: http://localhost:8989
- **Jellyfin**: http://localhost:8096
- **Jellyseerr**: http://localhost:5055

---

**Última actualización**: 25 de Octubre de 2025
