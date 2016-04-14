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
	// ErrDeleteTeamPermission ...
	ErrDeleteTeamPermission = errors.New("Need permission to delete team")
)

// Handles calls to delete a team (DELETE /v1/teams/{id:[0-9]+})
func (s *Service) deleteTeamHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch the team we want to delete
	team, err := s.FindTeamByID(uint(teamID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkDeleteTeamPermissions(authenticatedUser, team); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Soft delete the team
	if err := s.db.Delete(team).Error; err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 204 no content response
	response.NoContent(w)
}

func checkDeleteTeamPermissions(authenticatedUser *accounts.User, team *Team) error {
	// Superusers can delete any teams
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can delete their own teams
	if authenticatedUser.ID == team.Owner.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrDeleteTeamPermission
}
