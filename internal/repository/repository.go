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
