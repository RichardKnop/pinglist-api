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

	"github.com/RichardKnop/pinglist-api/util"
	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *TeamsTestSuite) TestInviteUserRequiresUserAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	w := httptest.NewRecorder()

	suite.service.inviteUserHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *TeamsTestSuite) TestInviteUserWithoutPermission() {
	// Prepare a request
	invitationRequest := &accounts.InvitationRequest{
		Email: "test@user",
	}
	payload, err := json.Marshal(invitationRequest)
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://1.2.3.4/v1/teams/%d/invitations", suite.teams[0].ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "invite_user", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[2])

	// Count before
	var countBefore int
	suite.db.Table("team_team_members").Count(&countBefore)

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
	suite.db.Table("team_team_members").Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrInviteUserPermission.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *TeamsTestSuite) TestInviteUser() {
	// Prepare a request
	invitationRequest := &InvitationRequest{"test@user"}
	payload, err := json.Marshal(invitationRequest)
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://1.2.3.4/v1/teams/%d/invitations", suite.teams[0].ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "invite_user", match.Route.GetName())
	}

	// Mock authentication
	suite.mockUserAuth(suite.users[0])

	// Mock invite user call to accounts service
	invitation := &accounts.Invitation{
		Model: gorm.Model{
			ID:        123,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Reference:     "invitation_reference",
		InvitedByUser: suite.users[0],
		InvitedUser:   suite.users[1],
	}
	suite.mockInviteUserTx(
		suite.users[0],
		"test@user",
		invitation,
		nil,
	)

	// Count before
	var countBefore int
	suite.db.Table("team_team_members").Count(&countBefore)

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
	var countAfter int
	suite.db.Table("team_team_members").Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the updated team
	team := new(Team)
	notFound := suite.db.Preload("Owner.OauthUser").Preload("Members.OauthUser").
		First(team, suite.teams[0].ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), 1, len(team.Members))
	assert.Equal(suite.T(), suite.users[1].ID, team.Members[0].ID)

	// Check the response body
	expected := &accounts.InvitationResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/accounts/invitations/%d", invitation.ID),
				},
			},
		},
		ID:              invitation.ID,
		Reference:       invitation.Reference,
		InvitedUserID:   invitation.InvitedUser.ID,
		InvitedByUserID: invitation.InvitedByUser.ID,
		CreatedAt:       util.FormatTime(invitation.CreatedAt),
		UpdatedAt:       util.FormatTime(invitation.UpdatedAt),
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
