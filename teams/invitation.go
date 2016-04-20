package teams

import (
	"github.com/RichardKnop/pinglist-api/accounts"
)

// inviteUser invites a new user to a team and sends an invitation email
func (s *Service) inviteUser(team *Team, invitedByUser *accounts.User, invitationRequest *accounts.InvitationRequest) (*accounts.Invitation, error) {
	// Begin a transaction
	tx := s.db.Begin()

	invitation, err := s.GetAccountsService().InviteUserTx(tx, invitedByUser, invitationRequest)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Update owners association
	membersAssoc := tx.Model(team).Association("Members")
	if err := membersAssoc.Append(invitation.InvitedUser).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return invitation, nil
}
