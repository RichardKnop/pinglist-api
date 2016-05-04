package alarms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/pagination"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestListAlarmResponseTimesRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.listAlarmResponseTimesHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AlarmsTestSuite) TestListAlarmResponseTimesWithoutPermission() {
	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://1.2.3.4/v1/alarms/%d/response-times",
			suite.alarms[0].ID,
		),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_alarm_response_times", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[2])

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 403, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrListAlarmResponseTimesPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AlarmsTestSuite) TestListAlarmResponseTimesNoResults() {
	var (
		okAlarmState *AlarmState
		testAlarm    *Alarm
		err          error
	)

	okAlarmState, err = suite.service.findAlarmStateByID(alarmstates.OK)
	assert.NoError(suite.T(), err, "Failed to fetch OK alarm state")

	// Insert a test alarm
	testAlarm = NewAlarm(
		suite.users[1],
		suite.regions[0],
		okAlarmState,
		&AlarmRequest{
			Region:                 "us-west-2",
			EndpointURL:            "http://new-endpoint",
			ExpectedHTTPCode:       200,
			MaxResponseTime:        1000,
			Interval:               60,
			EmailAlerts:            true,
			PushNotificationAlerts: true,
			Active:                 true,
		},
	)
	err = suite.db.Create(testAlarm).Error
	assert.NoError(suite.T(), err, "Failed to insert a test alarm")

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d/response-times", testAlarm.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_alarm_response_times", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Mock paginated response time metrics
	suite.mockResponseTimesCount(
		int(testAlarm.ID), // reference ID
		"",                // date_trunc
		nil,               // from
		nil,               // to
		0,                 // returned count
		nil,               // returned error
	)
	suite.mockFindPaginatedResponseTimes(
		0, // offset
		pagination.DefaultLimit,
		"",                // order by
		int(testAlarm.ID), // reference ID
		"",                // date_trunc
		nil,               // from
		nil,               // to
		[]*metrics.ResponseTime{}, // returned metrics
		nil, // returned error
	)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	expected := &ListResponseTimesResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/response-times", testAlarm.ID),
				},
				"first": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/response-times?page=1", testAlarm.ID),
				},
				"last": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/response-times?page=1", testAlarm.ID),
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				"response_times": jsonhal.Embedded([]*metrics.MetricResponse{}),
			},
		},
		Uptime:  0,
		Average: 0,
		IncidentTypeCounts: map[string]int{
			incidenttypes.SlowResponse: 0,
			incidenttypes.Timeout:      0,
			incidenttypes.BadCode:      0,
			incidenttypes.Other:        0,
		},
		Count: 0,
		Page:  1,
	}

	expectedJSON, err := json.Marshal(expected)
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
		)
	}
}

func (suite *AlarmsTestSuite) TestListAlarmResponseTimes() {
	var (
		today             = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName = "metrics_request_times_2016_02_09"
	)

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://1.2.3.4/v1/alarms/%d/response-times",
			suite.alarms[0].ID,
		),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_alarm_response_times", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Mock paginated response time metrics
	suite.mockResponseTimesCount(
		int(suite.alarms[0].ID), // reference ID
		"",  // date_trunc
		nil, // from
		nil, // to
		4,   // returned count
		nil, // returned error
	)
	testMetrics := []*metrics.ResponseTime{
		metrics.NewResponseTime(todaySubTableName, suite.alarms[0].ID, today, 123),
		metrics.NewResponseTime(todaySubTableName, suite.alarms[0].ID, today.Add(1*time.Hour), 234),
		metrics.NewResponseTime(todaySubTableName, suite.alarms[0].ID, today.Add(2*time.Hour), 345),
		metrics.NewResponseTime(todaySubTableName, suite.alarms[0].ID, today.Add(3*time.Hour), 456),
	}
	suite.mockFindPaginatedResponseTimes(
		0, // offset
		pagination.DefaultLimit,
		"", // order by
		int(suite.alarms[0].ID), // reference ID
		"",          // date_trunc
		nil,         // from
		nil,         // to
		testMetrics, // returned metrics
		nil,         // returned error
	)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	uptime, _, err := suite.service.getUptimeDowntime(suite.alarms[0])
	assert.NoError(suite.T(), err, "Failed to fetch uptime")
	expectedUptime, err := strconv.ParseFloat(fmt.Sprintf("%.4f", uptime), 64)
	assert.NoError(suite.T(), err, "Failed formatting uptime to 4 decimals")

	expectedAverage := 289.5

	responseTimeResponses := make([]*metrics.MetricResponse, len(testMetrics))
	for i, testMetric := range testMetrics {
		responseTimeResponse, err := metrics.NewMetricResponse(
			testMetric.Timestamp,
			testMetric.Value,
		)
		assert.NoError(suite.T(), err, "Creating response object failed")
		responseTimeResponses[i] = responseTimeResponse
	}

	expected := &ListResponseTimesResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/response-times", suite.alarms[0].ID),
				},
				"first": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/response-times?page=1", suite.alarms[0].ID),
				},
				"last": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/response-times?page=1", suite.alarms[0].ID),
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				"response_times": jsonhal.Embedded(responseTimeResponses),
			},
		},
		Uptime:  expectedUptime,
		Average: expectedAverage,
		IncidentTypeCounts: map[string]int{
			incidenttypes.SlowResponse: 0,
			incidenttypes.Timeout:      2,
			incidenttypes.BadCode:      1,
			incidenttypes.Other:        1,
		},
		Count: 4,
		Page:  1,
	}

	expectedJSON, err := json.Marshal(expected)
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
		)
	}
}
