[![Codeship Status for AreaHQ/go-slack](https://codeship.com/projects/6812efd0-14f0-0134-4f8d-12348d1f3442/status?branch=master)](https://codeship.com/projects/157933)

[![GoDoc](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/RichardKnop/go-slack)
[![Travis Status for RichardKnop/go-slack](https://travis-ci.org/RichardKnop/go-slack.svg?branch=master)](https://travis-ci.org/RichardKnop/go-slack)
[![Donate Bitcoin](https://img.shields.io/badge/donate-bitcoin-orange.svg)](https://richardknop.github.io/donate/)

# go-slack

A simple Golang SDK for Slack.

## Usage

```go
package main

import (
	"log"

	"github.com/AreaHQ/go-slack"
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
