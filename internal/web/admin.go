package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	auth "github.com/datasektionen/GO-auth-yourself"
)

type Item struct {
	ID          int
	Title       string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
}

func (s *Server) adminPage(w http.ResponseWriter, r *http.Request) error {
	info, loggedIn := auth.FromContext(r.Context())

	var items []Item
	if s.db != nil {
		rows, err := s.db.QueryContext(r.Context(), "SELECT id, title, coalesce(description, ''), created_by, created_at FROM items ORDER BY id DESC")
		if err != nil {
			return fmt.Errorf("query items: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Title, &it.Description, &it.CreatedBy, &it.CreatedAt); err != nil {
				return fmt.Errorf("scan item: %w", err)
			}
			items = append(items, it)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("read items: %w", err)
		}
	}

	return s.tmpl.ExecuteTemplate(w, "admin.gohtml", map[string]any{
		"title":    "Admin",
		"auth":     info,
		"loggedIn": loggedIn,
		"items":    items,
		"hasDB":    s.db != nil,
	})
}

func (s *Server) adminCreateItem(w http.ResponseWriter, r *http.Request) error {
	if s.db == nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return nil
	}

	info, _ := auth.FromContext(r.Context())
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))

	if title != "" {
		_, err := s.db.ExecContext(r.Context(), "INSERT INTO items (title, description, created_by) VALUES ($1, $2, $3)", title, description, info.Username)
		if err != nil {
			return fmt.Errorf("insert item: %w", err)
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	return nil
}

func (s *Server) adminDeleteItem(w http.ResponseWriter, r *http.Request) error {
	if s.db == nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return nil
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid item id", http.StatusBadRequest)
		return nil
	}
	if _, err := s.db.ExecContext(r.Context(), "DELETE FROM items WHERE id = $1", id); err != nil {
		return fmt.Errorf("delete item %d: %w", id, err)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	return nil
}
