package teams

// TeamMemberRequest ...
type TeamMemberRequest struct {
	Email string `json:"email"`
}

// TeamRequest ...
type TeamRequest struct {
	Name    string               `json:"name"`
	Members []*TeamMemberRequest `json:"members"`
}
