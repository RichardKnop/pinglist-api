package alarms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestUpdateAlarmRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.updateAlarmHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AlarmsTestSuite) TestUpdateAlarmWithoutPermission() {
	testAlarm, err := suite.insertTestAlarm(true)
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", testAlarm.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_alarm", match.Route.GetName())
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
		map[string]string{"error": ErrUpdateAlarmPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AlarmsTestSuite) TestUpdateAlarmMaxLimitReached() {
	_, err := suite.insertTestAlarm(true)
	assert.NoError(suite.T(), err, "Inserting test data failed")

	testAlarm, err := suite.insertTestAlarm(false)
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		Region:                 "us-west-2",
		EndpointURL:            "http://foobar-updated",
		ExpectedHTTPCode:       201,
		MaxResponseTime:        2000,
		Interval:               90,
		EmailAlerts:            false,
		PushNotificationAlerts: false,
		SlackAlerts:            false,
		Active:                 true,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", testAlarm.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

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

func (suite *AlarmsTestSuite) TestUpdateAlarmIntervalTooSmall() {
	testAlarm, err := suite.insertTestAlarm(true)
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		Region:                 "us-west-2",
		EndpointURL:            "http://foobar-updated",
		ExpectedHTTPCode:       201,
		MaxResponseTime:        2000,
		Interval:               5,
		EmailAlerts:            false,
		PushNotificationAlerts: false,
		SlackAlerts:            false,
		Active:                 true,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", testAlarm.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

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
				MaxAlarms:        10,
				MinAlarmInterval: 50,
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
	if !assert.Equal(suite.T(), 400, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	expectedJSON, err := json.Marshal(
		map[string]string{"error": NewErrIntervalTooSmall(uint(50)).Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AlarmsTestSuite) TestUpdateAlarmMaxResponseTimeTooBig() {
	testAlarm, err := suite.insertTestAlarm(true)
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		Region:                 "us-west-2",
		EndpointURL:            "http://foobar-updated",
		ExpectedHTTPCode:       201,
		MaxResponseTime:        10001,
		Interval:               60,
		EmailAlerts:            false,
		PushNotificationAlerts: false,
		SlackAlerts:            false,
		Active:                 true,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", testAlarm.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

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
	if !assert.Equal(suite.T(), 400, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrMaxResponseTimeTooBig.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AlarmsTestSuite) TestUpdateAlarmRegionNotFound() {
	testAlarm, err := suite.insertTestAlarm(true)
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		Region:                 "transylvania",
		EndpointURL:            "http://foobar-updated",
		ExpectedHTTPCode:       201,
		MaxResponseTime:        2000,
		Interval:               90,
		EmailAlerts:            false,
		PushNotificationAlerts: false,
		SlackAlerts:            false,
		Active:                 true,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", testAlarm.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

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
	if !assert.Equal(suite.T(), 400, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrRegionNotFound.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AlarmsTestSuite) TestUpdateAlarm() {
	testAlarm, err := suite.insertTestAlarm(true)
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		Region:                 "us-west-2",
		EndpointURL:            "http://foobar-updated",
		ExpectedHTTPCode:       201,
		MaxResponseTime:        2000,
		Interval:               90,
		EmailAlerts:            false,
		PushNotificationAlerts: false,
		SlackAlerts:            false,
		Active:                 true,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", testAlarm.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

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
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated alarm
	alarm := new(Alarm)
	notFound := suite.db.Preload("User").Preload("Incidents").
		Find(alarm, testAlarm.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), suite.users[1].ID, uint(alarm.UserID.Int64))
	assert.Equal(suite.T(), "http://foobar-updated", alarm.EndpointURL)
	assert.Equal(suite.T(), uint(201), alarm.ExpectedHTTPCode)
	assert.Equal(suite.T(), uint(2000), alarm.MaxResponseTime)
	assert.Equal(suite.T(), uint(90), alarm.Interval)
	assert.False(suite.T(), alarm.EmailAlerts)
	assert.False(suite.T(), alarm.PushNotificationAlerts)
	assert.False(suite.T(), alarm.SlackAlerts)
	assert.True(suite.T(), alarm.Active)
	assert.Equal(suite.T(), 0, len(alarm.Incidents))

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
		EndpointURL:            "http://foobar-updated",
		ExpectedHTTPCode:       uint(201),
		MaxResponseTime:        uint(2000),
		Interval:               uint(90),
		EmailAlerts:            false,
		PushNotificationAlerts: false,
		SlackAlerts:            false,
		Active:                 true,
		State:                  alarmstates.InsufficientData,
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

func (suite *AlarmsTestSuite) insertTestAlarm(active bool) (*Alarm, error) {
	// Insert a test alarm
	testAlarm := &Alarm{
		User:                   suite.users[1],
		Region:                 &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:             &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:            "http://foobar",
		ExpectedHTTPCode:       200,
		MaxResponseTime:        1000,
		Interval:               60,
		EmailAlerts:            true,
		PushNotificationAlerts: true,
		SlackAlerts:            true,
		Active:                 active,
	}
	return testAlarm, suite.db.Create(testAlarm).Error
}
