package alarms

import (
	"log"
	"testing"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/database"
	"github.com/RichardKnop/pinglist-api/email"
	"github.com/RichardKnop/pinglist-api/metrics"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
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
	teamsServiceMock         *teams.ServiceMock
	metricsServiceMock       *metrics.ServiceMock
	emailServiceMock         *email.ServiceMock
	emailFactoryMock         *EmailFactoryMock
	service                  *Service
	accounts                 []*accounts.Account
	users                    []*accounts.User
	regions                  []*Region
	alarms                   []*Alarm
	incidents                []*Incident
	router                   *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *AlarmsTestSuite) SetupSuite() {
	// NOTE: using Postgres test database instead of sqlite here as
	// we need to test a Postgres specific functionality (interval '1 second')

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

	// Fetch test regions
	suite.regions = make([]*Region, 0)
	if suite.db.Order("id").Find(&suite.regions).Error != nil {
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
	suite.teamsServiceMock = new(teams.ServiceMock)
	suite.metricsServiceMock = new(metrics.ServiceMock)
	suite.emailServiceMock = new(email.ServiceMock)
	suite.emailFactoryMock = new(EmailFactoryMock)

	// Initialise the service
	suite.service = NewService(
		suite.cnf,
		suite.db,
		suite.accountsServiceMock,
		suite.subscriptionsServiceMock,
		suite.teamsServiceMock,
		suite.metricsServiceMock,
		suite.emailServiceMock,
		suite.emailFactoryMock,
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

	// Reset mocks
	suite.oauthServiceMock.ExpectedCalls = suite.oauthServiceMock.ExpectedCalls[:0]
	suite.oauthServiceMock.Calls = suite.oauthServiceMock.Calls[:0]
	suite.accountsServiceMock.ExpectedCalls = suite.accountsServiceMock.ExpectedCalls[:0]
	suite.accountsServiceMock.Calls = suite.accountsServiceMock.Calls[:0]
	suite.subscriptionsServiceMock.ExpectedCalls = suite.subscriptionsServiceMock.ExpectedCalls[:0]
	suite.subscriptionsServiceMock.Calls = suite.subscriptionsServiceMock.Calls[:0]
	suite.teamsServiceMock.ExpectedCalls = suite.teamsServiceMock.ExpectedCalls[:0]
	suite.teamsServiceMock.Calls = suite.teamsServiceMock.Calls[:0]
	suite.metricsServiceMock.ExpectedCalls = suite.metricsServiceMock.ExpectedCalls[:0]
	suite.metricsServiceMock.Calls = suite.metricsServiceMock.Calls[:0]
	suite.emailServiceMock.ExpectedCalls = suite.emailServiceMock.ExpectedCalls[:0]
	suite.emailServiceMock.Calls = suite.emailServiceMock.Calls[:0]
	suite.emailFactoryMock.ExpectedCalls = suite.emailFactoryMock.ExpectedCalls[:0]
	suite.emailFactoryMock.Calls = suite.emailFactoryMock.Calls[:0]
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
func (suite *AlarmsTestSuite) mockFindActiveSubscriptionByUserID(userID uint, subscription *subscriptions.Subscription, err error) {
	suite.subscriptionsServiceMock.On(
		"FindActiveSubscriptionByUserID",
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

// Mock alarm down notification email
func (suite *AlarmsTestSuite) mockAlarmDownEmail() {
	emailMock := new(email.Email)
	suite.emailFactoryMock.On(
		"NewAlarmDownEmail",
		mock.AnythingOfType("*alarms.Alarm"),
	).Return(emailMock)
	suite.emailServiceMock.On("Send", emailMock).Return(nil)
}

// Mock alarm up notification email
func (suite *AlarmsTestSuite) mockAlarmUpEmail() {
	emailMock := new(email.Email)
	suite.emailFactoryMock.On(
		"NewAlarmUpEmail",
		mock.AnythingOfType("*alarms.Alarm"),
	).Return(emailMock)
	suite.emailServiceMock.On("Send", emailMock).Return(nil)
}

// Mock logging of request time metric
func (suite *AlarmsTestSuite) mockLogRequestTime(timestamp time.Time, referenceID uint, err error) {
	suite.metricsServiceMock.On(
		"LogRequestTime",
		timestamp,
		referenceID,
		mock.AnythingOfType("int64"),
	).Return(err)
}

// Mock counting of paginated request time metrics
func (suite *AlarmsTestSuite) mockPaginatedRequestTimesCount(alarmID uint, count int, err error) {
	suite.metricsServiceMock.On(
		"PaginatedRequestTimesCount",
		alarmID,
	).Return(count, err)
}

// Mock finding paginated request time metrics
func (suite *AlarmsTestSuite) mockFindPaginatedRequestTimes(offset, limit int, orderBy string, alarmID uint, requestTimes []*metrics.RequestTime, err error) {
	suite.metricsServiceMock.On(
		"FindPaginatedRequestTimes",
		offset,
		limit,
		orderBy,
		alarmID,
	).Return(requestTimes, err)
}
