package alarms

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestListRegions() {
	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		"http://1.2.3.4/v1/alarms/regions",
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_regions", match.Route.GetName())
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

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	regionResponses := make([]*RegionResponse, len(suite.regions))
	for i, region := range suite.regions {
		regionResponse, err := NewRegionResponse(region)
		assert.NoError(suite.T(), err, "Creating response object failed")
		regionResponses[i] = regionResponse
	}

	expected := &ListRegionsResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: "/v1/alarms/regions",
				},
			},
			Embedded: map[string]jsonhal.Embedded{
				"regions": jsonhal.Embedded(regionResponses),
			},
		},
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
