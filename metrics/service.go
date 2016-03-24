package metrics

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/jinzhu/gorm"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf             *config.Config
	db              *gorm.DB
	accountsService accounts.ServiceInterface
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, accountsService accounts.ServiceInterface) *Service {
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
