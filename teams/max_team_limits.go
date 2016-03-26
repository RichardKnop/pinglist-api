package teams

import (
	"github.com/RichardKnop/pinglist-api/accounts"
)

// getMaxTeamLimits finds out max number of teams as well as upper limit of members per one team
func (s *Service) getMaxTeamLimits(user *accounts.User) (int, int) {
	var (
		maxTeams          int
		maxMembersPerTeam int
	)

	// Fetch active user subscription
	subscription, err := s.subscriptionsService.FindActiveSubscriptionByUserID(user.ID)

	// If subscribed, take the max values from the subscription
	if err == nil && subscription != nil {
		maxTeams = int(subscription.Plan.MaxTeams)
		maxMembersPerTeam = int(subscription.Plan.MaxMembersPerTeam)
	}

	return maxTeams, maxMembersPerTeam
}
