package alarms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/pagination"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestListAlarmRequestTimesRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.listAlarmRequestTimesHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AlarmsTestSuite) TestListAlarmRequestTimesWithoutPermission() {
	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d/request-times", suite.alarms[0].ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_alarm_request_times", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[2])

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
		map[string]string{"error": ErrListAlarmRequestTimesPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AlarmsTestSuite) TestListAlarmRequestTimes() {
	var (
		today             = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName = "metrics_request_times_2016_02_09"
	)

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d/request-times", suite.alarms[0].ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_alarm_request_times", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Mock paginated request time metrics
	suite.mockPaginatedRequestTimesCount(
		int(suite.alarms[0].ID), // reference ID
		"",  // date_trunc
		nil, // from
		nil, // to
		4,   // returned count
		nil, // returned error
	)
	testMetrics := []*metrics.RequestTime{
		metrics.NewRequestTime(todaySubTableName, suite.alarms[0].ID, today, 123),
		metrics.NewRequestTime(todaySubTableName, suite.alarms[0].ID, today.Add(1*time.Hour), 234),
		metrics.NewRequestTime(todaySubTableName, suite.alarms[0].ID, today.Add(2*time.Hour), 345),
		metrics.NewRequestTime(todaySubTableName, suite.alarms[0].ID, today.Add(3*time.Hour), 456),
	}
	suite.mockFindPaginatedRequestTimes(
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
	requestTimeResponses := make([]*metrics.MetricResponse, len(testMetrics))
	for i, testMetric := range testMetrics {
		requestTimeResponse, err := metrics.NewMetricResponse(
			testMetric.Timestamp,
			testMetric.Value,
		)
		assert.NoError(suite.T(), err, "Creating response object failed")
		requestTimeResponses[i] = requestTimeResponse
	}
	expected := &ListIncidentsResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/request-times", suite.alarms[0].ID),
				},
				"first": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/request-times?page=1", suite.alarms[0].ID),
				},
				"last": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/request-times?page=1", suite.alarms[0].ID),
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				"request_times": jsonhal.Embedded(requestTimeResponses),
			},
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
