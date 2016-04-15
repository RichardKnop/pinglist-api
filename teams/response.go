package teams

import (
	"fmt"
	"time"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/accounts"
)

// TeamResponse ...
type TeamResponse struct {
	jsonhal.Hal
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListTeamsResponse ...
type ListTeamsResponse struct {
	jsonhal.Hal
	Count uint `json:"count"`
	Page  uint `json:"page"`
}

// NewTeamResponse creates new TeamResponse instance
func NewTeamResponse(team *Team) (*TeamResponse, error) {
	response := &TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		CreatedAt: team.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: team.UpdatedAt.UTC().Format(time.RFC3339),
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/teams/%d", team.ID), // href
		"", // title
	)

	// Create slice of member responses
	memberResponses := make([]*accounts.UserResponse, len(team.Members))
	for i, member := range team.Members {
		memberResponse, err := accounts.NewUserResponse(member)
		if err != nil {
			return nil, err
		}
		memberResponses[i] = memberResponse
	}

	// Set embedded members
	response.SetEmbedded(
		"members",
		jsonhal.Embedded(memberResponses),
	)

	return response, nil
}

// NewListTeamsResponse creates new ListTeamsResponse instance
func NewListTeamsResponse(count, page int, self, first, last, previous, next string, teams []*Team) (*ListTeamsResponse, error) {
	response := &ListTeamsResponse{
		Count: uint(count),
		Page:  uint(page),
	}

	// Set the self link
	response.SetLink("self", self, "")

	// Set the first link
	response.SetLink("first", first, "")

	// Set the last link
	response.SetLink("last", last, "")

	// Set the previous link
	response.SetLink("prev", previous, "")

	// Set the next link
	response.SetLink("next", next, "")

	// Create slice of team responses
	teamResponses := make([]*TeamResponse, len(teams))
	for i, team := range teams {
		teamResponse, err := NewTeamResponse(team)
		if err != nil {
			return nil, err
		}
		teamResponses[i] = teamResponse
	}

	// Set embedded teams
	response.SetEmbedded(
		"teams",
		jsonhal.Embedded(teamResponses),
	)

	return response, nil
}
