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
	// ErrCancelSubscriptionPermission ...
	ErrCancelSubscriptionPermission = errors.New("Need permission to cancel subscriptions")
)

// Handles calls to cancel a subscription (DELETE /v1/subscriptions/{id:[0-9]+})
func (s *Service) cancelSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch the subscription we want to cancel
	subscription, err := s.FindSubscriptionByID(uint(subscriptionID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkCancelSubscriptionPermissions(authenticatedUser, subscription); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Cancel the subscription
	if err := s.cancelSubscription(subscription); err != nil {
		logger.Errorf("Cancel subscription error: %s", err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 204 no content response
	response.NoContent(w)
}

func checkCancelSubscriptionPermissions(authenticatedUser *accounts.User, subscription *Subscription) error {
	// Superusers can cancel any subscriptions
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can cancel their own subscriptions
	if authenticatedUser.ID == subscription.Customer.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrCancelSubscriptionPermission
}
