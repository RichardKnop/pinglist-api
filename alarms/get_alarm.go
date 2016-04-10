package alarms

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
	// ErrGetAlarmPermission ...
	ErrGetAlarmPermission = errors.New("Need permission to get alarm")
)

// Handles calls to get an alarm (GET /v1/alarms/{id:[0-9]+})
func (s *Service) getAlarmHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := checkGetAlarmPermissions(authenticatedUser, alarm); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Create response
	alarmResponse, err := NewAlarmResponse(alarm)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, alarmResponse, http.StatusOK)
}

func checkGetAlarmPermissions(authenticatedUser *accounts.User, alarm *Alarm) error {
	// Superusers can get any alarms
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can get their own alarms
	if authenticatedUser.ID == alarm.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrGetAlarmPermission
}
