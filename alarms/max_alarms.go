package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
)

// getMaxAlarms finds out how many alarms user is allowed to have
func (s *Service) getUserMaxAlarms(user *accounts.User) int {
	var maxAlarms int

	// If user is in a free trial, allow one alarm
	if user.IsInFreeTrial() {
		maxAlarms = 1
	}

	// Fetch active user subscription
	subscription, err := s.subscriptionsService.FindActiveUserSubscription(user.ID)

	// If subscribed, take max allowed alarms from the subscription plan
	if err == nil && subscription != nil {
		maxAlarms = int(subscription.Plan.MaxAlarms)
	}

	return maxAlarms
}
