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

func (suite *SubscriptionsTestSuite) TestCreateCardRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.createCardHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *SubscriptionsTestSuite) TestCreateCardExistingValidCustomer() {
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
	testStripeCustomer, err := suite.service.stripeAdapter.CreateCustomer(
		suite.users[1].OauthUser.Username,
		"",
	)
	assert.NoError(suite.T(), err, "Creating test Stripe customer failed")

	// Create a test customer
	testCustomer := NewCustomer(suite.users[1], testStripeCustomer.ID)
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Creating test customer failed")

	// Prepare a request
	payload, err := json.Marshal(&CardRequest{
		Token: testStripeToken.ID,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/cards",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_card", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Count before
	var (
		countBefore         int
		customerCountBefore int
	)
	suite.db.Model(new(Card)).Count(&countBefore)
	suite.db.Model(new(Customer)).Count(&customerCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter         int
		customerCountAfter int
	)
	suite.db.Model(new(Card)).Count(&countAfter)
	suite.db.Model(new(Customer)).Count(&customerCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), customerCountBefore, customerCountAfter)

	var (
		customer *Customer
		card     *Card
		notFound bool
	)

	// Fetch the customer
	customer = new(Customer)
	notFound = suite.db.Preload("User").First(customer, testCustomer.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the created card
	card = new(Card)
	notFound = suite.db.Preload("Customer.User").Last(card).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), testCustomer.ID, customer.ID)
	assert.Equal(suite.T(), customer.ID, card.Customer.ID)
	assert.Equal(suite.T(), "Visa", card.Brand)
	assert.Equal(suite.T(), "4242", card.LastFour)
	assert.Equal(suite.T(), uint(10), card.ExpMonth)
	assert.Equal(suite.T(), uint(2020), card.ExpYear)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/cards/%d", card.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	expected := &CardResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/cards/%d", card.ID),
				},
			},
		},
		ID:        card.ID,
		Brand:     card.Brand,
		LastFour:  card.LastFour,
		ExpMonth:  card.ExpMonth,
		ExpYear:   card.ExpYear,
		CreatedAt: card.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: card.UpdatedAt.UTC().Format(time.RFC3339),
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

func (suite *SubscriptionsTestSuite) TestCreateCardExistingInvalidCustomer() {
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
	testCustomer := NewCustomer(suite.users[1], "bogus_customer_id")
	err = suite.db.Create(testCustomer).Error
	assert.NoError(suite.T(), err, "Creating test customer failed")

	// Prepare a request
	payload, err := json.Marshal(&CardRequest{
		Token: testStripeToken.ID,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/cards",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_card", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Count before
	var (
		countBefore         int
		customerCountBefore int
	)
	suite.db.Model(new(Card)).Count(&countBefore)
	suite.db.Model(new(Customer)).Count(&customerCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter         int
		customerCountAfter int
	)
	suite.db.Model(new(Card)).Count(&countAfter)
	suite.db.Model(new(Customer)).Count(&customerCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), customerCountBefore, customerCountAfter)

	var (
		customer *Customer
		card     *Card
		notFound bool
	)

	// Fetch the customer
	customer = new(Customer)
	notFound = suite.db.Preload("User").Last(customer).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the created card
	card = new(Card)
	notFound = suite.db.Preload("Customer.User").Last(card).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.NotEqual(suite.T(), testCustomer.ID, customer.ID)
	assert.Equal(suite.T(), customer.ID, card.Customer.ID)
	assert.Equal(suite.T(), "Visa", card.Brand)
	assert.Equal(suite.T(), "4242", card.LastFour)
	assert.Equal(suite.T(), uint(10), card.ExpMonth)
	assert.Equal(suite.T(), uint(2020), card.ExpYear)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/cards/%d", card.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	expected := &CardResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/cards/%d", card.ID),
				},
			},
		},
		ID:        card.ID,
		Brand:     card.Brand,
		LastFour:  card.LastFour,
		ExpMonth:  card.ExpMonth,
		ExpYear:   card.ExpYear,
		CreatedAt: card.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: card.UpdatedAt.UTC().Format(time.RFC3339),
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

func (suite *SubscriptionsTestSuite) TestCreateCardNewCustomer() {
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
	payload, err := json.Marshal(&CardRequest{
		Token: theStripeToken.ID,
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/cards",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_card", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[1])

	// Count before
	var (
		countBefore         int
		customerCountBefore int
	)
	suite.db.Model(new(Card)).Count(&countBefore)
	suite.db.Model(new(Customer)).Count(&customerCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter         int
		customerCountAfter int
	)
	suite.db.Model(new(Card)).Count(&countAfter)
	suite.db.Model(new(Customer)).Count(&customerCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), customerCountBefore+1, customerCountAfter)

	var (
		customer *Customer
		card     *Card
		notFound bool
	)

	// Fetch the customer
	customer = new(Customer)
	notFound = suite.db.Preload("User").Last(customer).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the created card
	card = new(Card)
	notFound = suite.db.Preload("Customer.User").Last(card).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the correct data was saved
	assert.Equal(suite.T(), suite.users[1].ID, uint(customer.UserID.Int64))
	assert.Equal(suite.T(), customer.ID, card.Customer.ID)
	assert.Equal(suite.T(), "Visa", card.Brand)
	assert.Equal(suite.T(), "4242", card.LastFour)
	assert.Equal(suite.T(), uint(10), card.ExpMonth)
	assert.Equal(suite.T(), uint(2020), card.ExpYear)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/cards/%d", card.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	expected := &CardResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/cards/%d", card.ID),
				},
			},
		},
		ID:        card.ID,
		Brand:     card.Brand,
		LastFour:  card.LastFour,
		ExpMonth:  card.ExpMonth,
		ExpYear:   card.ExpYear,
		CreatedAt: card.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: card.UpdatedAt.UTC().Format(time.RFC3339),
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
