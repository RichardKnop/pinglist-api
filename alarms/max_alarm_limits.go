package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/teams"
	"github.com/RichardKnop/pinglist-api/subscriptions"
)

// getMaxAlarms finds out max number of alarms
func (s *Service) getMaxAlarms(user *accounts.User) int {
	var (
		maxAlarms int
		team *teams.Team
		subscription *subscriptions.Subscription
		err error
	)

	// If user is in a free trial, allow one alarm
	if user.IsInFreeTrial() {
		maxAlarms = 1
	}

	// Is the user member of a team?
	team, err = s.teamsService.FindTeamByMemberID(user.ID)
	if err == nil {
		// User is member of a team, look for a team owner subscription
		subscription, err = s.subscriptionsService.FindActiveSubscriptionByUserID(team.Owner.ID)
	}

	// No subscription found yet, look for this user's subscription
	if subscription == nil {
		subscription, err = s.subscriptionsService.FindActiveSubscriptionByUserID(user.ID)
	}

	// If subscription found, take the max values from the subscription
	if subscription != nil {
		maxAlarms = int(subscription.Plan.MaxAlarms)
	}

	return maxAlarms
}
