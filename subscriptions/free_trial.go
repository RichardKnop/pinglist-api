package subscriptions

import (
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
)

var (
	// FreeTrialPeriod ...
	FreeTrialPeriod = time.Duration(30 * 24 * time.Hour) // 30 days
	// FreeTrialMaxAlarms ...
	FreeTrialMaxAlarms = 1
)

// IsInFreeTrial returns true if user is in a trial period
func IsInFreeTrial(user *accounts.User) bool {
	return time.Now().Before(user.CreatedAt.Add(FreeTrialPeriod))
}
