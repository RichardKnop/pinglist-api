package alarms

import (
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/jinzhu/gorm"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf                  *config.Config
	db                   *gorm.DB
	accountsService      accounts.ServiceInterface      // accounts service dependency injection
	subscriptionsService subscriptions.ServiceInterface // accounts service dependency injection
	client               *http.Client                   // clients are safe for concurrent use by multiple goroutines
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, accountsService accounts.ServiceInterface, subscriptionsService subscriptions.ServiceInterface, client *http.Client) *Service {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second, // 10 seconds timeout
		}
	}
	return &Service{
		cnf:                  cnf,
		db:                   db,
		accountsService:      accountsService,
		subscriptionsService: subscriptionsService,
		client:               client,
	}
}

// GetAccountsService returns accounts.Service instance
func (s *Service) GetAccountsService() accounts.ServiceInterface {
	return s.accountsService
}
