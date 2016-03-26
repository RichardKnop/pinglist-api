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

func (suite *TeamsTestSuite) TestUpdateTeamRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.updateTeamHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *TeamsTestSuite) TestUpdateTeamFailsWithoutPermission() {
	// Prepare a request
	payload, err := json.Marshal(&TeamRequest{
		Name:    "Test Team 1",
		Members: []*TeamMemberRequest{},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/teams/%d", suite.teams[0].ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_team", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[1])

	// Count before
	var countBefore int
	suite.db.Model(new(Team)).Count(&countBefore)

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
	suite.db.Model(new(Team)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrUpdateTeamPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *TeamsTestSuite) TestUpdateTeamMaxTeamsLimitReached() {
	// Prepare a request
	payload, err := json.Marshal(&TeamRequest{
		Name:    "Test Team 1",
		Members: []*TeamMemberRequest{},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/teams/%d", suite.teams[0].ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_team", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[0])

	// Mock find active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[0].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxTeams:          0,
				MaxMembersPerTeam: 0,
			},
		},
		nil,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Team)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 400, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Team)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrMaxTeamsLimitReached.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *TeamsTestSuite) TestUpdateTeamMaxMembersPerTeamLimitReached() {
	// Prepare a request
	payload, err := json.Marshal(&TeamRequest{
		Name: "Test Team 1",
		Members: []*TeamMemberRequest{
			&TeamMemberRequest{suite.users[1].ID},
			&TeamMemberRequest{suite.users[2].ID},
		},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/teams/%d", suite.teams[0].ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_team", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[0])

	// Mock find active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[0].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxTeams:          5,
				MaxMembersPerTeam: 1,
			},
		},
		nil,
	)

	// Count before
	var countBefore int
	suite.db.Model(new(Team)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 400, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(Team)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrMaxMembersPerTeamLimitReached.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *TeamsTestSuite) TestUpdateTeam() {
	// Prepare a request
	payload, err := json.Marshal(&TeamRequest{
		Name: "Test Team 1 Updated",
		Members: []*TeamMemberRequest{
			&TeamMemberRequest{suite.users[1].ID},
			&TeamMemberRequest{suite.users[2].ID},
		},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/teams/%d", suite.teams[0].ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_team", match.Route.GetName())
	}

	// Mock authentication
	suite.mockAuthentication(suite.users[0])

	// Mock find active subscription
	suite.mockFindActiveSubscriptionByUserID(
		suite.users[0].ID,
		&subscriptions.Subscription{
			Plan: &subscriptions.Plan{
				MaxTeams:          5,
				MaxMembersPerTeam: 10,
			},
		},
		nil,
	)

	// Mock find users
	suite.mockFindUser(suite.users[1].ID, suite.users[1], nil)
	suite.mockFindUser(suite.users[2].ID, suite.users[2], nil)

	// Count before
	var countBefore int
	suite.db.Model(new(Team)).Count(&countBefore)

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
	suite.db.Model(new(Team)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the created team
	team := new(Team)
	notFound := suite.db.Preload("Owner.OauthUser").Preload("Members.OauthUser").
		First(team, suite.teams[0].ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), "test@superuser", team.Owner.OauthUser.Username)
	assert.Equal(suite.T(), "Test Team 1 Updated", team.Name)
	assert.Equal(suite.T(), 2, len(team.Members))
	assert.Equal(suite.T(), "test@user", team.Members[0].OauthUser.Username)
	assert.Equal(suite.T(), "test@user2", team.Members[1].OauthUser.Username)

	// Check the response body
	memberResponses := make([]*accounts.UserResponse, len(team.Members))
	for i, member := range team.Members {
		memberResponse, err := accounts.NewUserResponse(member)
		assert.NoError(suite.T(), err, "Creating response object failed")
		memberResponses[i] = memberResponse
	}
	expected := &TeamResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/teams/%d", team.ID),
				},
			},
			Embedded: map[string]jsonhal.Embedded{
				"members": jsonhal.Embedded(memberResponses),
			},
		},
		ID:        team.ID,
		Name:      "Test Team 1 Updated",
		CreatedAt: team.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: team.UpdatedAt.UTC().Format(time.RFC3339),
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
