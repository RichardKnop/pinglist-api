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
		fmt.Sprintf("/v1/accounts/teams/%d", team.ID), // href
		"", // title
	)

	// Create slice of member responses
	memberResponses := make([]*accounts.UserResponse, len(team.Members))
	for i, user := range team.Members {
		memberResponse, err := accounts.NewUserResponse(user)
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
