package teams

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/response"
)

// Handles requests to create a new team (POST /v1/teams)
func (s *Service) createTeamHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
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
		logger.Errorf("Failed to unmarshal team request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new team
	team, err := s.createTeam(authenticatedUser, teamRequest)
	if err != nil {
		logger.Errorf("Create team error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// Create response
	teamResponse, err := NewTeamResponse(team)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set Location header to the newly created resource
	w.Header().Set("Location", fmt.Sprintf("/v1/teams/%d", team.ID))
	// Write JSON response
	response.WriteJSON(w, teamResponse, http.StatusCreated)
}
