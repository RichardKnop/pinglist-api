package subscriptions

import (
	"errors"

	"github.com/RichardKnop/pinglist-api/accounts"
	stripe "github.com/stripe/stripe-go"
)

var (
	// ErrCardNotFound ...
	ErrCardNotFound = errors.New("Card not found")
)

// FindCardByID looks up a card by an ID and returns it
func (s *Service) FindCardByID(cardID uint) (*Card, error) {
	// Fetch the card from the database
	card := new(Card)
	notFound := s.db.Preload("Customer.User.OauthUser").
		First(card, cardID).RecordNotFound()

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
	notFound := s.db.Preload("Customer.User.OauthUser").
		Where("card_id = ?", cardID).
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
		// Unexpected server error
		if err != ErrCustomerNotFound {
			tx.Rollback() // rollback the transaction
			return nil, err
		}

		// Create a new Stripe customer
		stripeCustomer, err = s.stripeAdapter.CreateCustomer(
			user.OauthUser.Username,
			"",
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
			"",
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

	// Create a new card object
	card := NewCard(
		customer,
		stripeCard.ID,
		string(stripeCard.Brand),
		stripeCard.LastFour,
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
	// Delete the card
	stripeCard, err := s.stripeAdapter.DeleteCard(
		card.CardID,
		card.Customer.CustomerID,
	)
	if err != nil {
		return err
	}

	logger.Infof("Deleted card: %s", stripeCard.ID)

	// Delete the record from our database
	if err := s.db.Delete(card).Error; err != nil {
		return err
	}

	return nil
}
