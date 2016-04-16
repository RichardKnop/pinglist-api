package teams

import (
	"errors"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/jinzhu/gorm"
)

var (
	// ErrTeamNotFound ...
	ErrTeamNotFound = errors.New("Team not found")
	// ErrMaxTeamsLimitReached ...
	ErrMaxTeamsLimitReached = errors.New("Max teams limit reached")
	// ErrMaxMembersPerTeamLimitReached ...
	ErrMaxMembersPerTeamLimitReached = errors.New("Max members per team limit reached")
	// ErrUserCanOnlyBeMemberOfOneTeam ...
	ErrUserCanOnlyBeMemberOfOneTeam = errors.New("User can only be member of one team")
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

// FindTeamByMemberID looks up a team by a member ID
func (s *Service) FindTeamByMemberID(memberID uint) (*Team, error) {
	// Fetch the team from the database
	team := new(Team)
	notFound := s.db.
		Joins("inner join team_team_members on team_team_members.team_id = team_teams.id").
		Where("team_team_members.user_id = ?", memberID).
		Preload("Owner.OauthUser").Preload("Members.OauthUser").
		First(team).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrTeamNotFound
	}

	return team, nil
}

// createTeam creates a new team
func (s *Service) createTeam(owner *accounts.User, teamRequest *TeamRequest) (*Team, error) {
	maxTeams, maxMembersPerTeam := s.getMaxTeamLimits(owner)

	// Limit teams to the max number defined as per subscription plan
	var teamsCount int
	err := s.db.Model(new(Team)).Where("owner_id = ?", owner.ID).Count(&teamsCount).Error
	if err != nil {
		return nil, err
	}
	if teamsCount >= maxTeams {
		return nil, ErrMaxTeamsLimitReached
	}

	// Limit team members to the max number defined as per subscription plan
	if len(teamRequest.Members) >= maxMembersPerTeam {
		return nil, ErrMaxMembersPerTeamLimitReached
	}

	// Members
	members := make([]*accounts.User, len(teamRequest.Members))
	for i, teamMemberRequest := range teamRequest.Members {
		// Fetch the member from the database
		member, err := s.GetAccountsService().FindUserByEmail(teamMemberRequest.Email)
		if err != nil {
			return nil, err
		}

		// Users can only be members of a single team
		_, err = s.FindTeamByMemberID(member.ID)
		if err == nil {
			return nil, ErrUserCanOnlyBeMemberOfOneTeam
		}

		members[i] = member
	}

	// Create a new team
	team := NewTeam(owner, members, teamRequest.Name)

	// Save the team to the database
	if err := s.db.Create(team).Error; err != nil {
		return nil, err
	}

	return team, nil
}

// updateTeam updates an existing team
func (s *Service) updateTeam(team *Team, teamRequest *TeamRequest) error {
	maxTeams, maxMembersPerTeam := s.getMaxTeamLimits(team.Owner)

	// Limit teams to the max number defined as per subscription plan
	var teamsCount int
	err := s.db.Model(new(Team)).Where("owner_id = ?", team.Owner.ID).Count(&teamsCount).Error
	if err != nil {
		return err
	}
	if teamsCount >= maxTeams {
		return ErrMaxTeamsLimitReached
	}

	// Limit team members to the max number defined as per subscription plan
	if len(teamRequest.Members) >= maxMembersPerTeam {
		return ErrMaxMembersPerTeamLimitReached
	}

	// Members
	members := make([]*accounts.User, len(teamRequest.Members))
	for i, teamMemberRequest := range teamRequest.Members {
		// Fetch the member from the database
		member, err := s.GetAccountsService().FindUserByEmail(teamMemberRequest.Email)
		if err != nil {
			return err
		}

		// Users can only be members of a single team
		memberTeam, err := s.FindTeamByMemberID(member.ID)
		if err == nil && memberTeam.ID != team.ID {
			return ErrUserCanOnlyBeMemberOfOneTeam
		}

		members[i] = member
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Update basic metadata
	if err := s.db.Model(team).UpdateColumns(Team{
		Name:  teamRequest.Name,
		Model: gorm.Model{UpdatedAt: time.Now()},
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

// paginatedTeamsCount returns a total count of teams
// Can be optionally filtered by owner
func (s *Service) paginatedTeamsCount(owner *accounts.User) (int, error) {
	var count int
	if err := s.paginatedTeamsQuery(owner).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedTeams returns paginated team records
// Results can optionally be filtered by owner
func (s *Service) findPaginatedTeams(offset, limit int, orderBy string, owner *accounts.User) ([]*Team, error) {
	var teams []*Team

	// Get the pagination query
	teamsQuery := s.paginatedTeamsQuery(owner)

	// Default ordering
	if orderBy == "" {
		orderBy = "id"
	}

	// Retrieve paginated results from the database
	err := teamsQuery.Offset(offset).Limit(limit).Order(orderBy).
		Preload("Owner.OauthUser").Preload("Members.OauthUser").
		Find(&teams).Error
	if err != nil {
		return teams, err
	}

	return teams, nil
}

// paginatedTeamsQuery returns a db query for paginated teams
func (s *Service) paginatedTeamsQuery(owner *accounts.User) *gorm.DB {
	// Basic query
	teamsQuery := s.db.Model(new(Team))

	// Optionally filter by user
	if owner != nil {
		teamsQuery = teamsQuery.Where("owner_id = ?", owner.ID)
	}

	return teamsQuery
}
