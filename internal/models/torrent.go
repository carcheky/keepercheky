package models

// TorrentInfo represents torrent information for media
type TorrentInfo struct {
	Hash      string  `json:"hash"`
	Name      string  `json:"name"`
	Size      int64   `json:"size"`
	Progress  float64 `json:"progress"`
	State     string  `json:"state"`
	Seeders   int     `json:"seeders"`
	Leechers  int     `json:"leechers"`
	IsSeeding bool    `json:"is_seeding"`
}
