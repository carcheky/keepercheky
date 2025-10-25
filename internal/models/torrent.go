package models

// TorrentInfo represents torrent information for media
type TorrentInfo struct {
	Hash        string  `json:"hash"`
	Name        string  `json:"name"`
	Size        int64   `json:"size"`
	Progress    float64 `json:"progress"`
	State       string  `json:"state"`
	Ratio       float64 `json:"ratio"`
	UpSpeed     int64   `json:"upspeed"`
	DlSpeed     int64   `json:"dlspeed"`
	SeedingTime int64   `json:"seeding_time"`
	Category    string  `json:"category"`
	Tags        string  `json:"tags"`
	SavePath    string  `json:"save_path"`
	Seeders     int     `json:"seeders"`
	Leechers    int     `json:"leechers"`
	IsSeeding   bool    `json:"is_seeding"`
	IsComplete  bool    `json:"is_complete"`
}
