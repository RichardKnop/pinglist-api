package subscriptions

import (
	"bytes"
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
	stripeToken "github.com/stripe/stripe-go/token"
)

func (suite *SubscriptionsTestSuite) TestCreateSubscriptionRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.createSubscriptionHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *SubscriptionsTestSuite) TestCreateSubscriptionExistingValidCustomer() {
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

	// Create a test Stripe customer
	teststripeCustomer, err := suite.service.stripeAdapter.CreateCustomer(
		testStripeToken.Email,
		testStripeToken.ID,
	)
	assert.NoError(suite.T(), err, "Creating test Stripe customer failed")

	// Create a test customer
	testCustomer := newCustomer(suite.users[1], teststripeCustomer.ID)
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Creating test customer failed")

	// Prepare a request
	payload, err := json.Marshal(&SubscriptionRequest{
		StripeToken: testStripeToken.ID,
		StripeEmail: testStripeToken.Email,
		PlanID:      suite.plans[0].ID,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/subscriptions",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_subscription", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Count before
	var (
		countBefore         int
		customerCountBefore int
	)
	suite.db.Model(new(Subscription)).Count(&countBefore)
	suite.db.Model(new(Customer)).Count(&customerCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter         int
		customerCountAfter int
	)
	suite.db.Model(new(Subscription)).Count(&countAfter)
	suite.db.Model(new(Customer)).Count(&customerCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), customerCountBefore, customerCountAfter)

	var (
		customer     *Customer
		subscription *Subscription
		notFound     bool
	)

	// Fetch the customer
	customer = new(Customer)
	notFound = suite.db.Preload("User").
		First(customer, testCustomer.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the created subscription
	subscription = new(Subscription)
	notFound = suite.db.Preload("Customer.User").Preload("Plan").
		Last(subscription).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), testCustomer.ID, customer.ID)
	assert.Equal(suite.T(), customer.ID, subscription.Customer.ID)
	assert.Equal(suite.T(), suite.plans[0].ID, uint(subscription.PlanID.Int64))
	assert.True(suite.T(), subscription.StartedAt.Valid)
	assert.False(suite.T(), subscription.CancelledAt.Valid)
	assert.False(suite.T(), subscription.EndedAt.Valid)
	assert.True(suite.T(), subscription.PeriodStart.Valid)
	assert.True(suite.T(), subscription.PeriodEnd.Valid)
	assert.True(suite.T(), subscription.TrialStart.Valid)
	assert.True(suite.T(), subscription.TrialEnd.Valid)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/subscriptions/%d", subscription.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	planResponse, err := NewPlanResponse(subscription.Plan)
	assert.NoError(suite.T(), err, "Creating response object failed")
	expected := &SubscriptionResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/subscriptions/%d", subscription.ID),
				},
			},
			Embedded: map[string]jsonhal.Embedded{
				"plan": jsonhal.Embedded(planResponse),
			},
		},
		ID:             subscription.ID,
		SubscriptionID: subscription.SubscriptionID,
		StartedAt:      subscription.StartedAt.Time.UTC().Format(time.RFC3339),
		PeriodStart:    subscription.PeriodStart.Time.UTC().Format(time.RFC3339),
		PeriodEnd:      subscription.PeriodEnd.Time.UTC().Format(time.RFC3339),
		TrialStart:     subscription.TrialStart.Time.UTC().Format(time.RFC3339),
		TrialEnd:       subscription.TrialEnd.Time.UTC().Format(time.RFC3339),
		CreatedAt:      subscription.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      subscription.UpdatedAt.UTC().Format(time.RFC3339),
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

