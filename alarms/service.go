package alarms

import (
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/email"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/jinzhu/gorm"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf                  *config.Config
	db                   *gorm.DB
	accountsService      accounts.ServiceInterface
	subscriptionsService subscriptions.ServiceInterface
	teamsService         teams.ServiceInterface
	metricsService       metrics.ServiceInterface
	emailService         email.ServiceInterface
	emailFactory         EmailFactoryInterface
	client               *http.Client
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, accountsService accounts.ServiceInterface, subscriptionsService subscriptions.ServiceInterface, teamsService teams.ServiceInterface, metricsService metrics.ServiceInterface, emailService email.ServiceInterface, emailFactory EmailFactoryInterface, client *http.Client) *Service {
	if emailFactory == nil {
		emailFactory = NewEmailFactory(cnf)
	}
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
		teamsService:         teamsService,
		metricsService:       metricsService,
		emailService:         emailService,
		emailFactory:         emailFactory,
		client:               client,
	}
}

// GetAccountsService returns accounts.Service instance
func (s *Service) GetAccountsService() accounts.ServiceInterface {
	return s.accountsService
}
