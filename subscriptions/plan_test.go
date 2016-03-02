package subscriptions

import (
	"github.com/stretchr/testify/assert"
)

func (suite *SubscriptionsTestSuite) TestFindPlanByID() {
	var (
		plan *Plan
		err  error
	)

	// When we try to find a plan with a bogus ID
	plan, err = suite.service.FindPlanByID(12345)

	// Plan object should be nil
	assert.Nil(suite.T(), plan)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errPlanNotFound, err)
	}

	// When we try to find a plan with a valid ID
	plan, err = suite.service.FindPlanByID(suite.plans[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct plan object should be returned
	if assert.NotNil(suite.T(), plan) {
		assert.Equal(suite.T(), suite.plans[0].ID, plan.ID)
	}
}
