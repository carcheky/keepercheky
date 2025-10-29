package handler

import (
	"testing"

	"github.com/carcheky/keepercheky/internal/repository"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestNewFileActionsHandler verifies that the handler can be created successfully
func TestNewFileActionsHandler(t *testing.T) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Create repositories
	mediaRepo := repository.NewMediaRepository(db)
	historyRepo := repository.NewHistoryRepository(db)

	// Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create handler (with nil clients for this basic test)
	handler := NewFileActionsHandler(
		mediaRepo,
		historyRepo,
		nil, // radarrClient
		nil, // sonarrClient
		nil, // qbitClient
		nil, // jellyfinClient
		logger,
	)

	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.mediaRepo == nil {
		t.Error("Expected mediaRepo to be set")
	}

	if handler.historyRepo == nil {
		t.Error("Expected historyRepo to be set")
	}

	if handler.logger == nil {
		t.Error("Expected logger to be set")
	}
}

// TestGetBoolParam tests the helper function
func TestGetBoolParam(t *testing.T) {
	tests := []struct {
		name         string
		params       map[string]interface{}
		key          string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "returns true when param is true",
			params:       map[string]interface{}{"test": true},
			key:          "test",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "returns false when param is false",
			params:       map[string]interface{}{"test": false},
			key:          "test",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "returns default when key missing",
			params:       map[string]interface{}{},
			key:          "test",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "returns default when value is not bool",
			params:       map[string]interface{}{"test": "not a bool"},
			key:          "test",
			defaultValue: true,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBoolParam(tt.params, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getBoolParam() = %v, want %v", result, tt.expected)
			}
		})
	}
}
