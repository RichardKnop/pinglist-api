package alarms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/response"
)

// Handles calls to create an alarm (POST /v1/alarms)
func (s *Service) createAlarmHandler(w http.ResponseWriter, r *http.Request) {
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
	alarmRequest := new(AlarmRequest)
	if err := json.Unmarshal(payload, alarmRequest); err != nil {
		logger.Errorf("Failed to unmarshal alarm request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new alarm
	alarm, err := s.createAlarm(authenticatedUser, alarmRequest)
	if err != nil {
		logger.Errorf("Create alarm error: %s", err)
		switch err {
		case ErrMaxAlarmsLimitReached:
			response.Error(w, err.Error(), http.StatusBadRequest)
		case ErrRegionNotFound:
			response.Error(w, err.Error(), http.StatusBadRequest)
		case ErrAlarmStateNotFound:
			response.Error(w, err.Error(), http.StatusBadRequest)
		default:
			response.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create response
	alarmResponse, err := NewAlarmResponse(alarm)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set Location header to the newly created resource
	w.Header().Set("Location", fmt.Sprintf("/v1/alarms/%d", alarm.ID))
	// Write JSON response
	response.WriteJSON(w, alarmResponse, http.StatusCreated)
}
