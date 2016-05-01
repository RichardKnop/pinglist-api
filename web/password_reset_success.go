package web

import (
	"net/http"
)

func (s *Service) passwordResetSuccess(w http.ResponseWriter, r *http.Request) {
	// Render the template
	s.renderTemplate(w, "password-reset-success.html", map[string]interface{}{
		"hideLogout": true,
	})
}
