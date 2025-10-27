// Package mocks provides mock implementations of client interfaces for testing
package mocks

import (
	"context"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockRadarrClient is a mock implementation of MediaClient for Radarr
type MockRadarrClient struct {
	mock.Mock
}

// TestConnection mocks the TestConnection method
func (m *MockRadarrClient) TestConnection(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// GetLibrary mocks the GetLibrary method
func (m *MockRadarrClient) GetLibrary(ctx context.Context) ([]*models.Media, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Media), args.Error(1)
}

// GetItem mocks the GetItem method
func (m *MockRadarrClient) GetItem(ctx context.Context, id int) (*models.Media, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Media), args.Error(1)
}

// DeleteItem mocks the DeleteItem method
func (m *MockRadarrClient) DeleteItem(ctx context.Context, id int, deleteFiles bool) error {
	args := m.Called(ctx, id, deleteFiles)
	return args.Error(0)
}

// GetTags mocks the GetTags method
func (m *MockRadarrClient) GetTags(ctx context.Context) ([]models.Tag, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Tag), args.Error(1)
}

// AddMovie mocks adding a movie to Radarr
func (m *MockRadarrClient) AddMovie(ctx context.Context, media *models.Media) (*models.Media, error) {
	args := m.Called(ctx, media)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Media), args.Error(1)
}

// MockSonarrClient is a mock implementation of MediaClient for Sonarr
type MockSonarrClient struct {
	mock.Mock
}

// TestConnection mocks the TestConnection method
func (m *MockSonarrClient) TestConnection(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// GetLibrary mocks the GetLibrary method
func (m *MockSonarrClient) GetLibrary(ctx context.Context) ([]*models.Media, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Media), args.Error(1)
}

// GetItem mocks the GetItem method
func (m *MockSonarrClient) GetItem(ctx context.Context, id int) (*models.Media, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Media), args.Error(1)
}

// DeleteItem mocks the DeleteItem method
func (m *MockSonarrClient) DeleteItem(ctx context.Context, id int, deleteFiles bool) error {
	args := m.Called(ctx, id, deleteFiles)
	return args.Error(0)
}

// GetTags mocks the GetTags method
func (m *MockSonarrClient) GetTags(ctx context.Context) ([]models.Tag, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Tag), args.Error(1)
}

// AddSeries mocks adding a series to Sonarr
func (m *MockSonarrClient) AddSeries(ctx context.Context, media *models.Media) (*models.Media, error) {
	args := m.Called(ctx, media)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Media), args.Error(1)
}

// MockJellyfinClient is a mock implementation of StreamingClient
type MockJellyfinClient struct {
	mock.Mock
}

// TestConnection mocks the TestConnection method
func (m *MockJellyfinClient) TestConnection(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// GetLibrary mocks the GetLibrary method
func (m *MockJellyfinClient) GetLibrary(ctx context.Context) ([]*models.Media, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Media), args.Error(1)
}

// GetPlaybackInfo mocks the GetPlaybackInfo method
func (m *MockJellyfinClient) GetPlaybackInfo(ctx context.Context, mediaID string) (*models.PlaybackInfo, error) {
	args := m.Called(ctx, mediaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlaybackInfo), args.Error(1)
}

// DeleteItem mocks the DeleteItem method
func (m *MockJellyfinClient) DeleteItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockQBittorrentClient is a mock implementation of qBittorrent client
type MockQBittorrentClient struct {
	mock.Mock
}

// TestConnection mocks the TestConnection method
func (m *MockQBittorrentClient) TestConnection(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// GetAllTorrentsMap mocks the GetAllTorrentsMap method
func (m *MockQBittorrentClient) GetAllTorrentsMap(ctx context.Context) (map[string]*models.TorrentInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*models.TorrentInfo), args.Error(1)
}

// DeleteTorrent mocks deleting a torrent
func (m *MockQBittorrentClient) DeleteTorrent(ctx context.Context, hash string, deleteFiles bool) error {
	args := m.Called(ctx, hash, deleteFiles)
	return args.Error(0)
}
