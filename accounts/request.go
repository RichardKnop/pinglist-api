package accounts

// UserRequest ...
type UserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
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
