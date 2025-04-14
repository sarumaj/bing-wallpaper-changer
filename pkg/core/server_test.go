package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupController(t *testing.T, cfg *Config, executed *bool) *Controller {
	t.Helper()
	img := &Image{}
	return &Controller{
		img: img,
		cfg: cfg,
		execute: func(cfg *Config) *Image {
			if executed != nil {
				*executed = true
			}
			return img
		},
	}
}

func TestHandleConfigGET(t *testing.T) {
	cfg := &Config{ApiPort: 8080}
	controller := setupController(t, cfg, nil)
	server := NewServer(cfg, controller)

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	w := httptest.NewRecorder()

	server.handleConfig(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response Config
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ApiPort != cfg.ApiPort {
		t.Errorf("Expected ApiPort %d, got %d", cfg.ApiPort, response.ApiPort)
	}
}

func TestHandleConfigPATCH(t *testing.T) {
	cfg := &Config{ApiPort: 8080}
	controller := setupController(t, cfg, nil)
	server := NewServer(cfg, controller)

	newConfig := Config{
		ApiPort: 9090,
	}
	body, _ := json.Marshal(newConfig)

	req := httptest.NewRequest(http.MethodPatch, "/config", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	server.handleConfig(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, got %d", http.StatusAccepted, w.Code)
	}

	if cfg.ApiPort != newConfig.ApiPort {
		t.Errorf("Expected ApiPort to be updated to %d, got %d", newConfig.ApiPort, cfg.ApiPort)
	}
}

func TestHandleConfigPATCHWithRefresh(t *testing.T) {
	cfg := &Config{ApiPort: 8080}
	executed := false
	controller := setupController(t, cfg, &executed)
	server := NewServer(cfg, controller)

	newConfig := Config{
		ApiPort: 9090,
	}
	body, _ := json.Marshal(newConfig)

	req := httptest.NewRequest(http.MethodPatch, "/config?refresh=true", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	server.handleConfig(w, req)

	// Give some time for the goroutine to execute
	time.Sleep(100 * time.Millisecond)

	if !executed {
		t.Error("Expected OnReady to be called")
	}
}

func TestHandleConfigInvalidMethod(t *testing.T) {
	cfg := &Config{}
	controller := setupController(t, cfg, nil)
	server := NewServer(cfg, controller)

	req := httptest.NewRequest(http.MethodPost, "/config", nil)
	w := httptest.NewRecorder()

	server.handleConfig(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHandleRoot(t *testing.T) {
	cfg := &Config{}
	controller := setupController(t, cfg, nil)
	server := NewServer(cfg, controller)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.handleRoot(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
