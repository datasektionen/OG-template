package web

import (
	"io"
	"testing"

	auth "github.com/datasektionen/GO-auth-yourself"
)

func TestTemplatesRender(t *testing.T) {
	tmpl, err := parseTemplates()
	if err != nil {
		t.Fatal(err)
	}

	users := map[string]auth.User{
		"anonymous": {},
		"admin":     {Username: "turetek", Name: "Ture Teknolog", Permissions: []auth.Permission{{ID: "admin"}}},
	}
	for name, info := range users {
		for _, page := range []string{"index.gohtml", "admin.gohtml"} {
			err := tmpl.ExecuteTemplate(io.Discard, page, map[string]any{
				"title":    "Test",
				"auth":     info,
				"loggedIn": info.Username != "",
				"hasDB":    false,
			})
			if err != nil {
				t.Errorf("%s as %s: %v", page, name, err)
			}
		}
	}
}
