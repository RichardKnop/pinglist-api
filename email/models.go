package email

// Recipient represents a single email recipient
type Recipient struct {
	Email string
	Name  string
}

// Email represents an email message
type Email struct {
	Subject    string
	Recipients []*Recipient
	From       string
	Text       string
}
