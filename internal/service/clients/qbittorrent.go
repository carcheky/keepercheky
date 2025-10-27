package clients

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// QBittorrentClient implements the TorrentClient interface for qBittorrent.
type QBittorrentClient struct {
	client   *resty.Client
	baseURL  string
	username string
	password string
	logger   *zap.Logger
	cookie   string // SID cookie for authentication
}

// NewQBittorrentClient creates a new qBittorrent client.
func NewQBittorrentClient(baseURL, username, password string, logger *zap.Logger) *QBittorrentClient {
	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetTimeout(30 * time.Second)

	return &QBittorrentClient{
		client:   client,
		baseURL:  baseURL,
		username: username,
		password: password,
		logger:   logger,
	}
}

// qbTorrent represents a torrent from qBittorrent API.
type qbTorrent struct {
	Hash        string  `json:"hash"`
	Name        string  `json:"name"`
	State       string  `json:"state"`
	Progress    float64 `json:"progress"`
	Ratio       float64 `json:"ratio"`
	Size        int64   `json:"size"`
	UpSpeed     int64   `json:"upspeed"`
	DlSpeed     int64   `json:"dlspeed"`
	SeedingTime int64   `json:"seeding_time"`
	Category    string  `json:"category"`
	Tags        string  `json:"tags"`
	SavePath    string  `json:"save_path"`
	ContentPath string  `json:"content_path"`
	AmountLeft  int64   `json:"amount_left"`
}

// qbBuildInfo represents build info from qBittorrent API.
type qbBuildInfo struct {
	Qt         string `json:"qt"`
	Libtorrent string `json:"libtorrent"`
	Boost      string `json:"boost"`
	Openssl    string `json:"openssl"`
	Bitness    int    `json:"bitness"`
}

// qbAppVersion represents app version from qBittorrent API.
// nolint:unused // Reserved for future use
type qbAppVersion struct {
	Version string `json:"version"` // Returns just the version string
}

// QBittorrentSystemInfo representa toda la información del sistema de qBittorrent
type QBittorrentSystemInfo struct {
	Version    string `json:"version"`
	Qt         string `json:"qt"`
	Libtorrent string `json:"libtorrent"`
	Boost      string `json:"boost"`
	Openssl    string `json:"openssl"`
	Bitness    int    `json:"bitness"`
}

// QBittorrentPreferences represents the application preferences from qBittorrent API.
// We only include the fields we need for path configuration.
type QBittorrentPreferences struct {
	SavePath       string                 `json:"save_path"`         // Default save path for torrents
	TempPath       string                 `json:"temp_path"`         // Path for incomplete torrents
	TempPathEnable bool                   `json:"temp_path_enabled"` // True if temp_path is enabled
	ExportDir      string                 `json:"export_dir"`        // Path to copy .torrent files
	ExportDirFin   string                 `json:"export_dir_fin"`    // Path to copy completed .torrent files
	ScanDirs       map[string]interface{} `json:"scan_dirs"`         // Directories to watch for torrents
}

// QBittorrentCategory represents a qBittorrent category with its save path.
type QBittorrentCategory struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

// login authenticates with qBittorrent and stores the SID cookie.
func (c *QBittorrentClient) login(ctx context.Context) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetFormData(map[string]string{
			"username": c.username,
			"password": c.password,
		}).
		Post("/api/v2/auth/login")

	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	// Extract SID cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SID" {
			c.cookie = cookie.Value
			c.client.SetHeader("Cookie", fmt.Sprintf("SID=%s", c.cookie))
			break
		}
	}

	if c.cookie == "" {
		return fmt.Errorf("no SID cookie received")
	}

	return nil
}

// TestConnection verifies the connection to qBittorrent.
func (c *QBittorrentClient) TestConnection(ctx context.Context) error {
	// Login first
	if err := c.login(ctx); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	var buildInfo qbBuildInfo
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&buildInfo).
		Get("/api/v2/app/buildInfo")

	if err != nil {
		return fmt.Errorf("failed to get build info: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	c.logger.Info("qBittorrent connection successful",
		zap.String("qt", buildInfo.Qt),
		zap.String("libtorrent", buildInfo.Libtorrent),
		zap.Int("bitness", buildInfo.Bitness),
		zap.String("url", c.baseURL),
	)

	return nil
}

