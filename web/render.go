package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/oxtoacart/bpool"
)

var (
	templates map[string]*template.Template
	bufpool   *bpool.BufferPool
)

// Load templates on program initialisation
func init() {
	templates = make(map[string]*template.Template)

	bufpool = bpool.NewBufferPool(64)

	layoutTemplates := map[string][]string{
		"web/layouts/layout.html": []string{
			"web/includes/register.html",
			"web/includes/login.html",
			"web/includes/password-reset.html",
			"web/includes/confirm-invitation.html",
			"web/includes/authorize.html",
			"web/includes/confirm-email.html",
			"web/includes/password-reset-success.html",
			"web/includes/confirm-invitation-success.html",
		},
	}

	for layout, includes := range layoutTemplates {
		for _, include := range includes {
			files := []string{include, layout}
			templates[filepath.Base(include)] = template.Must(template.ParseFiles(files...))
		}
	}
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
// It writes into a bytes.Buffer before writing to the http.ResponseWriter to catch
// any errors resulting from populating the template.
func (s *Service) renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	// Add app link to the data
	data["appLink"] = fmt.Sprintf(
		"%s://%s",
		s.cnf.Web.AppScheme,
		s.cnf.Web.AppHost,
	)

	// Execute the template
	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	// The X-Frame-Options HTTP response header can be used to indicate whether
	// or not a browser should be allowed to render a page in a <frame>,
	// <iframe> or <object> . Sites can use this to avoid clickjacking attacks,
	// by ensuring that their content is not embedded into other sites.
	w.Header().Set("X-Frame-Options", "deny")
	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
