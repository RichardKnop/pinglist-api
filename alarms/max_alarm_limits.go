package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
)

var (
	// FreeTierMaxAlarms ...
	FreeTierMaxAlarms = 1
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

	// Users in free tier get 1 free alarm all the time
	maxAlarms = FreeTierMaxAlarms

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
