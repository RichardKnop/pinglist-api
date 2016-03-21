package subscriptions

import (
	"errors"

	"github.com/RichardKnop/pinglist-api/util"
)

var (
	// ErrCustomerNotFound ...
	ErrCustomerNotFound = errors.New("Customer not found")
)

// FindCustomerByID looks up a customer by an ID and returns it
func (s *Service) FindCustomerByID(customerID uint) (*Customer, error) {
	// Fetch the subscription from the database
	customer := new(Customer)
	notFound := s.db.Preload("User").First(customer, customerID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrCustomerNotFound
	}

	return customer, nil
}

// FindCustomerByUserID looks up a customer by a user ID and returns it
func (s *Service) FindCustomerByUserID(userID uint) (*Customer, error) {
	// Fetch the subscription from the database
	customer := new(Customer)
	notFound := s.db.Where(Customer{
		UserID: util.PositiveIntOrNull(int64(userID)),
	}).Preload("User").First(customer).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrCustomerNotFound
	}

	return customer, nil
}

// FindCustomerByCustomerID looks up a customer by a customer ID and returns it
func (s *Service) FindCustomerByCustomerID(customerID string) (*Customer, error) {
	// Fetch the subscription from the database
	customer := new(Customer)
	notFound := s.db.Preload("User").Where("customer_id = ?", customerID).
		First(customer).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrCustomerNotFound
	}

	return customer, nil
}
