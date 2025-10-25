package repository

import (
	"github.com/carcheky/keepercheky/internal/models"
	"gorm.io/gorm"
)

// MediaRepository handles media data access.
type MediaRepository struct {
	db *gorm.DB
}

// NewMediaRepository creates a new media repository.
func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

// GetAll retrieves all media items.
func (r *MediaRepository) GetAll() ([]models.Media, error) {
	var media []models.Media
	result := r.db.Find(&media)
	return media, result.Error
}

// GetByID retrieves a media item by ID.
func (r *MediaRepository) GetByID(id uint) (*models.Media, error) {
	var media models.Media
	result := r.db.First(&media, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &media, nil
}

// Create creates a new media item.
func (r *MediaRepository) Create(media *models.Media) error {
	return r.db.Create(media).Error
}

// Update updates a media item.
func (r *MediaRepository) Update(media *models.Media) error {
	return r.db.Save(media).Error
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

// Delete deletes a media item.
func (r *MediaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Media{}, id).Error
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
