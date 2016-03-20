package accounts

import (
	"errors"
)

var (
	// ErrTeamNotFound ...
	ErrTeamNotFound = errors.New("Team not found")
	// ErrUserCanCreateOnlyOneTeam ...
	ErrUserCanCreateOnlyOneTeam = errors.New("User can create only one team")
)

// FindTeamByID looks up a team by ID
func (s *Service) FindTeamByID(teamID uint) (*Team, error) {
	// Fetch the team from the database
	team := new(Team)
	notFound := s.db.Preload("Owner.OauthUser").Preload("Members.OauthUser").
		First(team, teamID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrTeamNotFound
	}

	return team, nil
}

// FindTeamByOwnerID looks up a team by an owner ID
func (s *Service) FindTeamByOwnerID(ownerID uint) (*Team, error) {
	// Fetch the team from the database
	team := new(Team)
	notFound := s.db.Where("owner_id = ?", ownerID).
		Preload("Owner.OauthUser").Preload("Members.OauthUser").
		First(team).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrTeamNotFound
	}

	return team, nil
}

// userOwnsTeam returns true if the user is an owner of a team already
func (s *Service) userOwnsTeam(user *User) bool {
	_, err := s.FindTeamByOwnerID(user.ID)
	return err == nil
}

// createTeam creates a new team
func (s *Service) createTeam(owner *User, teamRequest *TeamRequest) (*Team, error) {
	// Users can only be owners of a single team
	if s.userOwnsTeam(owner) {
		return nil, ErrUserCanCreateOnlyOneTeam
	}

	// Members
	members := make([]*User, len(teamRequest.Members))
	for i, teamMemberRequest := range teamRequest.Members {
		// Fetch the member from the database
		member, err := s.FindUserByID(teamMemberRequest.ID)
		if err != nil {
			return nil, err
		}
		members[i] = member
	}

	// Create a new team
	team := newTeam(owner, members, teamRequest)

	// Save the team to the database
	if err := s.db.Create(team).Error; err != nil {
		return nil, err
	}

	return team, nil
}

// updateTeam updates an existing team
func (s *Service) updateTeam(team *Team, teamRequest *TeamRequest) error {
	// Members
	members := make([]*User, len(teamRequest.Members))
	for i, teamMemberRequest := range teamRequest.Members {
		// Fetch the member from the database
		member, err := s.FindUserByID(teamMemberRequest.ID)
		if err != nil {
			return err
		}
		members[i] = member
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Update basic metadata
	if err := s.db.Model(team).UpdateColumns(Team{
		Name: teamRequest.Name,
	}).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Update owners association
	membersAssoc := tx.Model(team).Association("Members")
	if err := membersAssoc.Replace(members).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}
