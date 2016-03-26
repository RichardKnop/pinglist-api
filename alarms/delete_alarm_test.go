package alarms

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestDeleteAlarmRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.deleteAlarmHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AlarmsTestSuite) TestDeleteAlarm() {
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
		"DELETE",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d", alarm.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "delete_alarm", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

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
	suite.teamsServiceMock.AssertExpectations(suite.T())
	suite.metricsServiceMock.AssertExpectations(suite.T())
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Alarm)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore-1, countAfter)

	// Check the response body
	assert.Equal(
		suite.T(),
		"", // empty string
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)
}
