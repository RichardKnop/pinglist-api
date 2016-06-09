package accounts

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/uuid"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestConfirmEmailFailsWithoutAccountAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()

	suite.service.confirmEmailHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated account")
}

func (suite *AccountsTestSuite) TestConfirmEmailReferenceNotFound() {
	bogusUUID := uuid.New()
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/confirmations/%s", bogusUUID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"Basic %s",
			b64.StdEncoding.EncodeToString([]byte("test_client_1:test_secret")),
		),
	)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	if !assert.Equal(suite.T(), 404, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrConfirmationNotFound.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AccountsTestSuite) TestConfirmEmail() {
	var (
		testOauthUser    *oauth.User
		testUser         *User
		testConfirmation *Confirmation
		err              error
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
		"some_slack_incoming_webhook", // slack incoming webhook
		"some_slack_channel",          // slack channel
	)
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser
	testUser.Role = suite.userRole

	// Insert a test confirmation
	testConfirmation = NewConfirmation(testUser)
	err = suite.db.Create(testConfirmation).Error
	assert.NoError(suite.T(), err, "Failed to insert a test confirmation")
	testConfirmation.User = testUser

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/confirmations/%s", testConfirmation.Reference),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"Basic %s",
			b64.StdEncoding.EncodeToString([]byte("test_client_1:test_secret")),
		),
	)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Fetch the updated user
	user := new(User)
	notFound := suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role.Permissions").First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the confirmation
	confirmation := new(Confirmation)
	assert.False(suite.T(), suite.db.Preload("User.OauthUser").
		Last(confirmation).RecordNotFound())

	// And correct data was saved
	assert.True(suite.T(), user.Confirmed)

	// Check the response body
	assert.Equal(
		suite.T(),
		"", // empty string
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)
}
