package teams

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/gorilla/mux"
)

var (
	// ErrGetTeamPermission ...
	ErrGetTeamPermission = errors.New("Need permission to get team")
)

// Handles calls to get a team (GET /v1/teams/{id:[0-9]+})
func (s *Service) getTeamHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Get the id from request URI and type assert it
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the team we want to get
	team, err := s.FindTeamByID(uint(teamID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkGetTeamPermissions(authenticatedUser, team); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
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

func checkGetTeamPermissions(authenticatedUser *accounts.User, team *Team) error {
	// Superusers can get any teams
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can get their own teams
	if authenticatedUser.ID == team.Owner.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrGetTeamPermission
}
