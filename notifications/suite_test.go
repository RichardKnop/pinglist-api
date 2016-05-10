package notifications

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

var testDbPath = "/tmp/notifications_testdb.sqlite"

var testFixtures = []string{
	"../oauth/fixtures/test_clients.yml",
	"../oauth/fixtures/test_users.yml",
	"../accounts/fixtures/roles.yml",
	"../accounts/fixtures/test_accounts.yml",
	"../accounts/fixtures/test_users.yml",
	"fixtures/test_endpoints.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
	accounts.MigrateAll,
	MigrateAll,
}

// NotificationsTestSuite needs to be exported so the tests run
type NotificationsTestSuite struct {
	suite.Suite
	cnf                 *config.Config
	db                  *gorm.DB
	oauthServiceMock    *oauth.ServiceMock
	accountsServiceMock *accounts.ServiceMock
	snsAdapterMock      *SNSAdapterMock
	service             *Service
	accounts            []*accounts.Account
	users               []*accounts.User
	endpoints           []*Endpoint
	router              *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *NotificationsTestSuite) SetupSuite() {

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
	err = suite.db.Order("id").Find(&suite.accounts).Error
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
	suite.snsAdapterMock = new(SNSAdapterMock)

	// Initialise the service
	suite.service = NewService(
		suite.cnf,
		suite.db,
		suite.accountsServiceMock,
		suite.snsAdapterMock,
	)

	// Register routes
	suite.router = mux.NewRouter()
	RegisterRoutes(suite.router, suite.service)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *NotificationsTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *NotificationsTestSuite) SetupTest() {
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(Endpoint))

	// Fetch test endpoints
	suite.endpoints = make([]*Endpoint, 0)
	err := suite.db.Preload("User").Order("id").Find(&suite.endpoints).Error
	if err != nil {
		log.Fatal(err)
	}

	// Reset mocks
	suite.oauthServiceMock.ExpectedCalls = suite.oauthServiceMock.ExpectedCalls[:0]
	suite.oauthServiceMock.Calls = suite.oauthServiceMock.Calls[:0]
	suite.accountsServiceMock.ExpectedCalls = suite.accountsServiceMock.ExpectedCalls[:0]
	suite.accountsServiceMock.Calls = suite.accountsServiceMock.Calls[:0]
	suite.snsAdapterMock.ExpectedCalls = suite.snsAdapterMock.ExpectedCalls[:0]
	suite.snsAdapterMock.Calls = suite.snsAdapterMock.Calls[:0]
}

// The TearDownTest method will be run after every test in the suite.
func (suite *NotificationsTestSuite) TearDownTest() {
	//
}

// TestNotificationsTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNotificationsTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationsTestSuite))
}

// Checks that the mock object expectations were met
func (suite *NotificationsTestSuite) assertMockExpectations() {
	suite.oauthServiceMock.AssertExpectations(suite.T())
	suite.accountsServiceMock.AssertExpectations(suite.T())
	suite.snsAdapterMock.AssertExpectations(suite.T())
}

// Mock authentication
func (suite *NotificationsTestSuite) mockUserAuth(user *accounts.User) {
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

// Mock create endpoint
func (suite *NotificationsTestSuite) mockCreateEndpoint(applicationARN, deviceToken, endpointARN string, err error) {
	suite.snsAdapterMock.On(
		"CreateEndpoint",
		applicationARN,
		deviceToken,
	).Return(endpointARN, err)
}

// Mock get endpoint attributes
func (suite *NotificationsTestSuite) mockGetAttributes(endpointARN string, endpointAttributes *EndpointAttributes, err error) {
	suite.snsAdapterMock.On("GetEndpointAttributes", endpointARN).
		Return(endpointAttributes, err)
}

// Mock set endpoint attributes
func (suite *NotificationsTestSuite) mockSetAttributes(endpointARN string, endpointAttributes *EndpointAttributes, err error) {
	suite.snsAdapterMock.On("SetEndpointAttributes", endpointARN, endpointAttributes).
		Return(err)
}
