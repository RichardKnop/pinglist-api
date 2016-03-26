package subscriptions

import (
	"testing"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestUserIsInFreeTrial(t *testing.T) {
	var user *accounts.User

	user = &accounts.User{Model: gorm.Model{CreatedAt: time.Now()}}
	assert.True(t, IsInFreeTrial(user))

	user.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)
	assert.False(t, IsInFreeTrial(user))
}
