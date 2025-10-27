package service

import (
	"testing"
	"time"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHealthAnalyzer_DetectOrphanDownloads(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewHealthAnalyzer(logger, 180)

	tests := []struct {
		name           string
		file           *models.MediaFileInfo
		expectedStatus HealthStatus
		expectedSeverity string
	}{
		{
			name: "Archivo en qBT pero no en Jellyfin ni Radarr - orphan_download",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				InQBittorrent: true,
				InJellyfin:    false,
				InRadarr:      false,
				InSonarr:      false,
			},
			expectedStatus:   HealthStatusOrphanDownload,
			expectedSeverity: "warning",
		},
		{
			name: "Archivo en qBT y Radarr pero no en Jellyfin - orphan_download",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				InQBittorrent: true,
				InJellyfin:    false,
				InRadarr:      true,
				InSonarr:      false,
			},
			expectedStatus:   HealthStatusOrphanDownload,
			expectedSeverity: "warning",
		},
		{
			name: "Archivo en todos los servicios - ok",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				InQBittorrent: true,
				InJellyfin:    true,
				InRadarr:      true,
				InSonarr:      false,
			},
			expectedStatus:   HealthStatusOK,
			expectedSeverity: "info",
		},
		{
			name: "Archivo solo en Jellyfin - ok",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				InQBittorrent: false,
				InJellyfin:    true,
				InRadarr:      false,
				InSonarr:      false,
			},
			expectedStatus:   HealthStatusOK,
			expectedSeverity: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := analyzer.AnalyzeFile(tt.file)
			assert.Equal(t, tt.expectedStatus, report.Status)
			assert.Equal(t, tt.expectedSeverity, report.Severity)
			
			if tt.expectedStatus == HealthStatusOrphanDownload {
				// Just check that there are suggestions, not the specific text
				// since different scenarios have different suggestions
				assert.NotEmpty(t, report.Suggestions)
				assert.NotEmpty(t, report.Actions)
			}
		})
	}
}

func TestHealthAnalyzer_DetectOnlyHardlinks(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewHealthAnalyzer(logger, 180)

	tests := []struct {
		name           string
		file           *models.MediaFileInfo
		expectedStatus HealthStatus
	}{
		{
			name: "Hardlink con torrent activo - ok",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				IsHardlink:    true,
				HardlinkPaths: "/jellyfin/movies/file.mkv|/downloads/file.mkv",
				InQBittorrent: true,
				InJellyfin:    true,
			},
			expectedStatus: HealthStatusOK,
		},
		{
			name: "Hardlink sin torrent - only_hardlink",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				IsHardlink:    true,
				HardlinkPaths: "/jellyfin/movies/file.mkv|/downloads/file.mkv",
				InQBittorrent: false,
				InJellyfin:    true,
			},
			expectedStatus: HealthStatusOnlyHardlink,
		},
		{
			name: "Un solo archivo (no hardlink) - ok",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				IsHardlink:    false,
				InQBittorrent: false,
				InJellyfin:    true,
			},
			expectedStatus: HealthStatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := analyzer.AnalyzeFile(tt.file)
			assert.Equal(t, tt.expectedStatus, report.Status)
			
			if tt.expectedStatus == HealthStatusOnlyHardlink {
				assert.Contains(t, report.Actions, "cleanup_hardlink")
				assert.Equal(t, "info", report.Severity)
			}
		})
	}
}

func TestHealthAnalyzer_DetectDeadTorrents(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewHealthAnalyzer(logger, 180)

	tests := []struct {
		name           string
		file           *models.MediaFileInfo
		expectedStatus HealthStatus
		expectedSeverity string
	}{
		{
			name: "Estado 'error' - dead_torrent",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				InQBittorrent: true,
				TorrentState:  "error",
				InJellyfin:    true,
			},
			expectedStatus:   HealthStatusDeadTorrent,
			expectedSeverity: "critical",
		},
		{
			name: "Estado 'missingFiles' - dead_torrent",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				InQBittorrent: true,
				TorrentState:  "missingFiles",
				InJellyfin:    true,
			},
			expectedStatus:   HealthStatusDeadTorrent,
			expectedSeverity: "critical",
		},
		{
			name: "Estado 'uploading' - ok",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				InQBittorrent: true,
				TorrentState:  "uploading",
				InJellyfin:    true,
			},
			expectedStatus:   HealthStatusOK,
			expectedSeverity: "info",
		},
		{
			name: "Estado 'pausedUP' - ok (por ahora)",
			file: &models.MediaFileInfo{
				Title:         "Test Movie",
				InQBittorrent: true,
				TorrentState:  "pausedUP",
				InJellyfin:    true,
			},
			expectedStatus:   HealthStatusOK,
			expectedSeverity: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := analyzer.AnalyzeFile(tt.file)
			assert.Equal(t, tt.expectedStatus, report.Status)
			assert.Equal(t, tt.expectedSeverity, report.Severity)
			
			if tt.expectedStatus == HealthStatusDeadTorrent {
				assert.NotEmpty(t, report.Suggestions)
				assert.Contains(t, report.Actions, "remove_torrent")
			}
		})
	}
}

