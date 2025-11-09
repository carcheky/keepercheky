# 🏗️ Arquitectura del Sistema de Progreso en Tiempo Real

## 📊 Vista General del Sistema

```
┌─────────────────────────────────────────────────────────────────────┐
│                           BROWSER (Alpine.js)                        │
│                                                                       │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │ filesPage Component                                           │ │
│  │                                                                │ │
│  │  Properties:                                                   │ │
│  │  ├─ syncing: false                                            │ │
│  │  ├─ syncMessage: ''                                           │ │
│  │  ├─ syncProgress: 0  ← NUEVO                                  │ │
│  │  └─ syncMessageType: 'info'                                   │ │
│  │                                                                │ │
│  │  Methods:                                                      │ │
│  │  └─ syncFiles()                                               │ │
│  │      ├─ EventSource('/api/sync/files')                        │ │
│  │      ├─ onmessage: Update syncProgress & syncMessage          │ │
│  │      └─ onerror: Handle errors                                │ │
│  └───────────────────────────────────────────────────────────────┘ │
│                                    ↑                                  │
│                                    │ SSE Events                       │
│                                    │ data: {step, message, percent}   │
└────────────────────────────────────┼─────────────────────────────────┘
                                     │
                                     │ HTTP/1.1 text/event-stream
                                     │
┌────────────────────────────────────┼─────────────────────────────────┐
│                                    ↓                                  │
│                      FIBER WEB SERVER (Go)                            │
│                                                                       │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │ SyncHandler                                                    │ │
│  │                                                                │ │
│  │  SyncFiles(c *fiber.Ctx) error                                │ │
│  │  ├─ Set SSE headers                                           │ │
│  │  │   ├─ Content-Type: text/event-stream                       │ │
│  │  │   ├─ Cache-Control: no-cache                               │ │
│  │  │   ├─ Connection: keep-alive                                │ │
│  │  │   └─ X-Accel-Buffering: no                                 │ │
│  │  │                                                             │ │
│  │  ├─ Create progressChan (buffered, cap=100)                   │ │
│  │  │                                                             │ │
│  │  ├─ Launch goroutine                                          │ │
│  │  │   └─ syncService.SyncAllWithProgress(ctx, progressChan)    │ │
│  │  │                                                             │ │
│  │  └─ Stream events                                             │ │
│  │      └─ for progress := range progressChan                    │ │
│  │          ├─ json.Marshal(progress)                            │ │
│  │          ├─ fmt.Fprintf(w, "data: %s\n\n", jsonData)         │ │
│  │          └─ w.Flush()                                         │ │
│  └───────────────────────────────────────────────────────────────┘ │
│                                    ↑                                  │
│                                    │ progressChan                     │
│                                    │ (chan SyncProgress)              │
│                                    │                                  │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │ SyncService                                                    │ │
│  │                                                                │ │
│  │  SyncAllWithProgress(ctx, progressChan)                       │ │
│  │  │                                                             │ │
│  │  ├─ Step 1: Clear DB (5-10%)                                  │ │
│  │  │   progressChan <- {step: "clear_db", percent: 5, ...}     │ │
│  │  │   mediaRepo.DeleteAll()                                    │ │
│  │  │   progressChan <- {step: "clear_db_complete", percent: 10}│ │
│  │  │                                                             │ │
│  │  ├─ Step 2: Invalidate Cache (15-20%)                         │ │
│  │  │   progressChan <- {step: "invalidate_cache", percent: 15} │ │
│  │  │   invalidateAllCaches(ctx, progressChan)                   │ │
│  │  │   progressChan <- {..., percent: 20}                       │ │
│  │  │                                                             │ │
│  │  ├─ Step 3: Sync Radarr (25-35%)                              │ │
│  │  │   progressChan <- {step: "sync_radarr", percent: 25}      │ │
│  │  │   media := radarrClient.GetLibrary(ctx)                    │ │
│  │  │   progressChan <- {..., percent: 35}                       │ │
│  │  │                                                             │ │
│  │  ├─ Step 4: Sync Sonarr (40-50%)                              │ │
│  │  │   progressChan <- {step: "sync_sonarr", percent: 40}      │ │
│  │  │   media := sonarrClient.GetLibrary(ctx)                    │ │
│  │  │   progressChan <- {..., percent: 50}                       │ │
│  │  │                                                             │ │
│  │  ├─ Step 5: Sync Jellyfin (55-65%)                            │ │
│  │  │   progressChan <- {step: "sync_jellyfin", percent: 55}    │ │
│  │  │   syncJellyfin(ctx, mediaMap, progressChan)                │ │
│  │  │   ├─ Sub-progress: 57% → 65% (incremental)                │ │
│  │  │   └─ progressChan <- {..., percent: 57, 60, 62, 65}       │ │
│  │  │                                                             │ │
│  │  ├─ Step 6: Enrich Torrents (70-75%)                          │ │
│  │  │   progressChan <- {step: "enrich_torrents", percent: 70}  │ │
│  │  │   enrichWithSeedingStatus(ctx, mediaMap)                   │ │
│  │  │   progressChan <- {..., percent: 75}                       │ │
│  │  │                                                             │ │
│  │  ├─ Step 7: Save to DB (80-95%)                               │ │
│  │  │   progressChan <- {step: "save_db", percent: 80}          │ │
│  │  │   for each media in mediaMap:                              │ │
│  │  │   ├─ mediaRepo.Create(media)                               │ │
│  │  │   └─ Sub-progress: 80% → 95% (incremental)                │ │
│  │  │       progressChan <- {..., percent: 82, 85, 90, 95}      │ │
│  │  │                                                             │ │
│  │  └─ Step 8: Complete (100%)                                   │ │
│  │      progressChan <- {step: "complete", percent: 100}        │ │
│  │      close(progressChan)                                       │ │
│  └───────────────────────────────────────────────────────────────┘ │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Flujo de Datos Detallado

### 1. Inicio de Sincronización

```
Usuario             Frontend            Backend              Service
  │                   │                   │                    │
  ├─[Click Sync]─────>│                   │                    │
  │                   │                   │                    │
  │                   ├─[EventSource]────>│                    │
  │                   │   /api/sync/files │                    │
  │                   │                   │                    │
  │                   │                   ├─[Create chan]─────>│
  │                   │                   │                    │
  │                   │                   ├─[Launch goroutine]>│
  │                   │                   │                    │
  │                   │<──[SSE Headers]───┤                    │
  │                   │   200 OK          │                    │
  │                   │   Content-Type:   │                    │
  │                   │   text/event-stream                    │
