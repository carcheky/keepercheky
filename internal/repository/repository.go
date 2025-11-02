package repository

import (
	"github.com/carcheky/keepercheky/internal/models"
	"gorm.io/gorm"
)

type Repositories struct {
	Media    *MediaRepository
	Schedule *ScheduleRepository
	History  *HistoryRepository
	Settings *SettingsRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Media:    NewMediaRepository(db),
		Schedule: NewScheduleRepository(db),
		History:  NewHistoryRepository(db),
		Settings: NewSettingsRepository(db),
	}
}

// MediaRepository handles media data access
type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *MediaRepository) GetAll() ([]models.Media, error) {
	var media []models.Media
	result := r.db.Find(&media)
	return media, result.Error
}

func (r *MediaRepository) GetByID(id uint) (*models.Media, error) {
	var media models.Media
	result := r.db.First(&media, id)
	return &media, result.Error
}

func (r *MediaRepository) Create(media *models.Media) error {
	return r.db.Create(media).Error
}

func (r *MediaRepository) Update(media *models.Media) error {
	return r.db.Save(media).Error
}

func (r *MediaRepository) Delete(id uint) error {
	// Use Unscoped() to permanently delete (hard delete) instead of soft delete
	return r.db.Unscoped().Delete(&models.Media{}, id).Error
}

// DeleteAll permanently deletes ALL media from the database.
// This is used for full sync to ensure clean state.
// WARNING: This is a destructive operation!
func (r *MediaRepository) DeleteAll() error {
	// Use Unscoped() to permanently delete all records (hard delete)
	return r.db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Media{}).Error
}

// CreateOrUpdate creates or updates a media item based on external IDs.
func (r *MediaRepository) CreateOrUpdate(media *models.Media) error {
	// Try to find existing media by external IDs
	var existing models.Media
	var found bool

	if media.RadarrID != nil {
		result := r.db.Where("radarr_id = ?", *media.RadarrID).First(&existing)
		found = result.Error == nil
	} else if media.SonarrID != nil {
		result := r.db.Where("sonarr_id = ?", *media.SonarrID).First(&existing)
		found = result.Error == nil
	} else if media.JellyfinID != nil {
		result := r.db.Where("jellyfin_id = ?", *media.JellyfinID).First(&existing)
		found = result.Error == nil
	}

	if found {
		// Update existing
		media.ID = existing.ID
		media.CreatedAt = existing.CreatedAt
		return r.db.Save(media).Error
	}

	// Create new
	return r.db.Create(media).Error
}

// GetStats retrieves media statistics.
func (r *MediaRepository) GetStats() (map[string]interface{}, error) {
	var stats models.GlobalStats

	// Count total media
	var totalMedia int64
	r.db.Model(&models.Media{}).Count(&totalMedia)
	stats.TotalMedia = totalMedia

	// Count movies
	var totalMovies int64
	r.db.Model(&models.Media{}).Where("type = ?", "movie").Count(&totalMovies)
	stats.TotalMovies = totalMovies

	// Count series
	var totalSeries int64
	r.db.Model(&models.Media{}).Where("type = ?", "series").Count(&totalSeries)
	stats.TotalSeries = totalSeries

	_ = r.db.Model(&models.Media{}).Select("COALESCE(SUM(size), 0)").Row().Scan(&stats.TotalSize)
	_ = r.db.Model(&models.Media{}).Where("type = ?", "series").Select("COALESCE(SUM(episode_count), 0)").Row().Scan(&stats.TotalEpisodes)
	_ = r.db.Model(&models.Media{}).Where("type = ?", "series").Select("COALESCE(SUM(episode_file_count), 0)").Row().Scan(&stats.TotalEpisodesDownload)

	return map[string]interface{}{
		"total_media":             stats.TotalMedia,
		"total_movies":            stats.TotalMovies,
		"total_series":            stats.TotalSeries,
		"total_size":              stats.TotalSize,
		"total_episodes":          stats.TotalEpisodes,
		"total_episodes_download": stats.TotalEpisodesDownload,
	}, nil
}

