package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
	"github.com/RichardKnop/pinglist-api/teams"
)

var (
	// FreeTierMaxAlarms ...
	FreeTierMaxAlarms = uint(1)
	// FreeTierMinAlarmInterval ...
	FreeTierMinAlarmInterval = uint(60)
)

// countActiveAlarms counts active alarms of the current user or his/her team
func (s *Service) countActiveAlarms(team *teams.Team, user *accounts.User) uint {
	var (
		userIDs = []uint{user.ID}
		count   uint
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

type alarmLimitsConfig struct {
	maxAlarms        uint
	minAlarmInterval uint
}

// getAlarmLimits returns a struct containing different alarm limits based on
// a plan (e.g. max number of active alarms or min alarm interval)
func (s *Service) getAlarmLimits(team *teams.Team, user *accounts.User) *alarmLimitsConfig {
	var (
		alarmLimits = &alarmLimitsConfig{
			// Users in free tier get 1 free alarm all the time
			maxAlarms: FreeTierMaxAlarms,
			// Users in free tier have minimum alarm check interval of 60s
			minAlarmInterval: FreeTierMinAlarmInterval,
		}
		subscription *subscriptions.Subscription
	)

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
		alarmLimits.maxAlarms = subscription.Plan.MaxAlarms
		alarmLimits.minAlarmInterval = subscription.Plan.MinAlarmInterval
	}

	return alarmLimits
}
