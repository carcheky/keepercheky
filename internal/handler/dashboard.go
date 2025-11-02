package handler

import (
	"context"
	"time"

	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

const (
	// Dashboard simulation constants
	defaultDiskCapacityGB = 5000.0 // 5TB default for simulation
	dailyGrowthRate       = 0.01   // 1% daily growth for simulation
	historyDays           = 30     // Days of history to simulate
	trendThreshold        = 1.0    // Percentage threshold for trend detection
)

type DashboardHandler struct {
	repos       *repository.Repositories
	logger      *logger.Logger
	syncService *service.SyncService
}

func NewDashboardHandler(repos *repository.Repositories, logger *logger.Logger, syncService *service.SyncService) *DashboardHandler {
	return &DashboardHandler{
		repos:       repos,
		logger:      logger,
		syncService: syncService,
	}
}

func (h *DashboardHandler) Index(c *fiber.Ctx) error {
	return c.Render("pages/dashboard", fiber.Map{
		"Title": "Dashboard - KeeperCheky",
	}, "layouts/main")
}

func (h *DashboardHandler) Stats(c *fiber.Ctx) error {
	// Get basic stats from repository
	basicStats, err := h.repos.Media.GetStats()
	if err != nil {
		h.logger.Error("Failed to get stats", "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get statistics",
		})
	}

	// Build enhanced stats
	enhancedStats := h.buildEnhancedStats(basicStats)

	return c.JSON(enhancedStats)
}

// buildEnhancedStats creates comprehensive dashboard statistics
func (h *DashboardHandler) buildEnhancedStats(basicStats map[string]interface{}) map[string]interface{} {
	stats := make(map[string]interface{})

	// Copy basic stats
	for k, v := range basicStats {
		stats[k] = v
	}

	// Add placeholder values for features not yet implemented
	stats["to_delete"] = 0
	stats["leaving_soon"] = 0

	// Get disk usage history (simulated for last 30 days)
	stats["disk_usage_history"] = h.getDiskUsageHistory()

	// Get distribution by quality
	stats["distribution_by_quality"] = h.getDistributionByQuality()

	// Get distribution by size
	stats["distribution_by_size"] = h.getDistributionBySize()

	// Get recent activity from history
	stats["recent_activity"] = h.getRecentActivity(10)

	// Get top media by size
	stats["top_media_by_size"] = h.getTopMediaBySize(10)

	// Get never watched stats
	neverWatchedCount, neverWatchedSize := h.getNeverWatchedStats()
	stats["never_watched_count"] = neverWatchedCount
	stats["never_watched_size"] = neverWatchedSize

	// Get torrent stats
	activeTorrents, totalSeedRatio, avgSeedRatio := h.getTorrentStats()
	stats["active_torrents"] = activeTorrents
	stats["total_seed_ratio"] = totalSeedRatio
	stats["avg_seed_ratio"] = avgSeedRatio

	// Calculate trend
	trend, trendPercent := h.calculateTrend(stats["disk_usage_history"])
	stats["trend"] = trend
	stats["trend_percent"] = trendPercent

	return stats
}

// DiskUsagePoint represents disk usage at a point in time
type DiskUsagePoint struct {
	Date        time.Time `json:"date"`
	UsedGB      float64   `json:"used_gb"`
	FreeGB      float64   `json:"free_gb"`
	UsedPercent float64   `json:"used_percent"`
}

// getDiskUsageHistory simulates disk usage history for the last 30 days
// In a real implementation, this would pull from actual historical data
func (h *DashboardHandler) getDiskUsageHistory() []DiskUsagePoint {
	var history []DiskUsagePoint
	now := time.Now()

	// Get current total size from repository
	totalSize, err := h.repos.Media.GetTotalSize()
	if err != nil {
		h.logger.Error("Failed to get total size", "error", err)
		return history
	}
	currentUsedGB := float64(totalSize) / (1024 * 1024 * 1024)

	// Simulate historyDays of history with slight variations
	// In production, this should come from actual historical tracking
	for i := historyDays - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)

		// Simulate gradual growth with dailyGrowthRate
		growthFactor := 1.0 - (float64(i) * dailyGrowthRate)
		usedGB := currentUsedGB * growthFactor

		// Use defaultDiskCapacityGB for simulation
		freeGB := defaultDiskCapacityGB - usedGB
		usedPercent := (usedGB / defaultDiskCapacityGB) * 100

		history = append(history, DiskUsagePoint{
			Date:        date,
			UsedGB:      usedGB,
			FreeGB:      freeGB,
			UsedPercent: usedPercent,
		})
	}

	return history
}

// getDistributionByQuality returns media count grouped by quality
func (h *DashboardHandler) getDistributionByQuality() map[string]int {
	results, err := h.repos.Media.GetDistributionByQuality()
	if err != nil {
		h.logger.Error("Failed to get distribution by quality", "error", err)
		return make(map[string]int)
	}

	distribution := make(map[string]int)
	for _, r := range results {
		quality := r.Quality
		if quality == "" {
			quality = "Unknown"
		}
		distribution[quality] = r.Count
	}

	return distribution
}