// GetPaginated retrieves media with pagination and optional filters.
func (r *MediaRepository) GetPaginated(page, pageSize int, mediaType string) ([]models.Media, int64, error) {
	var media []models.Media
	var total int64

	query := r.db.Model(&models.Media{})

	// Apply type filter if provided
	if mediaType != "" && mediaType != "all" {
		query = query.Where("type = ?", mediaType)
	}

	// Get total count
	query.Count(&total)

	// Apply pagination
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("added_date DESC").Find(&media)

	return media, total, result.Error
}

// MediaFilters holds all available filter options
type MediaFilters struct {
	Type              string // movie, series, torrent, all
	Status            string // active, excluded, seeding, all
	Search            string // search by title
	Service           string // radarr, sonarr, jellyfin, orphan, all
	SizeRange         string // small, medium, large, xlarge, all
	AddedDate         string // week, month, quarter, older, all
	SeedRatio         string // low, medium, high, all
	Quality           string // search by quality
	EpisodeCompletion string // complete, incomplete, empty, all
}

// GetPaginatedWithFilters retrieves media with advanced filtering and pagination
func (r *MediaRepository) GetPaginatedWithFilters(page, pageSize int, filters MediaFilters) ([]models.Media, int64, error) {
	var media []models.Media
	var total int64

	query := r.db.Model(&models.Media{})

	// Apply type filter
	if filters.Type != "" && filters.Type != "all" {
		query = query.Where("type = ?", filters.Type)
	}

	// Apply status filter
	switch filters.Status {
	case "excluded":
		query = query.Where("excluded = ?", true)
	case "active":
		query = query.Where("excluded = ?", false)
	case "seeding":
		query = query.Where("is_seeding = ?", true)
	}

	// Apply search filter
	if filters.Search != "" {
		query = query.Where("title LIKE ?", "%"+filters.Search+"%")
	}

	// Apply service filter
	switch filters.Service {
	case "radarr":
		query = query.Where("radarr_id IS NOT NULL")
	case "sonarr":
		query = query.Where("sonarr_id IS NOT NULL")
	case "jellyfin":
		query = query.Where("jellyfin_id IS NOT NULL")
	case "orphan":
		query = query.Where("type = ?", "torrent")
	}

	// Apply size range filter (in bytes)
	switch filters.SizeRange {
	case "small": // < 5 GB
		query = query.Where("size < ?", int64(5*1024*1024*1024))
	case "medium": // 5-20 GB
		query = query.Where("size >= ? AND size < ?", int64(5*1024*1024*1024), int64(20*1024*1024*1024))
	case "large": // 20-50 GB
		query = query.Where("size >= ? AND size < ?", int64(20*1024*1024*1024), int64(50*1024*1024*1024))
	case "xlarge": // > 50 GB
		query = query.Where("size >= ?", int64(50*1024*1024*1024))
	}

	// Apply date filter
	now := r.db.NowFunc()
	switch filters.AddedDate {
	case "week": // Last 7 days
		query = query.Where("added_date >= ?", now.AddDate(0, 0, -7))
	case "month": // Last 30 days
		query = query.Where("added_date >= ?", now.AddDate(0, 0, -30))
	case "quarter": // Last 90 days
		query = query.Where("added_date >= ?", now.AddDate(0, 0, -90))
	case "older": // Older than 90 days
		query = query.Where("added_date < ?", now.AddDate(0, 0, -90))
	}

	// Apply seed ratio filter
	switch filters.SeedRatio {
	case "low": // < 1.0
		query = query.Where("seed_ratio < ?", 1.0)
	case "medium": // 1.0-2.0
		query = query.Where("seed_ratio >= ? AND seed_ratio < ?", 1.0, 2.0)
	case "high": // > 2.0
		query = query.Where("seed_ratio >= ?", 2.0)
	}

	// Apply quality filter
	if filters.Quality != "" {
		query = query.Where("quality LIKE ?", "%"+filters.Quality+"%")
	}

	// Apply episode completion filter (for series)
	switch filters.EpisodeCompletion {
	case "complete": // Has all episodes
		query = query.Where("type = ? AND episode_count > 0 AND episode_count = episode_file_count", "series")
	case "incomplete": // Missing episodes
		query = query.Where("type = ? AND episode_count > 0 AND episode_count > episode_file_count", "series")
	case "empty": // No episode info
		query = query.Where("type = ? AND (episode_count IS NULL OR episode_count = 0)", "series")
	}

	// Get total count
	query.Count(&total)

	// Apply pagination and ordering
	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("added_date DESC").Find(&media)

	return media, total, result.Error
}

