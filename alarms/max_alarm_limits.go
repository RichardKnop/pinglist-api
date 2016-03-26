package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
)

// getMaxAlarms finds out max number of alarms
func (s *Service) getMaxAlarms(user *accounts.User) int {
	var maxAlarms int

	// If user is in a free trial, allow one alarm
	if user.IsInFreeTrial() {
		maxAlarms = 1
	}

	// Fetch active user subscription
	subscription, err := s.subscriptionsService.FindActiveSubscriptionByUserID(user.ID)

	// If subscribed, take the max values from the subscription
	if err == nil && subscription != nil {
		maxAlarms = int(subscription.Plan.MaxAlarms)
	}

	return maxAlarms
}
