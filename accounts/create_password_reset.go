package accounts

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/pinglist-api/response"
)

// Handles requests to reset a password (POST /v1/accounts/passwordreset)
func (s *Service) createPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated client from the request context
	_, err := GetAuthenticatedAccount(r)
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
	passwordResetRequest := new(PasswordResetRequest)
	if err := json.Unmarshal(payload, passwordResetRequest); err != nil {
		logger.Errorf("Failed to unmarshal password reset request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user who wants to reset his/her password based on the email
	user, err := s.FindUserByEmail(passwordResetRequest.Email)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create a new password reset
	_, err = s.createPasswordReset(user)
	if err != nil {
		logger.Errorf("Create password reset error: %s", err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 204 no content response
	response.NoContent(w)
}
