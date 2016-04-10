package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
)

// countActiveAlarms counts active alarms of the current user or his/her team
func (s *Service) countActiveAlarms(team *teams.Team, user *accounts.User) int {
	var (
		userIDs = []uint{user.ID}
		count   int
	)
	if team != nil {
		for _, member := range team.Members {
			userIDs = append(userIDs, member.ID)
		}
	}
	s.db.Model(new(Alarm)).Where("user_id IN (?)", userIDs).
		Where("active = ?", true).Count(&count)
	return count
}

// getMaxAlarms finds out max number of alarms
func (s *Service) getMaxAlarms(team *teams.Team, user *accounts.User) int {
	var (
		maxAlarms    int
		subscription *subscriptions.Subscription
	)

	// If user is in a free trial, allow one alarm
	if subscriptions.IsInFreeTrial(user) {
		maxAlarms = subscriptions.FreeTrialMaxAlarms
	}

	// If the user is member of a team, look for a team owner subscription
	if team != nil {
		subscription, _ = s.subscriptionsService.FindActiveSubscriptionByUserID(team.Owner.ID)
	}

	// No subscription found yet, look for this user's subscription
	if subscription == nil {
		subscription, _ = s.subscriptionsService.FindActiveSubscriptionByUserID(user.ID)
	}

	// If subscription found, take the max values from the subscription
	if subscription != nil {
		maxAlarms = int(subscription.Plan.MaxAlarms)
	}

	return maxAlarms
}
