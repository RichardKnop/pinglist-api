package subscriptions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/gorilla/mux"
)

var (
	// ErrUpdateSubscriptionPermission ...
	ErrUpdateSubscriptionPermission = errors.New("Need permission to update subscription")
)

// Handles calls to update a subscription (PUT /v1/subscriptions/{id:[0-9]+})
func (s *Service) updateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Get the id from request URI and type assert it
	vars := mux.Vars(r)
	subscriptionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the subscription we want to update
	subscription, err := s.FindSubscriptionByID(uint(subscriptionID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkUpdateSubscriptionPermissions(authenticatedUser, subscription); err != nil {
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
	subscriptionRequest := new(SubscriptionRequest)
	if err := json.Unmarshal(payload, subscriptionRequest); err != nil {
		log.Printf("Failed to unmarshal subscription request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the subscription
	if err := s.updateSubscription(subscription, subscriptionRequest); err != nil {
		log.Printf("Update subscription error: %s", err)
		code, ok := errStatusCodeMap[err]
		if !ok {
			code = http.StatusInternalServerError
		}
		response.Error(w, err.Error(), code)
		return
	}

	// Create response
	subscriptionResponse, err := NewSubscriptionResponse(subscription)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, subscriptionResponse, http.StatusOK)
}

func checkUpdateSubscriptionPermissions(authenticatedUser *accounts.User, subscription *Subscription) error {
	// Superusers can update any subscriptions
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can update their own subscriptions
	if authenticatedUser.ID == subscription.Customer.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrUpdateSubscriptionPermission
}
