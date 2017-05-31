package teams

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/gorilla/mux"
)

var (
	// ErrUpdateTeamPermission ...
	ErrUpdateTeamPermission = errors.New("Need permission to update team")
)

// Handles requests to update a team (PUT /v1/teams/{id:[0-9]+})
func (s *Service) updateTeamHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch the team we want to update
	team, err := s.FindTeamByID(uint(teamID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkUpdateTeamPermissions(authenticatedUser, team); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Request body cannot be nil
	if r.Body == nil {
		response.Error(w, "Request body cannot be nil", http.StatusBadRequest)
		return
	}

	// Read the request body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into the request prototype
	teamRequest := new(TeamRequest)
	if err := json.Unmarshal(payload, teamRequest); err != nil {
		logger.ERROR.Printf("Failed to unmarshal team request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the team
	if err := s.updateTeam(team, teamRequest); err != nil {
		logger.ERROR.Printf("Update team error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
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

func checkUpdateTeamPermissions(authenticatedUser *accounts.User, team *Team) error {
	// Superusers can update any team
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Owners can update their own team
	if authenticatedUser.ID == team.Owner.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrUpdateTeamPermission
}
