package subscriptions

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/pagination"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/RichardKnop/pinglist-api/util"
)

var (
	// ErrListSubscriptionsPermission ...
	ErrListSubscriptionsPermission = errors.New("Need permission to list subscriptions")
)

// Handles calls to list subscriptions (GET /v1/subscriptions)
func (s *Service) listSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Optional filtering by user
	user, err := s.GetAccountsService().GetUserFromQueryString(r)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Default to the authenticated user unless logged in as superuser
	if authenticatedUser.Role.Name != roles.Superuser && user == nil {
		user = authenticatedUser
	}

	// Check permissions
	if err := checkListSubscriptionsPermissions(authenticatedUser, user); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Get page, limit and order by
	page, limit, orderBy, err := pagination.GetPageLimitOrderBy(r)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Count total number of results
	count, err := s.subscriptionsCount(user)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get pagination links
	first, last, previous, next, err := pagination.GetPaginationLinks(
		r.URL,
		count,
		page,
		limit,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get paginated results
	subscriptions, err := s.findPaginatedSubscriptions(
		pagination.GetOffsetForPage(count, page, limit),
		limit,
		orderBy,
		user,
	)
	if err != nil {
		logger.ERROR.Printf("Find paginated subscriptions error: %s", err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	self := util.GetCurrentURL(r)
	listSubscriptionsResponse, err := NewListSubscriptionsResponse(
		count, page,
		self, first, last, previous, next,
		subscriptions,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, listSubscriptionsResponse, http.StatusOK)
}

func checkListSubscriptionsPermissions(authenticatedUser *accounts.User, user *accounts.User) error {
	// Superusers can list any subscriptions
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can list their own subscriptions
	if user != nil && authenticatedUser.ID == user.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrListSubscriptionsPermission
}
