package web

import (
	"encoding/json"
	"net/http"
	"time"

	auth "github.com/datasektionen/GO-auth-yourself"
)

func (s *Server) apiHealthHandler(w http.ResponseWriter, r *http.Request) error {
	dbStatus := "disabled"
	if s.db != nil {
		if err := s.db.Ping(); err == nil {
			dbStatus = "connected"
		} else {
			dbStatus = "error"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]any{
		"status":   "ok",
		"time":     time.Now().UTC().Format(time.RFC3339),
		"database": dbStatus,
	})
}

func (s *Server) apiUserHandler(w http.ResponseWriter, r *http.Request) error {
	info, loggedIn := auth.FromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]any{
		"loggedIn":    loggedIn,
		"user":        info.Username,
		"name":        info.Name,
		"email":       info.Email,
		"permissions": info.Permissions,
	})
}