func TestHealthAnalyzer_DetectNeverWatched(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewHealthAnalyzer(logger, 180) // 180 days threshold

	tests := []struct {
		name           string
		file           *models.MediaFileInfo
		addedDate      time.Time
		expectedStatus HealthStatus
	}{
		{
			name: "Nunca visto y > 180 días - never_watched",
			file: &models.MediaFileInfo{
				Title:          "Test Movie",
				InJellyfin:     true,
				HasBeenWatched: false,
			},
			addedDate:      time.Now().AddDate(0, 0, -200), // 200 days ago
			expectedStatus: HealthStatusNeverWatched,
		},
		{
			name: "Nunca visto y < 180 días - ok",
			file: &models.MediaFileInfo{
				Title:          "Test Movie",
				InJellyfin:     true,
				HasBeenWatched: false,
			},
			addedDate:      time.Now().AddDate(0, 0, -100), // 100 days ago
			expectedStatus: HealthStatusOK,
		},
		{
			name: "Visto recientemente - ok",
			file: &models.MediaFileInfo{
				Title:          "Test Movie",
				InJellyfin:     true,
				HasBeenWatched: true,
			},
			addedDate:      time.Now().AddDate(0, 0, -200), // 200 days ago
			expectedStatus: HealthStatusOK,
		},
		{
			name: "Visto hace mucho - ok (ya fue visto)",
			file: &models.MediaFileInfo{
				Title:          "Test Movie",
				InJellyfin:     true,
				HasBeenWatched: true,
			},
			addedDate:      time.Now().AddDate(0, 0, -300), // 300 days ago
			expectedStatus: HealthStatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Using the helper method that accepts addedDate
			report := analyzer.AnalyzeFileWithDate(tt.file, tt.addedDate)
			assert.Equal(t, tt.expectedStatus, report.Status)
			
			if tt.expectedStatus == HealthStatusNeverWatched {
				assert.Equal(t, "info", report.Severity)
				assert.NotEmpty(t, report.Suggestions)
			}
		})
	}
}

func TestHealthAnalyzer_GetHealthSummary(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewHealthAnalyzer(logger, 180)

	// Create test files with different states
	files := []*models.MediaFileInfo{
		{
			Title:         "Orphan 1",
			InQBittorrent: true,
			InJellyfin:    false,
			InRadarr:      false,
		},
		{
			Title:         "Orphan 2",
			InQBittorrent: true,
			InJellyfin:    false,
			InRadarr:      false,
		},
		{
			Title:         "Healthy",
			InQBittorrent: true,
			InJellyfin:    true,
			InRadarr:      true,
		},
		{
			Title:         "Dead Torrent",
			InQBittorrent: true,
			InJellyfin:    true,
			TorrentState:  "error",
		},
		{
			Title:         "Only Hardlink",
			IsHardlink:    true,
			InQBittorrent: false,
			InJellyfin:    true,
		},
	}

	summary := analyzer.GetHealthSummary(files)

	assert.Equal(t, 2, summary["orphan_downloads"])
	assert.Equal(t, 1, summary["healthy"])
	assert.Equal(t, 1, summary["dead_torrents"])
	assert.Equal(t, 1, summary["only_hardlinks"])
	assert.Equal(t, 0, summary["never_watched"])
}

func TestHealthAnalyzer_DefaultThreshold(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewHealthAnalyzer(logger, 0) // Should default to 180
	
	// Verify default was set
	assert.Equal(t, 180, analyzer.neverWatchedDays)
}

func TestHealthAnalyzer_CustomThreshold(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewHealthAnalyzer(logger, 90) // Custom 90 days
	
	assert.Equal(t, 90, analyzer.neverWatchedDays)
	
	// Test with custom threshold
	file := &models.MediaFileInfo{
		Title:          "Test Movie",
		InJellyfin:     true,
		HasBeenWatched: false,
	}
	
	// File older than 90 days should be flagged
	report := analyzer.AnalyzeFileWithDate(file, time.Now().AddDate(0, 0, -100))
	assert.Equal(t, HealthStatusNeverWatched, report.Status)
	
	// File newer than 90 days should be OK
	report = analyzer.AnalyzeFileWithDate(file, time.Now().AddDate(0, 0, -50))
	assert.Equal(t, HealthStatusOK, report.Status)
}

func TestHealthAnalyzer_PriorityOrder(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewHealthAnalyzer(logger, 180)

	// File with multiple issues - dead torrent should take priority
	file := &models.MediaFileInfo{
		Title:          "Test Movie",
		InQBittorrent:  true,
		InJellyfin:     false, // Also orphan
		TorrentState:   "error", // Dead torrent
		HasBeenWatched: false,
	}

	report := analyzer.AnalyzeFile(file)
	
	// Dead torrent should be detected first (highest priority)
	assert.Equal(t, HealthStatusDeadTorrent, report.Status)
	assert.Equal(t, "critical", report.Severity)
}
