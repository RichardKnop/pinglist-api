package subscriptions

import (
	"errors"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/jinzhu/gorm"
	stripe "github.com/stripe/stripe-go"
)

var (
	// ErrCardNotFound ...
	ErrCardNotFound = errors.New("Card not found")
	// ErrCardCanOnlyBeDeletedFromCancelledSubsription ...
	ErrCardCanOnlyBeDeletedFromCancelledSubsription = errors.New("Card can only be deleted from cancelled subscription")
)

// FindCardByID looks up a card by an ID and returns it
func (s *Service) FindCardByID(cardID uint) (*Card, error) {
	// Fetch the card from the database
	card := new(Card)
	notFound := s.db.Preload("Customer.User").First(card, cardID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrCardNotFound
	}

	return card, nil
}

// FindCardByCardID looks up a card by a card ID and returns it
func (s *Service) FindCardByCardID(cardID string) (*Card, error) {
	// Fetch the card from the database
	card := new(Card)
	notFound := s.db.Preload("Customer.User").Where("card_id = ?", cardID).
		First(card).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrCardNotFound
	}

	return card, nil
}

// createCard creates a new Stripe card payment source
func (s *Service) createCard(user *accounts.User, cardRequest *CardRequest) (*Card, error) {
	var (
		customer       *Customer
		stripeCustomer *stripe.Customer
		created        bool
		err            error
	)

	// Do we already store a customer recors for this user?
	customer, err = s.FindCustomerByUserID(user.ID)

	// Begin a transaction
	tx := s.db.Begin()

	if err != nil {
		// Create a new Stripe customer
		stripeCustomer, err = s.stripeAdapter.CreateCustomer(
			user.OauthUser.Username,
			"", // token
		)
		if err != nil {
			tx.Rollback() // rollback the transaction
			return nil, err
		}

		logger.Infof("Created customer: %s", stripeCustomer.ID)

		// Create a new customer object
		customer = NewCustomer(user, stripeCustomer.ID)

		// Save the customer to the database
		if err := tx.Create(customer).Error; err != nil {
			tx.Rollback() // rollback the transaction
			return nil, err
		}
	} else {
		// Get an existing Stripe customer or create a new one
		stripeCustomer, created, err = s.stripeAdapter.GetOrCreateCustomer(
			customer.CustomerID,
			user.OauthUser.Username,
			"", // token
		)
		if err != nil {
			tx.Rollback() // rollback the transaction
			return nil, err
		}

		if created {
			logger.Infof("Created customer: %s", stripeCustomer.ID)

			// Our customer record is not valid so delete it
			if err := tx.Delete(customer).Error; err != nil {
				tx.Rollback() // rollback the transaction
				return nil, err
			}

			// Create a new customer object
			customer = NewCustomer(user, stripeCustomer.ID)

			// Save the customer to the database
			if err := tx.Create(customer).Error; err != nil {
				tx.Rollback() // rollback the transaction
				return nil, err
			}
		}
	}

	// Create a new Stripe card
	stripeCard, err := s.stripeAdapter.CreateCard(
		customer.CustomerID,
		cardRequest.Token,
	)
	if err != nil {
		return nil, err
	}

	logger.Infof("Created card: %s", stripeCard.ID)

	var lastFour string
	if stripeCard.DynLastFour != "" {
		lastFour = stripeCard.DynLastFour
	} else {
		lastFour = stripeCard.LastFour
	}

	// Create a new card object
	card := NewCard(
		customer,
		stripeCard.ID,
		string(stripeCard.Brand),
		lastFour,
		uint(stripeCard.Month),
		uint(stripeCard.Year),
	)

	// Save the card to the database
	if err := tx.Create(card).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return card, nil
}

// deleteCard deletes a card payment source
func (s *Service) deleteCard(card *Card) error {
	// Fetch subscription using this card
	subscription, err := s.FindSubscriptionByCardID(card.ID)

	// Begin a transaction
	tx := s.db.Begin()

	// There is a subscription using this card
	if err == nil && subscription != nil {
		// Subscription must be cancelled first before the card can be removed
		if !subscription.IsCancelled() {
			tx.Rollback() // rollback the transaction
			return ErrCardCanOnlyBeDeletedFromCancelledSubsription
		}

		// Update the subscription (need to use map here because card_id field is
		// changing to nil which would not work with struct)
		if err := tx.Model(subscription).UpdateColumns(map[string]interface{}{
			"card_id":    nil,
			"updated_at": time.Now(),
		}).Error; err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}
	}

	// Delete the card
	stripeCard, err := s.stripeAdapter.DeleteCard(
		card.CardID,
		card.Customer.CustomerID,
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	logger.Infof("Deleted card: %s", stripeCard.ID)

	// Delete the record from our database
	if err := tx.Delete(card).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}

// paginatedCardsCount returns a total count of cards
// Can be optionally filtered by user
func (s *Service) paginatedCardsCount(user *accounts.User) (int, error) {
	var count int
	if err := s.paginatedCardsQuery(user).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedCards returns paginated card records
// Results can optionally be filtered by user
func (s *Service) findPaginatedCards(offset, limit int, orderBy string, user *accounts.User) ([]*Card, error) {
	var cards []*Card

	// Get the pagination query
	cardsQuery := s.paginatedCardsQuery(user)

	// Default ordering
	if orderBy == "" {
		orderBy = "id"
	}

	// Retrieve paginated results from the database
	err := cardsQuery.Offset(offset).Limit(limit).Order(orderBy).
		Preload("Customer.User").Find(&cards).Error
	if err != nil {
		return cards, err
	}

	return cards, nil
}

// paginatedCardsQuery returns a db query for paginated cards
func (s *Service) paginatedCardsQuery(user *accounts.User) *gorm.DB {
	// Basic query
	cardsQuery := s.db.Model(new(Card))

	// Optionally filter by user
	if user != nil {
		cardsQuery = cardsQuery.
			Joins("inner join subscription_customers on subscription_customers.id = subscription_cards.customer_id").
			Joins("inner join account_users on account_users.id = subscription_customers.user_id").
			Where("subscription_customers.user_id = ?", user.ID)
	}

	return cardsQuery
}