// GetSystemInfo retrieves complete system information from qBittorrent.
func (c *QBittorrentClient) GetSystemInfo(ctx context.Context) (*QBittorrentSystemInfo, error) {
	// Login first
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Get version
	var version string
	versionResp, err := c.client.R().
		SetContext(ctx).
		Get("/api/v2/app/version")

	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	if versionResp.StatusCode() == 200 {
		version = strings.Trim(versionResp.String(), "\"")
	}

	// Get build info
	var buildInfo qbBuildInfo
	buildResp, err := c.client.R().
		SetContext(ctx).
		SetResult(&buildInfo).
		Get("/api/v2/app/buildInfo")

	if err != nil {
		return nil, fmt.Errorf("failed to get build info: %w", err)
	}

	if buildResp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", buildResp.StatusCode())
	}

	// Convertir a nuestro modelo
	systemInfo := &QBittorrentSystemInfo{
		Version:    version,
		Qt:         buildInfo.Qt,
		Libtorrent: buildInfo.Libtorrent,
		Boost:      buildInfo.Boost,
		Openssl:    buildInfo.Openssl,
		Bitness:    buildInfo.Bitness,
	}

	c.logger.Info("Retrieved qBittorrent system info",
		zap.String("version", systemInfo.Version),
		zap.String("qt", systemInfo.Qt),
		zap.String("libtorrent", systemInfo.Libtorrent),
		zap.Int("bitness", systemInfo.Bitness),
	)

	return systemInfo, nil
}

// GetPreferences retrieves the application preferences from qBittorrent.
// This includes important path configuration like save_path, temp_path, etc.
func (c *QBittorrentClient) GetPreferences(ctx context.Context) (*QBittorrentPreferences, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	var prefs QBittorrentPreferences
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&prefs).
		Get("/api/v2/app/preferences")

	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), resp.String())
	}

	c.logger.Info("Retrieved qBittorrent preferences",
		zap.String("save_path", prefs.SavePath),
		zap.String("temp_path", prefs.TempPath),
		zap.Bool("temp_path_enabled", prefs.TempPathEnable),
		zap.String("export_dir", prefs.ExportDir),
		zap.String("export_dir_fin", prefs.ExportDirFin),
		zap.Int("scan_dirs_count", len(prefs.ScanDirs)),
	)

	return &prefs, nil
}

// GetCategories retrieves all torrent categories with their save paths from qBittorrent.
func (c *QBittorrentClient) GetCategories(ctx context.Context) (map[string]QBittorrentCategory, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	var categories map[string]QBittorrentCategory
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&categories).
		Get("/api/v2/torrents/categories")

	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), resp.String())
	}

	c.logger.Info("Retrieved qBittorrent categories",
		zap.Int("category_count", len(categories)),
	)

	return categories, nil
}

