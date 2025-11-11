package models

// TorrentInfo represents comprehensive torrent information for media.
// This struct captures all relevant data from qBittorrent API for display and cleanup logic.
type TorrentInfo struct {
	// Basic identification
	Hash       string `json:"hash"`        // Torrent hash (unique identifier)
	InfoHashV1 string `json:"infohash_v1"` // InfoHash v1
	InfoHashV2 string `json:"infohash_v2"` // InfoHash v2
	Name       string `json:"name"`        // Torrent name
	MagnetURI  string `json:"magnet_uri"`  // Magnet URI

	// Size and progress
	Size       int64   `json:"size"`        // Total size (bytes)
	TotalSize  int64   `json:"total_size"`  // Total size including unselected files
	Progress   float64 `json:"progress"`    // Download progress (0.0 to 1.0)
	AmountLeft int64   `json:"amount_left"` // Bytes left to download

	// State and activity
	State        string  `json:"state"`        // Torrent state (downloading, uploading, etc.)
	IsSeeding    bool    `json:"is_seeding"`   // True if in seeding state
	IsComplete   bool    `json:"is_complete"`  // True if download complete
	Availability float64 `json:"availability"` // Piece availability

	// Speeds
	UpSpeed int64 `json:"upspeed"` // Upload speed (bytes/s)
	DlSpeed int64 `json:"dlspeed"` // Download speed (bytes/s)

	// Transfer statistics
	Ratio             float64 `json:"ratio"`                        // Share ratio
	TotalUploaded     int64   `json:"total_uploaded,omitempty"`     // Total bytes uploaded
	TotalDownloaded   int64   `json:"total_downloaded,omitempty"`   // Total bytes downloaded
	UploadedSession   int64   `json:"uploaded_session,omitempty"`   // Uploaded this session
	DownloadedSession int64   `json:"downloaded_session,omitempty"` // Downloaded this session

	// Peers and seeds
	NumSeeds      int `json:"num_seeds,omitempty"`      // Seeds connected
	NumComplete   int `json:"num_complete,omitempty"`   // Total seeds in swarm
	NumPeers      int `json:"num_peers,omitempty"`      // Peers/leechers connected
	NumIncomplete int `json:"num_incomplete,omitempty"` // Total leechers in swarm
	Seeders       int `json:"seeders"`                  // Alias for NumSeeds (backward compat)
	Leechers      int `json:"leechers"`                 // Alias for NumPeers (backward compat)

	// Time information
	AddedOn      int64 `json:"added_on,omitempty"`      // Unix timestamp when added
	CompletedOn  int64 `json:"completed_on,omitempty"`  // Unix timestamp when completed
	LastActivity int64 `json:"last_activity,omitempty"` // Unix timestamp of last activity
	SeenComplete int64 `json:"seen_complete,omitempty"` // Unix timestamp last seen complete
	TimeActive   int64 `json:"time_active,omitempty"`   // Total active time (seconds)
	SeedingTime  int64 `json:"seeding_time"`            // Total seeding time (seconds)
	ETA          int64 `json:"eta,omitempty"`           // Estimated time to completion (seconds)

	// Paths
	SavePath     string `json:"save_path"`               // Save path
	ContentPath  string `json:"content_path,omitempty"`  // Content path (full path)
	DownloadPath string `json:"download_path,omitempty"` // Temporary download path

	// Organization
	Category string `json:"category"` // Category
	Tags     string `json:"tags"`     // Comma-separated tags

	// Tracker
	Tracker       string `json:"tracker,omitempty"`        // Primary tracker URL
	TrackersCount int    `json:"trackers_count,omitempty"` // Number of trackers

	// Priority and limits
	Priority int   `json:"priority,omitempty"` // Torrent priority
	DlLimit  int64 `json:"dl_limit,omitempty"` // Download limit (-1 = unlimited)
	UpLimit  int64 `json:"up_limit,omitempty"` // Upload limit (-1 = unlimited)

	// Ratio and time limits
	MaxRatio         float64 `json:"max_ratio,omitempty"`          // Max ratio (-1 = global, -2 = unlimited)
	MaxSeedingTime   int64   `json:"max_seeding_time,omitempty"`   // Max seeding time (-1 = global, -2 = unlimited)
	RatioLimit       float64 `json:"ratio_limit,omitempty"`        // Applied ratio limit
	SeedingTimeLimit int64   `json:"seeding_time_limit,omitempty"` // Applied seeding time limit

	// Flags
	SeqDl        bool `json:"seq_dl,omitempty"`         // Sequential download enabled
	FlPiecePrio  bool `json:"f_l_piece_prio,omitempty"` // First/last piece priority
	ForceStart   bool `json:"force_start,omitempty"`    // Force start enabled
	SuperSeeding bool `json:"super_seeding,omitempty"`  // Super seeding enabled
	AutoTMM      bool `json:"auto_tmm,omitempty"`       // Automatic Torrent Management
}

// QBittorrentTransferInfo represents global transfer statistics
type QBittorrentTransferInfo struct {
	DLInfoSpeed      int64  `json:"dl_info_speed"`     // Global download rate (bytes/s)
	DLInfoData       int64  `json:"dl_info_data"`      // Data downloaded this session (bytes)
	UPInfoSpeed      int64  `json:"up_info_speed"`     // Global upload rate (bytes/s)
	UPInfoData       int64  `json:"up_info_data"`      // Data uploaded this session (bytes)
	DHTNodes         int    `json:"dht_nodes"`         // DHT nodes connected to
	ConnectionStatus string `json:"connection_status"` // Connection status (connected, firewalled, disconnected)
}

// QBittorrentServerState represents server state from sync/maindata
type QBittorrentServerState struct {
	DLInfoSpeed      int64  `json:"dl_info_speed"`
	UPInfoSpeed      int64  `json:"up_info_speed"`
	DLInfoData       int64  `json:"dl_info_data"`
	UPInfoData       int64  `json:"up_info_data"`
	DHTNodes         int    `json:"dht_nodes"`
	ConnectionStatus string `json:"connection_status"`
	FreeSpaceOnDisk  int64  `json:"free_space_on_disk"`
}

// QBittorrentTorrentProperties represents detailed torrent properties
type QBittorrentTorrentProperties struct {
	AdditionDate    int64  `json:"addition_date"`
	CompletionDate  int64  `json:"completion_date"`
	TotalUploaded   int64  `json:"total_uploaded"`
	TotalDownloaded int64  `json:"total_downloaded"`
	NbConnections   int    `json:"nb_connections"`
	Seeds           int    `json:"seeds"`
	SeedsTotal      int    `json:"seeds_total"`
	Peers           int    `json:"peers"`
	PeersTotal      int    `json:"peers_total"`
	ETA             int64  `json:"eta"`
	SavePath        string `json:"save_path"`
	CreationDate    int64  `json:"creation_date"`
	PieceSize       int64  `json:"piece_size"`
	Comment         string `json:"comment"`
	TotalWasted     int64  `json:"total_wasted"`
	TotalSize       int64  `json:"total_size"`
	UpLimit         int64  `json:"up_limit"`
	DlLimit         int64  `json:"dl_limit"`
}

// QBittorrentTracker represents a tracker for a torrent
type QBittorrentTracker struct {
	URL      string `json:"url"`
	Status   int    `json:"status"` // 0=disabled, 1=not contacted, 2=working, 3=updating, 4=not working
	NumPeers int    `json:"num_peers"`
	NumSeeds int    `json:"num_seeds"`
	Msg      string `json:"msg"` // Tracker message
}
