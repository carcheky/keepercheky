# 🔍 Mejoras de Búsqueda y Ordenamiento - Files View

## Resumen

Se han implementado mejoras significativas en la navegación de la vista de lista de archivos, específicamente para el Issue #92 (sub-issues #92.2 y #92.4).

## ✅ Funcionalidades Implementadas

### 1. Búsqueda del Lado del Servidor (Issue #92.2)

**Backend (`/api/files`):**
- Nuevo parámetro query: `search`
- Búsqueda en múltiples campos:
  - `title` - Título del archivo
  - `file_path` - Ruta completa del archivo
  - `torrent_hash` - Hash del torrent asociado
- Uso de búsqueda SQL optimizada con `LIKE`
- Compatible con paginación y filtros por tab

**Frontend:**
- Input de búsqueda con debounce de 500ms
- Indicador visual cuando hay búsqueda activa
- Búsqueda integrada con el servidor (no solo filtrado de cliente)
- Botón para limpiar búsqueda rápidamente
- Compatibilidad con tabs (OK, Atención, Críticos, etc.)

### 2. Ordenamiento Mejorado (Issue #92.4)

**Nuevos Campos de Ordenamiento:**
- `file_path` - Ruta del archivo (por defecto)
- `title` - Título
- `size` - Tamaño
- `type` - Tipo (película/serie)
- `created_at` - ⭐ **NUEVO**: Fecha de agregado
- `updated_at` - ⭐ **NUEVO**: Última actualización
- `seed_ratio` - ⭐ **NUEVO**: Ratio de seeds
- `torrent_state` - ⭐ **NUEVO**: Estado del torrent

**Dirección de Ordenamiento:**
- Ascendente (A → Z, 0 → 9)
- Descendente (Z → A, 9 → 0)

### 3. UI/UX Mejorada

**Sección de Filtros Colapsable:**
```
🔍 Búsqueda y Filtros [activos]
└── 🔎 Buscar archivos
    └── Input con placeholder: "Buscar por título, ruta o hash de torrent..."
    └── Botón "Limpiar"
└── 📊 Ordenar por
    └── Selector con 8 opciones
└── 🔄 Orden
    └── Selector ASC/DESC
└── 🔄 Restablecer todo
```

**Indicadores Visuales:**
- Badge "activos" cuando hay filtros aplicados
- Indicador de número de resultados encontrados
- Información contextual sobre funcionamiento

## 📊 Paginación (Ya Existía)

La paginación avanzada ya estaba implementada:
- ✅ Botones Primera/Última página
- ✅ Páginas numeradas con ellipsis (...)
- ✅ Información de rango (Mostrando X - Y de Z)
- ✅ Selector de items por página (10, 25, 50, 100)

## 🔧 Cambios Técnicos

### Backend

**Archivo:** `internal/handler/files.go`

```go
// Nuevos parámetros
search := c.Query("search", "")

// Aplicar búsqueda
if search != "" {
    searchPattern := "%" + search + "%"
    query = query.Where(
        "title LIKE ? OR file_path LIKE ? OR torrent_hash LIKE ?",
        searchPattern, searchPattern, searchPattern,
    )
}
```

### Frontend

**Archivo:** `web/templates/pages/files.html`

```javascript
// Nuevas propiedades
searchQuery: '',
sortBy: 'file_path',
sortOrder: 'asc',

// Nuevos métodos
handleSearchChange() {
    this.currentPage = 1;
    this.loadFiles();
},

handleSortChange() {
    this.currentPage = 1;
    this.loadFiles();
},

clearSearch() {
    this.searchQuery = '';
    this.currentPage = 1;
    this.loadFiles();
},

resetFiltersAndSort() {
    this.searchQuery = '';
    this.sortBy = 'file_path';
    this.sortOrder = 'asc';
    this.currentPage = 1;
    this.loadFiles();
}
```

## 🎯 Casos de Uso

### Ejemplo 1: Buscar película específica
1. Click en "🔍 Búsqueda y Filtros"
2. Escribir título en el input de búsqueda
3. Esperar 500ms (debounce automático)
4. Ver resultados filtrados en tiempo real

### Ejemplo 2: Encontrar archivos grandes
1. Click en "🔍 Búsqueda y Filtros"
2. Ordenar por: "Tamaño"
3. Orden: "Descendente"
4. Ver archivos más grandes primero

### Ejemplo 3: Ver archivos recientes
1. Click en "🔍 Búsqueda y Filtros"
2. Ordenar por: "Fecha de agregado"
3. Orden: "Descendente"
4. Ver archivos más recientes primero

### Ejemplo 4: Combinar búsqueda con filtros
1. Click en tab "Atención" (orphan downloads)
2. Click en "🔍 Búsqueda y Filtros"
3. Buscar: "1080p"
4. Ver solo archivos huérfanos en 1080p

## 🚀 Performance

- **Búsqueda servidor-side**: Más rápida que filtrado de cliente
- **Paginación**: Solo carga items necesarios (no todos los resultados)
- **Debounce**: Evita búsquedas excesivas mientras el usuario escribe
- **Cache de contadores**: Los totales por categoría se cachean

## 📝 Notas Técnicas

### Validación de Seguridad
- Campos de ordenamiento validados en whitelist
- Protección contra SQL injection con placeholders
- Validación de parámetros de paginación

### Compatibilidad
- Compatible con todos los tabs existentes
- Compatible con vista organizada
- Compatible con selección masiva
- No afecta funcionalidad existente

## 🐛 Testing

### Tests Ejecutados
```bash
go test ./internal/handler/... -v
```

**Resultado:** ✅ Todos los tests pasan

### Validación Manual Pendiente
- [ ] Verificar búsqueda en diferentes campos
- [ ] Verificar ordenamiento por cada campo
- [ ] Verificar combinación búsqueda + tabs
- [ ] Verificar combinación búsqueda + ordenamiento
- [ ] Verificar performance con muchos archivos
- [ ] Tomar screenshots del antes/después

## 📚 Referencias

- **Issue Principal:** #92 - Navegación vista lista confusa
- **Sub-Issues Implementados:**
  - #92.2 - Agregar búsqueda de texto
  - #92.4 - Optimizar ordenamiento
- **Sub-Issue Ya Existente:**
  - #92.1 - Mejorar paginación (ya estaba implementado)

## 🎨 Capturas de Pantalla

_Pendiente: Agregar screenshots del UI con los cambios_

### Antes
- Búsqueda solo en cliente (página actual)
- Ordenamiento limitado
- Sin indicadores visuales de filtros activos

### Después
- ✅ Búsqueda en servidor (todos los resultados)
- ✅ Ordenamiento extendido (8 opciones)
- ✅ Indicadores visuales claros
- ✅ Botón para restablecer todo