// GetAllTorrentsMap retrieves all torrents and returns them indexed by content_path for fast lookup.
// This is much more efficient than calling GetTorrentByPath() for each media item.
func (c *QBittorrentClient) GetAllTorrentsMap(ctx context.Context) (map[string]*models.TorrentInfo, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, err
		}
	}

	var torrents []qbTorrent
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&torrents).
		Get("/api/v2/torrents/info")

	if err != nil {
		return nil, fmt.Errorf("failed to get torrents: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// Create map indexed by content_path for fast lookup
	torrentMap := make(map[string]*models.TorrentInfo, len(torrents))

	for i, t := range torrents {
		info := &models.TorrentInfo{
			Hash:        t.Hash,
			Name:        t.Name,
			State:       t.State,
			Progress:    t.Progress,
			Ratio:       t.Ratio,
			Size:        t.Size,
			UpSpeed:     t.UpSpeed,
			DlSpeed:     t.DlSpeed,
			SeedingTime: t.SeedingTime,
			Category:    t.Category,
			Tags:        t.Tags,
			SavePath:    t.SavePath,
			IsSeeding:   c.isStateSeeding(t.State),
			IsComplete:  t.Progress >= 1.0,
		}

		// Log first 5 torrents for debugging path structure
		if i < 5 {
			c.logger.Debug("Sample torrent",
				zap.String("name", t.Name),
				zap.String("content_path", t.ContentPath),
				zap.String("save_path", t.SavePath),
				zap.String("hash", t.Hash),
			)
		}

		// Index by both content_path and save_path for better matching
		if t.ContentPath != "" {
			torrentMap[t.ContentPath] = info
		}
		if t.SavePath != "" && t.SavePath != t.ContentPath {
			torrentMap[t.SavePath] = info
		}
	}

	c.logger.Info("Retrieved all torrents from qBittorrent",
		zap.Int("total_torrents", len(torrents)),
		zap.Int("indexed_paths", len(torrentMap)),
	)

	return torrentMap, nil
}

// GetTorrentByPath finds a torrent by its file path.
func (c *QBittorrentClient) GetTorrentByPath(ctx context.Context, filePath string) (*models.TorrentInfo, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, err
		}
	}

	var torrents []qbTorrent
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&torrents).
		Get("/api/v2/torrents/info")

	if err != nil {
		return nil, fmt.Errorf("failed to get torrents: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// Search for torrent containing the file path
	for _, t := range torrents {
		// Check if the file path matches this torrent's save path or content path
		if strings.Contains(filePath, t.SavePath) || strings.Contains(filePath, t.ContentPath) {
			info := &models.TorrentInfo{
				Hash:        t.Hash,
				Name:        t.Name,
				State:       t.State,
				Progress:    t.Progress,
				Ratio:       t.Ratio,
				Size:        t.Size,
				UpSpeed:     t.UpSpeed,
				DlSpeed:     t.DlSpeed,
				SeedingTime: t.SeedingTime,
				Category:    t.Category,
				Tags:        t.Tags,
				SavePath:    t.SavePath,
				IsSeeding:   c.isStateSeeding(t.State),
				IsComplete:  t.Progress >= 1.0,
			}

			c.logger.Info("Found torrent for path",
				zap.String("path", filePath),
				zap.String("hash", info.Hash),
				zap.String("name", info.Name),
				zap.String("state", info.State),
				zap.Bool("is_seeding", info.IsSeeding),
			)

			return info, nil
		}
	}

	return nil, fmt.Errorf("no torrent found for path: %s", filePath)
}

// IsSeeding checks if a file is currently seeding.
func (c *QBittorrentClient) IsSeeding(ctx context.Context, filePath string) (bool, error) {
	torrent, err := c.GetTorrentByPath(ctx, filePath)
	if err != nil {
		// If torrent not found, it's not seeding
		if strings.Contains(err.Error(), "no torrent found") {
			return false, nil
		}
		return false, err
	}

	c.logger.Info("Checked seeding status",
		zap.String("path", filePath),
		zap.Bool("is_seeding", torrent.IsSeeding),
		zap.Float64("ratio", torrent.Ratio),
		zap.Int64("seed_time", torrent.SeedingTime),
	)

	return torrent.IsSeeding, nil
}

// isStateSeeding determines if a torrent state indicates seeding.
func (c *QBittorrentClient) isStateSeeding(state string) bool {
	seedingStates := map[string]bool{
		"uploading":  true,
		"stalledUP":  true,
		"checkingUP": true,
		"forcedUP":   true,
		"queuedUP":   true,
	}

	return seedingStates[state]
}

// GetAllTorrents retrieves all torrents as Media items (for orphaned torrents not in Radarr/Sonarr).
func (c *QBittorrentClient) GetAllTorrents(ctx context.Context) ([]*models.Media, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, err
		}
	}

	var torrents []qbTorrent
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&torrents).
		Get("/api/v2/torrents/info")

	if err != nil {
		return nil, fmt.Errorf("failed to get torrents: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	mediaList := make([]*models.Media, 0, len(torrents))
	for _, t := range torrents {
		media := &models.Media{
			Title:       t.Name,
			Type:        "torrent", // Special type for standalone torrents
			FilePath:    t.ContentPath,
			Size:        t.Size,
			AddedDate:   time.Now(), // qBittorrent API doesn't provide this easily
			TorrentHash: t.Hash,
			IsSeeding:   c.isStateSeeding(t.State),
			SeedRatio:   t.Ratio,
			Quality:     t.Category, // Use category as quality indicator
			Tags:        []string{t.Tags},
		}

		mediaList = append(mediaList, media)
	}

	c.logger.Info("Retrieved all torrents from qBittorrent",
		zap.Int("total_torrents", len(mediaList)),
	)

	return mediaList, nil
}

// GetSeedingStatus checks the seeding status of a torrent by hash.
func (c *QBittorrentClient) GetSeedingStatus(ctx context.Context, hash string) (bool, float64, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return false, 0, err
		}
	}

	var torrents []qbTorrent
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&torrents).
		SetQueryParam("hashes", hash).
		Get("/api/v2/torrents/info")

	if err != nil {
		return false, 0, fmt.Errorf("failed to get torrent info: %w", err)
	}

	if resp.StatusCode() != 200 {
		return false, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	if len(torrents) == 0 {
		return false, 0, fmt.Errorf("torrent not found: %s", hash)
	}

	torrent := torrents[0]
	isSeeding := c.isStateSeeding(torrent.State)

	return isSeeding, torrent.Ratio, nil
}

// DeleteTorrent removes a torrent from qBittorrent.
func (c *QBittorrentClient) DeleteTorrent(ctx context.Context, hash string, deleteFiles bool) error {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return err
		}
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetFormData(map[string]string{
			"hashes":      hash,
			"deleteFiles": fmt.Sprintf("%t", deleteFiles),
		}).
		Post("/api/v2/torrents/delete")

	if err != nil {
		return fmt.Errorf("failed to delete torrent: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	c.logger.Info("Deleted torrent from qBittorrent",
		zap.String("hash", hash),
		zap.Bool("deleted_files", deleteFiles),
	)

	return nil
}
