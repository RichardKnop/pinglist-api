package subscriptions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/response"
)

// Handles calls to create a subscription (POST /v1/subscriptions)
func (s *Service) createSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
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
	subscriptionRequest := new(SubscriptionRequest)
	if err := json.Unmarshal(payload, subscriptionRequest); err != nil {
		logger.ERROR.Printf("Failed to unmarshal subscription request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a subscription
	subscription, err := s.createSubscription(authenticatedUser, subscriptionRequest)
	if err != nil {
		logger.ERROR.Printf("Create subscription error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// Create response
	subscriptionResponse, err := NewSubscriptionResponse(subscription)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set Location header to the newly created resource
	w.Header().Set("Location", fmt.Sprintf("/v1/subscriptions/%d", subscription.ID))
	// Write JSON response
	response.WriteJSON(w, subscriptionResponse, http.StatusCreated)
}
