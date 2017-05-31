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
	// ErrInviteUserPermission ...
	ErrInviteUserPermission = errors.New("Need permission to invite user")
)

// Handles requests to invite a new user to a team (POST /v1/teams/{id:[0-9]+}/invitations)
func (s *Service) inviteUserHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := checkInviteUserPermissions(authenticatedUser, team); err != nil {
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
	invitationRequest := new(InvitationRequest)
	if err := json.Unmarshal(payload, invitationRequest); err != nil {
		logger.ERROR.Printf("Failed to unmarshal invitation request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new invited user account
	invitation, err := s.inviteUser(
		team,
		invitationRequest.Email,
		true, // update members assoc
	)
	if err != nil {
		logger.ERROR.Printf("Invite user error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// Create response
	invitationResponse, err := accounts.NewInvitationResponse(invitation)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, invitationResponse, http.StatusCreated)
}

func checkInviteUserPermissions(authenticatedUser *accounts.User, team *Team) error {
	// Superusers can invite users to any team
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Owners can invite users to their own team
	if authenticatedUser.ID == team.Owner.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrInviteUserPermission
}
