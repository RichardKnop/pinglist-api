package subscriptions

import (
	"github.com/stretchr/testify/assert"
)

func (suite *SubscriptionsTestSuite) TestFindCustomerByID() {
	var (
		customer *Customer
		err      error
	)

	// When we try to find a customer with a bogus ID
	customer, err = suite.service.FindCustomerByID(12345)

	// Customer object should be nil
	assert.Nil(suite.T(), customer)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrCustomerNotFound, err)
	}

	// When we try to find a plan with a valid ID
	customer, err = suite.service.FindCustomerByID(suite.customers[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct customer object should be returned
	if assert.NotNil(suite.T(), customer) {
		assert.Equal(suite.T(), suite.customers[0].ID, customer.ID)
		assert.Equal(suite.T(), suite.users[0].ID, customer.User.ID)
	}
}

func (suite *SubscriptionsTestSuite) TestFindCustomerByUserID() {
	var (
		customer *Customer
		err      error
	)

	// When we try to find a customer with a bogus user ID
	customer, err = suite.service.FindCustomerByUserID(12345)

	// Customer object should be nil
	assert.Nil(suite.T(), customer)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrCustomerNotFound, err)
	}

	// When we try to find a plan with a valid customer ID
	customer, err = suite.service.FindCustomerByUserID(suite.users[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct customer object should be returned
	if assert.NotNil(suite.T(), customer) {
		assert.Equal(suite.T(), suite.customers[0].ID, customer.ID)
		assert.Equal(suite.T(), suite.users[0].ID, customer.User.ID)
	}
}

func (suite *SubscriptionsTestSuite) TestFindCustomerByCustomerID() {
	var (
		customer *Customer
		err      error
	)

	// When we try to find a customer with a bogus customer ID
	customer, err = suite.service.FindCustomerByCustomerID("bogus")

	// Customer object should be nil
	assert.Nil(suite.T(), customer)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrCustomerNotFound, err)
	}

	// When we try to find a plan with a valid customer ID
	customer, err = suite.service.FindCustomerByCustomerID(suite.customers[0].CustomerID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct customer object should be returned
	if assert.NotNil(suite.T(), customer) {
		assert.Equal(suite.T(), suite.customers[0].ID, customer.ID)
		assert.Equal(suite.T(), suite.users[0].ID, customer.User.ID)
	}
}
