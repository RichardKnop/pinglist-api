package metrics

import (
	"log"
	"testing"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/database"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

var (
	testDbUser = "pinglist"
	testDbName = "pinglist_metrics_test"
)

var testFixtures = []string{
	"../oauth/fixtures/test_clients.yml",
	"../oauth/fixtures/test_users.yml",
	"../accounts/fixtures/roles.yml",
	"../accounts/fixtures/test_accounts.yml",
	"../accounts/fixtures/test_users.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
	accounts.MigrateAll,
	MigrateAll,
}

// AlarmsTestSuite needs to be exported so the tests run
type MetricsTestSuite struct {
	suite.Suite
	cnf                 *config.Config
	db                  *gorm.DB
	oauthServiceMock    *oauth.ServiceMock
	accountsServiceMock *accounts.ServiceMock
	service             *Service
	accounts            []*accounts.Account
	users               []*accounts.User
	router              *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *MetricsTestSuite) SetupSuite() {
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
	suite.db.LogMode(true)

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

	// Initialise mocks
	suite.oauthServiceMock = new(oauth.ServiceMock)
	suite.accountsServiceMock = new(accounts.ServiceMock)

	// Initialise the service
	suite.service = NewService(
		suite.cnf,
		suite.db,
		suite.accountsServiceMock,
	)

	// Register routes
	suite.router = mux.NewRouter()
	RegisterRoutes(suite.router, suite.service)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *MetricsTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *MetricsTestSuite) SetupTest() {
	suite.db.Unscoped().Delete(new(RequestTime))

	// Delete sub tables
	var subTables []*SubTable
	if err := suite.db.Order("id").Find(&subTables).Error; err != nil {
		log.Fatal(err)
	}
	for _, subTable := range subTables {
		if err := suite.db.DropTable(subTable.Name).Error; err != nil {
			log.Fatal(err)
		}
	}
	suite.db.Unscoped().Delete(new(SubTable))

	// Reset mocks
	suite.oauthServiceMock.ExpectedCalls = suite.oauthServiceMock.ExpectedCalls[:0]
	suite.oauthServiceMock.Calls = suite.oauthServiceMock.Calls[:0]
	suite.accountsServiceMock.ExpectedCalls = suite.accountsServiceMock.ExpectedCalls[:0]
	suite.accountsServiceMock.Calls = suite.accountsServiceMock.Calls[:0]
}

// The TearDownTest method will be run after every test in the suite.
func (suite *MetricsTestSuite) TearDownTest() {
	//
}

// TestMetricsTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsTestSuite))
}

// Mock authentication
func (suite *MetricsTestSuite) mockAuthentication(user *accounts.User) {
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
