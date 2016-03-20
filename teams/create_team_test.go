package teams

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
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *TeamsTestSuite) TestCreateTeamRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.createTeamHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *TeamsTestSuite) TestCreateTeamFailsWhenUserAlreadyOwnsOne() {
	// Prepare a request
	payload, err := json.Marshal(&TeamRequest{
		Name:    "Test Team 2",
		Members: []*TeamMemberRequest{},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/teams",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_team", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[0])

	// Count before
	var countBefore int
	suite.db.Model(new(Team)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 400, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Team)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrUserCanOnlyCreateOneTeam.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *TeamsTestSuite) TestCreateTeamWithoutMembers() {
	// Prepare a request
	payload, err := json.Marshal(&TeamRequest{
		Name:    "Test Team 2",
		Members: []*TeamMemberRequest{},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/teams",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_team", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Mock find active subscription
	suite.mockFindActiveSubscription(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Team)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Team)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created team
	team := new(Team)
	notFound := suite.db.Preload("Owner.OauthUser").Preload("Members.OauthUser").
		Last(team).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), "test@user", team.Owner.OauthUser.Username)
	assert.Equal(suite.T(), "Test Team 2", team.Name)
	assert.Equal(suite.T(), 0, len(team.Members))

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/accounts/teams/%d", team.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	expected := &TeamResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/accounts/teams/%d", team.ID),
				},
			},
			Embedded: map[string]jsonhal.Embedded{
				"members": jsonhal.Embedded([]*accounts.UserResponse{}),
			},
		},
		ID:        team.ID,
		Name:      "Test Team 2",
		CreatedAt: team.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: team.CreatedAt.UTC().Format(time.RFC3339),
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

func (suite *TeamsTestSuite) TestCreateTeamMaxMemberLimitReached() {
	// Prepare a request
	payload, err := json.Marshal(&TeamRequest{
		Name: "Test Team 2",
		Members: []*TeamMemberRequest{
			&TeamMemberRequest{suite.users[2].ID},
		},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/teams",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_team", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Mock find active subscription
	suite.mockFindActiveSubscription(
		suite.users[1].ID,
		nil,
		subscriptions.ErrUserHasNoActiveSubscription,
	)

	// Mock find user
	suite.mockFindUser(suite.users[2].ID, suite.users[2], nil)

	// Count before
	var countBefore int
	suite.db.Model(new(Team)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 400, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Team)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrMaxTeamMembersLimitReached.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *TeamsTestSuite) TestCreateTeam() {
	// Prepare a request
	payload, err := json.Marshal(&TeamRequest{
		Name: "Test Team 2",
		Members: []*TeamMemberRequest{
			&TeamMemberRequest{suite.users[2].ID},
		},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/teams",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_team", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Mock find active subscription
	suite.mockFindActiveSubscription(
		suite.users[1].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxTeamMembers: 10,
			},
		},
		nil,
	)

	// Mock find user
	suite.mockFindUser(suite.users[2].ID, suite.users[2], nil)

	// Count before
	var countBefore int
	suite.db.Model(new(Team)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.subscriptionsServiceMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Team)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created team
	team := new(Team)
	notFound := suite.db.Preload("Owner.OauthUser").Preload("Members.OauthUser").
		Last(team).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), "test@user", team.Owner.OauthUser.Username)
	assert.Equal(suite.T(), "Test Team 2", team.Name)
	assert.Equal(suite.T(), 1, len(team.Members))
	assert.Equal(suite.T(), "test@user2", team.Members[0].OauthUser.Username)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/accounts/teams/%d", team.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	memberResponse, err := accounts.NewUserResponse(suite.users[2])
	assert.Nil(suite.T(), err, "Error creating user response")
	expected := &TeamResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/accounts/teams/%d", team.ID),
				},
			},
			Embedded: map[string]jsonhal.Embedded{
				"members": jsonhal.Embedded([]*accounts.UserResponse{memberResponse}),
			},
		},
		ID:        team.ID,
		Name:      "Test Team 2",
		CreatedAt: team.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: team.CreatedAt.UTC().Format(time.RFC3339),
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