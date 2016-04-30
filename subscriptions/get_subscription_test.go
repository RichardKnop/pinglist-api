package subscriptions

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
	stripe "github.com/stripe/stripe-go"
	stripeSub "github.com/stripe/stripe-go/sub"
	stripeToken "github.com/stripe/stripe-go/token"
)

func (suite *SubscriptionsTestSuite) TestGetSubscriptionRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.getSubscriptionHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *SubscriptionsTestSuite) TestGetSubscriptionWithoutPermission() {
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
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/subscriptions/%d", testSubscription.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "get_subscription", match.Route.GetName())
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

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrGetSubscriptionPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *SubscriptionsTestSuite) TestGetSubscription() {
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
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/subscriptions/%d", testSubscription.ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "get_subscription", match.Route.GetName())
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
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Subscription)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	planResponse, err := NewPlanResponse(testSubscription.Plan)
	assert.NoError(suite.T(), err, "Creating response object failed")
	expected := &SubscriptionResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/subscriptions/%d", testSubscription.ID),
				},
			},
			Embedded: map[string]jsonhal.Embedded{
				"plan": jsonhal.Embedded(planResponse),
			},
		},
		ID:             testSubscription.ID,
		SubscriptionID: testSubscription.SubscriptionID,
		StartedAt:      testSubscription.StartedAt.Time.UTC().Format(time.RFC3339),
		PeriodStart:    testSubscription.PeriodStart.Time.UTC().Format(time.RFC3339),
		PeriodEnd:      testSubscription.PeriodEnd.Time.UTC().Format(time.RFC3339),
		TrialStart:     testSubscription.TrialStart.Time.UTC().Format(time.RFC3339),
		TrialEnd:       testSubscription.TrialEnd.Time.UTC().Format(time.RFC3339),
		Status:         string(stripeSub.Trialing),
		CreatedAt:      testSubscription.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      testSubscription.UpdatedAt.UTC().Format(time.RFC3339),
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
