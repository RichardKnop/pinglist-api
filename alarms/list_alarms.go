package alarms

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/pagination"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/RichardKnop/pinglist-api/util"
)

var (
	// ErrListAlarmsPermission ...
	ErrListAlarmsPermission = errors.New("Need permission to list alarms")
)

// Handles calls to list alarms (GET /v1/alarms)
func (s *Service) listAlarmsHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := checkListAlarmsPermissions(authenticatedUser, user); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Get page and limit
	page, limit, err := pagination.GetPageLimit(r)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Count total number of results
	count, err := s.alarmsCount(user)
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
	alarms, err := s.findPaginatedAlarms(
		pagination.GetOffsetForPage(count, page, limit),
		limit,
		r.URL.Query().Get("order_by"),
		user,
	)
	if err != nil {
		logger.Errorf("Find paginated alarms error: %s", err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	self := util.GetCurrentURL(r)
	listAlarmsResponse, err := NewListAlarmsResponse(
		count, page,
		self, first, last, previous, next,
		alarms,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, listAlarmsResponse, http.StatusOK)
}

func checkListAlarmsPermissions(authenticatedUser *accounts.User, user *accounts.User) error {
	// Superusers can list any alarms
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can list their own alarms
	if user != nil && authenticatedUser.ID == user.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrListAlarmsPermission
}
