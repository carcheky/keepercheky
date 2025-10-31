package filesystem

import (
	"context"
	"testing"

	"github.com/carcheky/keepercheky/internal/models"
	"go.uber.org/zap"
)

func TestEnricher_EnrichWithRadarr_HardlinkMatching(t *testing.T) {
	logger := zap.NewNop()
	enricher := NewEnricher(logger)

	// Simulate a file with hardlinks
	// File is in /Descargas but has hardlink in /Peliculas
	files := map[string]*EnrichedFile{
		"/Descargas/Movie.mkv": {
			FileEntry: &FileEntry{
				Path:          "/Descargas/Movie.mkv",
				Size:          1024 * 1024 * 1024,
				Inode:         12345,
				IsHardlink:    true,
				HardlinkPaths: []string{"/Descargas/Movie.mkv", "/Peliculas/Movie.mkv"},
				PrimaryPath:   "/Peliculas/Movie.mkv", // Library path is primary
			},
		},
	}

	// Radarr reports the library path
	radarrID := 1
	radarrMedia := []*models.Media{
		{
			FilePath: "/Peliculas/Movie.mkv",
			Title:    "Test Movie",
			RadarrID: &radarrID,
			Quality:  "Bluray-1080p",
		},
	}

	// Enrich
	count := enricher.EnrichWithRadarr(context.Background(), files, radarrMedia)

	// Verify enrichment happened
	if count != 1 {
		t.Errorf("Expected 1 file enriched, got %d", count)
	}

	file := files["/Descargas/Movie.mkv"]
	if !file.InRadarr {
		t.Error("File should be marked as InRadarr")
	}
	if file.Title != "Test Movie" {
		t.Errorf("Expected title 'Test Movie', got '%s'", file.Title)
	}
	if file.Quality != "Bluray-1080p" {
		t.Errorf("Expected quality 'Bluray-1080p', got '%s'", file.Quality)
	}
}

func TestEnricher_EnrichWithSonarr_HardlinkMatching(t *testing.T) {
	logger := zap.NewNop()
	enricher := NewEnricher(logger)

	// Simulate a series file with hardlinks
	files := map[string]*EnrichedFile{
		"/Descargas/Series/Show.S01E01.mkv": {
			FileEntry: &FileEntry{
				Path:          "/Descargas/Series/Show.S01E01.mkv",
				Size:          500 * 1024 * 1024,
				Inode:         54321,
				IsHardlink:    true,
				HardlinkPaths: []string{"/Descargas/Series/Show.S01E01.mkv", "/Series/Show/S01/Show.S01E01.mkv"},
				PrimaryPath:   "/Series/Show/S01/Show.S01E01.mkv",
			},
		},
	}

	// Sonarr reports the organized library path
	sonarrID := 10
	sonarrMedia := []*models.Media{
		{
			FilePath: "/Series/Show/S01/Show.S01E01.mkv",
			Title:    "Test Show",
			SonarrID: &sonarrID,
			Quality:  "WEBDL-720p",
		},
	}

	// Enrich
	count := enricher.EnrichWithSonarr(context.Background(), files, sonarrMedia)

	// Verify
	if count != 1 {
		t.Errorf("Expected 1 file enriched, got %d", count)
	}

	file := files["/Descargas/Series/Show.S01E01.mkv"]
	if !file.InSonarr {
		t.Error("File should be marked as InSonarr")
	}
	if file.Title != "Test Show" {
		t.Errorf("Expected title 'Test Show', got '%s'", file.Title)
	}
}

