package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app" yaml:"app"`
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Clients  ClientsConfig  `mapstructure:"clients" yaml:"clients"`
	Cleanup  CleanupConfig  `mapstructure:"cleanup" yaml:"cleanup"`
}

type AppConfig struct {
	Environment      string `mapstructure:"environment" yaml:"environment"`
	LogLevel         string `mapstructure:"log_level" yaml:"log_level"`
	SchedulerEnabled bool   `mapstructure:"scheduler_enabled" yaml:"scheduler_enabled"`
}

type CleanupConfig struct {
	DryRun            bool     `mapstructure:"dry_run" yaml:"dry_run"`
	DaysToKeep        int      `mapstructure:"days_to_keep" yaml:"days_to_keep"`
	LeavingSoonDays   int      `mapstructure:"leaving_soon_days" yaml:"leaving_soon_days"`
	ExclusionTags     []string `mapstructure:"exclusion_tags" yaml:"exclusion_tags"`
	DeleteUnmonitored bool     `mapstructure:"delete_unmonitored" yaml:"delete_unmonitored"`
}

type ServerConfig struct {
	Port string `mapstructure:"port" yaml:"port"`
	Host string `mapstructure:"host" yaml:"host"`
}

type DatabaseConfig struct {
	Type string `mapstructure:"type" yaml:"type"` // sqlite or postgres
	Path string `mapstructure:"path" yaml:"path"` // for sqlite
	Host string `mapstructure:"host" yaml:"host"` // for postgres
	Port string `mapstructure:"port" yaml:"port"`
	User string `mapstructure:"user" yaml:"user"`
	Pass string `mapstructure:"password" yaml:"password"`
	Name string `mapstructure:"name" yaml:"name"`
}

type ClientsConfig struct {
	Radarr      ServiceClient     `mapstructure:"radarr" yaml:"radarr"`
	Sonarr      ServiceClient     `mapstructure:"sonarr" yaml:"sonarr"`
	Jellyfin    ServiceClient     `mapstructure:"jellyfin" yaml:"jellyfin"`
	Jellyseerr  ServiceClient     `mapstructure:"jellyseerr" yaml:"jellyseerr"`
	QBittorrent QBittorrentClient `mapstructure:"qbittorrent" yaml:"qbittorrent"`
}

type ServiceClient struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
	URL     string `mapstructure:"url" yaml:"url"`
	APIKey  string `mapstructure:"api_key" yaml:"api_key"`
}

type QBittorrentClient struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	URL      string `mapstructure:"url" yaml:"url"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/config") // Primary location (mounted volume)
	viper.AddConfigPath("./config")    // Fallback for local development
	viper.AddConfigPath(".")

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

// Save writes the current configuration to a YAML file
func Save(cfg *Config) error {
	viper.Set("app", cfg.App)
	viper.Set("server", cfg.Server)
	viper.Set("database", cfg.Database)
	viper.Set("clients", cfg.Clients)
	viper.Set("cleanup", cfg.Cleanup)

	// Write to /app/config/config.yaml (mounted volume in container)
	configPath := "/app/config/config.yaml"
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file to %s: %w", configPath, err)
	}

	return nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.scheduler_enabled", false)

	// Server defaults
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("server.host", "0.0.0.0")

	// Database defaults
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.path", "./data/keepercheky.db")

	// Cleanup defaults
	viper.SetDefault("cleanup.dry_run", true)
	viper.SetDefault("cleanup.days_to_keep", 90)
	viper.SetDefault("cleanup.leaving_soon_days", 7)
	viper.SetDefault("cleanup.delete_unmonitored", false)

	// Clients defaults (all disabled by default)
	viper.SetDefault("clients.radarr.enabled", false)
	viper.SetDefault("clients.sonarr.enabled", false)
	viper.SetDefault("clients.jellyfin.enabled", false)
	viper.SetDefault("clients.jellyseerr.enabled", false)
	viper.SetDefault("clients.qbittorrent.enabled", false)
}
