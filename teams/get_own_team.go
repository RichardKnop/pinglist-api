package teams

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/response"
)

// Handles requests to get own team data (GET /v1/ownteam)
func (s *Service) getOwnTeamHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := accounts.GetAuthenticatedUser(r)
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
