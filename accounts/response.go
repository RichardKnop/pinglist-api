package accounts

import (
	"fmt"
	"time"

	"github.com/RichardKnop/jsonhal"
)

// UserResponse ...
type UserResponse struct {
	jsonhal.Hal
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Confirmed bool   `json:"confirmed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// InvitationResponse ...
type InvitationResponse struct {
	jsonhal.Hal
	ID              uint   `json:"id"`
	Reference       string `json:"reference"`
	InvitedUserID   uint   `json:"invited_user_id"`
	InvitedByUserID uint   `json:"invited_by_user_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// NewUserResponse creates new UserResponse instance
func NewUserResponse(user *User) (*UserResponse, error) {
	response := &UserResponse{
		ID:        user.OauthUser.ID,
		Email:     user.OauthUser.Username,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Role:      user.RoleID.String,
		Confirmed: user.Confirmed,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/accounts/users/%d", user.ID), // href
		"", // title
	)

	return response, nil
}

// NewInvitationResponse creates new InvitationResponse instance
func NewInvitationResponse(invitation *Invitation) (*InvitationResponse, error) {
	response := &InvitationResponse{
		ID:              invitation.ID,
		Reference:       invitation.Reference,
		InvitedUserID:   invitation.InvitedUser.ID,
		InvitedByUserID: invitation.InvitedByUser.ID,
		CreatedAt:       invitation.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       invitation.UpdatedAt.UTC().Format(time.RFC3339),
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/accounts/invitations/%d", invitation.ID), // href
		"", // title
	)

	return response, nil
}
