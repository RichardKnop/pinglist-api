package subscriptions

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	stripe "github.com/stripe/stripe-go"
)

func (suite *SubscriptionsTestSuite) TestStripeWebhookNoPayload() {
	// Prepare a request
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/stripe-webhook",
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "stripe_webhook", match.Route.GetName())
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())

	assert.Equal(suite.T(), 400, w.Code, "Expected a 400 (Bad Request) response")
}

func (suite *SubscriptionsTestSuite) TestStripeWebhookBogusEventID() {
	// Prepare a request
	payload, err := json.Marshal(&stripe.Event{ID: "bogus"})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/stripe-webhook",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "stripe_webhook", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(StripeEventLog)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 404, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(StripeEventLog)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)
}
