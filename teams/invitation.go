package teams

import (
	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/jinzhu/gorm"
)

// inviteUser invites a new user to a team and sends an invitation email
func (s *Service) inviteUser(team *Team, invitedByUser *accounts.User, invitationRequest *accounts.InvitationRequest) (*accounts.Invitation, error) {
	// Begin a transaction
	tx := s.db.Begin()

	invitation, err := s.inviteUserCommon(tx, team, invitedByUser, invitationRequest)
	if err != nil {
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

// inviteUser invites a new user to a team and sends an invitation email in a transaction
func (s *Service) inviteUserTx(tx *gorm.DB, team *Team, invitedByUser *accounts.User, invitationRequest *accounts.InvitationRequest) (*accounts.Invitation, error) {
	return s.inviteUserCommon(tx, team, invitedByUser, invitationRequest)
}

func (s *Service) inviteUserCommon(db *gorm.DB, team *Team, invitedByUser *accounts.User, invitationRequest *accounts.InvitationRequest) (*accounts.Invitation, error) {
	invitation, err := s.GetAccountsService().InviteUserTx(db, invitedByUser, invitationRequest)
	if err != nil {
		return nil, err
	}

	// Update owners association
	membersAssoc := db.Model(team).Association("Members")
	if err := membersAssoc.Append(invitation.InvitedUser).Error; err != nil {
		return nil, err
	}

	return invitation, nil
}
