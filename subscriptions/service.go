package subscriptions

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/jinzhu/gorm"
	"github.com/stripe/stripe-go"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf             *config.Config
	db              *gorm.DB
	accountsService accounts.ServiceInterface // accounts service dependency injection
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, accountsService accounts.ServiceInterface) *Service {
	// Assign secret key from configuration to Stripe
	stripe.Key = cnf.Stripe.SecretKey

	return &Service{
		cnf:             cnf,
		db:              db,
		accountsService: accountsService,
	}
}

// GetAccountsService returns accounts.Service instance
func (s *Service) GetAccountsService() accounts.ServiceInterface {
	return s.accountsService
}
