package accounts

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
	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestUpdateUser() {
	payload, err := json.Marshal(&UserRequest{
		FirstName: "John",
		LastName:  "Reese",
	})
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/users/%d", suite.users[1].ID),
		bytes.NewBuffer(payload),
	)
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Set("Authorization", "Bearer test_user_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_user", match.Route.GetName())
	}

	// Count before
	var (
		countBefore int
	)
	suite.db.Model(new(User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter int
	)
	suite.db.Model(new(User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated user
	user := new(User)
	assert.False(suite.T(), suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role").First(user, suite.users[1].ID).RecordNotFound())

	// Check that the correct data was saved
	assert.Equal(suite.T(), "John", user.FirstName.String)
	assert.Equal(suite.T(), "Reese", user.LastName.String)
	assert.Equal(suite.T(), roles.User, user.Role.Name)
	assert.True(suite.T(), user.Confirmed)
	assert.Equal(suite.T(), uint(1), user.MaxAlarms)

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
		Email:     "test@user",
		FirstName: "John",
		LastName:  "Reese",
		Role:      roles.User,
		MaxAlarms: 1,
		Confirmed: true,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
	}
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		suite.T(),
		string(expectedJSON),
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)
}
