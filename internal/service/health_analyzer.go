package service

import (
	"time"

	"github.com/carcheky/keepercheky/internal/models"
	"go.uber.org/zap"
)

// HealthStatus represents the health status of a media file
type HealthStatus string

const (
	HealthStatusOK              HealthStatus = "ok"
	HealthStatusOrphanDownload  HealthStatus = "orphan_download"
	HealthStatusOnlyHardlink    HealthStatus = "only_hardlink"
	HealthStatusDeadTorrent     HealthStatus = "dead_torrent"
	HealthStatusNeverWatched    HealthStatus = "never_watched"
)

// HealthReport represents a health analysis report for a media file
type HealthReport struct {
	Status      HealthStatus `json:"status"`
	Severity    string       `json:"severity"` // "info", "warning", "critical"
	Message     string       `json:"message"`
	Suggestions []string     `json:"suggestions"`
	Actions     []string     `json:"actions"`
}

// HealthAnalyzer analyzes media files for health issues
type HealthAnalyzer struct {
	logger              *zap.Logger
	neverWatchedDays    int // Threshold for never watched files
}

// NewHealthAnalyzer creates a new health analyzer
func NewHealthAnalyzer(logger *zap.Logger, neverWatchedDays int) *HealthAnalyzer {
	if neverWatchedDays == 0 {
		neverWatchedDays = 180 // Default: 180 days
	}
	return &HealthAnalyzer{
		logger:              logger,
		neverWatchedDays:    neverWatchedDays,
	}
}

// AnalyzeFile analyzes a single media file and returns a health report
func (a *HealthAnalyzer) AnalyzeFile(file *models.MediaFileInfo) HealthReport {
	// Check for dead torrents first (highest priority)
	if file.InQBittorrent && a.isDeadTorrent(file) {
		return HealthReport{
			Status:      HealthStatusDeadTorrent,
			Severity:    "critical",
			Message:     "Torrent is in error state or missing files",
			Suggestions: []string{"Remove torrent from qBittorrent", "Re-download if needed"},
			Actions:     []string{"remove_torrent", "delete_file"},
		}
	}

	// Check for orphan downloads (in qBittorrent but not in media servers)
	if file.InQBittorrent && !file.InJellyfin && !file.InRadarr && !file.InSonarr {
		return HealthReport{
			Status:      HealthStatusOrphanDownload,
			Severity:    "warning",
			Message:     "File downloaded but not imported to any media server",
			Suggestions: []string{"Importar a Radarr", "Importar a Sonarr", "Delete if not needed"},
			Actions:     []string{"import_to_radarr", "import_to_sonarr", "delete_file"},
		}
	}

	// Check if it's in qBittorrent, in Radarr/Sonarr but not in Jellyfin
	if file.InQBittorrent && (file.InRadarr || file.InSonarr) && !file.InJellyfin {
		return HealthReport{
			Status:      HealthStatusOrphanDownload,
			Severity:    "warning",
			Message:     "File in download manager and *arr but not in Jellyfin library",
			Suggestions: []string{"Trigger Jellyfin library scan", "Check file permissions"},
			Actions:     []string{"scan_jellyfin"},
		}
	}

	// Check for hardlinks without active torrent
	if file.IsHardlink && !file.InQBittorrent && file.InJellyfin {
		return HealthReport{
			Status:      HealthStatusOnlyHardlink,
			Severity:    "info",
			Message:     "Hardlink exists but torrent has been removed",
			Suggestions: []string{"Clean up hardlink if space is needed"},
			Actions:     []string{"cleanup_hardlink"},
		}
	}

	// Check for never watched files (if in Jellyfin)
	if file.InJellyfin && !file.HasBeenWatched {
		// Calculate age based on file mod time (we'd need to add AddedDate to MediaFileInfo)
		// For now, we'll use a simple check if the file is old enough
		// This would need proper implementation with actual dates
		if a.isOldFile(file) {
			return HealthReport{
				Status:      HealthStatusNeverWatched,
				Severity:    "info",
				Message:     "File has never been watched and is old",
				Suggestions: []string{"Consider removing if not interested"},
				Actions:     []string{"delete_from_jellyfin", "delete_file"},
			}
		}
	}

	// Everything looks OK
	return HealthReport{
		Status:      HealthStatusOK,
		Severity:    "info",
		Message:     "File is healthy",
		Suggestions: []string{},
		Actions:     []string{},
	}
}

