package accounts

import "errors"

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
	notFound := s.db.Preload("Owner").Preload("Members").
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
		Preload("Owner").Preload("Members").First(team).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrTeamNotFound
	}

	return team, nil
}

// createTeam creates a new team
func (s *Service) createTeam(user *User, teamRequest *TeamRequest) (*Team, error) {
	// Users are allowed to create only one team
	_, err := s.FindTeamByOwnerID(user.ID)
	if err == nil {
		return nil, ErrUserCanCreateOnlyOneTeam
	}

	return nil, nil
}

// updateTeam updates an existing team
func (s *Service) updateTeam(team *Team, teamRequest *TeamRequest) error {
	// TODO

	return nil
}