// Search searches media by title.
func (r *MediaRepository) Search(query string) ([]models.Media, error) {
	var media []models.Media
	result := r.db.Where("title LIKE ?", "%"+query+"%").Order("title ASC").Find(&media)
	return media, result.Error
}

// GetByType retrieves all media of a specific type.
func (r *MediaRepository) GetByType(mediaType string) ([]models.Media, error) {
	var media []models.Media
	result := r.db.Where("type = ?", mediaType).Order("added_date DESC").Find(&media)
	return media, result.Error
}

// GetExcluded retrieves all excluded media.
func (r *MediaRepository) GetExcluded() ([]models.Media, error) {
	var media []models.Media
	result := r.db.Where("excluded = ?", true).Order("title ASC").Find(&media)
	return media, result.Error
}

// GetOlderThan retrieves media added before a specific date.
func (r *MediaRepository) GetOlderThan(days int) ([]models.Media, error) {
	var media []models.Media
	cutoffDate := r.db.NowFunc().AddDate(0, 0, -days)
	result := r.db.Where("added_date < ?", cutoffDate).Order("added_date ASC").Find(&media)
	return media, result.Error
}

// SetExcluded marks media as excluded or not excluded.
func (r *MediaRepository) SetExcluded(id uint, excluded bool) error {
	return r.db.Model(&models.Media{}).Where("id = ?", id).Update("excluded", excluded).Error
}

// CountByType counts media by type.
func (r *MediaRepository) CountByType() (map[string]int64, error) {
	counts := make(map[string]int64)

	var movies, series int64
	r.db.Model(&models.Media{}).Where("type = ?", "movie").Count(&movies)
	r.db.Model(&models.Media{}).Where("type = ?", "series").Count(&series)

	counts["movies"] = movies
	counts["series"] = series

	return counts, nil
}

// GetTotalSize returns the total size of all media
func (r *MediaRepository) GetTotalSize() (int64, error) {
	var totalSize int64
	err := r.db.Model(&models.Media{}).Select("COALESCE(SUM(size), 0)").Row().Scan(&totalSize)
	return totalSize, err
}

// QualityCount represents quality distribution data
type QualityCount struct {
	Quality string
	Count   int
}

// GetDistributionByQuality returns media count grouped by quality
func (r *MediaRepository) GetDistributionByQuality() ([]QualityCount, error) {
	var results []QualityCount
	err := r.db.
		Model(&models.Media{}).
		Select("COALESCE(quality, 'Unknown') as quality, COUNT(*) as count").
		Group("quality").
		Order("count DESC").
		Scan(&results).Error

	return results, err
}

// SizeDistribution represents size range distribution
type SizeDistribution struct {
	Range string
	Count int64
}

