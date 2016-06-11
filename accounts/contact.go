package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/pinglist-api/email"
	"github.com/RichardKnop/pinglist-api/response"
)

// Handles requests to send a contact email (POST /v1/accounts/contact)
func (s *Service) contactHandler(w http.ResponseWriter, r *http.Request) {
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
	contactRequest := new(ContactRequest)
	if err := json.Unmarshal(payload, contactRequest); err != nil {
		logger.Errorf("Failed to unmarshal contact request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send contact email
	go func() {
		contactEmail := &email.Email{
			Subject: contactRequest.Subject,
			Recipients: []*email.Recipient{&email.Recipient{
				Email: s.cnf.Pinglist.ContactEmail,
				Name:  "Pinglist Admin",
			}},
			From: fmt.Sprintf("%s <%s>", contactRequest.Name, contactRequest.Email),
			Text: contactRequest.Message,
		}

		// Try to send the contact email
		if err := s.emailService.Send(contactEmail); err != nil {
			logger.Errorf("Send email error: %s", err)
			return
		}
	}()

	// 204 no content response
	response.NoContent(w)
}
