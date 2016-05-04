package alarms

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/pagination"
	"github.com/RichardKnop/pinglist-api/response"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/gorilla/mux"
)

var (
	// ErrListAlarmResponseTimesPermission ...
	ErrListAlarmResponseTimesPermission = errors.New("Need permission to list alarm response times")
)

// Handles calls to list alarm response times (GET /v1/alarms/{id:[0-9]+}/response-times)
func (s *Service) listAlarmResponseTimesHandler(w http.ResponseWriter, r *http.Request) {
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

	// Fetch the alarm
	alarm, err := s.FindAlarmByID(uint(alarmID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check permissions
	if err := checkListAlarmResponseTimesPermissions(authenticatedUser, alarm); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Get page and limit
	page, limit, err := pagination.GetPageLimit(r)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get other params
	dateTrunc, from, to, err := metrics.GetParamsFromQueryString(r)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Count total number of metric records
	count, err := s.metricsService.ResponseTimesCount(
		int(alarm.ID),
		dateTrunc,
		from,
		to,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		wg                     sync.WaitGroup
		responseTimesChan      = make(chan []*metrics.ResponseTime)
		incidentTypeCountsChan = make(chan map[string]int)
		uptimeChan             = make(chan float64)
		errChan                = make(chan error)
	)

	wg.Add(1)
	go func() {
		// Get paginated metric records
		responseTimes, err := s.metricsService.FindPaginatedResponseTimes(
			pagination.GetOffsetForPage(count, page, limit),
			limit,
			r.URL.Query().Get("order_by"),
			int(alarm.ID),
			dateTrunc,
			from,
			to,
		)
		if err != nil {
			errChan <- err
			return
		}
		responseTimesChan <- responseTimes
	}()

	wg.Add(1)
	go func() {
		// Get aggregate incident counts based on type
		incidentTypeCounts, err := s.incidentTypeCounts(
			nil, // user
			alarm,
			from,
			to,
		)
		if err != nil {
			errChan <- err
			return
		}
		incidentTypeCountsChan <- incidentTypeCounts
	}()

	wg.Add(1)
	go func() {
		// Calculate uptime
		uptime, _, err := s.getUptimeDowntime(alarm)
		if err != nil {
			errChan <- err
			return
		}
		uptimeChan <- uptime
	}()

	var (
		responseTimes      []*metrics.ResponseTime
		incidentTypeCounts map[string]int
		uptime             float64
		errs               []error
	)

	for i := 0; i < 3; i++ {
		select {
		case responseTimes = <-responseTimesChan:
		case incidentTypeCounts = <-incidentTypeCountsChan:
		case uptime = <-uptimeChan:
		case err := <-errChan:
			errs = append(errs, err)
		}
	}

	// If one of the goroutines failed
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error(err)
		}
		response.Error(
			w,
			"Something went wrong while fetching metrics data",
			http.StatusInternalServerError,
		)
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

	// Create response
	self := util.GetCurrentURL(r)
	listResponseTimesResponse, err := NewListResponseTimesResponse(
		count, page,
		self, first, last, previous, next,
		responseTimes, uptime, incidentTypeCounts,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Infof("%v", listResponseTimesResponse)
	// Write JSON response
	response.WriteJSON(w, listResponseTimesResponse, http.StatusOK)
}

func checkListAlarmResponseTimesPermissions(authenticatedUser *accounts.User, alarm *Alarm) error {
	// Superusers can list any alarm response times
	if authenticatedUser.Role.Name == roles.Superuser {
		return nil
	}

	// Users can list their own alarm response times
	if authenticatedUser.ID == alarm.User.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrListAlarmResponseTimesPermission
}
