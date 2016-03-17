package alarms

import (
	"log"
	"testing"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/database"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	testDbUser = "pinglist"
	testDbName = "pinglist_alarms_test"
)

var testFixtures = []string{
	"../oauth/fixtures/test_clients.yml",
	"../oauth/fixtures/test_users.yml",
	"../accounts/fixtures/roles.yml",
	"../accounts/fixtures/test_accounts.yml",
	"../accounts/fixtures/test_users.yml",
	"../subscriptions/fixtures/plans.yml",
	"../subscriptions/fixtures/test_customers.yml",
	"../subscriptions/fixtures/test_subscriptions.yml",
	"fixtures/regions.yml",
	"fixtures/alarm_states.yml",
	"fixtures/incident_types.yml",
	"fixtures/test_alarms.yml",
	"fixtures/test_incidents.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
	accounts.MigrateAll,
	subscriptions.MigrateAll,
	MigrateAll,
}

// AlarmsTestSuite needs to be exported so the tests run
type AlarmsTestSuite struct {
	suite.Suite
	cnf                      *config.Config
	db                       *gorm.DB
	oauthServiceMock         *oauth.ServiceMock
	accountsServiceMock      *accounts.ServiceMock
	subscriptionsServiceMock *subscriptions.ServiceMock
	service                  *Service
	accounts                 []*accounts.Account
	users                    []*accounts.User
	alarms                   []*Alarm
	incidents                []*Incident
	router                   *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *AlarmsTestSuite) SetupSuite() {
	// NOTE: using Postgres test database instead of sqlite here as
	// we need to test a Postgres specific functionality (table inheritance)

	// Initialise the config
	suite.cnf = config.NewConfig(false, false)

	// Create the test database
	db, err := database.CreateTestDatabasePostgres(
		testDbUser,
		testDbName,
		testMigrations,
		testFixtures,
	)
	if err != nil {
		log.Fatal(err)
	}
	suite.db = db

	// Fetch test accounts
	suite.accounts = make([]*accounts.Account, 0)
	err = suite.db.Preload("OauthClient").Order("id").Find(&suite.accounts).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*accounts.User, 0)
	err = suite.db.Preload("Account").Preload("OauthUser").Preload("Role").
		Order("id").Find(&suite.users).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test alarms
	suite.alarms = make([]*Alarm, 0)
	err = suite.db.Preload("User").Preload("Incidents").
		Order("id").Find(&suite.alarms).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test incidents
	suite.incidents = make([]*Incident, 0)
	err = suite.db.Preload("Alarm").Order("id").
		Order("id").Find(&suite.incidents).Error
	if err != nil {
		log.Fatal(err)
	}

	// Initialise mocks
	suite.oauthServiceMock = new(oauth.ServiceMock)
	suite.accountsServiceMock = new(accounts.ServiceMock)
	suite.subscriptionsServiceMock = new(subscriptions.ServiceMock)

	// Initialise the service
	suite.service = NewService(
		suite.cnf,
		suite.db,
		suite.accountsServiceMock,
		suite.subscriptionsServiceMock,
		nil, // HTTP client
	)

	// Register routes
	suite.router = mux.NewRouter()
	RegisterRoutes(suite.router, suite.service)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *AlarmsTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *AlarmsTestSuite) SetupTest() {
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4}).Delete(new(Alarm))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4}).Delete(new(Incident))
	suite.db.Unscoped().Delete(new(Result))

	// Delete result sub tables
	var resultSubTables []*ResultSubTable
	if err := suite.db.Order("id").Find(&resultSubTables).Error; err != nil {
		log.Fatal(err)
	}
	for _, resultSubTable := range resultSubTables {
		if err := suite.db.DropTable(resultSubTable.Name).Error; err != nil {
			log.Fatal(err)
		}
	}
	suite.db.Unscoped().Delete(new(ResultSubTable))

	// Reset mocks
	suite.oauthServiceMock.ExpectedCalls = suite.oauthServiceMock.ExpectedCalls[:0]
	suite.oauthServiceMock.Calls = suite.oauthServiceMock.Calls[:0]
	suite.accountsServiceMock.ExpectedCalls = suite.accountsServiceMock.ExpectedCalls[:0]
	suite.accountsServiceMock.Calls = suite.accountsServiceMock.Calls[:0]
	suite.subscriptionsServiceMock.ExpectedCalls = suite.subscriptionsServiceMock.ExpectedCalls[:0]
	suite.subscriptionsServiceMock.Calls = suite.subscriptionsServiceMock.Calls[:0]
}

// The TearDownTest method will be run after every test in the suite.
func (suite *AlarmsTestSuite) TearDownTest() {
	//
}

// TestAlarmsTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAlarmsTestSuite(t *testing.T) {
	suite.Run(t, new(AlarmsTestSuite))
}

// Mock authentication
func (suite *AlarmsTestSuite) mockAuthentication(user *accounts.User) {
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

// Mock find active subscription
func (suite *AlarmsTestSuite) mockFindActiveSubscription(userID uint, subscription *subscriptions.Subscription, err error) {
	suite.subscriptionsServiceMock.On(
		"FindActiveUserSubscription",
		userID,
	).Return(subscription, err)
}

// Mock user querystring filtering
func (suite *AlarmsTestSuite) mockUserFiltering(user *accounts.User) {
	suite.accountsServiceMock.On(
		"GetUserFromQueryString",
		mock.AnythingOfType("*http.Request"),
	).Return(user, nil)
}
