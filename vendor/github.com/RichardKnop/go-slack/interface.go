package slack

// AdapterInterface defines exported methods
type AdapterInterface interface {
	// Exported methods
	SendMessage(channel, username, text, emoji string) error
}
