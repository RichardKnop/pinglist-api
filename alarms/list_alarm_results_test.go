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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestListAlarmResults() {
	var (
		today             = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName = "alarm_results_2016_02_09"
		err               error
	)

	// Partition the results table
	err = suite.service.PartitionTable(ResultParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test results
	testResults := []*Result{
		newResult(todaySubTableName, suite.alarms[0], today, 123),
		newResult(todaySubTableName, suite.alarms[0], today.Add(1*time.Hour), 234),
		newResult(todaySubTableName, suite.alarms[0], today.Add(2*time.Hour), 345),
		newResult(todaySubTableName, suite.alarms[0], today.Add(3*time.Hour), 456),
	}
	for _, result := range testResults {
		err := suite.db.Create(result).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d/results", suite.alarms[0].ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_alarm_results", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	var results []*Result
	err = suite.db.Order("id").Find(&results).Error
	assert.NoError(suite.T(), err, "Fetching data failed")

	resultResponses := make([]*ResultResponse, len(results))
	for i, result := range results {
		resultResponse, err := NewResultResponse(result)
		assert.NoError(suite.T(), err, "Creating response object failed")
		resultResponses[i] = resultResponse
	}

	expected := &ListIncidentsResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/results", suite.alarms[0].ID),
				},
				"first": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/results?page=1", suite.alarms[0].ID),
				},
				"last": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/results?page=1", suite.alarms[0].ID),
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				"results": jsonhal.Embedded(resultResponses),
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
