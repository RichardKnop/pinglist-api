package web

import (
	"net/http"
)

func (s *Service) confirmInvitationForm(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	errMsg, _ := sessionService.GetFlashMessage()
	s.renderTemplate(w, "confirm-invitation.html", map[string]interface{}{
		"error":       errMsg,
		"queryString": getQueryString(r.URL.Query()),
	})
}

func (s *Service) confirmInvitation(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the invitation from the request context
	invitation, err := getInvitation(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check both password fields were submitted with the same password
	if r.Form.Get("password1") != r.Form.Get("password2") {
		sessionService.SetFlashMessage("Passwords must match")
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Confirm the invitation
	if err := s.accountsService.ConfirmInvitation(
		invitation,
		r.Form.Get("password1"),
	); err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Redirect to the success page
	redirectWithQueryString("/web/confirm-invitation-success", r.URL.Query(), w, r)
}
