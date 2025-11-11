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

// Torrent state constants from qBittorrent API
const (
	TorrentStateDownloading = "downloading"
	TorrentStateUploading   = "uploading"
	TorrentStateStalledUP   = "stalledUP"
	TorrentStateCheckingUP  = "checkingUP"
	TorrentStateForcedUP    = "forcedUP"
	TorrentStateQueuedUP    = "queuedUP"
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
// This struct captures all available fields from /api/v2/torrents/info endpoint.
type qbTorrent struct {
	// Basic identification
	Hash       string `json:"hash"`        // Torrent hash (unique identifier)
	InfoHashV1 string `json:"infohash_v1"` // InfoHash v1 (BitTorrent v1)
	InfoHashV2 string `json:"infohash_v2"` // InfoHash v2 (BitTorrent v2)
	Name       string `json:"name"`        // Torrent name
	MagnetURI  string `json:"magnet_uri"`  // Magnet URI

	// Size and progress
	Size       int64   `json:"size"`        // Total size (bytes)
	TotalSize  int64   `json:"total_size"`  // Total size including unselected files
	Progress   float64 `json:"progress"`    // Download progress (0.0 to 1.0)
	AmountLeft int64   `json:"amount_left"` // Bytes left to download
	Completed  int64   `json:"completed"`   // Bytes completed

	// State
	State string `json:"state"` // Torrent state (see TorrentState constants)

	// Speeds
	DlSpeed int64 `json:"dlspeed"` // Download speed (bytes/s)
	UpSpeed int64 `json:"upspeed"` // Upload speed (bytes/s)

	// Transfer statistics
	Downloaded        int64   `json:"downloaded"`         // Total downloaded (bytes)
	Uploaded          int64   `json:"uploaded"`           // Total uploaded (bytes)
	DownloadedSession int64   `json:"downloaded_session"` // Downloaded this session
	UploadedSession   int64   `json:"uploaded_session"`   // Uploaded this session
	Ratio             float64 `json:"ratio"`              // Share ratio

	// Peers and seeds
	NumSeeds      int     `json:"num_seeds"`      // Seeds connected
	NumComplete   int     `json:"num_complete"`   // Total seeds in swarm
	NumLeechs     int     `json:"num_leechs"`     // Leechers connected
	NumIncomplete int     `json:"num_incomplete"` // Total leechers in swarm
	Availability  float64 `json:"availability"`   // Piece availability

	// Time information
	AddedOn      int64 `json:"added_on"`      // Unix timestamp when added
	CompletionOn int64 `json:"completion_on"` // Unix timestamp when completed
	CompletedOn  int64 `json:"completed_on"`  // Alias for completion_on (some API versions)
	LastActivity int64 `json:"last_activity"` // Unix timestamp of last activity
	SeenComplete int64 `json:"seen_complete"` // Unix timestamp last seen complete
	TimeActive   int64 `json:"time_active"`   // Total active time (seconds)
	SeedingTime  int64 `json:"seeding_time"`  // Total seeding time (seconds)
	ETA          int64 `json:"eta"`           // Estimated time to completion (seconds)

	// Paths
	SavePath     string `json:"save_path"`     // Save path
	ContentPath  string `json:"content_path"`  // Content path (full path)
	DownloadPath string `json:"download_path"` // Temporary download path

	// Organization
	Category string `json:"category"` // Category
	Tags     string `json:"tags"`     // Comma-separated tags

	// Tracker
	Tracker       string `json:"tracker"`        // Primary tracker URL
	TrackersCount int    `json:"trackers_count"` // Number of trackers

	// Priority and limits
	Priority int   `json:"priority"` // Torrent priority
	DlLimit  int64 `json:"dl_limit"` // Download limit (-1 = unlimited)
	UpLimit  int64 `json:"up_limit"` // Upload limit (-1 = unlimited)

	// Ratio and time limits
	MaxRatio         float64 `json:"max_ratio"`          // Max ratio (-1 = global, -2 = unlimited)
	MaxSeedingTime   int64   `json:"max_seeding_time"`   // Max seeding time (-1 = global, -2 = unlimited)
	RatioLimit       float64 `json:"ratio_limit"`        // Applied ratio limit
	SeedingTimeLimit int64   `json:"seeding_time_limit"` // Applied seeding time limit

	// Flags
	SeqDl        bool `json:"seq_dl"`         // Sequential download enabled
	FlPiecePrio  bool `json:"f_l_piece_prio"` // First/last piece priority
	ForceStart   bool `json:"force_start"`    // Force start enabled
	SuperSeeding bool `json:"super_seeding"`  // Super seeding enabled
	AutoTMM      bool `json:"auto_tmm"`       // Automatic Torrent Management
}

// qbBuildInfo represents build info from qBittorrent API.
type qbBuildInfo struct {
	Qt         string `json:"qt"`
	Libtorrent string `json:"libtorrent"`
	Boost      string `json:"boost"`
	Openssl    string `json:"openssl"`
	Bitness    int    `json:"bitness"`
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
		info := c.convertToTorrentInfo(&t)

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
			info := c.convertToTorrentInfo(&t)

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
		TorrentStateUploading:  true,
		TorrentStateStalledUP:  true,
		TorrentStateCheckingUP: true,
		TorrentStateForcedUP:   true,
		TorrentStateQueuedUP:   true,
	}

	return seedingStates[state]
}

