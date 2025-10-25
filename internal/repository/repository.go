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
	return r.db.Delete(&models.Media{}, id).Error
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
	var stats struct {
		TotalMedia  int64
		TotalMovies int64
		TotalSeries int64
		TotalSize   int64
	}
	
	r.db.Model(&models.Media{}).Count(&stats.TotalMedia)
	r.db.Model(&models.Media{}).Where("type = ?", "movie").Count(&stats.TotalMovies)
	r.db.Model(&models.Media{}).Where("type = ?", "series").Count(&stats.TotalSeries)
	r.db.Model(&models.Media{}).Select("COALESCE(SUM(size), 0)").Row().Scan(&stats.TotalSize)
	
	return map[string]interface{}{
		"total_media":  stats.TotalMedia,
		"total_movies": stats.TotalMovies,
		"total_series": stats.TotalSeries,
		"total_size":   stats.TotalSize,
	}, nil
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
	setting := models.Settings{Key: key, Value: value}
	return r.db.Save(&setting).Error
}