// GetDistributionBySize returns media count grouped by size ranges
func (r *MediaRepository) GetDistributionBySize() (map[string]int64, error) {
	distribution := make(map[string]int64)

	// Small (< 5 GB)
	var smallCount int64
	if err := r.db.Model(&models.Media{}).
		Where("size < ?", int64(5*1024*1024*1024)).
		Count(&smallCount).Error; err != nil {
		return nil, err
	}
	distribution["< 5 GB"] = smallCount

	// Medium (5-20 GB)
	var mediumCount int64
	if err := r.db.Model(&models.Media{}).
		Where("size >= ? AND size < ?", int64(5*1024*1024*1024), int64(20*1024*1024*1024)).
		Count(&mediumCount).Error; err != nil {
		return nil, err
	}
	distribution["5-20 GB"] = mediumCount

	// Large (20-50 GB)
	var largeCount int64
	if err := r.db.Model(&models.Media{}).
		Where("size >= ? AND size < ?", int64(20*1024*1024*1024), int64(50*1024*1024*1024)).
		Count(&largeCount).Error; err != nil {
		return nil, err
	}
	distribution["20-50 GB"] = largeCount

	// XLarge (> 50 GB)
	var xlargeCount int64
	if err := r.db.Model(&models.Media{}).
		Where("size >= ?", int64(50*1024*1024*1024)).
		Count(&xlargeCount).Error; err != nil {
		return nil, err
	}
	distribution["> 50 GB"] = xlargeCount

	return distribution, nil
}

// MediaSize represents a media item with size info for top lists
type MediaSize struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Size  int64  `json:"size"`
	Type  string `json:"type"`
}

// GetTopMediaBySize returns the largest media items
func (r *MediaRepository) GetTopMediaBySize(limit int) ([]MediaSize, error) {
	var topMedia []MediaSize
	err := r.db.
		Model(&models.Media{}).
		Select("id, title, size, type").
		Order("size DESC").
		Limit(limit).
		Scan(&topMedia).Error

	return topMedia, err
}

// GetNeverWatchedStats returns count and size of never watched content
func (r *MediaRepository) GetNeverWatchedStats() (int64, int64, error) {
	var count int64
	var totalSize int64

	if err := r.db.Model(&models.Media{}).
		Where("last_watched IS NULL").
		Count(&count).Error; err != nil {
		return 0, 0, err
	}

	if err := r.db.Model(&models.Media{}).
		Where("last_watched IS NULL").
		Select("COALESCE(SUM(size), 0)").
		Scan(&totalSize).Error; err != nil {
		return 0, 0, err
	}

	return count, totalSize, nil
}

// GetTorrentStats returns torrent statistics
func (r *MediaRepository) GetTorrentStats() (int64, float64, float64, error) {
	var activeTorrents int64
	var totalSeedRatio float64
	var avgSeedRatio float64

	// Count active torrents
	if err := r.db.Model(&models.Media{}).
		Where("is_seeding = ?", true).
		Count(&activeTorrents).Error; err != nil {
		return 0, 0, 0, err
	}

	// Calculate total and average seed ratio
	if err := r.db.Model(&models.Media{}).
		Where("is_seeding = ?", true).
		Select("COALESCE(SUM(seed_ratio), 0)").
		Scan(&totalSeedRatio).Error; err != nil {
		return 0, 0, 0, err
	}

	if activeTorrents > 0 {
		avgSeedRatio = totalSeedRatio / float64(activeTorrents)
	}

	return activeTorrents, totalSeedRatio, avgSeedRatio, nil
}

// ScheduleRepository handles schedule data access
type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) GetAll() ([]models.Schedule, error) {
	var schedules []models.Schedule
	result := r.db.Find(&schedules)
	return schedules, result.Error
}

func (r *ScheduleRepository) GetEnabled() ([]models.Schedule, error) {
	var schedules []models.Schedule
	result := r.db.Where("enabled = ?", true).Find(&schedules)
	return schedules, result.Error
}

// HistoryRepository handles history data access
type HistoryRepository struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) Create(history *models.History) error {
	return r.db.Create(history).Error
}

func (r *HistoryRepository) GetRecent(limit int) ([]models.History, error) {
	var history []models.History
	result := r.db.Order("created_at DESC").Limit(limit).Find(&history)
	return history, result.Error
}

// SettingsRepository handles settings data access
type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(key string) (string, error) {
	var setting models.Settings
	result := r.db.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		return "", result.Error
	}
	return setting.Value, nil
}

func (r *SettingsRepository) Set(key, value string) error {
	// Update if exists, create if not (upsert)
	return r.db.
		Where("key = ?", key).
		Assign(models.Settings{Value: value}).
		FirstOrCreate(&models.Settings{Key: key}).
		Error
}
