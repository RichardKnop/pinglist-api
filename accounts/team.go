package accounts

import "errors"

var (
	// ErrTeamNotFound ...
	ErrTeamNotFound = errors.New("Team not found")
)

// findTeamByOwnerID looks up a team by an owner ID
func (s *Service) findTeamByOwnerID(ownerID uint) (*Team, error) {
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
	// TODO

	return nil, nil
}

// updateTeam updates an existing team
func (s *Service) updateTeam(team *Team, teamRequest *TeamRequest) error {
	// TODO

	return nil
}
