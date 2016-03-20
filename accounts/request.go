package accounts

// UserRequest ...
type UserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// PasswordResetRequest ...
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// TeamMemberRequest ...
type TeamMemberRequest struct {
	ID uint `json:"id"`
}

// TeamRequest ...
type TeamRequest struct {
	Name    string               `json:"name"`
	Members []*TeamMemberRequest `json:"members"`
}
