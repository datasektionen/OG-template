package web

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	auth "github.com/datasektionen/GO-auth-yourself"
	"github.com/datasektionen/OG-template/internal/config"
)

//go:embed templates/*.gohtml
var templatesFS embed.FS

//go:embed static/*
var staticFiles embed.FS

// Server holds the application dependencies and handlers.
type Server struct {
	db   *sql.DB
	tmpl *template.Template
	env  config.EnvVar
}

// NewHandler initializes the Server and returns a handler serving all its routes.
func NewHandler(db *sql.DB, env config.EnvVar) (http.Handler, error) {
	ctx := context.Background()

	// OIDC settings and APP_SECRET_KEY are read and validated by GO-auth-yourself
	authCfg := auth.ConfigFromEnv()
	authenticator, err := auth.New(ctx, authCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth: %w", err)
	}
	slog.Info("OIDC authentication initialized", "provider", authCfg.ProviderURL)

	tmpl, err := parseTemplates()
	if err != nil {
		return nil, err
	}

	s := &Server{
		db:   db,
		tmpl: tmpl,
		env:  env,
	}

	mux := http.NewServeMux()

	// Static assets
	mux.Handle("/static/", http.FileServerFS(staticFiles))

	// Mount default auth routes (/login, /logout, /auth/callback)
	authenticator.MountAuthRoutes(mux)

	// Web pages
	mux.Handle("GET /{$}", s.route(s.indexPage))

	// Admin protected routes
	requireAdmin := authenticator.RequirePermissions("admin")
	mux.Handle("GET /admin", requireAdmin(s.route(s.adminPage)))
	mux.Handle("POST /admin/items", requireAdmin(s.route(s.adminCreateItem)))
	mux.Handle("POST /admin/items/delete", requireAdmin(s.route(s.adminDeleteItem)))

	// API endpoints
	mux.Handle("GET /api/health", s.route(s.apiHealthHandler))
	mux.Handle("GET /api/user", s.route(s.apiUserHandler))

	// Makes the logged-in user available to every route through auth.FromContext.
	handler := authenticator.SessionMiddleware(mux)
	// SameSite=Lax cookies still allow form posts from other *.datasektionen.se sites.
	return http.NewCrossOriginProtection().Handler(handler), nil
}

func parseTemplates() (*template.Template, error) {
	funcs := template.FuncMap{
		"year": func() int {
			return time.Now().Year()
		},
	}
	tmpl, err := template.New("").Funcs(funcs).ParseFS(templatesFS, "templates/*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return tmpl, nil
}

// --- Middleware & Route Helpers ---

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (s *Server) route(handler HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			slog.Error("Handler error", "path", r.URL.Path, "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
