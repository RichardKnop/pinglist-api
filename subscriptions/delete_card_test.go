package subscriptions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	stripe "github.com/stripe/stripe-go"
	stripeToken "github.com/stripe/stripe-go/token"
)

func (suite *SubscriptionsTestSuite) TestDeleteCardRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.deleteCardHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *SubscriptionsTestSuite) TestDeleteCardWithoutPermission() {
	// Create a test Stripe token
	testStripeToken, err := stripeToken.New(&stripe.TokenParams{
		Card: &stripe.CardParams{
			Number: "4242424242424242",
			Month:  "10",
			Year:   "20",
			CVC:    "123",
		},
		Email: suite.users[1].OauthUser.Username,
	})
	assert.NoError(suite.T(), err, "Creating test Stripe token failed")

	// Create a test card
	testCard, err := suite.service.createCard(
		suite.users[1],
		&CardRequest{
			Token: testStripeToken.ID,
		},
	)
	assert.NoError(suite.T(), err, "Creating test card failed")

	// Prepare a request
	r, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://1.2.3.4/v1/cards/%d", testCard.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "delete_card", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[2])

	// Count before
	var countBefore int
	suite.db.Model(new(Card)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 403, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Card)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrDeleteCardPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *SubscriptionsTestSuite) TestDeleteCard() {
	// Create a test Stripe token
	testStripeToken, err := stripeToken.New(&stripe.TokenParams{
		Card: &stripe.CardParams{
			Number: "4242424242424242",
			Month:  "10",
			Year:   "20",
			CVC:    "123",
		},
		Email: suite.users[1].OauthUser.Username,
	})
	assert.NoError(suite.T(), err, "Creating test Stripe token failed")

	// Create a test card
	testCard, err := suite.service.createCard(
		suite.users[1],
		&CardRequest{
			Token: testStripeToken.ID,
		},
	)
	assert.NoError(suite.T(), err, "Creating test card failed")

	// Prepare a request
	r, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://1.2.3.4/v1/cards/%d", testCard.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "delete_card", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Count before
	var countBefore int
	suite.db.Model(new(Card)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Card)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore-1, countAfter)

	// Try to fetch the deleted card
	notFound := suite.db.First(new(Card), testCard.ID).RecordNotFound()
	assert.True(suite.T(), notFound)

	// Check the response body
	assert.Equal(
		suite.T(),
		"", // empty string
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)
}