// isDeadTorrent checks if a torrent is in a dead state
func (a *HealthAnalyzer) isDeadTorrent(file *models.MediaFileInfo) bool {
	// Dead torrent states
	deadStates := map[string]bool{
		"error":        true,
		"missingFiles": true,
	}

	if deadStates[file.TorrentState] {
		return true
	}

	// PausedUP without seeders could also be considered dead
	// This would need additional info from qBittorrent about seeders
	// For now, we'll keep it simple
	
	return false
}

// isOldFile checks if a file is considered old (for never watched analysis)
// This is a placeholder - in real implementation, we'd need actual added date
func (a *HealthAnalyzer) isOldFile(file *models.MediaFileInfo) bool {
	// This is a simplified check
	// In production, we'd need to compare against actual file added date
	// For now, we'll return false to avoid false positives
	// The proper implementation would be:
	// return time.Since(file.AddedDate).Hours() > float64(a.neverWatchedDays * 24)
	return false
}

// GetHealthSummary returns a summary of health statuses for a list of files
func (a *HealthAnalyzer) GetHealthSummary(files []*models.MediaFileInfo) map[string]int {
	summary := map[string]int{
		"healthy":           0,
		"orphan_downloads":  0,
		"only_hardlinks":    0,
		"dead_torrents":     0,
		"never_watched":     0,
	}

	for _, file := range files {
		report := a.AnalyzeFile(file)
		switch report.Status {
		case HealthStatusOK:
			summary["healthy"]++
		case HealthStatusOrphanDownload:
			summary["orphan_downloads"]++
		case HealthStatusOnlyHardlink:
			summary["only_hardlinks"]++
		case HealthStatusDeadTorrent:
			summary["dead_torrents"]++
		case HealthStatusNeverWatched:
			summary["never_watched"]++
		}
	}

	return summary
}

// AnalyzeFileWithDate is a helper for testing that allows passing a custom added date
func (a *HealthAnalyzer) AnalyzeFileWithDate(file *models.MediaFileInfo, addedDate time.Time) HealthReport {
	// Check for dead torrents first
	if file.InQBittorrent && a.isDeadTorrent(file) {
		return HealthReport{
			Status:      HealthStatusDeadTorrent,
			Severity:    "critical",
			Message:     "Torrent is in error state or missing files",
			Suggestions: []string{"Remove torrent from qBittorrent", "Re-download if needed"},
			Actions:     []string{"remove_torrent", "delete_file"},
		}
	}

	// Check for orphan downloads
	if file.InQBittorrent && !file.InJellyfin && !file.InRadarr && !file.InSonarr {
		return HealthReport{
			Status:      HealthStatusOrphanDownload,
			Severity:    "warning",
			Message:     "File downloaded but not imported to any media server",
			Suggestions: []string{"Importar a Radarr", "Importar a Sonarr", "Delete if not needed"},
			Actions:     []string{"import_to_radarr", "import_to_sonarr", "delete_file"},
		}
	}

	if file.InQBittorrent && (file.InRadarr || file.InSonarr) && !file.InJellyfin {
		return HealthReport{
			Status:      HealthStatusOrphanDownload,
			Severity:    "warning",
			Message:     "File in download manager and *arr but not in Jellyfin library",
			Suggestions: []string{"Trigger Jellyfin library scan", "Check file permissions"},
			Actions:     []string{"scan_jellyfin"},
		}
	}

	// Check for hardlinks
	if file.IsHardlink && !file.InQBittorrent && file.InJellyfin {
		return HealthReport{
			Status:      HealthStatusOnlyHardlink,
			Severity:    "info",
			Message:     "Hardlink exists but torrent has been removed",
			Suggestions: []string{"Clean up hardlink if space is needed"},
			Actions:     []string{"cleanup_hardlink"},
		}
	}

	// Check for never watched with custom date
	if file.InJellyfin && !file.HasBeenWatched {
		daysSinceAdded := time.Since(addedDate).Hours() / 24
		if daysSinceAdded > float64(a.neverWatchedDays) {
			return HealthReport{
				Status:      HealthStatusNeverWatched,
				Severity:    "info",
				Message:     "File has never been watched and is old",
				Suggestions: []string{"Consider removing if not interested"},
				Actions:     []string{"delete_from_jellyfin", "delete_file"},
			}
		}
	}

	return HealthReport{
		Status:      HealthStatusOK,
		Severity:    "info",
		Message:     "File is healthy",
		Suggestions: []string{},
		Actions:     []string{},
	}
}