func TestEnricher_EnrichWithJellyfin_HardlinkMatching(t *testing.T) {
	logger := zap.NewNop()
	enricher := NewEnricher(logger)

	// Simulate file with multiple hardlinks
	files := map[string]*EnrichedFile{
		"/qBittorrent/downloads/Movie.mkv": {
			FileEntry: &FileEntry{
				Path:       "/qBittorrent/downloads/Movie.mkv",
				Size:       2 * 1024 * 1024 * 1024,
				Inode:      99999,
				IsHardlink: true,
				HardlinkPaths: []string{
					"/qBittorrent/downloads/Movie.mkv",
					"/media/Movies/Movie (2024)/Movie.mkv",
				},
				PrimaryPath: "/media/Movies/Movie (2024)/Movie.mkv",
			},
		},
	}

	// Jellyfin scans its library path
	jellyfinID := "abc123"
	jellyfinMedia := []*models.Media{
		{
			FilePath:   "/media/Movies/Movie (2024)/Movie.mkv",
			Title:      "Movie (2024)",
			JellyfinID: &jellyfinID,
			PosterURL:  "https://example.com/poster.jpg",
		},
	}

	// Enrich
	count := enricher.EnrichWithJellyfin(context.Background(), files, jellyfinMedia)

	// Verify
	if count != 1 {
		t.Errorf("Expected 1 file enriched, got %d", count)
	}

	file := files["/qBittorrent/downloads/Movie.mkv"]
	if !file.InJellyfin {
		t.Error("File should be marked as InJellyfin")
	}
	if file.Title != "Movie (2024)" {
		t.Errorf("Expected title 'Movie (2024)', got '%s'", file.Title)
	}
	if file.JellyfinID == nil || *file.JellyfinID != "abc123" {
		t.Error("JellyfinID should be set correctly")
	}
}

func TestEnricher_EnrichWithQBittorrent_HardlinkMatching(t *testing.T) {
	logger := zap.NewNop()
	enricher := NewEnricher(logger)

	// File in library but torrent in downloads
	files := map[string]*EnrichedFile{
		"/BibliotecaMultimedia/Peliculas/Movie.mkv": {
			FileEntry: &FileEntry{
				Path:       "/BibliotecaMultimedia/Peliculas/Movie.mkv",
				Size:       1500 * 1024 * 1024,
				Inode:      77777,
				IsHardlink: true,
				HardlinkPaths: []string{
					"/BibliotecaMultimedia/Descargas/Movie.mkv",
					"/BibliotecaMultimedia/Peliculas/Movie.mkv",
				},
				PrimaryPath: "/BibliotecaMultimedia/Peliculas/Movie.mkv",
			},
		},
	}

	// qBittorrent reports the download path
	torrentMap := map[string]*models.TorrentInfo{
		"/BibliotecaMultimedia/Descargas/Movie.mkv": {
			Hash:      "abc123def456",
			IsSeeding: true,
			Ratio:     2.5,
			Category:  "movies",
			State:     "uploading",
		},
	}

	// Enrich
	count := enricher.EnrichWithQBittorrent(context.Background(), files, torrentMap)

	// Verify
	if count != 1 {
		t.Errorf("Expected 1 file enriched, got %d", count)
	}

	file := files["/BibliotecaMultimedia/Peliculas/Movie.mkv"]
	if !file.InQBittorrent {
		t.Error("File should be marked as InQBittorrent")
	}
	if file.TorrentHash != "abc123def456" {
		t.Errorf("Expected hash 'abc123def456', got '%s'", file.TorrentHash)
	}
	if !file.IsSeeding {
		t.Error("File should be marked as seeding")
	}
	if file.SeedRatio != 2.5 {
		t.Errorf("Expected ratio 2.5, got %f", file.SeedRatio)
	}
}

