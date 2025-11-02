package handler

import (
	"context"
	"time"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
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
	
	// Get current total size
	var totalSize int64
	h.repos.Media.GetDB().Model(&models.Media{}).Select("COALESCE(SUM(size), 0)").Scan(&totalSize)
	currentUsedGB := float64(totalSize) / (1024 * 1024 * 1024)
	
	// Simulate 30 days of history with slight variations
	// In production, this should come from actual historical tracking
	for i := 29; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		
		// Simulate gradual growth with some randomness
		growthFactor := 1.0 - (float64(i) * 0.01) // Grow ~1% per day
		usedGB := currentUsedGB * growthFactor
		
		// Assume disk capacity of 5TB for simulation
		diskCapacityGB := 5000.0
		freeGB := diskCapacityGB - usedGB
		usedPercent := (usedGB / diskCapacityGB) * 100
		
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
	type QualityCount struct {
		Quality string
		Count   int
	}
	
	var results []QualityCount
	h.repos.Media.GetDB().
		Model(&models.Media{}).
		Select("COALESCE(quality, 'Unknown') as quality, COUNT(*) as count").
		Group("quality").
		Order("count DESC").
		Scan(&results)
	
	distribution := make(map[string]int)
	for _, r := range results {
		if r.Quality == "" {
			r.Quality = "Unknown"
		}
		distribution[r.Quality] = r.Count
	}
	
	return distribution
}

// getDistributionBySize returns media count grouped by size ranges
func (h *DashboardHandler) getDistributionBySize() map[string]int {
	distribution := make(map[string]int)
	
	// Small (< 5 GB)
	var smallCount int64
	h.repos.Media.GetDB().Model(&models.Media{}).
		Where("size < ?", int64(5*1024*1024*1024)).
		Count(&smallCount)
	distribution["< 5 GB"] = int(smallCount)
	
	// Medium (5-20 GB)
	var mediumCount int64
	h.repos.Media.GetDB().Model(&models.Media{}).
		Where("size >= ? AND size < ?", int64(5*1024*1024*1024), int64(20*1024*1024*1024)).
		Count(&mediumCount)
	distribution["5-20 GB"] = int(mediumCount)
	
	// Large (20-50 GB)
	var largeCount int64
	h.repos.Media.GetDB().Model(&models.Media{}).
		Where("size >= ? AND size < ?", int64(20*1024*1024*1024), int64(50*1024*1024*1024)).
		Count(&largeCount)
	distribution["20-50 GB"] = int(largeCount)
	
	// XLarge (> 50 GB)
	var xlargeCount int64
	h.repos.Media.GetDB().Model(&models.Media{}).
		Where("size >= ?", int64(50*1024*1024*1024)).
		Count(&xlargeCount)
	distribution["> 50 GB"] = int(xlargeCount)
	
	return distribution
}

// ActivityEvent represents a recent activity event
type ActivityEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "added", "deleted", "watched"
	Title     string    `json:"title"`
	Size      int64     `json:"size,omitempty"`
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
	for _, h := range history {
		eventType := "added"
		if h.Action == "deleted" {
			eventType = "deleted"
		} else if h.Action == "excluded" {
			eventType = "excluded"
		}
		
		events = append(events, ActivityEvent{
			Timestamp: h.CreatedAt,
			Type:      eventType,
			Title:     h.MediaTitle,
			Size:      0, // Size not tracked in history currently
		})
	}
	
	return events
}

// MediaSize represents a media item with just ID, title, and size for top lists
type MediaSize struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Size  int64  `json:"size"`
	Type  string `json:"type"`
}

// getTopMediaBySize returns the largest media items
func (h *DashboardHandler) getTopMediaBySize(limit int) []MediaSize {
	var topMedia []MediaSize
	
	h.repos.Media.GetDB().
		Model(&models.Media{}).
		Select("id, title, size, type").
		Order("size DESC").
		Limit(limit).
		Scan(&topMedia)
	
	return topMedia
}

// getNeverWatchedStats returns count and size of never watched content
func (h *DashboardHandler) getNeverWatchedStats() (int, int64) {
	var count int64
	var totalSize int64
	
	h.repos.Media.GetDB().
		Model(&models.Media{}).
		Where("last_watched IS NULL").
		Count(&count)
	
	h.repos.Media.GetDB().
		Model(&models.Media{}).
		Where("last_watched IS NULL").
		Select("COALESCE(SUM(size), 0)").
		Scan(&totalSize)
	
	return int(count), totalSize
}

// getTorrentStats returns torrent statistics
func (h *DashboardHandler) getTorrentStats() (int, float64, float64) {
	var activeTorrents int64
	var totalSeedRatio float64
	var avgSeedRatio float64
	
	// Count active torrents
	h.repos.Media.GetDB().
		Model(&models.Media{}).
		Where("is_seeding = ?", true).
		Count(&activeTorrents)
	
	// Calculate total and average seed ratio
	h.repos.Media.GetDB().
		Model(&models.Media{}).
		Where("is_seeding = ?", true).
		Select("COALESCE(SUM(seed_ratio), 0)").
		Scan(&totalSeedRatio)
	
	if activeTorrents > 0 {
		avgSeedRatio = totalSeedRatio / float64(activeTorrents)
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
	
	// Determine trend
	trend := "stable"
	if percentChange > 1.0 {
		trend = "growing"
	} else if percentChange < -1.0 {
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
