// server.go
package debugkit

import (
	"embed"
	"net/http"
)

// go runtime knows to embed everything in the ui folder into this variable
//go:embed ui/*
var uiAssets embed.FS

func DebugDashboardUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Read the index.html from the embedded file system
	htmlContent, err := uiAssets.ReadFile("ui/index.html")
	if err != nil {
		http.Error(w, "Dashboard UI not found", http.StatusInternalServerError)
		return
	}

	w.Write(htmlContent)
}
