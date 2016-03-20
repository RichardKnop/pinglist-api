package accounts

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/response"
)

// Handles requests to get owned team data (GET /v1/accounts/ownedteam)
func (s *Service) getOwnedTeamHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Fetch the team user owns
	team, err := s.FindTeamByOwnerID(authenticatedUser.ID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create response
	teamResponse, err := NewTeamResponse(team)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, teamResponse, http.StatusOK)
}
