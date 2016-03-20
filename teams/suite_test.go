package teams

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
	"github.com/stretchr/testify/suite"
)

var testDbPath = "/tmp/pinglist_teams_testdb.sqlite"

var testFixtures = []string{
	"../oauth/fixtures/test_clients.yml",
	"../oauth/fixtures/test_users.yml",
	"../oauth/fixtures/test_access_tokens.yml",
	"../accounts/fixtures/roles.yml",
	"../accounts/fixtures/test_accounts.yml",
	"../accounts/fixtures/test_users.yml",
	"fixtures/test_teams.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
	accounts.MigrateAll,
	MigrateAll,
}

// TeamsTestSuite needs to be exported so the tests run
type TeamsTestSuite struct {
	suite.Suite
	cnf                      *config.Config
	db                       *gorm.DB
	oauthServiceMock         *oauth.ServiceMock
	accountsServiceMock      *accounts.ServiceMock
	subscriptionsServiceMock *subscriptions.ServiceMock
	service                  *Service
	accounts                 []*accounts.Account
	users                    []*accounts.User
	teams                    []*Team
	router                   *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *TeamsTestSuite) SetupSuite() {

	// Initialise the config
	suite.cnf = config.NewConfig(false, false)

	// Create the test database
	db, err := database.CreateTestDatabase(testDbPath, testMigrations, testFixtures)
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

	// Fetch test teams
	suite.teams = make([]*Team, 0)
	err = suite.db.Preload("Owner.OauthUser").Preload("Members.OauthUser").
		Order("id").Find(&suite.teams).Error
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
	)

	// Register routes
	suite.router = mux.NewRouter()
	RegisterRoutes(suite.router, suite.service)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *TeamsTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *TeamsTestSuite) SetupTest() {
	suite.db.Exec("delete from team_team_members;")
	suite.db.Unscoped().Not("id", []int64{1}).Delete(new(Team))

	// Reset mocks
	suite.oauthServiceMock.ExpectedCalls = suite.oauthServiceMock.ExpectedCalls[:0]
	suite.oauthServiceMock.Calls = suite.oauthServiceMock.Calls[:0]
	suite.accountsServiceMock.ExpectedCalls = suite.accountsServiceMock.ExpectedCalls[:0]
	suite.accountsServiceMock.Calls = suite.accountsServiceMock.Calls[:0]
	suite.subscriptionsServiceMock.ExpectedCalls = suite.subscriptionsServiceMock.ExpectedCalls[:0]
	suite.subscriptionsServiceMock.Calls = suite.subscriptionsServiceMock.Calls[:0]
}

// The TearDownTest method will be run after every test in the suite.
func (suite *TeamsTestSuite) TearDownTest() {
	//
}

// TestTeamsTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTeamsTestSuite(t *testing.T) {
	suite.Run(t, new(TeamsTestSuite))
}

// Mock authentication
func (suite *TeamsTestSuite) mockAuthentication(user *accounts.User) {
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
func (suite *TeamsTestSuite) mockFindActiveSubscription(userID uint, subscription *subscriptions.Subscription, err error) {
	suite.subscriptionsServiceMock.On(
		"FindActiveUserSubscription",
		userID,
	).Return(subscription, err)
}

// Mock find user
func (suite *TeamsTestSuite) mockFindUser(userID uint, user *accounts.User, err error) {
	suite.accountsServiceMock.On(
		"FindUserByID",
		userID,
	).Return(user, err)
}
