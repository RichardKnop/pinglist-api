package notifications

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/jinzhu/gorm"
)

// Service struct keeps objects to avoid passing them around
type Service struct {
	cnf             *config.Config
	db              *gorm.DB
	accountsService accounts.ServiceInterface
	snsAdapter      SNSAdapterInterface
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, accountsService accounts.ServiceInterface, snsAdapter SNSAdapterInterface) *Service {
	if snsAdapter == nil {
		snsAdapter = NewSNSAdapter(cnf.AWS.Region)
	}
	return &Service{
		cnf:             cnf,
		db:              db,
		accountsService: accountsService,
		snsAdapter:      snsAdapter,
	}
}

// GetAccountsService returns accounts.Service instance
func (s *Service) GetAccountsService() accounts.ServiceInterface {
	return s.accountsService
}
