package web

import (
	"net/http"

	auth "github.com/datasektionen/GO-auth-yourself"
)

func (s *Server) indexPage(w http.ResponseWriter, r *http.Request) error {
	info, loggedIn := auth.FromContext(r.Context())

	return s.tmpl.ExecuteTemplate(w, "index.gohtml", map[string]any{
		"title":    "Hem",
		"auth":     info,
		"loggedIn": loggedIn,
	})
}
