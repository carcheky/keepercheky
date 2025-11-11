# Card de Películas - Especificación

## Formato: UNA LÍNEA HORIZONTAL por película

Cada película se muestra como **una sola fila horizontal** con todos los elementos alineados.

## Estructura

```
[ ] Título de la película - Año - Tamaño - ICONOS - Carátula
```

### Elementos (de izquierda a derecha):

1. **[ ] Checkbox** - Selección múltiple
2. **Título de la película** - Nombre completo
3. **Año de lanzamiento** - Entre paréntesis o separado con guión
4. **Tamaño** - Formato legible (GB, MB, etc.)
5. **ICONOS** - Servicios e información:
   - 🎞️ **Jellyfin** - Con popup/desplegable mostrando todos los datos disponibles:
     - ID
     - Sinopsis
     - Calificación
     - Reproducciones
     - Última reproducción
     - Fecha agregada
   - 🔽 **qBittorrent** (futuro) - Estado del torrent
   - 📡 **Radarr** (futuro) - Info de Radarr
   - 📺 **Sonarr** (futuro) - Para series
6. **Carátula de Jellyfin** - Poster pequeño al final de la línea

## Ejemplo Visual

```
┌───────────────────────────────────────────────────────────────────────────┐
│ [ ] The Matrix (1999) - 2.0 GB - 🎞️ 🔽 📡 - [📷 poster]                 │
│ [ ] Inception (2010) - 2.5 GB - 🎞️ - [📷 poster]                         │
│ [ ] Pulp Fiction (1994) - 2.3 GB - 🎞️ 🔽 - [📷 poster]                  │
└───────────────────────────────────────────────────────────────────────────┘
```

## Layout HTML/CSS

- **Display:** `flex` con `items-center`
- **Ancho:** 100% (responsive)
- **Gap:** Espaciado consistente entre elementos
- **Altura:** Fija por fila (~50-60px)
- **Poster:** Tamaño pequeño (40-50px altura, aspect ratio 2:3)

## Interacciones

- Click en checkbox → Selecciona/deselecciona
- Click en icono Jellyfin (🎞️) → Abre popup con metadatos
- Hover sobre fila → Resalta visualmente

---

**IMPORTANTE:** NO usar grid de tarjetas verticales. Cada película es UNA FILA HORIZONTAL en una lista.