// convertToTorrentInfo converts a qbTorrent to TorrentInfo model with all fields populated.
func (c *QBittorrentClient) convertToTorrentInfo(t *qbTorrent) *models.TorrentInfo {
	return &models.TorrentInfo{
		// Basic identification
		Hash:       t.Hash,
		InfoHashV1: t.InfoHashV1,
		InfoHashV2: t.InfoHashV2,
		Name:       t.Name,
		MagnetURI:  t.MagnetURI,

		// Size and progress
		Size:       t.Size,
		TotalSize:  t.TotalSize,
		Progress:   t.Progress,
		AmountLeft: t.AmountLeft,

		// State and activity
		State:        t.State,
		IsSeeding:    c.isStateSeeding(t.State),
		IsComplete:   t.Progress >= 1.0,
		Availability: t.Availability,

		// Speeds
		UpSpeed: t.UpSpeed,
		DlSpeed: t.DlSpeed,

		// Transfer statistics
		Ratio:             t.Ratio,
		TotalUploaded:     t.Uploaded,
		TotalDownloaded:   t.Downloaded,
		UploadedSession:   t.UploadedSession,
		DownloadedSession: t.DownloadedSession,

		// Peers and seeds
		NumSeeds:      t.NumSeeds,
		NumComplete:   t.NumComplete,
		NumPeers:      t.NumLeechs,
		NumIncomplete: t.NumIncomplete,
		Seeders:       t.NumSeeds,  // Backward compatibility
		Leechers:      t.NumLeechs, // Backward compatibility

		// Time information
		AddedOn:      t.AddedOn,
		CompletedOn:  t.CompletedOn,
		LastActivity: t.LastActivity,
		SeenComplete: t.SeenComplete,
		TimeActive:   t.TimeActive,
		SeedingTime:  t.SeedingTime,
		ETA:          t.ETA,

		// Paths
		SavePath:     t.SavePath,
		ContentPath:  t.ContentPath,
		DownloadPath: t.DownloadPath,

		// Organization
		Category: t.Category,
		Tags:     t.Tags,

		// Tracker
		Tracker:       t.Tracker,
		TrackersCount: t.TrackersCount,

		// Priority and limits
		Priority: t.Priority,
		DlLimit:  t.DlLimit,
		UpLimit:  t.UpLimit,

		// Ratio and time limits
		MaxRatio:         t.MaxRatio,
		MaxSeedingTime:   t.MaxSeedingTime,
		RatioLimit:       t.RatioLimit,
		SeedingTimeLimit: t.SeedingTimeLimit,

		// Flags
		SeqDl:        t.SeqDl,
		FlPiecePrio:  t.FlPiecePrio,
		ForceStart:   t.ForceStart,
		SuperSeeding: t.SuperSeeding,
		AutoTMM:      t.AutoTMM,
	}
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

// GetTransferInfo retrieves global transfer statistics from qBittorrent.
func (c *QBittorrentClient) GetTransferInfo(ctx context.Context) (*models.QBittorrentTransferInfo, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	var transferInfo models.QBittorrentTransferInfo
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&transferInfo).
		Get("/api/v2/transfer/info")

	if err != nil {
		return nil, fmt.Errorf("failed to get transfer info: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), resp.String())
	}

	c.logger.Debug("Retrieved qBittorrent transfer info",
		zap.Int64("dl_speed", transferInfo.DLInfoSpeed),
		zap.Int64("up_speed", transferInfo.UPInfoSpeed),
		zap.Int("dht_nodes", transferInfo.DHTNodes),
		zap.String("connection_status", transferInfo.ConnectionStatus),
	)

	return &transferInfo, nil
}

// GetServerState retrieves server state from qBittorrent sync/maindata endpoint.
func (c *QBittorrentClient) GetServerState(ctx context.Context) (*models.QBittorrentServerState, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	// The maindata endpoint returns a complex structure, we only need server_state
	var response struct {
		ServerState models.QBittorrentServerState `json:"server_state"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&response).
		Get("/api/v2/sync/maindata")

	if err != nil {
		return nil, fmt.Errorf("failed to get server state: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), resp.String())
	}

	c.logger.Debug("Retrieved qBittorrent server state",
		zap.Int64("dl_speed", response.ServerState.DLInfoSpeed),
		zap.Int64("up_speed", response.ServerState.UPInfoSpeed),
		zap.Int64("free_space", response.ServerState.FreeSpaceOnDisk),
	)

	return &response.ServerState, nil
}

// GetTorrentProperties retrieves detailed properties for a specific torrent.
func (c *QBittorrentClient) GetTorrentProperties(ctx context.Context, hash string) (*models.QBittorrentTorrentProperties, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	var props models.QBittorrentTorrentProperties
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&props).
		SetQueryParam("hash", hash).
		Get("/api/v2/torrents/properties")

	if err != nil {
		return nil, fmt.Errorf("failed to get torrent properties: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), resp.String())
	}

	c.logger.Debug("Retrieved torrent properties",
		zap.String("hash", hash),
		zap.Int64("eta", props.ETA),
		zap.Int("seeds", props.Seeds),
		zap.Int("peers", props.Peers),
	)

	return &props, nil
}

// GetTorrentTrackers retrieves tracker information for a specific torrent.
func (c *QBittorrentClient) GetTorrentTrackers(ctx context.Context, hash string) ([]models.QBittorrentTracker, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	var trackers []models.QBittorrentTracker
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&trackers).
		SetQueryParam("hash", hash).
		Get("/api/v2/torrents/trackers")

	if err != nil {
		return nil, fmt.Errorf("failed to get torrent trackers: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), resp.String())
	}

	c.logger.Debug("Retrieved torrent trackers",
		zap.String("hash", hash),
		zap.Int("tracker_count", len(trackers)),
	)

	return trackers, nil
}

// GetEnhancedTorrentInfo retrieves basic torrent info with enhanced details from properties endpoint.
// This provides a richer dataset for displaying in the UI.
func (c *QBittorrentClient) GetEnhancedTorrentInfo(ctx context.Context, hash string) (*models.TorrentInfo, error) {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Get basic torrent info
	var torrents []qbTorrent
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&torrents).
		SetQueryParam("hashes", hash).
		Get("/api/v2/torrents/info")

	if err != nil {
		return nil, fmt.Errorf("failed to get torrent info: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	if len(torrents) == 0 {
		return nil, fmt.Errorf("torrent not found: %s", hash)
	}

	torrent := torrents[0]

	// Get enhanced properties
	props, err := c.GetTorrentProperties(ctx, hash)
	if err != nil {
		// Log but don't fail - return basic info using helper
		c.logger.Warn("Failed to get torrent properties, using basic info only",
			zap.String("hash", hash),
			zap.Error(err),
		)
		return c.convertToTorrentInfo(&torrent), nil
	}

	// Start with complete basic info
	info := c.convertToTorrentInfo(&torrent)

	// Override/enhance with properties endpoint data if available
	// Properties endpoint may have more accurate data for some fields
	if props.AdditionDate > 0 {
		info.AddedOn = props.AdditionDate
	}
	if props.CompletionDate > 0 {
		info.CompletedOn = props.CompletionDate
	}
	if props.ETA > 0 {
		info.ETA = props.ETA
	}
	if props.TotalUploaded > 0 {
		info.TotalUploaded = props.TotalUploaded
	}
	if props.TotalDownloaded > 0 {
		info.TotalDownloaded = props.TotalDownloaded
	}
	if props.Seeds > 0 {
		info.NumSeeds = props.Seeds
		info.Seeders = props.Seeds
	}
	if props.Peers > 0 {
		info.NumPeers = props.Peers
		info.Leechers = props.Peers
	}

	return info, nil
}

// performTorrentAction is a generic helper method for torrent operations
func (c *QBittorrentClient) performTorrentAction(ctx context.Context, hash, endpoint, action string) error {
	// Ensure logged in
	if c.cookie == "" {
		if err := c.login(ctx); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetFormData(map[string]string{"hashes": hash}).
		Post(endpoint)

	if err != nil {
		return fmt.Errorf("failed to %s torrent: %w", action, err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), resp.String())
	}

	c.logger.Info(fmt.Sprintf("%s torrent", action), zap.String("hash", hash))
	return nil
}

// PauseTorrent pauses a torrent
func (c *QBittorrentClient) PauseTorrent(ctx context.Context, hash string) error {
	return c.performTorrentAction(ctx, hash, "/api/v2/torrents/pause", "Paused")
}

// ResumeTorrent resumes a paused torrent
func (c *QBittorrentClient) ResumeTorrent(ctx context.Context, hash string) error {
	return c.performTorrentAction(ctx, hash, "/api/v2/torrents/resume", "Resumed")
}

// RecheckTorrent rechecks a torrent
func (c *QBittorrentClient) RecheckTorrent(ctx context.Context, hash string) error {
	return c.performTorrentAction(ctx, hash, "/api/v2/torrents/recheck", "Rechecked")
}
