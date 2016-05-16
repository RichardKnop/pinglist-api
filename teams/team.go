package teams

import (
	"errors"
	"fmt"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/jinzhu/gorm"
)

// ErrUserCanOnlyBeMemberOfOneTeam ...
type ErrUserCanOnlyBeMemberOfOneTeam struct {
	email, teamName string
}

// NewErrUserCanOnlyBeMemberOfOneTeam returns new ErrUserCanOnlyBeMemberOfOneTeam
func NewErrUserCanOnlyBeMemberOfOneTeam(email, teamName string) ErrUserCanOnlyBeMemberOfOneTeam {
	return ErrUserCanOnlyBeMemberOfOneTeam{email, teamName}
}

// Error method so we implement the error interface
func (e ErrUserCanOnlyBeMemberOfOneTeam) Error() string {
	return fmt.Sprintf("%s is already member of the %s", e.email, e.teamName)
}

var (
	// ErrTeamNotFound ...
	ErrTeamNotFound = errors.New("Team not found")
	// ErrMaxTeamsLimitReached ...
	ErrMaxTeamsLimitReached = errors.New("Max teams limit reached")
	// ErrMaxMembersPerTeamLimitReached ...
	ErrMaxMembersPerTeamLimitReached = errors.New("Max members per team limit reached")
	// ErrCannotAddYourself ...
	ErrCannotAddYourself = errors.New("You cannot add yourself to the you have created")
)

// FindTeamByID looks up a team by ID
func (s *Service) FindTeamByID(teamID uint) (*Team, error) {
	// Fetch the team from the database
	team := new(Team)
	notFound := s.db.Preload("Owner.Account").Preload("Owner.OauthUser").
		Preload("Owner.Role").Preload("Members.Account").Preload("Members.OauthUser").
		Preload("Members.Role").First(team, teamID).RecordNotFound()

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
	notFound := s.db.Where("owner_id = ?", ownerID).Preload("Owner.Account").
		Preload("Owner.OauthUser").Preload("Owner.Role").Preload("Members.Account").
		Preload("Members.OauthUser").Preload("Members.Role").
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
		Preload("Owner.Account").Preload("Owner.OauthUser").Preload("Owner.Role").
		Preload("Members.Account").Preload("Members.OauthUser").Preload("Members.Role").
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

	// Begin a transaction
	tx := s.db.Begin()

	// Create a new team
	team := NewTeam(owner, []*accounts.User{}, teamRequest.Name)

	// Save the team to the database
	if err := tx.Create(team).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Members
	members := make([]*accounts.User, len(teamRequest.Members))
	for i, teamMemberRequest := range teamRequest.Members {
		// Owner cannot add himself / herself
		if teamMemberRequest.Email == owner.OauthUser.Username {
			tx.Rollback() // rollback the transaction
			return nil, ErrCannotAddYourself
		}

		var member *accounts.User

		// Fetch the member from the database
		member, err = s.GetAccountsService().FindUserByEmail(teamMemberRequest.Email)
		if err != nil {
			switch err {
			case accounts.ErrUserNotFound:
				invitation, err := s.inviteUserTx(
					tx,
					team,
					teamMemberRequest.Email,
					false, // update members assoc
				)
				if err != nil {
					tx.Rollback() // rollback the transaction
					return nil, err
				}
				member = invitation.InvitedUser
			default:
				tx.Rollback() // rollback the transaction
				return nil, err
			}
		}

		// Users can only be members of a single team
		memberTeam, err := s.FindTeamByMemberID(member.ID)
		if err == nil && memberTeam.ID != team.ID {
			tx.Rollback() // rollback the transaction
			return nil, NewErrUserCanOnlyBeMemberOfOneTeam(
				member.OauthUser.Username,
				memberTeam.Name,
			)
		}

		members[i] = member
	}

	// Update members association
	membersAssoc := tx.Model(team).Association("Members")
	if err := membersAssoc.Replace(members).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
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
	if teamsCount > maxTeams {
		return ErrMaxTeamsLimitReached
	}

	// Limit team members to the max number defined as per subscription plan
	if len(teamRequest.Members) >= maxMembersPerTeam {
		return ErrMaxMembersPerTeamLimitReached
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Update basic metadata
	if err := tx.Model(team).UpdateColumns(Team{
		Name:  teamRequest.Name,
		Model: gorm.Model{UpdatedAt: time.Now()},
	}).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Members
	members := make([]*accounts.User, len(teamRequest.Members))
	for i, teamMemberRequest := range teamRequest.Members {
		// Owner cannot add himself / herself
		if teamMemberRequest.Email == team.Owner.OauthUser.Username {
			tx.Rollback() // rollback the transaction
			return ErrCannotAddYourself
		}

		var member *accounts.User

		// Fetch the member from the database
		member, err = s.GetAccountsService().FindUserByEmail(teamMemberRequest.Email)
		if err != nil {
			switch err {
			case accounts.ErrUserNotFound:
				invitation, err := s.inviteUserTx(
					tx,
					team,
					teamMemberRequest.Email,
					false, // update members assoc
				)
				if err != nil {
					tx.Rollback() // rollback the transaction
					return err
				}
				member = invitation.InvitedUser
			default:
				tx.Rollback() // rollback the transaction
				return err
			}
		}

		// Users can only be members of a single team
		memberTeam, err := s.FindTeamByMemberID(member.ID)
		if err == nil && memberTeam.ID != team.ID {
			tx.Rollback() // rollback the transaction
			return NewErrUserCanOnlyBeMemberOfOneTeam(
				member.OauthUser.Username,
				memberTeam.Name,
			)
		}

		members[i] = member
	}

	// Update members association
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
		Preload("Owner.Account").Preload("Owner.OauthUser").Preload("Owner.Role").
		Preload("Members.Account").Preload("Members.OauthUser").Preload("Members.Role").
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