func TestEnricher_MultiServiceHardlinkMatching(t *testing.T) {
	logger := zap.NewNop()
	enricher := NewEnricher(logger)

	// Single file with hardlinks, enriched by multiple services
	files := map[string]*EnrichedFile{
		"/storage/downloads/Movie.2024.mkv": {
			FileEntry: &FileEntry{
				Path:       "/storage/downloads/Movie.2024.mkv",
				Size:       3 * 1024 * 1024 * 1024,
				Inode:      11111,
				IsHardlink: true,
				HardlinkPaths: []string{
					"/storage/downloads/Movie.2024.mkv",
					"/storage/movies/Movie (2024).mkv",
					"/media/jellyfin/Movies/Movie (2024).mkv",
				},
				PrimaryPath: "/storage/movies/Movie (2024).mkv",
			},
		},
	}

	// Each service sees different path
	radarrID := 42
	radarrMedia := []*models.Media{
		{
			FilePath: "/storage/movies/Movie (2024).mkv",
			Title:    "Movie",
			RadarrID: &radarrID,
			Quality:  "Bluray-2160p",
		},
	}

	jellyfinID := "jf-999"
	jellyfinMedia := []*models.Media{
		{
			FilePath:   "/media/jellyfin/Movies/Movie (2024).mkv",
			JellyfinID: &jellyfinID,
		},
	}

	torrentMap := map[string]*models.TorrentInfo{
		"/storage/downloads/Movie.2024.mkv": {
			Hash:      "torrent-hash-123",
			IsSeeding: true,
			Ratio:     3.0,
		},
	}

	// Enrich with all services
	radarrCount := enricher.EnrichWithRadarr(context.Background(), files, radarrMedia)
	jellyfinCount := enricher.EnrichWithJellyfin(context.Background(), files, jellyfinMedia)
	qbitCount := enricher.EnrichWithQBittorrent(context.Background(), files, torrentMap)

	// Verify all enrichments succeeded
	if radarrCount != 1 {
		t.Errorf("Expected Radarr to enrich 1 file, got %d", radarrCount)
	}
	if jellyfinCount != 1 {
		t.Errorf("Expected Jellyfin to enrich 1 file, got %d", jellyfinCount)
	}
	if qbitCount != 1 {
		t.Errorf("Expected qBittorrent to enrich 1 file, got %d", qbitCount)
	}

	// Verify file has all service flags
	file := files["/storage/downloads/Movie.2024.mkv"]
	if !file.InRadarr {
		t.Error("File should be in Radarr")
	}
	if !file.InJellyfin {
		t.Error("File should be in Jellyfin")
	}
	if !file.InQBittorrent {
		t.Error("File should be in qBittorrent")
	}

	// Verify metadata
	if file.Title != "Movie" {
		t.Errorf("Expected title 'Movie', got '%s'", file.Title)
	}
	if file.Quality != "Bluray-2160p" {
		t.Errorf("Expected quality 'Bluray-2160p', got '%s'", file.Quality)
	}
	if file.TorrentHash != "torrent-hash-123" {
		t.Errorf("Expected hash 'torrent-hash-123', got '%s'", file.TorrentHash)
	}
}

func TestEnricher_NoHardlinks_ExactMatch(t *testing.T) {
	logger := zap.NewNop()
	enricher := NewEnricher(logger)

	// File without hardlinks - should still work with exact match
	files := map[string]*EnrichedFile{
		"/movies/SingleFile.mkv": {
			FileEntry: &FileEntry{
				Path:        "/movies/SingleFile.mkv",
				Size:        1024 * 1024 * 1024,
				Inode:       33333,
				IsHardlink:  false,
				PrimaryPath: "/movies/SingleFile.mkv",
			},
		},
	}

	radarrID := 5
	radarrMedia := []*models.Media{
		{
			FilePath: "/movies/SingleFile.mkv",
			Title:    "Single File Movie",
			RadarrID: &radarrID,
		},
	}

	count := enricher.EnrichWithRadarr(context.Background(), files, radarrMedia)

	if count != 1 {
		t.Errorf("Expected 1 file enriched, got %d", count)
	}

	file := files["/movies/SingleFile.mkv"]
	if !file.InRadarr {
		t.Error("File should be in Radarr")
	}
	if file.Title != "Single File Movie" {
		t.Errorf("Expected title 'Single File Movie', got '%s'", file.Title)
	}
}
