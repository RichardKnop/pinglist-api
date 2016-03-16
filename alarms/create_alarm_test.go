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
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestCreateAlarm() {
	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		EndpointURL:      "http://endpoint-5",
		ExpectedHTTPCode: 200,
		Interval:         60,
		Active:           false,
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

	// Mock find active subscription
	suite.mockFindActiveSubscription(
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
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())

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
	notFound := suite.db.Preload("User").Preload("Incidents").Preload("Results").
		Last(alarm).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), suite.users[1].ID, uint(alarm.UserID.Int64))
	assert.Equal(suite.T(), "http://endpoint-5", alarm.EndpointURL)
	assert.Equal(suite.T(), uint(200), alarm.ExpectedHTTPCode)
	assert.Equal(suite.T(), uint(60), alarm.Interval)
	assert.False(suite.T(), alarm.Active)
	assert.Equal(suite.T(), 0, len(alarm.Incidents))
	assert.Equal(suite.T(), 0, len(alarm.Results))

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
		ID:               alarm.ID,
		UserID:           suite.users[1].ID,
		EndpointURL:      "http://endpoint-5",
		ExpectedHTTPCode: uint(200),
		Interval:         uint(60),
		Active:           false,
		State:            alarmstates.InsufficientData,
		CreatedAt:        alarm.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:        alarm.UpdatedAt.UTC().Format(time.RFC3339),
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

func (suite *AlarmsTestSuite) TestCreateAlarmMaxLimitReached() {
	// Prepare a request
	payload, err := json.Marshal(&AlarmRequest{
		EndpointURL:      "http://endpoint-5",
		ExpectedHTTPCode: 200,
		Interval:         60,
		Active:           false,
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

	// Mock find active subscription
	suite.mockFindActiveSubscription(
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
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())

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
