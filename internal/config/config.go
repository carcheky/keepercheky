package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// EnvSourceMap tracks which config values come from environment variables
type EnvSourceMap struct {
	Radarr struct {
		Enabled bool `json:"enabled"`
		APIKey  bool `json:"api_key"`
		URL     bool `json:"url"`
	} `json:"radarr"`
	Sonarr struct {
		Enabled bool `json:"enabled"`
		APIKey  bool `json:"api_key"`
		URL     bool `json:"url"`
	} `json:"sonarr"`
	Jellyfin struct {
		Enabled bool `json:"enabled"`
		APIKey  bool `json:"api_key"`
		URL     bool `json:"url"`
	} `json:"jellyfin"`
	Jellyseerr struct {
		Enabled bool `json:"enabled"`
		APIKey  bool `json:"api_key"`
		URL     bool `json:"url"`
	} `json:"jellyseerr"`
	QBittorrent struct {
		Enabled  bool `json:"enabled"`
		Username bool `json:"username"`
		Password bool `json:"password"`
		URL      bool `json:"url"`
	} `json:"qbittorrent"`
}

// GetEnvSourceMap returns a map indicating which config values come from environment variables
func GetEnvSourceMap() *EnvSourceMap {
	envMap := &EnvSourceMap{}

	// Check Radarr
	envMap.Radarr.Enabled = os.Getenv("KEEPERCHEKY_CLIENTS_RADARR_ENABLED") != ""
	envMap.Radarr.URL = os.Getenv("KEEPERCHEKY_CLIENTS_RADARR_URL") != ""
	envMap.Radarr.APIKey = os.Getenv("KEEPERCHEKY_CLIENTS_RADARR_API_KEY") != ""

	// Check Sonarr
	envMap.Sonarr.Enabled = os.Getenv("KEEPERCHEKY_CLIENTS_SONARR_ENABLED") != ""
	envMap.Sonarr.URL = os.Getenv("KEEPERCHEKY_CLIENTS_SONARR_URL") != ""
	envMap.Sonarr.APIKey = os.Getenv("KEEPERCHEKY_CLIENTS_SONARR_API_KEY") != ""

	// Check Jellyfin
	envMap.Jellyfin.Enabled = os.Getenv("KEEPERCHEKY_CLIENTS_JELLYFIN_ENABLED") != ""
	envMap.Jellyfin.URL = os.Getenv("KEEPERCHEKY_CLIENTS_JELLYFIN_URL") != ""
	envMap.Jellyfin.APIKey = os.Getenv("KEEPERCHEKY_CLIENTS_JELLYFIN_API_KEY") != ""

	// Check Jellyseerr
	envMap.Jellyseerr.Enabled = os.Getenv("KEEPERCHEKY_CLIENTS_JELLYSEERR_ENABLED") != ""
	envMap.Jellyseerr.URL = os.Getenv("KEEPERCHEKY_CLIENTS_JELLYSEERR_URL") != ""
	envMap.Jellyseerr.APIKey = os.Getenv("KEEPERCHEKY_CLIENTS_JELLYSEERR_API_KEY") != ""

	// Check qBittorrent
	envMap.QBittorrent.Enabled = os.Getenv("KEEPERCHEKY_CLIENTS_QBITTORRENT_ENABLED") != ""
	envMap.QBittorrent.URL = os.Getenv("KEEPERCHEKY_CLIENTS_QBITTORRENT_URL") != ""
	envMap.QBittorrent.Username = os.Getenv("KEEPERCHEKY_CLIENTS_QBITTORRENT_USERNAME") != ""
	envMap.QBittorrent.Password = os.Getenv("KEEPERCHEKY_CLIENTS_QBITTORRENT_PASSWORD") != ""

	return envMap
}

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

	// Environment variables with KEEPERCHEKY_ prefix
	// CRITICAL: SetEnvKeyReplacer allows KEEPERCHEKY_CLIENTS_RADARR_URL to map to clients.radarr.url
	viper.SetEnvPrefix("KEEPERCHEKY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Bind environment variables WITHOUT prefix (for external service credentials)
	// These take precedence over config file values
	viper.BindEnv("clients.radarr.api_key", "RADARR_API_KEY")
	viper.BindEnv("clients.sonarr.api_key", "SONARR_API_KEY")
	viper.BindEnv("clients.jellyfin.api_key", "JELLYFIN_API_KEY")
	viper.BindEnv("clients.jellyseerr.api_key", "JELLYSEERR_API_KEY")
	viper.BindEnv("clients.qbittorrent.username", "QBITTORRENT_USERNAME")
	viper.BindEnv("clients.qbittorrent.password", "QBITTORRENT_PASSWORD")

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

	// ¡NUEVA LÓGICA! Si hay variables de entorno definidas, SOBREESCRIBIR config.yaml
	envSources := GetEnvSourceMap()
	needsSave := false

	// Sobreescribir valores desde variables de entorno SOLO SI SON DIFERENTES
	if envSources.Radarr.Enabled && os.Getenv("KEEPERCHEKY_CLIENTS_RADARR_ENABLED") != "" {
		newValue := viper.GetBool("clients.radarr.enabled")
		if config.Clients.Radarr.Enabled != newValue {
			config.Clients.Radarr.Enabled = newValue
			needsSave = true
		}
	}
	if envSources.Radarr.URL && os.Getenv("KEEPERCHEKY_CLIENTS_RADARR_URL") != "" {
		newValue := viper.GetString("clients.radarr.url")
		if config.Clients.Radarr.URL != newValue {
			config.Clients.Radarr.URL = newValue
			needsSave = true
		}
	}
	if envSources.Radarr.APIKey && os.Getenv("KEEPERCHEKY_CLIENTS_RADARR_API_KEY") != "" {
		newValue := viper.GetString("clients.radarr.api_key")
		if config.Clients.Radarr.APIKey != newValue {
			config.Clients.Radarr.APIKey = newValue
			needsSave = true
		}
	}

	if envSources.Sonarr.Enabled && os.Getenv("KEEPERCHEKY_CLIENTS_SONARR_ENABLED") != "" {
		newValue := viper.GetBool("clients.sonarr.enabled")
		if config.Clients.Sonarr.Enabled != newValue {
			config.Clients.Sonarr.Enabled = newValue
			needsSave = true
		}
	}
	if envSources.Sonarr.URL && os.Getenv("KEEPERCHEKY_CLIENTS_SONARR_URL") != "" {
		newValue := viper.GetString("clients.sonarr.url")
		if config.Clients.Sonarr.URL != newValue {
			config.Clients.Sonarr.URL = newValue
			needsSave = true
		}
	}
	if envSources.Sonarr.APIKey && os.Getenv("KEEPERCHEKY_CLIENTS_SONARR_API_KEY") != "" {
		newValue := viper.GetString("clients.sonarr.api_key")
		if config.Clients.Sonarr.APIKey != newValue {
			config.Clients.Sonarr.APIKey = newValue
			needsSave = true
		}
	}

	if envSources.Jellyfin.Enabled && os.Getenv("KEEPERCHEKY_CLIENTS_JELLYFIN_ENABLED") != "" {
		newValue := viper.GetBool("clients.jellyfin.enabled")
		if config.Clients.Jellyfin.Enabled != newValue {
			config.Clients.Jellyfin.Enabled = newValue
			needsSave = true
		}
	}
	if envSources.Jellyfin.URL && os.Getenv("KEEPERCHEKY_CLIENTS_JELLYFIN_URL") != "" {
		newValue := viper.GetString("clients.jellyfin.url")
		if config.Clients.Jellyfin.URL != newValue {
			config.Clients.Jellyfin.URL = newValue
			needsSave = true
		}
	}
	if envSources.Jellyfin.APIKey && os.Getenv("KEEPERCHEKY_CLIENTS_JELLYFIN_API_KEY") != "" {
		newValue := viper.GetString("clients.jellyfin.api_key")
		if config.Clients.Jellyfin.APIKey != newValue {
			config.Clients.Jellyfin.APIKey = newValue
			needsSave = true
		}
	}

	if envSources.Jellyseerr.Enabled && os.Getenv("KEEPERCHEKY_CLIENTS_JELLYSEERR_ENABLED") != "" {
		newValue := viper.GetBool("clients.jellyseerr.enabled")
		if config.Clients.Jellyseerr.Enabled != newValue {
			config.Clients.Jellyseerr.Enabled = newValue
			needsSave = true
		}
	}
	if envSources.Jellyseerr.URL && os.Getenv("KEEPERCHEKY_CLIENTS_JELLYSEERR_URL") != "" {
		newValue := viper.GetString("clients.jellyseerr.url")
		if config.Clients.Jellyseerr.URL != newValue {
			config.Clients.Jellyseerr.URL = newValue
			needsSave = true
		}
	}
	if envSources.Jellyseerr.APIKey && os.Getenv("KEEPERCHEKY_CLIENTS_JELLYSEERR_API_KEY") != "" {
		newValue := viper.GetString("clients.jellyseerr.api_key")
		if config.Clients.Jellyseerr.APIKey != newValue {
			config.Clients.Jellyseerr.APIKey = newValue
			needsSave = true
		}
	}

	if envSources.QBittorrent.Enabled && os.Getenv("KEEPERCHEKY_CLIENTS_QBITTORRENT_ENABLED") != "" {
		newValue := viper.GetBool("clients.qbittorrent.enabled")
		if config.Clients.QBittorrent.Enabled != newValue {
			config.Clients.QBittorrent.Enabled = newValue
			needsSave = true
		}
	}
	if envSources.QBittorrent.URL && os.Getenv("KEEPERCHEKY_CLIENTS_QBITTORRENT_URL") != "" {
		newValue := viper.GetString("clients.qbittorrent.url")
		if config.Clients.QBittorrent.URL != newValue {
			config.Clients.QBittorrent.URL = newValue
			needsSave = true
		}
	}
	if envSources.QBittorrent.Username && os.Getenv("KEEPERCHEKY_CLIENTS_QBITTORRENT_USERNAME") != "" {
		newValue := viper.GetString("clients.qbittorrent.username")
		if config.Clients.QBittorrent.Username != newValue {
			config.Clients.QBittorrent.Username = newValue
			needsSave = true
		}
	}
	if envSources.QBittorrent.Password && os.Getenv("KEEPERCHEKY_CLIENTS_QBITTORRENT_PASSWORD") != "" {
		newValue := viper.GetString("clients.qbittorrent.password")
		if config.Clients.QBittorrent.Password != newValue {
			config.Clients.QBittorrent.Password = newValue
			needsSave = true
		}
	}

	// Si se detectaron CAMBIOS en variables de entorno, GUARDAR en config.yaml
	if needsSave {
		if err := Save(&config); err != nil {
			// No fallar si no se puede guardar, solo loguear
			fmt.Printf("Warning: Could not save env vars to config.yaml: %v\n", err)
		} else {
			fmt.Println("Config synced: Environment variables written to config.yaml")
		}
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
