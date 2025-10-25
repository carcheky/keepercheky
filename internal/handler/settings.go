package handler

import (
	"context"
	"time"

	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type SettingsHandler struct {
	repos       *repository.Repositories
	logger      *logger.Logger
	config      *config.Config
	syncService *service.SyncService
}

func NewSettingsHandler(repos *repository.Repositories, logger *logger.Logger, cfg *config.Config, syncService *service.SyncService) *SettingsHandler {
	return &SettingsHandler{
		repos:       repos,
		logger:      logger,
		config:      cfg,
		syncService: syncService,
	}
}

func (h *SettingsHandler) Index(c *fiber.Ctx) error {
	return c.Render("pages/settings", fiber.Map{
		"Title": "Settings - KeeperCheky",
	})
}

func (h *SettingsHandler) Get(c *fiber.Ctx) error {
	// Return config in the format expected by frontend
	return c.JSON(fiber.Map{
		"services": fiber.Map{
			"radarr": fiber.Map{
				"enabled": h.config.Services.Radarr.Enabled,
				"url":     h.config.Services.Radarr.URL,
				"api_key": h.config.Services.Radarr.APIKey,
			},
			"sonarr": fiber.Map{
				"enabled": h.config.Services.Sonarr.Enabled,
				"url":     h.config.Services.Sonarr.URL,
				"api_key": h.config.Services.Sonarr.APIKey,
			},
			"jellyfin": fiber.Map{
				"enabled": h.config.Services.Jellyfin.Enabled,
				"url":     h.config.Services.Jellyfin.URL,
				"api_key": h.config.Services.Jellyfin.APIKey,
			},
			"jellyseerr": fiber.Map{
				"enabled": h.config.Services.Jellyseerr.Enabled,
				"url":     h.config.Services.Jellyseerr.URL,
				"api_key": h.config.Services.Jellyseerr.APIKey,
			},
		},
		"cleanup": fiber.Map{
			"dry_run":            h.config.Cleanup.DryRun,
			"days_to_keep":       h.config.Cleanup.DaysToKeep,
			"leaving_soon_days":  h.config.Cleanup.LeavingSoonDays,
			"exclusion_tags":     "",
			"delete_unmonitored": false,
		},
	})
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	type ConfigUpdate struct {
		Services struct {
			Radarr struct {
				Enabled bool   `json:"enabled"`
				URL     string `json:"url"`
				APIKey  string `json:"api_key"`
			} `json:"radarr"`
			Sonarr struct {
				Enabled bool   `json:"enabled"`
				URL     string `json:"url"`
				APIKey  string `json:"api_key"`
			} `json:"sonarr"`
			Jellyfin struct {
				Enabled bool   `json:"enabled"`
				URL     string `json:"url"`
				APIKey  string `json:"api_key"`
			} `json:"jellyfin"`
			Jellyseerr struct {
				Enabled bool   `json:"enabled"`
				URL     string `json:"url"`
				APIKey  string `json:"api_key"`
			} `json:"jellyseerr"`
		} `json:"services"`
		Cleanup struct {
			DryRun           bool `json:"dry_run"`
			DaysToKeep       int  `json:"days_to_keep"`
			LeavingSoonDays  int  `json:"leaving_soon_days"`
			DeleteUnmonitored bool `json:"delete_unmonitored"`
		} `json:"cleanup"`
	}
	
	var update ConfigUpdate
	if err := c.BodyParser(&update); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}
	
	// Update config in memory
	h.config.Services.Radarr.Enabled = update.Services.Radarr.Enabled
	h.config.Services.Radarr.URL = update.Services.Radarr.URL
	h.config.Services.Radarr.APIKey = update.Services.Radarr.APIKey
	
	h.config.Services.Sonarr.Enabled = update.Services.Sonarr.Enabled
	h.config.Services.Sonarr.URL = update.Services.Sonarr.URL
	h.config.Services.Sonarr.APIKey = update.Services.Sonarr.APIKey
	
	h.config.Services.Jellyfin.Enabled = update.Services.Jellyfin.Enabled
	h.config.Services.Jellyfin.URL = update.Services.Jellyfin.URL
	h.config.Services.Jellyfin.APIKey = update.Services.Jellyfin.APIKey
	
	h.config.Services.Jellyseerr.Enabled = update.Services.Jellyseerr.Enabled
	h.config.Services.Jellyseerr.URL = update.Services.Jellyseerr.URL
	h.config.Services.Jellyseerr.APIKey = update.Services.Jellyseerr.APIKey
	
	h.config.Cleanup.DryRun = update.Cleanup.DryRun
	h.config.Cleanup.DaysToKeep = update.Cleanup.DaysToKeep
	h.config.Cleanup.LeavingSoonDays = update.Cleanup.LeavingSoonDays
	
	h.logger.Info("Configuration updated",
		"radarr_enabled", h.config.Services.Radarr.Enabled,
		"sonarr_enabled", h.config.Services.Sonarr.Enabled,
		"jellyfin_enabled", h.config.Services.Jellyfin.Enabled,
		"dry_run", h.config.Cleanup.DryRun,
	)
	
	return c.JSON(fiber.Map{
		"message": "Configuration saved successfully",
		"success": true,
	})
}

func (h *SettingsHandler) TestConnection(c *fiber.Ctx) error {
	service := c.Params("service")
	
	h.logger.Info("Testing connection", "service", service)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	err := h.syncService.TestConnection(ctx, service)
	if err != nil {
		h.logger.Error("Connection test failed",
			"service", service,
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Connection test successful",
		"version": "Connected", // TODO: Get actual version from service
	})
}
