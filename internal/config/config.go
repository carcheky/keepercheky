package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Clients  ClientsConfig  `mapstructure:"clients"`
}

type AppConfig struct {
	Environment      string   `mapstructure:"environment"`
	LogLevel         string   `mapstructure:"log_level"`
	DryRun           bool     `mapstructure:"dry_run"`
	LeavingSoonDays  int      `mapstructure:"leaving_soon_days"`
	ExclusionTags    []string `mapstructure:"exclusion_tags"`
	SchedulerEnabled bool     `mapstructure:"scheduler_enabled"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Type string `mapstructure:"type"` // sqlite or postgres
	Path string `mapstructure:"path"` // for sqlite
	Host string `mapstructure:"host"` // for postgres
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"password"`
	Name string `mapstructure:"name"`
}

type ClientsConfig struct {
	Radarr     ServiceClient `mapstructure:"radarr"`
	Sonarr     ServiceClient `mapstructure:"sonarr"`
	Jellyfin   ServiceClient `mapstructure:"jellyfin"`
	Jellyseerr ServiceClient `mapstructure:"jellyseerr"`
}

type ServiceClient struct {
	Enabled  bool   `mapstructure:"enabled"`
	URL      string `mapstructure:"url"`
	APIKey   string `mapstructure:"api_key"`
	Username string `mapstructure:"username"` // for Jellyfin
	Password string `mapstructure:"password"` // for Jellyfin
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/config")

	// Environment variables
	viper.SetEnvPrefix("KEEPERCHEKY")
	viper.AutomaticEnv()

	// Defaults
	setDefaults()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is OK, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.dry_run", true)
	viper.SetDefault("app.leaving_soon_days", 7)
	viper.SetDefault("app.scheduler_enabled", false)

	// Server defaults
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("server.host", "0.0.0.0")

	// Database defaults
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.path", "./data/keepercheky.db")
}