```

### 2. Streaming de Progreso

```
Service             Backend             Frontend           UI
  │                   │                   │                 │
  ├─[Send Progress]──>│                   │                 │
  │ {step: "clear_db",│                   │                 │
  │  percent: 5}      │                   │                 │
  │                   │                   │                 │
  │                   ├─[Marshal JSON]    │                 │
  │                   │                   │                 │
  │                   ├─[Write SSE]──────>│                 │
  │                   │ data: {...}\n\n   │                 │
  │                   │                   │                 │
  │                   │                   ├─[Parse JSON]    │
  │                   │                   │                 │
  │                   │                   ├─[Update State]──>│
  │                   │                   │ syncProgress=5  │
  │                   │                   │ syncMessage=... │
  │                   │                   │                 │
  │                   │                   │                 ├─[Render]
  │                   │                   │                 │ Barra: 5%
  │                   │                   │                 │ Msg: 🗑️...
```

### 3. Ciclo de Actualización

```
[Repetir 15+ veces con diferentes porcentajes y mensajes]

Service → Backend → Frontend → UI
  5%       5%        5%         █░░░░░░░░░░
 10%      10%       10%         ██░░░░░░░░░
 25%      25%       25%         █████░░░░░░
 50%      50%       50%         ██████████░
 75%      75%       75%         ███████████████░
100%     100%      100%         ████████████████
```

### 4. Finalización

```
Service             Backend             Frontend           UI
  │                   │                   │                 │
  ├─[Send Complete]──>│                   │                 │
  │ {step: "complete",│                   │                 │
  │  percent: 100}    │                   │                 │
  │                   │                   │                 │
  ├─[Close chan]      │                   │                 │
  │                   │                   │                 │
  │                   ├─[Write SSE]──────>│                 │
  │                   │                   │                 │
  │                   │                   ├─[EventSource.   │
  │                   │                   │  close()]       │
  │                   │                   │                 │
  │                   │                   ├─[Set syncing   │
  │                   │                   │  = false]       │
  │                   │                   │                 │
  │                   │                   ├─[Reload files]──>│
  │                   │                   │                 │
  │                   │                   │                 ├─[Hide bar]
  │                   │                   │                 │
  │                   │                   │                 ├─[Show ✅]
```

---

## 🎨 Estructura de Datos

### SyncProgress (Go)

```go
type SyncProgress struct {
    Step    string                 // Identificador de etapa
    Message string                 // Mensaje descriptivo
    Status  string                 // "processing", "success", "error"
    Percent int                    // 0-100
    Data    map[string]interface{} // Opcional, datos adicionales
}

// Ejemplo:
SyncProgress{
    Step:    "sync_radarr",
    Message: "🎬 Sincronizando películas desde Radarr...",
    Status:  "processing",
    Percent: 25,
    Data:    nil,
}
```

### SSE Event (Wire Format)

```
data: {"step":"sync_radarr","message":"🎬 Sincronizando películas desde Radarr...","status":"processing","percent":25}

