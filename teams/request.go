package teams

// TeamMemberRequest ...
type TeamMemberRequest struct {
	ID uint `json:"id"`
}

// TeamRequest ...
type TeamRequest struct {
	Name    string               `json:"name"`
	Members []*TeamMemberRequest `json:"members"`
}
