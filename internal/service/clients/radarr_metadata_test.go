package clients

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRadarrClient_ConvertToMetadata(t *testing.T) {
	logger := zap.NewNop()
	client := NewRadarrClient(ClientConfig{
		BaseURL: "http://localhost:7878",
		APIKey:  "test-api-key",
		Timeout: 5 * time.Second,
	}, logger)

	// Create a comprehensive mock movie with all fields
	inCinemas := time.Date(1999, 3, 24, 0, 0, 0, 0, time.UTC)
	physicalRelease := time.Date(1999, 9, 21, 0, 0, 0, 0, time.UTC)
	digitalRelease := time.Date(1999, 9, 21, 0, 0, 0, 0, time.UTC)
	added := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)

	movie := &radarrMovie{
		ID:              1,
		Title:           "The Matrix",
		OriginalTitle:   "The Matrix",
		SortTitle:       "matrix",
		Year:            1999,
		Status:          "released",
		Overview:        "Set in the 22nd century...",
		Monitored:       true,
		InCinemas:       &inCinemas,
		PhysicalRelease: &physicalRelease,
		DigitalRelease:  &digitalRelease,
		Added:           added,
		Path:            "/movies/The Matrix (1999)",
		SizeOnDisk:      2147483648,
		HasFile:         true,
		Runtime:         136,
		IMDbID:          "tt0133093",
		TMDbID:          603,
		Certification:   "R",
		Genres:          []string{"Action", "Science Fiction"},
		Popularity:      82.458,
	}

	// Add images
	movie.Images = []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
		RemoteURL string `json:"remoteUrl"`
	}{
		{CoverType: "poster", URL: "http://localhost/poster.jpg"},
		{CoverType: "fanart", URL: "http://localhost/fanart.jpg"},
	}

	// Add alternate titles
	movie.AlternateTitles = []struct {
		SourceType      string `json:"sourceType"`
		MovieMetadataID int    `json:"movieMetadataId"`
		Title           string `json:"title"`
		ID              int    `json:"id"`
	}{
		{Title: "Matrix"},
	}

	// Add ratings
	movie.Ratings.IMDb = &struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
		Type  string  `json:"type"`
	}{Votes: 1800000, Value: 8.7}
	movie.Ratings.TMDb = &struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
		Type  string  `json:"type"`
	}{Votes: 24000, Value: 8.2}
	movie.Ratings.Metacritic = &struct {
		Votes int    `json:"votes"`
		Value int    `json:"value"`
		Type  string `json:"type"`
	}{Value: 73}

	// Add movie file with complete details - we'll initialize the struct inline
	movieFileStruct := struct {
		MovieID      int       `json:"movieId"`
		RelativePath string    `json:"relativePath"`
		Path         string    `json:"path"`
		Size         int64     `json:"size"`
		DateAdded    time.Time `json:"dateAdded"`
		SceneName    string    `json:"sceneName"`
		IndexerFlags int       `json:"indexerFlags"`
		Quality      struct {
			Quality struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Source     string `json:"source"`
				Resolution int    `json:"resolution"`
				Modifier   string `json:"modifier"`
			} `json:"quality"`
			Revision struct {
				Version  int  `json:"version"`
				Real     int  `json:"real"`
				IsRepack bool `json:"isRepack"`
			} `json:"revision"`
		} `json:"quality"`
		CustomFormatScore int `json:"customFormatScore"`
		CustomFormats     []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"customFormats"`
		MediaInfo *struct {
			AudioBitrate          int     `json:"audioBitrate"`
			AudioChannels         float64 `json:"audioChannels"`
			AudioCodec            string  `json:"audioCodec"`
			AudioLanguages        string  `json:"audioLanguages"`
			AudioStreamCount      int     `json:"audioStreamCount"`
			VideoBitDepth         int     `json:"videoBitDepth"`
			VideoBitrate          int     `json:"videoBitrate"`
			VideoCodec            string  `json:"videoCodec"`
			VideoFps              float64 `json:"videoFps"`
			VideoDynamicRange     string  `json:"videoDynamicRange"`
			VideoDynamicRangeType string  `json:"videoDynamicRangeType"`
			Resolution            string  `json:"resolution"`
			RunTime               string  `json:"runTime"`
			ScanType              string  `json:"scanType"`
			Subtitles             string  `json:"subtitles"`
		} `json:"mediaInfo"`
		QualityCutoffNotMet bool `json:"qualityCutoffNotMet"`
		Languages           []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"languages"`
		ReleaseGroup string `json:"releaseGroup"`
		Edition      string `json:"edition"`
		ID           int    `json:"id"`
	}{
		MovieID:      1,
		Path:         "/movies/The Matrix (1999)/The Matrix (1999).mkv",
		Size:         2147483648,
		SceneName:    "The.Matrix.1999.1080p.BluRay.x264-GROUP",
		ReleaseGroup: "GROUP",
	}

	// Set quality info
	movieFileStruct.Quality.Quality.Name = "Bluray-1080p"
	movieFileStruct.Quality.Quality.Source = "bluray"
	movieFileStruct.Quality.Quality.Resolution = 1080
	movieFileStruct.Quality.Revision.Version = 1

	// Set media info
	movieFileStruct.MediaInfo = &struct {
		AudioBitrate          int     `json:"audioBitrate"`
		AudioChannels         float64 `json:"audioChannels"`
		AudioCodec            string  `json:"audioCodec"`
		AudioLanguages        string  `json:"audioLanguages"`
		AudioStreamCount      int     `json:"audioStreamCount"`
		VideoBitDepth         int     `json:"videoBitDepth"`
		VideoBitrate          int     `json:"videoBitrate"`
		VideoCodec            string  `json:"videoCodec"`
		VideoFps              float64 `json:"videoFps"`
		VideoDynamicRange     string  `json:"videoDynamicRange"`
		VideoDynamicRangeType string  `json:"videoDynamicRangeType"`
		Resolution            string  `json:"resolution"`
		RunTime               string  `json:"runTime"`
		ScanType              string  `json:"scanType"`
		Subtitles             string  `json:"subtitles"`
	}{
		AudioCodec:        "DTS",
		AudioChannels:     5.1,
		VideoCodec:        "h264",
		VideoFps:          23.976,
		VideoBitDepth:     8,
		VideoDynamicRange: "SDR",
		Resolution:        "1920x1080",
		ScanType:          "Progressive",
		Subtitles:         "eng, spa",
	}

	// Set languages
	movieFileStruct.Languages = []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{{ID: 1, Name: "English"}}

	movie.MovieFile = &movieFileStruct

	// Convert to metadata
	metadata := client.ConvertToMetadata(movie)

	// Verify basic info
	assert.Equal(t, 1, metadata.ID)
	assert.Equal(t, "The Matrix", metadata.Title)
	assert.Equal(t, "The Matrix", metadata.OriginalTitle)
	assert.Equal(t, 1999, metadata.Year)
	assert.Equal(t, "released", metadata.Status)
	assert.True(t, metadata.Monitored)

	// Verify dates
	assert.Equal(t, "1999-03-24", metadata.InCinemas)
	assert.Equal(t, "1999-09-21", metadata.PhysicalRelease)
	assert.Equal(t, "2023-01-15", metadata.DateAdded)

	// Verify external IDs
	assert.Equal(t, "tt0133093", metadata.IMDbID)
	assert.Equal(t, 603, metadata.TMDbID)

	// Verify genres
	assert.Contains(t, metadata.Genres, "Action")
	assert.Contains(t, metadata.Genres, "Science Fiction")

	// Verify images
	assert.Equal(t, "http://localhost/poster.jpg", metadata.PosterURL)
	assert.Equal(t, "http://localhost/fanart.jpg", metadata.FanartURL)

	// Verify alternate titles
	assert.Contains(t, metadata.AlternateTitles, "Matrix")

	// Verify ratings
	assert.Equal(t, 8.7, metadata.IMDbRating)
	assert.Equal(t, 1800000, metadata.IMDbVotes)
	assert.Equal(t, 8.2, metadata.TMDbRating)
	assert.Equal(t, 73, metadata.MetacriticScore)

	// Verify quality info
	assert.Equal(t, "Bluray-1080p", metadata.QualityName)
	assert.Equal(t, "bluray", metadata.QualitySource)
	assert.Equal(t, 1080, metadata.QualityResolution)
	assert.Equal(t, 1, metadata.QualityVersion)

	// Verify media info
	require.NotNil(t, metadata.MediaInfo)
	assert.Equal(t, "DTS", metadata.MediaInfo.AudioCodec)
	assert.Equal(t, 5.1, metadata.MediaInfo.AudioChannels)
	assert.Equal(t, "h264", metadata.MediaInfo.VideoCodec)
	assert.Equal(t, 23.976, metadata.MediaInfo.VideoFps)
	assert.Equal(t, "SDR", metadata.MediaInfo.VideoDynamicRange)
	assert.Equal(t, "1920x1080", metadata.MediaInfo.Resolution)
	assert.Equal(t, "Progressive", metadata.MediaInfo.ScanType)

	// Verify file info
	assert.Equal(t, "GROUP", metadata.ReleaseGroup)
	assert.Contains(t, metadata.Languages, "English")

	t.Log("ConvertToMetadata successfully captures all fields from Radarr API")
}
