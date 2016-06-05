package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/password"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestUpdateUserRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.updateUserHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AccountsTestSuite) TestUpdateUserFailsWithoutPermission() {
	payload, err := json.Marshal(&UserRequest{
		Email:     "test@user",
		FirstName: "John",
		LastName:  "Reese",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/users/%d", suite.users[2].ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_user_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_user", match.Route.GetName())
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 403, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrUpdateUserPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AccountsTestSuite) TestUpdateUserChangePasswordWhenPasswordEmpty() {
	var (
		testOauthUser   *oauth.User
		testUser        *User
		testAccessToken *oauth.AccessToken
		err             error
	)

	// Insert a test user
	testOauthUser, err = suite.service.oauthService.CreateUser(
		"harold@finch",
		"", // empty password
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser = NewUser(
		suite.accounts[0],
		testOauthUser,
		suite.userRole,
		"some_facebook_id", // facebook ID
		"Harold",
		"Finch",
		true,                 // confirmed
		"some_slack_api_key", // slack API key
		"some_slack_channel", // slack channel
	)
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")

	// Assign related objects
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser
	testUser.Role = suite.userRole

	// Login the test user
	testAccessToken, _, err = suite.service.oauthService.Login(
		suite.accounts[0].OauthClient,
		testUser.OauthUser,
		"read_write", // scope
	)
	assert.NoError(suite.T(), err, "Failed to login the test user")

	payload, err := json.Marshal(&UserRequest{
		Password:    "",
		NewPassword: "some_new_password",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/users/%d", testUser.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAccessToken.Token))

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_user", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(User)).Count(&countBefore)

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
	suite.db.Model(new(User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated user
	user := new(User)
	assert.False(suite.T(), suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role").First(user, testUser.ID).RecordNotFound())

	// Check that the password has changed
	assert.NoError(suite.T(), password.VerifyPassword(
		user.OauthUser.Password.String,
		"some_new_password",
	))

	// And the user meta data is unchanged
	assert.Equal(suite.T(), "harold@finch", user.OauthUser.Username)
	assert.Equal(suite.T(), "Harold", user.FirstName.String)
	assert.Equal(suite.T(), "Finch", user.LastName.String)
	assert.Equal(suite.T(), roles.User, user.Role.ID)
	assert.True(suite.T(), user.Confirmed)
	assert.Equal(suite.T(), "some_slack_api_key", user.SlackAPIKey.String)
	assert.Equal(suite.T(), "some_slack_channel", user.SlackChannel.String)

	// Check the response body
	expected := &UserResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/accounts/users/%d", user.ID),
				},
			},
		},
		ID:           user.ID,
		Email:        "harold@finch",
		FirstName:    "Harold",
		LastName:     "Finch",
		Role:         roles.User,
		Confirmed:    true,
		SlackAPIKey:  "some_slack_api_key",
		SlackChannel: "some_slack_channel",
		CreatedAt:    util.FormatTime(user.CreatedAt),
		UpdatedAt:    util.FormatTime(user.UpdatedAt),
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

func (suite *AccountsTestSuite) TestUpdateUserChangePassword() {
	var (
		testOauthUser   *oauth.User
		testUser        *User
		testAccessToken *oauth.AccessToken
		err             error
	)

	// Insert a test user
	testOauthUser, err = suite.service.oauthService.CreateUser(
		"harold@finch",
		"test_password",
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser = NewUser(
		suite.accounts[0],
		testOauthUser,
		suite.userRole,
		"", // facebook ID
		"Harold",
		"Finch",
		false, // confirmed
		"",    // slack API key
		"",    // slack channel
	)
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser
	testUser.Role = suite.userRole

	// Login the test user
	testAccessToken, _, err = suite.service.oauthService.Login(
		suite.accounts[0].OauthClient,
		testUser.OauthUser,
		"read_write", // scope
	)
	assert.NoError(suite.T(), err, "Failed to login the test user")

	payload, err := json.Marshal(&UserRequest{
		Password:    "test_password",
		NewPassword: "some_new_password",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/users/%d", testUser.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAccessToken.Token))

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_user", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(User)).Count(&countBefore)

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
	suite.db.Model(new(User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated user
	user := new(User)
	assert.False(suite.T(), suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role").First(user, testUser.ID).RecordNotFound())

	// Check that the password has changed
	assert.Error(suite.T(), password.VerifyPassword(
		user.OauthUser.Password.String,
		"test_password",
	))
	assert.NoError(suite.T(), password.VerifyPassword(
		user.OauthUser.Password.String,
		"some_new_password",
	))

	// And the user meta data is unchanged
	assert.Equal(suite.T(), "harold@finch", user.OauthUser.Username)
	assert.Equal(suite.T(), "Harold", user.FirstName.String)
	assert.Equal(suite.T(), "Finch", user.LastName.String)
	assert.Equal(suite.T(), roles.User, user.Role.ID)
	assert.False(suite.T(), user.Confirmed)
	assert.False(suite.T(), user.SlackAPIKey.Valid)
	assert.False(suite.T(), user.SlackChannel.Valid)

	// Check the response body
	expected := &UserResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/accounts/users/%d", user.ID),
				},
			},
		},
		ID:        user.ID,
		Email:     "harold@finch",
		FirstName: "Harold",
		LastName:  "Finch",
		Role:      roles.User,
		Confirmed: false,
		CreatedAt: util.FormatTime(user.CreatedAt),
		UpdatedAt: util.FormatTime(user.UpdatedAt),
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

func (suite *AccountsTestSuite) TestUpdateUser() {
	var (
		testOauthUser   *oauth.User
		testUser        *User
		testAccessToken *oauth.AccessToken
		err             error
	)

	// Insert a test user
	testOauthUser, err = suite.service.oauthService.CreateUser(
		"harold@finch",
		"test_password",
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser = NewUser(
		suite.accounts[0],
		testOauthUser,
		suite.userRole,
		"",    // facebook ID
		"",    // first name
		"",    // last name
		false, // confirmed
		"",    // slack API key
		"",    // slack channel
	)
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser
	testUser.Role = suite.userRole

	// Login the test user
	testAccessToken, _, err = suite.service.oauthService.Login(
		suite.accounts[0].OauthClient,
		testUser.OauthUser,
		"read_write", // scope
	)
	assert.NoError(suite.T(), err, "Failed to login the test user")

	payload, err := json.Marshal(&UserRequest{
		FirstName:    "Harold",
		LastName:     "Finch",
		SlackAPIKey:  "harolds_slack_api_key",
		SlackChannel: "harolds_slack_channel",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/users/%d", testUser.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAccessToken.Token))

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_user", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(User)).Count(&countBefore)

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
	suite.db.Model(new(User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated user
	user := new(User)
	assert.False(suite.T(), suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role").First(user, testUser.ID).RecordNotFound())

	// Check that the password has NOT changed
	assert.NoError(suite.T(), password.VerifyPassword(
		user.OauthUser.Password.String,
		"test_password",
	))

	// And correct data was saved
	assert.Equal(suite.T(), "harold@finch", user.OauthUser.Username)
	assert.Equal(suite.T(), "Harold", user.FirstName.String)
	assert.Equal(suite.T(), "Finch", user.LastName.String)
	assert.Equal(suite.T(), roles.User, user.Role.ID)
	assert.False(suite.T(), user.Confirmed)
	assert.Equal(suite.T(), "harolds_slack_api_key", user.SlackAPIKey.String)
	assert.Equal(suite.T(), "harolds_slack_channel", user.SlackChannel.String)

	// Check the response body
	expected := &UserResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/accounts/users/%d", user.ID),
				},
			},
		},
		ID:           user.ID,
		Email:        "harold@finch",
		FirstName:    "Harold",
		LastName:     "Finch",
		Role:         roles.User,
		Confirmed:    false,
		SlackAPIKey:  "harolds_slack_api_key",
		SlackChannel: "harolds_slack_channel",
		CreatedAt:    util.FormatTime(user.CreatedAt),
		UpdatedAt:    util.FormatTime(user.UpdatedAt),
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
