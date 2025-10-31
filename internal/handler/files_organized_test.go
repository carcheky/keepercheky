package handler

import (
	"testing"
)

func TestParseSeasonEpisode(t *testing.T) {
	tests := []struct {
		filename string
		wantSeason int
		wantEpisode int
		wantFound bool
	}{
		{
			filename: "Breaking.Bad.S01E01.720p.mkv",
			wantSeason: 1,
			wantEpisode: 1,
			wantFound: true,
		},
		{
			filename: "Game.of.Thrones.s05e08.1080p.mp4",
			wantSeason: 5,
			wantEpisode: 8,
			wantFound: true,
		},
		{
			filename: "The.Office.1x01.Pilot.mp4",
			wantSeason: 1,
			wantEpisode: 1,
			wantFound: true,
		},
		{
			filename: "Friends.2X10.The.One.mkv",
			wantSeason: 2,
			wantEpisode: 10,
			wantFound: true,
		},
		{
			filename: "Movie.Title.2023.1080p.mkv",
			wantSeason: 0,
			wantEpisode: 0,
			wantFound: false,
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
