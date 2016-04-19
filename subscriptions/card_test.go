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

func (suite *SubscriptionsTestSuite) TestPaginatedCardsCount() {
	var (
		count int
		err   error
	)

	// Without any filtering
	count, err = suite.service.cardsCount(nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by user with 4 cards
	count, err = suite.service.cardsCount(suite.users[0])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by user without cards
	count, err = suite.service.cardsCount(suite.users[1])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}
}

func (suite *SubscriptionsTestSuite) TestFindPaginatedCards() {
	var (
		cards []*Card
		err   error
	)

	// This should return all cards
	cards, err = suite.service.findPaginatedCards(0, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(cards))
		assert.Equal(suite.T(), suite.cards[0].ID, cards[0].ID)
		assert.Equal(suite.T(), suite.cards[1].ID, cards[1].ID)
		assert.Equal(suite.T(), suite.cards[2].ID, cards[2].ID)
		assert.Equal(suite.T(), suite.cards[3].ID, cards[3].ID)
	}

	// This should return all cards ordered by ID desc
	cards, err = suite.service.findPaginatedCards(0, 25, "id desc", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(cards))
		assert.Equal(suite.T(), suite.cards[3].ID, cards[0].ID)
		assert.Equal(suite.T(), suite.cards[2].ID, cards[1].ID)
		assert.Equal(suite.T(), suite.cards[1].ID, cards[2].ID)
		assert.Equal(suite.T(), suite.cards[0].ID, cards[3].ID)
	}

	// Test offset
	cards, err = suite.service.findPaginatedCards(2, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(cards))
		assert.Equal(suite.T(), suite.cards[2].ID, cards[0].ID)
		assert.Equal(suite.T(), suite.cards[3].ID, cards[1].ID)
	}

	// Test limit
	cards, err = suite.service.findPaginatedCards(2, 1, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(cards))
		assert.Equal(suite.T(), suite.cards[2].ID, cards[0].ID)
	}
}
