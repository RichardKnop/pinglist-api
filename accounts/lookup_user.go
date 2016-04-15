package accounts

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/response"
)

// Handles requests to lookup a user by email (GET /v1/accounts/user-lookup)
func (s *Service) userLookupHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	_, err := GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Get the email from query string
	email := r.URL.Query().Get("email")
	if email == "" {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user by email
	user, err := s.FindUserByEmail(email)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create response
	userResponse, err := NewUserResponse(user)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, userResponse, http.StatusOK)
}
