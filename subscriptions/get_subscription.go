package subscriptions

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
	// ErrGetSubscriptionPermission ...
	ErrGetSubscriptionPermission = errors.New("Need permission to get subscription")
)

// Handles calls to get a subscription (GET /v1/subscriptions/{id:[0-9]+})
func (s *Service) getSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch the subscription we want to get
	subscription, err := s.FindSubscriptionByID(uint(subscriptionID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkGetSubscriptionPermissions(authenticatedUser, subscription); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
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

func checkGetSubscriptionPermissions(authenticatedUser *accounts.User, subscription *Subscription) error {
	// Superusers can get any subscriptions
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can get their own subscriptions
	if authenticatedUser.ID == subscription.Customer.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrGetSubscriptionPermission
}
