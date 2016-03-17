package notifications

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedAPNSMessage = `{
  "aps": {
    "alert": "Hello world message.",
    "sound": "default"
  },
  "foo": "bar"
}`

var expectedGCMMessage = `{
  "data": {
    "foo": "bar",
    "message": "Hello world message."
  }
}`

func TestComposeMessageAPNS(t *testing.T) {
	actual, err := composeMessageAPNS(
		"Hello world message.",
		map[string]interface{}{
			"foo": "bar",
		},
	)

	expected := bytes.NewBuffer([]byte{})
	if err := json.Compact(expected, []byte(expectedAPNSMessage)); err != nil {
		log.Fatal(err)
	}

	if assert.Nil(t, err) {
		assert.Equal(t, expected.String(), actual)
	}
}

func TestComposeMessageGCM(t *testing.T) {
	actual, err := composeMessageGCM(
		"Hello world message.",
		map[string]interface{}{
			"foo": "bar",
		},
	)

	expected := bytes.NewBuffer([]byte{})
	if err := json.Compact(expected, []byte(expectedGCMMessage)); err != nil {
		log.Print(actual)
		log.Fatal(err)
	}

	if assert.Nil(t, err) {
		assert.Equal(t, expected.String(), actual)
	}
}
