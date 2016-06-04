package alarms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/RichardKnop/pinglist-api/config"
)

// SlackAdapter ...
type SlackAdapter struct {
	incomingWebhook string
}

// NewSlackAdapter starts a new Adapter instance
func NewSlackAdapter(cnf *config.Config) *SlackAdapter {
	return &SlackAdapter{
		incomingWebhook: cnf.Slack.IncomingWebhook,
	}
}

// SendMessage pushes a message to one of Slack's channels
func (a *SlackAdapter) SendMessage(channel, username, text, emoji string) error {
	payload := map[string]string{
		"channel":  channel,
		"username": username,
		"text":     text,
		"emoji":    emoji,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.PostForm(a.incomingWebhook,
		url.Values{
			"payload": {string(payloadJSON)},
		},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Send Slack Message Error: %s", string(body))
	}
	return nil
}
