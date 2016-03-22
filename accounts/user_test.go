package accounts

import (
	"testing"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts/roles"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestUserGetName(t *testing.T) {
	user := &User{}
	assert.Equal(t, "", user.GetName())

	user.FirstName = util.StringOrNull("John")
	user.LastName = util.StringOrNull("Reese")
	assert.Equal(t, "John Reese", user.GetName())
}

func TestUserIsInFreeTrial(t *testing.T) {
	user := &User{Model: gorm.Model{CreatedAt: time.Now()}}
	assert.True(t, user.IsInFreeTrial())

	user.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)
	assert.False(t, user.IsInFreeTrial())
}

func (suite *AccountsTestSuite) TestFindUserByOauthUserID() {
	var (
		user *User
		err  error
	)

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByOauthUserID(12345)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByOauthUserID(suite.users[1].OauthUser.ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.Role.ID)
	}
}

func (suite *AccountsTestSuite) TestFindUserByEmail() {
	var (
		user *User
		err  error
	)

	// Let's try to find a user by a bogus email
	user, err = suite.service.FindUserByEmail("bogus")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserNotFound, err)
	}

	// Now let's pass a valid email
	user, err = suite.service.FindUserByEmail("test@user")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.Role.ID)
	}
}

func (suite *AccountsTestSuite) TestFindUserByID() {
	var (
		user *User
		err  error
	)

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByID(12345)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByID(suite.users[1].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned with preloaded data
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.Role.ID)
	}
}

func (suite *AccountsTestSuite) TestFindUserByFacebookID() {
	var (
		user *User
		err  error
	)

	// Let's try to find a user by an empty string Facebook ID
	user, err = suite.service.FindUserByFacebookID("")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserNotFound, err)
	}

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByFacebookID("bogus")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByFacebookID(suite.users[1].FacebookID.String)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned with preloaded data
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.Role.ID)
	}
}

func (suite *AccountsTestSuite) TestGetOrCreateFacebookUser() {
	var (
		countBefore, countAfter int
		user                    *User
		err                     error
	)

	// Count before
	suite.db.Model(new(User)).Count(&countBefore)

	// Let's try passing an existing facebook ID
	user, err = suite.service.GetOrCreateFacebookUser(
		suite.accounts[0], // account
		"facebook_id_2",   // facebook ID
		new(UserRequest),  // not important in this case
	)

	// Count after
	suite.db.Model(new(User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), "facebook_id_2", user.FacebookID.String)
		assert.Equal(suite.T(), "test_first_name_2", user.FirstName.String)
		assert.Equal(suite.T(), "test_last_name_2", user.LastName.String)
		assert.True(suite.T(), user.Confirmed)
	}

	// Count before
	suite.db.Model(new(User)).Count(&countBefore)

	// Let's try passing an existing email
	user, err = suite.service.GetOrCreateFacebookUser(
		suite.accounts[0],  // account
		"user_facebook_id", // facebook ID
		&UserRequest{
			Email:     "test@user",
			FirstName: "Harold",
			LastName:  "Finch",
		},
	)

	// Count after
	suite.db.Model(new(User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), "user_facebook_id", user.FacebookID.String)
		assert.Equal(suite.T(), "Harold", user.FirstName.String)
		assert.Equal(suite.T(), "Finch", user.LastName.String)
		assert.True(suite.T(), user.Confirmed)
	}

	// Count before
	suite.db.Model(new(User)).Count(&countBefore)

	// We pass new facebook ID and new email
	user, err = suite.service.GetOrCreateFacebookUser(
		suite.accounts[0],     // account
		"newuser_facebook_id", // facebook ID
		&UserRequest{
			Email:     "test@newuser",
			FirstName: "John",
			LastName:  "Reese",
		},
	)

	// Count after
	suite.db.Model(new(User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@newuser", user.OauthUser.Username)
		assert.Equal(suite.T(), "newuser_facebook_id", user.FacebookID.String)
		assert.Equal(suite.T(), "John", user.FirstName.String)
		assert.Equal(suite.T(), "Reese", user.LastName.String)
		assert.True(suite.T(), user.Confirmed)
	}
}

func (suite *AccountsTestSuite) TestCreateSuperuser() {
	var (
		user *User
		err  error
	)

	// We try to insert a user with a non unique oauth user
	user, err = suite.service.CreateSuperuser(
		suite.accounts[0], // account
		"test@superuser",  // email
		"test_password",   // password
	)

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrEmailTaken, err)
	}

	// We try to insert a unique superuser
	user, err = suite.service.CreateSuperuser(
		suite.accounts[0],   // account
		"test@newsuperuser", // email
		"test_password",     // password
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@newsuperuser", user.OauthUser.Username)
		assert.False(suite.T(), user.FacebookID.Valid)
		assert.False(suite.T(), user.FirstName.Valid)
		assert.False(suite.T(), user.LastName.Valid)
		assert.True(suite.T(), user.Confirmed)
	}
}
