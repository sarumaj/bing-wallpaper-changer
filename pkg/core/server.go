package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
)

type Server struct {
	config     *Config
	controller *Controller
	updateLock sync.Mutex
	server     *http.Server
}

// NewServer creates a new server.
func NewServer(config *Config, controller *Controller) *Server {
	return &Server{
		config:     config,
		controller: controller,
	}
}

// Start starts the server.
func (s *Server) Start() error {
	router := http.NewServeMux()
	router.HandleFunc("/config", s.handleConfig)
	router.HandleFunc("/", s.handleRoot)

	logger.Logger.Printf("Starting API server on port %d", s.config.ApiPort)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.ApiPort),
		Handler: router,
	}
	return s.server.ListenAndServe()
}

// Stop stops the server.
func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

// handleConfig handles the config endpoint.
// It returns the current config when GET request is made.
// It updates the config when PATCH request is made.
// It refreshes the wallpaper when PATCH request with query parameter refresh=true is made.
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		_ = json.NewEncoder(w).Encode(s.config)

	case http.MethodPatch:
		s.updateLock.Lock()
		defer s.updateLock.Unlock()

		target := reflect.ValueOf(s.config).Elem()

		// Create a new struct with the same fields as the original config
		// but with pointers to the fields instead of the fields themselves
		// so that we know which fields were provided in the request
		fields := make([]reflect.StructField, 0, target.NumField())
		for i := range target.NumField() {
			field := target.Type().Field(i)
			fields = append(fields, reflect.StructField{
				Name: field.Name,
				Type: reflect.PointerTo(field.Type),
				Tag:  field.Tag,
			})
		}
		source := reflect.New(reflect.StructOf(fields))

		if err := json.NewDecoder(r.Body).Decode(source.Interface()); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update the original config with the new values
		updatedFields := 0
		for i := range target.NumField() {
			sourceField := source.Elem().Field(i)
			targetField := target.Field(i)
			if sourceField.IsNil() {
				continue
			}

			if !reflect.DeepEqual(targetField.Interface(), sourceField.Elem().Interface()) {
				targetField.Set(sourceField.Elem())
				updatedFields++
			}
		}

		query := r.URL.Query()
		if result, err := strconv.ParseBool(query.Get("refresh")); updatedFields > 0 && err == nil && result {
			go func() {
				s.updateLock.Lock()
				defer s.updateLock.Unlock()

				autoPlayAudio := s.config.AutoPlayAudio
				s.config.AutoPlayAudio = false
				s.controller.OnReady()
				s.config.AutoPlayAudio = autoPlayAudio
			}()
		}

		w.WriteHeader(http.StatusAccepted)

	default:
		logger.Logger.Printf("Method not allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed: " + r.Method})
	}

}

// handleRoot handles the root endpoint.
// It returns a 404 error when the request is not found.
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Printf("Not found: %s", r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "Not found: " + r.URL.Path})
}
