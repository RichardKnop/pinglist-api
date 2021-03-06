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
	stripeSub "github.com/stripe/stripe-go/sub"
	stripeToken "github.com/stripe/stripe-go/token"
)

func (suite *SubscriptionsTestSuite) TestCancelSubscriptionRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.cancelSubscriptionHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *SubscriptionsTestSuite) TestCancelSubscriptionWithoutPermission() {
	// Create a test Stripe customer
	testStripeCustomer, err := suite.service.stripeAdapter.CreateCustomer(
		suite.users[1].OauthUser.Username,
		"", // token
	)
	assert.NoError(suite.T(), err, "Creating test Stripe customer failed")

	// Create a test customer
	testCustomer := NewCustomer(suite.users[1], testStripeCustomer.ID)
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Failed to insert a test customer")

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
	_, err = suite.service.createCard(
		suite.users[1],
		&CardRequest{
			Token: testStripeToken.ID,
		},
	)
	assert.NoError(suite.T(), err, "Creating test card failed")

	// Create a test subscription
	testSubscription, err := suite.service.createSubscription(
		suite.users[1],
		&SubscriptionRequest{
			PlanID: suite.plans[0].ID,
		},
	)
	assert.NoError(suite.T(), err, "Creating test subscription failed")

	// Prepare a request
	r, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://1.2.3.4/v1/subscriptions/%d", testSubscription.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "cancel_subscription", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[2])

	// Count before
	var countBefore int
	suite.db.Model(new(Subscription)).Count(&countBefore)

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
	suite.db.Model(new(Subscription)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the subscription
	subscription := new(Subscription)
	notFound := suite.db.Preload("Customer.User").Preload("Plan").
		First(subscription, testSubscription.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the subscription has not been cancelled
	assert.False(suite.T(), subscription.CancelledAt.Valid)
	assert.False(suite.T(), subscription.EndedAt.Valid)
	assert.Equal(suite.T(), string(stripeSub.Trialing), subscription.Status)

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrCancelSubscriptionPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *SubscriptionsTestSuite) TestCancelSubscription() {
	// Create a test Stripe customer
	testStripeCustomer, err := suite.service.stripeAdapter.CreateCustomer(
		suite.users[1].OauthUser.Username,
		"", // token
	)
	assert.NoError(suite.T(), err, "Creating test Stripe customer failed")

	// Create a test customer
	testCustomer := NewCustomer(suite.users[1], testStripeCustomer.ID)
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Failed to insert a test customer")

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
	_, err = suite.service.createCard(
		suite.users[1],
		&CardRequest{
			Token: testStripeToken.ID,
		},
	)
	assert.NoError(suite.T(), err, "Creating test card failed")

	// Create a test subscription
	testSubscription, err := suite.service.createSubscription(
		suite.users[1],
		&SubscriptionRequest{
			PlanID: suite.plans[0].ID,
		},
	)
	assert.NoError(suite.T(), err, "Creating test subscription failed")

	assert.False(suite.T(), testSubscription.CancelledAt.Valid)
	assert.False(suite.T(), testSubscription.EndedAt.Valid)

	// Prepare a request
	r, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://1.2.3.4/v1/subscriptions/%d", testSubscription.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "cancel_subscription", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Count before
	var countBefore int
	suite.db.Model(new(Subscription)).Count(&countBefore)

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
	suite.db.Model(new(Subscription)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the cancelled subscription
	subscription := new(Subscription)
	notFound := suite.db.Preload("Customer.User").Preload("Plan").
		First(subscription, testSubscription.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the subscription has been cancelled
	// Cancelled at should be set, ended at should still be null
	// Status should remain trialing as subscription will run until the period end
	assert.True(suite.T(), subscription.CancelledAt.Valid)
	assert.False(suite.T(), subscription.EndedAt.Valid)
	assert.Equal(suite.T(), string(stripeSub.Trialing), subscription.Status)

	// Check the response body
	assert.Equal(
		suite.T(),
		"", // empty string
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)
}
