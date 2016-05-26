package alarms

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
	// ErrUpdateAlarmPermission ...
	ErrUpdateAlarmPermission = errors.New("Need permission to update alarm")
)

// Handles calls to update an alarm (PUT /v1/alarms/{id:[0-9]+})
func (s *Service) updateAlarmHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := checkUpdateAlarmPermissions(authenticatedUser, alarm); err != nil {
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
	alarmRequest := new(AlarmRequest)
	if err := json.Unmarshal(payload, alarmRequest); err != nil {
		log.Printf("Failed to unmarshal alarm request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update an alarm
	if err := s.updateAlarm(alarm, alarmRequest); err != nil {
		log.Printf("Update alarm error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
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

func checkUpdateAlarmPermissions(authenticatedUser *accounts.User, alarm *Alarm) error {
	// Superusers can update any alarms
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can update their own alarms
	if authenticatedUser.ID == alarm.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrUpdateAlarmPermission
}
