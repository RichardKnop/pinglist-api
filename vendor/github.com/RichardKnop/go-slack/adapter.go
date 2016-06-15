package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Adapter ...
type Adapter struct {
	cnf *Config
}

// NewAdapter starts a new Adapter instance
func NewAdapter(cnf *Config) *Adapter {
	return &Adapter{
		cnf: cnf,
	}
}

// SendMessage pushes a message to one of Slack's channels
func (a *Adapter) SendMessage(channel, username, text, emoji string) error {
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
	resp, err := http.PostForm(
		a.cnf.IncomingWebhook,
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
