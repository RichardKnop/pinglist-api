package alarms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestListAlarmIncidents() {
	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/alarms/%d/incidents", suite.alarms[0].ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_alarm_incidents", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	var incidents []*Incident
	err = suite.db.Order("id").Find(&incidents).Error
	assert.NoError(suite.T(), err, "Fetching data failed")

	incidentResponses := make([]*IncidentResponse, len(incidents))
	for i, incident := range incidents {
		incidentResponse, err := NewIncidentResponse(incident)
		assert.NoError(suite.T(), err, "Creating response object failed")
		incidentResponses[i] = incidentResponse
	}

	expected := &ListIncidentsResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/incidents", suite.alarms[0].ID),
				},
				"first": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/incidents?page=1", suite.alarms[0].ID),
				},
				"last": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/alarms/%d/incidents?page=1", suite.alarms[0].ID),
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				"incidents": jsonhal.Embedded(incidentResponses),
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
