[![Codeship Status for AreaHQ/go-slack](https://codeship.com/projects/6812efd0-14f0-0134-4f8d-12348d1f3442/status?branch=master)](https://codeship.com/projects/157933)

# Go Slack

A simple Golang SDK for Slack.

## Usage

```go
package main

import (
	"log"

	slack "github.com/AreaHQ/go-slack"
)

func main() {
	cnf := &slack.Config{IncomingWebhook: "incoming_webhook"}
	adapter := slack.NewAdapter(cnf)
	err = adapter.SendMessage(
		"#some-channel",
		"some-username",
		"message to send",
		"", // emmoji
	)
}
```
