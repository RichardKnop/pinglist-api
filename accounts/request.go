package accounts

// UserRequest ...
type UserRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	NewPassword          string `json:"new_password"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Role                 string `json:"role"`
	SlackIncomingWebhook string `json:"slack_incoming_webhook"`
	SlackChannel         string `json:"slack_channel"`
}

// InvitationRequest ...
type InvitationRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// PasswordResetRequest ...
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// PasswordRequest ...
type PasswordRequest struct {
	Password string `json:"password"`
}

// ConfirmInvitationRequest ...
type ConfirmInvitationRequest struct {
	PasswordRequest
}

// ConfirmPasswordResetRequest ...
type ConfirmPasswordResetRequest struct {
	PasswordRequest
}

// ContactRequest ...
type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
