package web

import (
	"net/http"
)

func (s *Service) confirmInvitationSuccess(w http.ResponseWriter, r *http.Request) {
	// Render the template
	renderTemplate(w, "confirm-invitation-success.html", map[string]interface{}{
		"hideLogout": true,
	})
}