func (suite *SubscriptionsTestSuite) TestCreateSubscriptionExistingInvalidCustomer() {
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

	// Create a test customer
	testCustomer := newCustomer(suite.users[1], "bogus_customer_id")
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Creating test customer failed")

	// Prepare a request
	payload, err := json.Marshal(&SubscriptionRequest{
		StripeToken: testStripeToken.ID,
		StripeEmail: testStripeToken.Email,
		PlanID:      suite.plans[0].ID,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/subscriptions",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_subscription", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Count before
	var (
		countBefore         int
		customerCountBefore int
	)
	suite.db.Model(new(Subscription)).Count(&countBefore)
	suite.db.Model(new(Customer)).Count(&customerCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter         int
		customerCountAfter int
	)
	suite.db.Model(new(Subscription)).Count(&countAfter)
	suite.db.Model(new(Customer)).Count(&customerCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), customerCountBefore, customerCountAfter)

	var (
		customer     *Customer
		subscription *Subscription
		notFound     bool
	)

	// Fetch the customer
	customer = new(Customer)
	notFound = suite.db.Preload("User").Last(customer).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the created subscription
	subscription = new(Subscription)
	notFound = suite.db.Preload("Customer.User").Preload("Plan").
		Last(subscription).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.NotEqual(suite.T(), testCustomer.ID, customer.ID)
	assert.Equal(suite.T(), customer.ID, subscription.Customer.ID)
	assert.Equal(suite.T(), suite.plans[0].ID, uint(subscription.PlanID.Int64))
	assert.True(suite.T(), subscription.StartedAt.Valid)
	assert.False(suite.T(), subscription.CancelledAt.Valid)
	assert.False(suite.T(), subscription.EndedAt.Valid)
	assert.True(suite.T(), subscription.PeriodStart.Valid)
	assert.True(suite.T(), subscription.PeriodEnd.Valid)
	assert.True(suite.T(), subscription.TrialStart.Valid)
	assert.True(suite.T(), subscription.TrialEnd.Valid)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/subscriptions/%d", subscription.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	planResponse, err := NewPlanResponse(subscription.Plan)
	assert.NoError(suite.T(), err, "Creating response object failed")
	expected := &SubscriptionResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/subscriptions/%d", subscription.ID),
				},
			},
			Embedded: map[string]jsonhal.Embedded{
				"plan": jsonhal.Embedded(planResponse),
			},
		},
		ID:             subscription.ID,
		SubscriptionID: subscription.SubscriptionID,
		StartedAt:      subscription.StartedAt.Time.UTC().Format(time.RFC3339),
		PeriodStart:    subscription.PeriodStart.Time.UTC().Format(time.RFC3339),
		PeriodEnd:      subscription.PeriodEnd.Time.UTC().Format(time.RFC3339),
		TrialStart:     subscription.TrialStart.Time.UTC().Format(time.RFC3339),
		TrialEnd:       subscription.TrialEnd.Time.UTC().Format(time.RFC3339),
		CreatedAt:      subscription.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      subscription.UpdatedAt.UTC().Format(time.RFC3339),
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

func (suite *SubscriptionsTestSuite) TestCreateSubscriptionNewCustomer() {
	// Create a test Stripe token
	theStripeToken, err := stripeToken.New(&stripe.TokenParams{
		Card: &stripe.CardParams{
			Number: "4242424242424242",
			Month:  "10",
			Year:   "20",
			CVC:    "123",
		},
		Email: suite.users[1].OauthUser.Username,
	})
	assert.NoError(suite.T(), err, "Creating test Stripe token failed")

	// Prepare a request
	payload, err := json.Marshal(&SubscriptionRequest{
		StripeToken: theStripeToken.ID,
		StripeEmail: theStripeToken.Email,
		PlanID:      suite.plans[0].ID,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/subscriptions",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_subscription", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Count before
	var (
		countBefore         int
		customerCountBefore int
	)
	suite.db.Model(new(Subscription)).Count(&countBefore)
	suite.db.Model(new(Customer)).Count(&customerCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter         int
		customerCountAfter int
	)
	suite.db.Model(new(Subscription)).Count(&countAfter)
	suite.db.Model(new(Customer)).Count(&customerCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), customerCountBefore+1, customerCountAfter)

	var (
		customer     *Customer
		subscription *Subscription
		notFound     bool
	)

	// Fetch the created customer
	customer = new(Customer)
	notFound = suite.db.Preload("User").Last(customer).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the created subscription
	subscription = new(Subscription)
	notFound = suite.db.Preload("Customer.User").Preload("Plan").
		Last(subscription).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), suite.users[1].ID, uint(customer.UserID.Int64))
	assert.Equal(suite.T(), customer.ID, subscription.Customer.ID)
	assert.Equal(suite.T(), suite.plans[0].ID, uint(subscription.PlanID.Int64))
	assert.True(suite.T(), subscription.StartedAt.Valid)
	assert.False(suite.T(), subscription.CancelledAt.Valid)
	assert.False(suite.T(), subscription.EndedAt.Valid)
	assert.True(suite.T(), subscription.PeriodStart.Valid)
	assert.True(suite.T(), subscription.PeriodEnd.Valid)
	assert.True(suite.T(), subscription.TrialStart.Valid)
	assert.True(suite.T(), subscription.TrialEnd.Valid)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/subscriptions/%d", subscription.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	planResponse, err := NewPlanResponse(subscription.Plan)
	assert.NoError(suite.T(), err, "Creating response object failed")
	expected := &SubscriptionResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/subscriptions/%d", subscription.ID),
				},
			},
			Embedded: map[string]jsonhal.Embedded{
				"plan": jsonhal.Embedded(planResponse),
			},
		},
		ID:             subscription.ID,
		SubscriptionID: subscription.SubscriptionID,
		StartedAt:      subscription.StartedAt.Time.UTC().Format(time.RFC3339),
		PeriodStart:    subscription.PeriodStart.Time.UTC().Format(time.RFC3339),
		PeriodEnd:      subscription.PeriodEnd.Time.UTC().Format(time.RFC3339),
		TrialStart:     subscription.TrialStart.Time.UTC().Format(time.RFC3339),
		TrialEnd:       subscription.TrialEnd.Time.UTC().Format(time.RFC3339),
		CreatedAt:      subscription.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      subscription.UpdatedAt.UTC().Format(time.RFC3339),
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