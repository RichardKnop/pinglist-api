package subscriptions

import (
	"github.com/stretchr/testify/assert"
)

func (suite *SubscriptionsTestSuite) TestFindCardByID() {
	var (
		card *Card
		err  error
	)

	// When we try to find a card with a bogus ID
	card, err = suite.service.FindCardByID(12345)

	// Card object should be nil
	assert.Nil(suite.T(), card)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrCardNotFound, err)
	}

	// When we try to find a card with a valid ID
	card, err = suite.service.FindCardByID(suite.cards[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct card object should be returned
	if assert.NotNil(suite.T(), card) {
		assert.Equal(suite.T(), suite.cards[0].ID, card.ID)
		assert.Equal(suite.T(), suite.users[0].ID, card.Customer.User.ID)
	}
}

func (suite *SubscriptionsTestSuite) TestFindCardByCardID() {
	var (
		card *Card
		err  error
	)

	// When we try to find a card with a bogus card ID
	card, err = suite.service.FindCardByCardID("bogus")

	// Card object should be nil
	assert.Nil(suite.T(), card)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrCardNotFound, err)
	}

	// When we try to find a card with a valid card ID
	card, err = suite.service.FindCardByCardID(suite.cards[0].CardID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct card object should be returned
	if assert.NotNil(suite.T(), card) {
		assert.Equal(suite.T(), suite.cards[0].ID, card.ID)
		assert.Equal(suite.T(), suite.users[0].ID, card.Customer.User.ID)
	}
}