```

### Frontend State (JavaScript)

```javascript
{
    syncing: true,
    syncProgress: 25,
    syncMessage: "🎬 Sincronizando películas desde Radarr...",
    syncMessageType: "info",
    syncMessageTimeout: null
}
```

---

## 📊 Distribución de Responsabilidades

### Backend (Go)
- ✅ **Calcular porcentajes**: Determina 0-100% en cada etapa
- ✅ **Generar mensajes**: Crea mensajes descriptivos
- ✅ **Ejecutar sync**: Realiza operaciones de sincronización
- ✅ **Enviar eventos**: Stream SSE al frontend
- ✅ **Manejar errores**: Reporta errores con contexto

### Frontend (JavaScript)
- ✅ **Recibir eventos**: EventSource para SSE
- ✅ **Actualizar estado**: Reactivamente con Alpine.js
- ✅ **Renderizar UI**: Barra de progreso y mensajes
- ✅ **Manejar ciclo de vida**: Abrir/cerrar conexión
- ✅ **Recargar datos**: Al completar sincronización

### UI (HTML/CSS)
- ✅ **Mostrar progreso**: Barra visual animada
- ✅ **Mostrar mensajes**: Texto contextual
- ✅ **Animaciones**: Transiciones suaves
- ✅ **Responsive**: Adapta a diferentes tamaños
- ✅ **Accesibilidad**: Mensajes para lectores de pantalla

---

## 🔐 Seguridad y Robustez

### Timeouts
```go
// Backend
ctx, cancel := context.WithTimeout(
    context.Background(), 
    10*time.Minute,  // Máximo 10 minutos
)
defer cancel()
```

### Buffer de Canal
```go
// Backend - Previene bloqueo
progressChan := make(chan SyncProgress, 100)
// Buffer de 100 permite burst de eventos sin bloquear
```

### Manejo de Errores
```go
// Backend
if err := syncRadarr(ctx); err != nil {
    progressChan <- SyncProgress{
        Step:    "sync_radarr_error",
        Message: fmt.Sprintf("⚠️ Error: %v", err),
        Status:  "error",
        Percent: 25,  // Mantiene último porcentaje
    }
    // Continúa con siguiente servicio
}
```

### Limpieza de Recursos
```javascript
// Frontend
eventSource.onerror = (error) => {
    console.error('SSE error:', error);
    eventSource.close();  // Libera conexión
    this.syncing = false;  // Resetea estado
};
```

---

## 🔄 Estados del Sistema

```
┌────────────────────────────────────────────────────┐
│                  Estado: IDLE                       │
│  ┌────────────────────────────────────────────┐   │
│  │ syncing: false                              │   │
│  │ syncProgress: 0                             │   │
│  │ syncMessage: ''                             │   │
│  │ [Botón: "Sincronizar" habilitado]          │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
            │
            │ [Click "Sincronizar"]
            ▼
┌────────────────────────────────────────────────────┐
│              Estado: CONNECTING                     │
│  ┌────────────────────────────────────────────┐   │
│  │ syncing: true                               │   │
│  │ syncProgress: 0                             │   │
│  │ syncMessage: ''                             │   │
│  │ [EventSource creado]                        │   │
│  │ [Botón: "Sincronizando..." deshabilitado]  │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
            │
            │ [Primer evento SSE]
            ▼
┌────────────────────────────────────────────────────┐
│              Estado: SYNCING                        │
│  ┌────────────────────────────────────────────┐   │
│  │ syncing: true                               │   │
│  │ syncProgress: 5-99                          │   │
│  │ syncMessage: "🎬 Sincronizando..."          │   │
│  │ [Barra de progreso visible]                 │   │
│  │ [Actualizaciones continuas via SSE]         │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
            │
            │ [Evento "complete"]
            ▼
┌────────────────────────────────────────────────────┐
│              Estado: COMPLETING                     │
│  ┌────────────────────────────────────────────┐   │
│  │ syncing: true                               │   │
│  │ syncProgress: 100                           │   │
│  │ syncMessage: "✅ Completado"                │   │
│  │ [EventSource cerrado]                       │   │
│  │ [Recargando datos...]                       │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
            │
            │ [Datos recargados]
            ▼
┌────────────────────────────────────────────────────┐
│              Estado: IDLE (again)                   │
│  ┌────────────────────────────────────────────┐   │
│  │ syncing: false                              │   │
│  │ syncProgress: 0                             │   │
│  │ syncMessage: ''                             │   │
│  │ [Barra oculta]                              │   │
│  │ [Mensaje de éxito temporal]                 │   │
│  └────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────┘
```

---

## 🎯 Puntos Clave de la Arquitectura

### 1. Desacoplamiento
- **Backend**: Calcula y envía progreso
- **Frontend**: Recibe y renderiza
- **Sin dependencia directa**: Comunicación via SSE

### 2. Escalabilidad
- **Buffer de canal**: 100 eventos sin bloqueo
- **Goroutines**: Procesamiento asíncrono
- **Múltiples clientes**: Cada uno con su stream SSE

### 3. Robustez
- **Timeouts**: Previene cuelgues infinitos
- **Error handling**: Errores específicos por etapa
- **Resource cleanup**: Cierre apropiado de conexiones

### 4. Performance
- **CSS animations**: No JavaScript para animaciones
- **Transiciones suaves**: 500ms ease-out
- **Minimal re-renders**: Solo actualiza cuando cambia

### 5. User Experience
- **Feedback continuo**: 15+ actualizaciones
- **Información clara**: Mensajes descriptivos
- **Progreso visible**: Barra de 0% a 100%
- **Sub-progreso**: Operaciones largas con incrementos

---

**Arquitectura diseñada para ser escalable, robusta y proporcionar la mejor experiencia de usuario posible.**
