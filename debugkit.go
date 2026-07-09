package debugkit

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
)

//go:embed ui/index.html
var uiFS embed.FS

// New initializes the global Collector instance
func New() {
	// collector.go instance initialization
	Instance = &Collector{
		Requests:   make([]RequestLog, 0),
		SQLQueries: make([]SQLQueryLog, 0),
	}
}

// RegisterUIHandlers mounts the premium dashboard and telemetry API
func RegisterUIHandlers(mux *http.ServeMux) {
	// 1. Extract the ui sub-directory from the embedded uiFS
	subFS, err := fs.Sub(uiFS, "ui")
	if err != nil {
		panic("DebugKit: Failed to load embedded UI filesystem: " + err.Error())
	}

	// 2. Create a file server (this now points directly to the root of the ui folder, so it can directly serve index.html)
	fileServer := http.FileServer(http.FS(subFS))

	// ৩. UI Route Dashboard: /debugkit বা /debugkit/ both should serve the index.html
	mux.HandleFunc("/debugkit", func(w http.ResponseWriter, r *http.Request) {
		// Redirect to /debugkit/ if the request is for /debugkit without a trailing slash
		if r.URL.Path == "/debugkit" {
			http.Redirect(w, r, "/debugkit/", http.StatusMovedPermanently)
			return
		}
		// Strip the /debugkit/ prefix and serve the file from the embedded filesystem
		http.StripPrefix("/debugkit/", fileServer).ServeHTTP(w, r)
	})

	// 4. Serve static assets
	mux.Handle("/debugkit/", http.StripPrefix("/debugkit/", fileServer))

	// ৫. Telemetry Metrics API Endpoint
	mux.HandleFunc("/debugkit/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if Instance == nil {
			http.Error(w, `{"error": "DebugKit instance not initialized. Call debugkit.New() first."}`, http.StatusInternalServerError)
			return
		}

		stats := Instance.GetFullStats()
		json.NewEncoder(w).Encode(stats)
	})
}

