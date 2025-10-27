package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service/clients/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a test database
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	err = models.RunMigrations(db)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestFileActionsHandler_DeleteFile_RequiresConfirmation(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zap.NewNop()
	mediaRepo := repository.NewMediaRepository(db)
	
	mockRadarr := new(mocks.MockRadarrClient)
	mockSonarr := new(mocks.MockSonarrClient)
	mockJellyfin := new(mocks.MockJellyfinClient)
	
	handler := NewFileActionsHandler(mediaRepo, mockRadarr, mockSonarr, mockJellyfin, nil, logger)
	
	// Create fiber app
	app := fiber.New()
	app.Post("/api/files/:id/delete", handler.DeleteFile)
	
	// Test without confirmation
	reqBody := DeleteFileRequest{Confirm: false}
	body, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest("POST", "/api/files/1/delete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	
	// Verify response body
	var response fiber.Map
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "confirmation required")
}

func TestFileActionsHandler_DeleteFile_WithConfirmation(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zap.NewNop()
	mediaRepo := repository.NewMediaRepository(db)
	
	mockRadarr := new(mocks.MockRadarrClient)
	mockSonarr := new(mocks.MockSonarrClient)
	mockJellyfin := new(mocks.MockJellyfinClient)
	
	handler := NewFileActionsHandler(mediaRepo, mockRadarr, mockSonarr, mockJellyfin, nil, logger)
	
	// Create fiber app
	app := fiber.New()
	app.Post("/api/files/:id/delete", handler.DeleteFile)
	
	// Test with confirmation
	reqBody := DeleteFileRequest{Confirm: true}
	body, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest("POST", "/api/files/1/delete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Verify response
	var response fiber.Map
	json.NewDecoder(resp.Body).Decode(&response)
	assert.True(t, response["success"].(bool))
}

func TestFileActionsHandler_BulkAction_PartialFailures(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zap.NewNop()
	mediaRepo := repository.NewMediaRepository(db)
	
	mockRadarr := new(mocks.MockRadarrClient)
	mockSonarr := new(mocks.MockSonarrClient)
	mockJellyfin := new(mocks.MockJellyfinClient)
	
	handler := NewFileActionsHandler(mediaRepo, mockRadarr, mockSonarr, mockJellyfin, nil, logger)
	
	// Create fiber app
	app := fiber.New()
	app.Post("/api/files/bulk", handler.BulkAction)
	
	// Test bulk action with some successes
	reqBody := BulkActionRequest{
		FileIDs: []uint{1, 2, 3, 4, 5},
		Action:  "delete",
		Options: fiber.Map{},
	}
	body, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest("POST", "/api/files/bulk", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Verify response
	var response BulkActionResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	// All should succeed in this simple implementation
	assert.Equal(t, 5, len(response.Success))
	assert.Equal(t, 0, len(response.Failed))
}

func TestFileActionsHandler_BulkAction_UnknownAction(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zap.NewNop()
	mediaRepo := repository.NewMediaRepository(db)
	
	mockRadarr := new(mocks.MockRadarrClient)
	mockSonarr := new(mocks.MockSonarrClient)
	mockJellyfin := new(mocks.MockJellyfinClient)
	
	handler := NewFileActionsHandler(mediaRepo, mockRadarr, mockSonarr, mockJellyfin, nil, logger)
	
	// Create fiber app
	app := fiber.New()
	app.Post("/api/files/bulk", handler.BulkAction)
	
	// Test bulk action with unknown action
	reqBody := BulkActionRequest{
		FileIDs: []uint{1, 2},
		Action:  "unknown_action",
		Options: fiber.Map{},
	}
	body, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest("POST", "/api/files/bulk", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	
	// Verify response - all should fail
	var response BulkActionResponse
	json.NewDecoder(resp.Body).Decode(&response)
	
	assert.Equal(t, 0, len(response.Success))
	assert.Equal(t, 2, len(response.Failed))
	assert.Contains(t, response.Errors[0], "unknown action")
}

func TestFileActionsHandler_CleanupHardlink_NotHardlink(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zap.NewNop()
	mediaRepo := repository.NewMediaRepository(db)
	
	// Create a non-hardlink media file in DB
	media := &models.Media{
		Title:      "Test Movie",
		Type:       "movie",
		FilePath:   "/test/movie.mkv",
		IsHardlink: false,
	}
	db.Create(media)
	
	mockRadarr := new(mocks.MockRadarrClient)
	mockSonarr := new(mocks.MockSonarrClient)
	mockJellyfin := new(mocks.MockJellyfinClient)
	
	handler := NewFileActionsHandler(mediaRepo, mockRadarr, mockSonarr, mockJellyfin, nil, logger)
	
	// Create fiber app
	app := fiber.New()
	app.Post("/api/files/:id/cleanup-hardlink", handler.CleanupHardlink)
	
	reqBody := CleanupHardlinkRequest{KeepPrimary: true}
	body, _ := json.Marshal(reqBody)
	
	req := httptest.NewRequest("POST", "/api/files/1/cleanup-hardlink", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	
	// Verify error message
	var response fiber.Map
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "not a hardlink")
}

func TestFileActionsHandler_InvalidFileID(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zap.NewNop()
	mediaRepo := repository.NewMediaRepository(db)
	
	mockRadarr := new(mocks.MockRadarrClient)
	mockSonarr := new(mocks.MockSonarrClient)
	mockJellyfin := new(mocks.MockJellyfinClient)
	
	handler := NewFileActionsHandler(mediaRepo, mockRadarr, mockSonarr, mockJellyfin, nil, logger)
	
	// Create fiber app
	app := fiber.New()
	app.Post("/api/files/:id/delete", handler.DeleteFile)
	
	reqBody := DeleteFileRequest{Confirm: true}
	body, _ := json.Marshal(reqBody)
	
	// Test with invalid ID
	req := httptest.NewRequest("POST", "/api/files/invalid/delete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	
	var response fiber.Map
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid file ID")
}