// getDistributionBySize returns media count grouped by size ranges
func (h *DashboardHandler) getDistributionBySize() map[string]int {
	distribution, err := h.repos.Media.GetDistributionBySize()
	if err != nil {
		h.logger.Error("Failed to get distribution by size", "error", err)
		return make(map[string]int)
	}

	// Convert int64 to int for consistency
	result := make(map[string]int)
	for k, v := range distribution {
		result[k] = int(v)
	}

	return result
}

// ActivityEvent represents a recent activity event
type ActivityEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "added", "deleted", "excluded"
	Title     string    `json:"title"`
}

// getRecentActivity returns recent activity from history
func (h *DashboardHandler) getRecentActivity(limit int) []ActivityEvent {
	var events []ActivityEvent

	// Get recent history entries
	history, err := h.repos.History.GetRecent(limit)
	if err != nil {
		h.logger.Error("Failed to get recent history", "error", err)
		return events
	}

	// Convert history to activity events
	for _, entry := range history {
		eventType := "added"
		if entry.Action == "deleted" {
			eventType = "deleted"
		} else if entry.Action == "excluded" {
			eventType = "excluded"
		}

		events = append(events, ActivityEvent{
			Timestamp: entry.CreatedAt,
			Type:      eventType,
			Title:     entry.MediaTitle,
		})
	}

	return events
}

// getTopMediaBySize returns the largest media items
func (h *DashboardHandler) getTopMediaBySize(limit int) []repository.MediaSize {
	topMedia, err := h.repos.Media.GetTopMediaBySize(limit)
	if err != nil {
		h.logger.Error("Failed to get top media by size", "error", err)
		return []repository.MediaSize{}
	}
	return topMedia
}

// getNeverWatchedStats returns count and size of never watched content
func (h *DashboardHandler) getNeverWatchedStats() (int, int64) {
	count, totalSize, err := h.repos.Media.GetNeverWatchedStats()
	if err != nil {
		h.logger.Error("Failed to get never watched stats", "error", err)
		return 0, 0
	}
	return int(count), totalSize
}

// getTorrentStats returns torrent statistics
func (h *DashboardHandler) getTorrentStats() (int, float64, float64) {
	activeTorrents, totalSeedRatio, avgSeedRatio, err := h.repos.Media.GetTorrentStats()
	if err != nil {
		h.logger.Error("Failed to get torrent stats", "error", err)
		return 0, 0, 0
	}
	return int(activeTorrents), totalSeedRatio, avgSeedRatio
}

// calculateTrend analyzes disk usage history to determine growth trend
func (h *DashboardHandler) calculateTrend(historyInterface interface{}) (string, float64) {
	history, ok := historyInterface.([]DiskUsagePoint)
	if !ok || len(history) < 2 {
		return "stable", 0.0
	}

	// Compare first and last day
	first := history[0].UsedGB
	last := history[len(history)-1].UsedGB

	if first == 0 {
		return "stable", 0.0
	}

	percentChange := ((last - first) / first) * 100

	// Determine trend using trendThreshold constant
	trend := "stable"
	if percentChange > trendThreshold {
		trend = "growing"
	} else if percentChange < -trendThreshold {
		trend = "shrinking"
	}

	return trend, percentChange
}

// GetJellyseerrStats returns detailed Jellyseerr statistics.
func (h *DashboardHandler) GetJellyseerrStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestStats, err := h.syncService.GetJellyseerrRequestStats(ctx)
	if err != nil {
		h.logger.Error("Failed to get Jellyseerr stats",
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(requestStats)
}

// GetJellyseerrRequests returns recent Jellyseerr requests.
func (h *DashboardHandler) GetJellyseerrRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	requests, err := h.syncService.GetJellyseerrRequests(ctx)
	if err != nil {
		h.logger.Error("Failed to get Jellyseerr requests",
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"requests": requests,
		"count":    len(requests),
	})
}

// GetJellystatStats returns detailed Jellystat statistics for the dashboard.
func (h *DashboardHandler) GetJellystatStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get days from query params, default to 7 days for dashboard
	days := c.QueryInt("days", 7)

	stats, err := h.syncService.GetJellystatStatistics(ctx, days)
	if err != nil {
		h.logger.Error("Failed to get Jellystat stats",
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stats)
}

// GetJellystatViewsByType returns views by library type for the dashboard.
func (h *DashboardHandler) GetJellystatViewsByType(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get days from query params, default to 7 days for dashboard
	days := c.QueryInt("days", 7)

	views, err := h.syncService.GetJellystatViewsByLibraryType(ctx, days)
	if err != nil {
		h.logger.Error("Failed to get Jellystat views by type",
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(views)
}
