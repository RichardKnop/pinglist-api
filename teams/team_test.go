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
	}
}

func (suite *TeamsTestSuite) TestFindTeamByOwnerID() {
	var (
		team *Team
		err  error
	)

	// Let's try to find a team by a bogus owner ID
	team, err = suite.service.FindTeamByOwnerID(12345)

	// Team should be nil
	assert.Nil(suite.T(), team)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrTeamNotFound, err)
	}

	// Now let's pass a valid owner ID
	team, err = suite.service.FindTeamByOwnerID(suite.users[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct team should be returned
	if assert.NotNil(suite.T(), team) {
		assert.Equal(suite.T(), suite.teams[0].ID, team.ID)
		assert.Equal(suite.T(), "test@superuser", team.Owner.OauthUser.Username)
	}
}
