package alarms

// SlackAdapterInterface defines exported methods
type SlackAdapterInterface interface {
	// Exported methods
	SendMessage(incomingWebhook, channel, username, emoji, text string) error
}
