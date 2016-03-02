package subscriptions

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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

func (suite *SubscriptionsTestSuite) TestStripeWebhook() {
	// Prepare a request
	var payload []byte
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

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())

	// TODO
}
