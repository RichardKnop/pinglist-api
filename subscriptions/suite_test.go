package subscriptions

import (
	"log"
	"testing"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/database"
	"github.com/RichardKnop/pinglist-api/logger"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	stripeCustomer "github.com/stripe/stripe-go/customer"
	stripeUtils "github.com/stripe/stripe-go/utils"
)

var testDbPath = "/tmp/subscriptions_testdb.sqlite"

var testFixtures = []string{
	"../oauth/fixtures/test_clients.yml",
	"../oauth/fixtures/test_users.yml",
	"../accounts/fixtures/roles.yml",
	"../accounts/fixtures/test_accounts.yml",
	"../accounts/fixtures/test_users.yml",
	"fixtures/plans.yml",
	"fixtures/test_customers.yml",
	"fixtures/test_cards.yml",
	"fixtures/test_subscriptions.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
	accounts.MigrateAll,
	MigrateAll,
}

// SubscriptionsTestSuite needs to be exported so the tests run
type SubscriptionsTestSuite struct {
	suite.Suite
	cnf                 *config.Config
	db                  *gorm.DB
	oauthServiceMock    *oauth.ServiceMock
	accountsServiceMock *accounts.ServiceMock
	service             *Service
	accounts            []*accounts.Account
	users               []*accounts.User
	plans               []*Plan
	customers           []*Customer
	cards               []*Card
	subscriptions       []*Subscription
	router              *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *SubscriptionsTestSuite) SetupSuite() {

	// Initialise the config
	suite.cnf = config.NewConfig(false, false)

	// Overwrite Stripe secret key
	suite.cnf.Stripe.SecretKey = stripeUtils.GetTestKey()

	// Create the test database
	db, err := database.CreateTestDatabase(testDbPath, testMigrations, testFixtures)
	if err != nil {
		log.Fatal(err)
	}
	suite.db = db

	// Fetch test accounts
	suite.accounts = make([]*accounts.Account, 0)
	if suite.db.Preload("OauthClient").Order("id").Find(&suite.accounts).Error != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*accounts.User, 0)
	err = suite.db.Preload("Account").Preload("OauthUser").Preload("Role").
		Order("id").Find(&suite.users).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test plans
	suite.plans = make([]*Plan, 0)
	if suite.db.Order("id").Find(&suite.plans).Error != nil {
		log.Fatal(err)
	}

	// Fetch test customers
	suite.customers = make([]*Customer, 0)
	err = suite.db.Preload("User.OauthUser").Preload("User.Role").
		Order("id").Find(&suite.customers).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test cards
	suite.cards = make([]*Card, 0)
	err = suite.db.Preload("Customer.User").Order("id").Find(&suite.cards).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test subscriptions
	suite.subscriptions = make([]*Subscription, 0)
	err = suite.db.Preload("Customer.User").Preload("Plan").Order("id").
		Find(&suite.subscriptions).Error
	if err != nil {
		log.Fatal(err)
	}

	// Initialise mocks
	suite.oauthServiceMock = new(oauth.ServiceMock)
	suite.accountsServiceMock = new(accounts.ServiceMock)

	// Initialise the service
	suite.service = NewService(
		suite.cnf,
		suite.db,
		suite.accountsServiceMock,
		nil, // Stripe adapter... TODO replace with mock
	)

	// Register routes
	suite.router = mux.NewRouter()
	RegisterRoutes(suite.router, suite.service)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *SubscriptionsTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *SubscriptionsTestSuite) SetupTest() {
	// Delete Stripe customers
	var customersToDelete []*Customer
	if err := suite.db.Not("id", []int64{1}).Find(&customersToDelete).Error; err != nil {
		log.Fatal(err)
	}
	for _, customer := range customersToDelete {
		_, err := stripeCustomer.Del(customer.CustomerID)
		logger.ERROR.Print(err)
	}

	// Delete data inserted by tests
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4}).Delete(new(Subscription))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4}).Delete(new(Card))
	suite.db.Unscoped().Not("id", []int64{1}).Delete(new(Customer))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4, 5}).Delete(new(Plan))

	// Reset mocks
	suite.resetMocks()
}

// The TearDownTest method will be run after every test in the suite.
func (suite *SubscriptionsTestSuite) TearDownTest() {
	//
}

// TestSubscriptionsTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSubscriptionsTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionsTestSuite))
}

// Reset mocks
func (suite *SubscriptionsTestSuite) resetMocks() {
	suite.oauthServiceMock.ExpectedCalls = suite.oauthServiceMock.ExpectedCalls[:0]
	suite.oauthServiceMock.Calls = suite.oauthServiceMock.Calls[:0]
	suite.accountsServiceMock.ExpectedCalls = suite.accountsServiceMock.ExpectedCalls[:0]
	suite.accountsServiceMock.Calls = suite.accountsServiceMock.Calls[:0]
}

// Checks that the mock object expectations were met
func (suite *SubscriptionsTestSuite) assertMockExpectations() {
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.resetMocks()
}

// Mock resource owner credentials grant auth
func (suite *SubscriptionsTestSuite) mockUserAuth(user *accounts.User) {
	// Mock GetConfig call to return the config object
	suite.accountsServiceMock.On("GetConfig").Return(suite.cnf)

	// Mock GetOauthService to return a mock oauth service
	suite.accountsServiceMock.On("GetOauthService").Return(suite.oauthServiceMock)

	// Mock Authenticate to return a mock access token
	mockOauthAccessToken := &oauth.AccessToken{User: user.OauthUser}
	suite.oauthServiceMock.On("Authenticate", "test_token").
		Return(mockOauthAccessToken, nil)

	// Mock FindUserByOauthUserID to return the wanted user
	suite.accountsServiceMock.On("FindUserByOauthUserID", user.OauthUser.ID).
		Return(user, nil)
}

// Mock client credentials grant auth
func (suite *SubscriptionsTestSuite) mockClientAuth(account *accounts.Account) {
	// Mock GetConfig call to return the config object
	suite.accountsServiceMock.On("GetConfig").Return(suite.cnf)

	// Mock GetOauthService to return a mock oauth service
	suite.accountsServiceMock.On("GetOauthService").Return(suite.oauthServiceMock)

	// Mock AuthClient to return a mock client
	suite.oauthServiceMock.On("AuthClient", "test_client_1", "test_secret").
		Return(account.OauthClient, nil)

	// Mock FindAccountByOauthClientID to return the wanted account
	suite.accountsServiceMock.
		On("FindAccountByOauthClientID", account.OauthClient.ID).
		Return(account, nil)
}

// Mock user querystring filtering
func (suite *SubscriptionsTestSuite) mockUserFiltering(user *accounts.User) {
	suite.accountsServiceMock.On(
		"GetUserFromQueryString",
		mock.AnythingOfType("*http.Request"),
	).Return(user, nil)
}
