package handler

import (
	"testing"

	"go.uber.org/zap"
)

func TestParseSeasonEpisode(t *testing.T) {
	tests := []struct {
		filename    string
		wantSeason  int
		wantEpisode int
		wantFound   bool
	}{
		{
			filename:    "Breaking.Bad.S01E01.720p.mkv",
			wantSeason:  1,
			wantEpisode: 1,
			wantFound:   true,
		},
		{
			filename:    "Game.of.Thrones.s05e08.1080p.mp4",
			wantSeason:  5,
			wantEpisode: 8,
			wantFound:   true,
		},
		{
			filename:    "The.Office.1x01.Pilot.mp4",
			wantSeason:  1,
			wantEpisode: 1,
			wantFound:   true,
		},
		{
			filename:    "Friends.2X10.The.One.mkv",
			wantSeason:  2,
			wantEpisode: 10,
			wantFound:   true,
		},
		{
			filename:    "Movie.Title.2023.1080p.mkv",
			wantSeason:  0,
			wantEpisode: 0,
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			season, episode, found := parseSeasonEpisode(tt.filename)

			if found != tt.wantFound {
				t.Errorf("parseSeasonEpisode(%q) found = %v, want %v", tt.filename, found, tt.wantFound)
			}

			if found && (season != tt.wantSeason || episode != tt.wantEpisode) {
				t.Errorf("parseSeasonEpisode(%q) = (%d, %d), want (%d, %d)",
					tt.filename, season, episode, tt.wantSeason, tt.wantEpisode)
			}
		})
	}
}

func TestExtractSeriesName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{
			path: "/media/tv/Breaking.Bad.S01E01.720p.mkv",
			want: "Breaking Bad",
		},
		{
			path: "/media/tv/Game_of_Thrones_s05e08_1080p.mp4",
			want: "Game of Thrones",
		},
		{
			path: "/downloads/The.Office.1x01.Pilot.1080p.BluRay.x264.mp4",
			want: "The Office",
		},
		{
			path: "/series/Friends.2X10.The.One.720p.WEB-DL.mkv",
			want: "Friends",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := extractSeriesName(tt.path)
			if got != tt.want {
				t.Errorf("extractSeriesName(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestOrganizeFiles_TypeDetection(t *testing.T) {
	// Create a mock logger
	logger, _ := zap.NewDevelopment()
	handler := &FilesHandler{logger: logger}

	tests := []struct {
		name           string
		files          []MediaFileInfo
		expectedSeries int
		expectedMovies int
	}{
		{
			name: "explicit types",
			files: []MediaFileInfo{
				{FilePath: "/data/movies/Movie1.mkv", Type: "movie", Title: "Movie 1"},
				{FilePath: "/data/tv/Series1.S01E01.mkv", Type: "series", Title: "Series 1"},
				{FilePath: "/data/tv/Series1.S01E02.mkv", Type: "episode", Title: "Series 1"},
			},
			expectedSeries: 1, // Series1 with 2 episodes
			expectedMovies: 1,
		},
		{
			name: "episode pattern detection",
			files: []MediaFileInfo{
				{FilePath: "/downloads/Breaking.Bad.S01E01.720p.mkv", Type: "", Title: ""},
				{FilePath: "/downloads/Breaking.Bad.S01E02.720p.mkv", Type: "", Title: ""},
				{FilePath: "/downloads/Movie.2023.1080p.mkv", Type: "", Title: ""},
			},
			expectedSeries: 1, // Breaking Bad
			expectedMovies: 1,
		},
		{
			name: "service flags detection",
			files: []MediaFileInfo{
				{FilePath: "/media/file1.mkv", Type: "", InSonarr: true, Title: "TV Show"},
				{FilePath: "/media/file2.mkv", Type: "", InRadarr: true, Title: "Movie"},
			},
			expectedSeries: 1,
			expectedMovies: 1,
		},
		{
			name: "path-based detection",
			files: []MediaFileInfo{
				{FilePath: "/media/tv/show.mkv", Type: "", Title: "Show"},
				{FilePath: "/media/series/another.mkv", Type: "", Title: "Another"},
				{FilePath: "/media/movies/film.mkv", Type: "", Title: "Film"},
			},
			expectedSeries: 2,
			expectedMovies: 1,
		},
		{
			name: "mixed scenarios",
			files: []MediaFileInfo{
				// Explicit series type
				{FilePath: "/data/shows/Show1.mkv", Type: "series", Title: "Show 1"},
				// Episode pattern
				{FilePath: "/downloads/Show2.S01E01.mkv", Type: "", Title: ""},
				// In Sonarr
				{FilePath: "/library/show3.mkv", Type: "", InSonarr: true, Title: "Show 3"},
				// Explicit movie type
				{FilePath: "/data/films/Movie1.mkv", Type: "movie", Title: "Movie 1"},
				// In Radarr
				{FilePath: "/library/movie2.mkv", Type: "", InRadarr: true, Title: "Movie 2"},
				// Default to movie
				{FilePath: "/random/file.mkv", Type: "", Title: "Random"},
			},
			expectedSeries: 3,
			expectedMovies: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.organizeFiles(tt.files)

			if len(result.Series) != tt.expectedSeries {
				t.Errorf("Expected %d series, got %d", tt.expectedSeries, len(result.Series))
			}

			if len(result.Movies) != tt.expectedMovies {
				t.Errorf("Expected %d movies, got %d", tt.expectedMovies, len(result.Movies))
			}
		})
	}
}
