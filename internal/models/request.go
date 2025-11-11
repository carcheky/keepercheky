package models

import (
	"time"

	"gorm.io/gorm"
)

// Request represents a media request from Jellyseerr/Overseerr with complete metadata.
type Request struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// External service information
	ServiceID int `json:"service_id" gorm:"not null;index"` // ID in Jellyseerr

	// Request details
	MediaType   string    `json:"media_type" gorm:"not null"` // "movie" or "tv"
	MediaTitle  string    `json:"media_title" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null;index"` // "pending", "approved", "available", "declined"
	RequestedBy string    `json:"requested_by"`
	RequestedAt time.Time `json:"requested_at" gorm:"index"`
	ModifiedBy  string    `json:"modified_by,omitempty"`

	// Request configuration
	Is4K              bool   `json:"is_4k"`
	ProfileID         *int   `json:"profile_id,omitempty"`
	RootFolder        string `json:"root_folder,omitempty"`
	LanguageProfileID *int   `json:"language_profile_id,omitempty"`
	IsAutoRequest     bool   `json:"is_auto_request"`

	// Media metadata from Jellyseerr
	TMDBID              int        `json:"tmdb_id" gorm:"index"`
	TVDBID              *int       `json:"tvdb_id,omitempty" gorm:"index"`
	IMDBID              string     `json:"imdb_id,omitempty" gorm:"index"`
	MediaStatus         int        `json:"media_status"` // Status in Jellyseerr media system
	MediaStatus4K       int        `json:"media_status_4k"`
	MediaAddedAt        *time.Time `json:"media_added_at,omitempty"`
	ExternalServiceSlug string     `json:"external_service_slug,omitempty"`

	// Associated media services
	RadarrID          *int `json:"radarr_id" gorm:"index"`
	SonarrID          *int `json:"sonarr_id" gorm:"index"`
	ServiceInstanceID *int `json:"service_instance_id,omitempty"` // Which Radarr/Sonarr instance

	// User information
	RequestedByUserID int    `json:"requested_by_user_id" gorm:"index"`
	RequestedByEmail  string `json:"requested_by_email,omitempty"`
	ModifiedByUserID  *int   `json:"modified_by_user_id,omitempty"`
}

// TableName specifies the table name for Request model.
func (Request) TableName() string {
	return "requests"
}
