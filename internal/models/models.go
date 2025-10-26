package models

import (
	"time"

	"gorm.io/gorm"
)

// Media represents a media item (movie or TV show)
type Media struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Title       string     `json:"title" gorm:"not null;index"`
	Type        string     `json:"type" gorm:"not null;index"` // "movie" or "series"
	PosterURL   string     `json:"poster_url"`
	FilePath    string     `json:"file_path" gorm:"not null"`
	Size        int64      `json:"size"`
	AddedDate   time.Time  `json:"added_date" gorm:"index"`
	LastWatched *time.Time `json:"last_watched"`

	// Series specific
	EpisodeCount     int `json:"episode_count"`
	SeasonCount      int `json:"season_count"`
	EpisodeFileCount int `json:"episode_file_count"` // Downloaded episodes

	// Torrent status
	IsSeeding   bool    `json:"is_seeding" gorm:"default:false"`
	TorrentHash string  `json:"torrent_hash" gorm:"index"` // Hash from qBittorrent
	SeedRatio   float64 `json:"seed_ratio"`

	// Service IDs
	RadarrID     *int    `json:"radarr_id" gorm:"index"`
	SonarrID     *int    `json:"sonarr_id" gorm:"index"`
	JellyfinID   *string `json:"jellyfin_id" gorm:"index"`
	JellyseerrID *int    `json:"jellyseerr_id" gorm:"index"`

	// Metadata
	Tags     []string `json:"tags" gorm:"type:json"`
	Quality  string   `json:"quality"`
	Excluded bool     `json:"excluded" gorm:"default:false;index"`

	// Relationships
	History []History `json:"history,omitempty" gorm:"foreignKey:MediaID"`
}

func (Media) TableName() string {
	return "media"
}

// Schedule represents a cleanup schedule
type Schedule struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name        string     `json:"name" gorm:"not null"`
	Enabled     bool       `json:"enabled" gorm:"default:true"`
	CronExpr    string     `json:"cron_expr" gorm:"not null"`
	LastRun     *time.Time `json:"last_run"`
	NextRun     *time.Time `json:"next_run"`
	Description string     `json:"description"`

	// Rules (stored as JSON for now, can be normalized later)
	Rules string `json:"rules" gorm:"type:json"`
}

func (Schedule) TableName() string {
	return "schedules"
}

// History represents an action log entry
type History struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`

	MediaID    uint   `json:"media_id" gorm:"index"`
	MediaTitle string `json:"media_title"`
	Action     string `json:"action" gorm:"not null"` // "deleted", "excluded", "unmonitored"
	Status     string `json:"status" gorm:"not null"` // "success", "failed"
	Message    string `json:"message"`
	DryRun     bool   `json:"dry_run" gorm:"default:false"`
}

func (History) TableName() string {
	return "history"
}

// Settings represents application settings
type Settings struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UpdatedAt time.Time `json:"updated_at"`

	Key   string `json:"key" gorm:"uniqueIndex;not null"`
	Value string `json:"value" gorm:"type:text"`
}

func (Settings) TableName() string {
	return "settings"
}

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&Media{},
		&Schedule{},
		&History{},
		&Settings{},
	)
}
