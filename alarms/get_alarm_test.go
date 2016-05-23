package alarms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/pinglist-api/util"
	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestGetAlarmRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.getAlarmHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AlarmsTestSuite) TestGetAlarmWithoutPermission() {
	// Insert a test alarm
	alarm := &Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://endpoint-5",
		ExpectedHTTPCode: 200,
		Interval:         60,
		Active:           false,
	}
	assert.NoError(suite.T(), suite.db.Create(alarm).Error, "Inserting test data failed")

	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", alarm.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "get_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[2])

	// Count before
	var countBefore int
	suite.db.Model(new(Alarm)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 403, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrGetAlarmPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AlarmsTestSuite) TestGetAlarm() {
	// Insert a test alarm
	alarm := &Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://endpoint-5",
		ExpectedHTTPCode: 200,
		Interval:         60,
		Active:           false,
	}
	assert.NoError(suite.T(), suite.db.Create(alarm).Error, "Inserting test data failed")

	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", alarm.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "get_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Count before
	var countBefore int
	suite.db.Model(new(Alarm)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	expected := &AlarmResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d", alarm.ID),
				},
			},
		},
		ID:                     alarm.ID,
		UserID:                 uint(alarm.UserID.Int64),
		Region:                 regions.USWest2,
		EndpointURL:            alarm.EndpointURL,
		ExpectedHTTPCode:       alarm.ExpectedHTTPCode,
		MaxResponseTime:        alarm.MaxResponseTime,
		Interval:               alarm.Interval,
		EmailAlerts:            alarm.EmailAlerts,
		PushNotificationAlerts: alarm.PushNotificationAlerts,
		Active:                 alarm.Active,
		State:                  alarm.AlarmStateID.String,
		CreatedAt:              util.FormatTime(alarm.CreatedAt),
		UpdatedAt:              util.FormatTime(alarm.UpdatedAt),
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
