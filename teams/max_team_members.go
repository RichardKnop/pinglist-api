package teams

import (
	"github.com/RichardKnop/pinglist-api/accounts"
)

// getUserMaxTeamMembers finds out how many members can user add to his/her team
func (s *Service) getUserMaxTeamMembers(user *accounts.User) int {
	var maxTeamMembers int

	// Fetch active user subscription
	subscription, err := s.subscriptionsService.FindActiveUserSubscription(user.ID)

	// If subscribed, take max allowed team members from the subscription plan
	if err == nil && subscription != nil {
		maxTeamMembers = int(subscription.Plan.MaxTeamMembers)
	}

	return maxTeamMembers
}
