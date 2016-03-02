package accounts

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestGetAccountFromQueryString() {
	var (
		account *Account
		err     error
	)

	// Create a test request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/foobar?account_id=%d", suite.accounts[0].ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// Let's try without user authentication first
	account, err = suite.service.GetAccountFromQueryString(r)

	// Account should be nil
	assert.Nil(suite.T(), account)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errUserAuthenticationRequired, err)
	}

	// Now, let's set an authenticated user and try again
	context.Set(r, authenticatedUserKey, suite.users[1])
	account, err = suite.service.GetAccountFromQueryString(r)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct account object should be returned
	if assert.NotNil(suite.T(), account) {
		assert.Equal(suite.T(), "Test Account 1", account.Name)
	}
}

func (suite *AccountsTestSuite) TestGetUserFromQueryString() {
	var (
		user *User
		err  error
	)

	// Create a test request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/foobar?user_id=%d", suite.users[1].ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// Let's try without user authentication first
	user, err = suite.service.GetUserFromQueryString(r)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errUserAuthenticationRequired, err)
	}

	// Now, let's set an authenticated user and try again
	context.Set(r, authenticatedUserKey, suite.users[1])
	user, err = suite.service.GetUserFromQueryString(r)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
	}
}
