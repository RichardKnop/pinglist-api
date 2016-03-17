package subscriptions

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

func (suite *SubscriptionsTestSuite) TestListSubscriptions() {
	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		"http://1.2.3.4/v1/subscriptions",
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "list_subscriptions", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[0])
	suite.mockUserFiltering(suite.users[0])

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
	var subscriptions []*Subscription
	err = suite.db.Preload("Customer").Preload("Plan").Find(&subscriptions).Error
	assert.NoError(suite.T(), err, "Fetching data failed")

	subscriptionResponses := make([]*SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		subscriptionResponse, err := NewSubscriptionResponse(subscription)
		assert.NoError(suite.T(), err, "Creating response object failed")
		subscriptionResponses[i] = subscriptionResponse
	}

	expected := &ListSubscriptionsResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: "/v1/subscriptions",
				},
				"first": &jsonhal.Link{
					Href: "/v1/subscriptions?page=1",
				},
				"last": &jsonhal.Link{
					Href: "/v1/subscriptions?page=1",
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				"subscriptions": jsonhal.Embedded(subscriptionResponses),
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
