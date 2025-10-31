# Vista Organizada de Archivos

## Descripción

La nueva "Vista Organizada" permite visualizar los archivos de media de forma jerárquica, organizados por series/temporadas o películas, facilitando la gestión de grandes bibliotecas de contenido.

## Características

### 📺 Series
- **Organización jerárquica de 3 niveles**: Serie → Temporada → Episodios
- **Desplegables en cada nivel**: Puedes expandir/colapsar series, temporadas y episodios
- **Información de tamaño**: Muestra el tamaño total en cada nivel (serie completa, temporada, episodio individual)
- **Detección automática**: Reconoce patrones de nombres como:
  - `S01E01` / `s01e01` (Breaking Bad S01E01.mkv)
  - `1x01` / `1X01` (The Office 1x01.mp4)
- **Especiales**: Los episodios con temporada 0 se muestran como "Especiales"

### 🎬 Películas
- **Agrupación por título**: Todas las versiones de una película en un solo lugar
- **Versión principal**: Prioriza automáticamente la versión en Jellyfin
- **Otras versiones**: Desplegable para mostrar versiones alternativas (diferentes calidades, ubicaciones)
- **Información detallada**: Calidad y tamaño para cada versión

## Cómo usar

### Cambiar entre vistas

1. En la página de Files, encontrarás el toggle de vista en la esquina superior derecha
2. Haz clic en "Organizado" para cambiar a la vista jerárquica
3. Haz clic en "Lista" para volver a la vista clásica

### Navegar en la vista organizada

#### Series
1. Haz clic en el nombre de una serie para expandir/colapsar sus temporadas
2. Haz clic en una temporada para ver sus episodios
3. Cada episodio muestra:
   - Número de episodio (E01, E02, etc.)
   - Título (si está disponible)
   - Calidad del archivo
   - Tamaño del archivo

#### Películas
1. Cada película muestra su versión principal
2. Si hay múltiples versiones, aparece el botón "Otras versiones"
3. Haz clic para ver todas las versiones alternativas con sus detalles

### Filtros

La vista organizada respeta los mismos filtros que la vista de lista:
- 🟢 OK (Archivos saludables)
- 🟡 Atención (Huérfanos en descargas)
- 🔴 Críticos (Torrents muertos)
- 🔗 Hardlinks (Solo hardlinks)
- 👁️ Sin Ver (Sin reproducir)

## API

### Endpoint

```
GET /api/files/organized
```

### Parámetros de consulta

| Parámetro | Tipo | Descripción | Por defecto |
|-----------|------|-------------|-------------|
| `page` | int | Número de página | 1 |
| `perPage` | int | Elementos por página | 25 |
| `tab` | string | Filtro de categoría (healthy, attention, critical, etc.) | "" |
| `type` | string | Tipo de media (series, movie) | "" |

### Respuesta

```json
{
  "data": {
    "series": [
      {
        "series_title": "Breaking Bad",
        "total_size": 52428800000,
        "season_count": 5,
        "episode_count": 62,
        "seasons": [
          {
            "season_number": 1,
            "episode_count": 7,
            "total_size": 7516192768,
            "episodes": [
              {
                "episode_number": 1,
                "title": "Pilot",
                "file_path": "/media/tv/Breaking.Bad.S01E01.720p.mkv",
                "size": 1073741824,
                "quality": "720p",
                "versions": []
              }
            ]
          }
        ],
        "poster_url": "https://...",
        "primary_path": "/media/tv/Breaking Bad/",
        "metadata": {
          "in_jellyfin": true,
          "in_sonarr": true,
          "in_qbittorrent": true,
          "sonarr_id": 123,
          "jellyfin_id": "abc-123"
        }
      }
    ],
    "movies": [
      {
        "title": "Inception",
        "total_size": 10737418240,
        "primary_file": {
          "id": 1,
          "title": "Inception",
          "file_path": "/media/movies/Inception.2010.1080p.mkv",
          "size": 8589934592,
          "quality": "1080p"
        },
        "other_versions": [
          {
            "file_path": "/downloads/Inception.2010.2160p.mkv",
            "size": 2147483648,
            "quality": "2160p"
          }
        ],
        "metadata": {
          "in_jellyfin": true,
          "in_radarr": true,
          "radarr_id": 456
        }
      }
    ]
  },
  "page": 1,
  "perPage": 25,
  "totalCount": 150
}
```

## Patrones de nombres soportados

### Series
- `Serie.S01E01.720p.mkv` → Serie: "Serie", Temporada: 1, Episodio: 1
- `The.Office.1x05.mp4` → Serie: "The Office", Temporada: 1, Episodio: 5
- `Game.of.Thrones.s05e08.1080p.mp4` → Serie: "Game of Thrones", Temporada: 5, Episodio: 8

### Películas
- Cualquier archivo que no contenga patrones de serie/episodio
- Se agrupa por título extraído del nombre del archivo

## Notas técnicas

- La detección de series se basa en el nombre del archivo, no en metadatos
- Si el título ya está en la base de datos, se usa ese título
- La versión principal se prioriza según este orden:
  1. Archivos en Jellyfin
  2. Archivos en Radarr/Sonarr
  3. Otros archivos
- Los tamaños se calculan sumando todos los archivos relacionados
