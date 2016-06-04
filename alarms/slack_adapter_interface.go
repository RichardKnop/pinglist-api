package alarms

// SlackAdapterInterface defines exported methods
type SlackAdapterInterface interface {
	// Exported methods
	SendMessage(channel, username, text, emoji string) error
}
