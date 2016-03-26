package alarms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestCreateAlarmRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.createAlarmHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AlarmsTestSuite) TestCreateAlarmMaxLimitReached() {
	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		Region:                 "us-west-2",
		EndpointURL:            "http://endpoint-5",
		ExpectedHTTPCode:       200,
		Interval:               60,
		EmailAlerts:            true,
		PushNotificationAlerts: true,
		Active:                 true,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/alarms",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Mock find team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		nil,
		teams.ErrTeamNotFound,
	)

	// Mock find active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Alarm)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 400, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrMaxAlarmsLimitReached.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AlarmsTestSuite) TestCreateAlarm() {
	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		Region:                 "us-west-2",
		EndpointURL:            "http://endpoint-5",
		ExpectedHTTPCode:       200,
		Interval:               60,
		EmailAlerts:            true,
		PushNotificationAlerts: true,
		Active:                 true,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/alarms",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Mock find team
	suite.mockFindTeamByMemberID(
		suite.users[1].ID,
		nil,
		teams.ErrTeamNotFound,
	)

	// Mock find active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[1].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxAlarms: 10,
			},
		},
		nil,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Alarm)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created alarm
	alarm := new(Alarm)
	notFound := suite.db.Preload("User").Preload("Incidents").
		Last(alarm).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), suite.users[1].ID, uint(alarm.UserID.Int64))
	assert.Equal(suite.T(), "http://endpoint-5", alarm.EndpointURL)
	assert.Equal(suite.T(), uint(200), alarm.ExpectedHTTPCode)
	assert.Equal(suite.T(), uint(60), alarm.Interval)
	assert.True(suite.T(), alarm.EmailAlerts)
	assert.True(suite.T(), alarm.PushNotificationAlerts)
	assert.True(suite.T(), alarm.Active)
	assert.Equal(suite.T(), 0, len(alarm.Incidents))

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/alarms/%d", alarm.ID),
		w.Header().Get("Location"),
	)

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
		UserID:                 suite.users[1].ID,
		Region:                 regions.USWest2,
		EndpointURL:            "http://endpoint-5",
		ExpectedHTTPCode:       uint(200),
		Interval:               uint(60),
		EmailAlerts:            true,
		PushNotificationAlerts: true,
		Active:                 true,
		State:                  alarmstates.InsufficientData,
		CreatedAt:              alarm.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:              alarm.UpdatedAt.UTC().Format(time.RFC3339),
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
