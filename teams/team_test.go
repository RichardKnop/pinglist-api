package teams

import (
	"github.com/stretchr/testify/assert"
)

func (suite *TeamsTestSuite) TestFindTeamByID() {
	var (
		team *Team
		err  error
	)

	// Let's try to find a team by a bogus ID
	team, err = suite.service.FindTeamByID(12345)

	// Team should be nil
	assert.Nil(suite.T(), team)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrTeamNotFound, err)
	}

	// Now let's pass a valid ID
	team, err = suite.service.FindTeamByID(suite.teams[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct team should be returned
	if assert.NotNil(suite.T(), team) {
		assert.Equal(suite.T(), suite.teams[0].ID, team.ID)
		assert.Equal(suite.T(), "test@superuser", team.Owner.OauthUser.Username)
		assert.Equal(suite.T(), 0, len(team.Members))
	}
}

func (suite *TeamsTestSuite) TestFindTeamByMemberID() {
	var (
		team *Team
		err  error
	)

	// Insert a test team member
	err = suite.db.Model(&suite.teams[0]).Association("Members").Append(suite.users[1]).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Let's try to find a team by a bogus member ID
	team, err = suite.service.FindTeamByMemberID(12345)

	// Team should be nil
	assert.Nil(suite.T(), team)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrTeamNotFound, err)
	}

	// Now let's pass a valid member ID
	team, err = suite.service.FindTeamByMemberID(suite.users[1].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct team should be returned
	if assert.NotNil(suite.T(), team) {
		assert.Equal(suite.T(), suite.teams[0].ID, team.ID)
		assert.Equal(suite.T(), "test@superuser", team.Owner.OauthUser.Username)
		assert.Equal(suite.T(), 1, len(team.Members))
		assert.Equal(suite.T(), "test@user", team.Members[0].OauthUser.Username)
	}
}

func (suite *TeamsTestSuite) TestPaginatedTeamsCount() {
	var (
		count int
		err   error
	)

	count, err = suite.service.paginatedTeamsCount(nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	count, err = suite.service.paginatedTeamsCount(suite.users[0])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	count, err = suite.service.paginatedTeamsCount(suite.users[1])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}
}

func (suite *TeamsTestSuite) TestFindPaginatedTeams() {
	var (
		teams []*Team
		err   error
	)

	// This should return all teams
	teams, err = suite.service.findPaginatedTeams(0, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(teams))
		assert.Equal(suite.T(), suite.teams[0].ID, teams[0].ID)
		assert.Equal(suite.T(), suite.teams[1].ID, teams[1].ID)
		assert.Equal(suite.T(), suite.teams[2].ID, teams[2].ID)
		assert.Equal(suite.T(), suite.teams[3].ID, teams[3].ID)
	}

	// This should return all teams ordered by ID desc
	teams, err = suite.service.findPaginatedTeams(0, 25, "id desc", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(teams))
		assert.Equal(suite.T(), suite.teams[3].ID, teams[0].ID)
		assert.Equal(suite.T(), suite.teams[2].ID, teams[1].ID)
		assert.Equal(suite.T(), suite.teams[1].ID, teams[2].ID)
		assert.Equal(suite.T(), suite.teams[0].ID, teams[3].ID)
	}

	// Test offset
	teams, err = suite.service.findPaginatedTeams(2, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(teams))
		assert.Equal(suite.T(), suite.teams[2].ID, teams[0].ID)
		assert.Equal(suite.T(), suite.teams[3].ID, teams[1].ID)
	}

	// Test limit
	teams, err = suite.service.findPaginatedTeams(2, 1, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(teams))
		assert.Equal(suite.T(), suite.teams[2].ID, teams[0].ID)
	}
}
