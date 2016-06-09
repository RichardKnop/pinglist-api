package alarms

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/subscriptions"
)

var (
	// FreeTierMaxAlarms ...
	FreeTierMaxAlarms = uint(1)
	// FreeTierMinAlarmInterval ...
	FreeTierMinAlarmInterval = uint(60)
	// FreeTierMaxEmailsPerInterval ...
	FreeTierMaxEmailsPerInterval = uint(50)
)

// countActiveAlarms counts active alarms of the current user or his/her team
func (s *Service) countActiveAlarms(user *accounts.User) uint {
	var (
		userIDs = []uint{user.ID}
		count   uint
	)

	// Fetch the user team
	team, _ := s.teamsService.FindTeamByMemberID(user.ID)

	if team != nil {
		// Aggregate IDs of all team members
		for _, member := range team.Members {
			userIDs = append(userIDs, member.ID)
		}
	}

	// Count active alarms
	s.db.Model(new(Alarm)).Where("user_id IN (?)", userIDs).
		Where("active = ?", true).Count(&count)
	return count
}

type alarmLimitsConfig struct {
	maxAlarms            uint
	minAlarmInterval     uint
	unlimitedEmails      bool
	maxEmailsPerInterval uint
	slackAlerts          bool
}

// getAlarmLimits returns a struct containing different alarm limits based on
// a plan (e.g. max number of active alarms or min alarm interval)
func (s *Service) getAlarmLimits(user *accounts.User) *alarmLimitsConfig {
	var (
		alarmLimits = &alarmLimitsConfig{
			// Users in free tier get 1 free alarm all the time
			maxAlarms: FreeTierMaxAlarms,
			// Users in free tier have minimum alarm check interval of 60s
			minAlarmInterval: FreeTierMinAlarmInterval,
			// Emails are not unlimited by default
			unlimitedEmails: false,
			// Users in free tier have limit of 50 emails per interval
			maxEmailsPerInterval: FreeTierMaxEmailsPerInterval,
			// Slack alerts are not enabled by default
			slackAlerts: false,
		}
		subscription *subscriptions.Subscription
	)

	// Fetch the user team
	team, _ := s.teamsService.FindTeamByMemberID(user.ID)

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
		alarmLimits.unlimitedEmails = subscription.Plan.UnlimitedEmails
		alarmLimits.maxEmailsPerInterval = uint(subscription.Plan.MaxEmailsPerInterval.Int64)
		alarmLimits.slackAlerts = subscription.Plan.SlackAlerts
	}

	return alarmLimits
}
