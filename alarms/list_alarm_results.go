package alarms

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/pagination"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/gorilla/mux"
)

var (
	// ErrListAlarmResultsPermission ...
	ErrListAlarmResultsPermission = errors.New("Need permission to list alarm results")
)

// Handles calls to list alarm results (GET /v1/alarms/{id:[0-9]+}/results)
func (s *Service) listAlarmResultsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := accounts.GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Get the id from request URI and type assert it
	vars := mux.Vars(r)
	alarmID, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the alarm we want to update
	alarm, err := s.FindAlarmByID(uint(alarmID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkListAlarmResultsPermissions(authenticatedUser, alarm); err != nil {
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
	count, err := s.paginatedResultsCount(alarm)
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
	results, err := s.findPaginatedResults(
		pagination.GetOffsetForPage(count, page, limit),
		limit,
		r.URL.Query().Get("order_by"),
		alarm,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	self := util.GetCurrentURL(r)
	listResultsResponse, err := NewListResultsResponse(
		count, page,
		self, first, last, next, previous,
		results,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, listResultsResponse, http.StatusOK)
}

func checkListAlarmResultsPermissions(authenticatedUser *accounts.User, alarm *Alarm) error {
	// Superusers can list any alarm results
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can list their own alarm results
	if authenticatedUser.ID == alarm.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrListAlarmResultsPermission
}
