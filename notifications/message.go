package notifications

import (
	"encoding/json"
	"unicode/utf8"
)

const (
	// MessageBodyLimit is maximum allowed message payload size
	MessageBodyLimit = 2000
	gcmKeyMessage    = "message"
	apnsKeyMessage   = "alert"
	apnsKeySound     = "sound"
	apnsKeyBadge     = "badge"
)

// Make sns message for Apple Push Notification Service
// For full documentation of APNS message payload keys, see:
// https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/TheNotificationPayload.html#//apple_ref/doc/uid/TP40008194-CH107-SW1
func composeMessageAPNS(msg string, opt map[string]interface{}) (string, error) {
	aps := make(map[string]interface{})
	aps[apnsKeyMessage] = msg

	aps[apnsKeySound] = "default"
	if v, ok := opt[apnsKeySound]; ok {
		aps[apnsKeySound] = v
	}

	if v, ok := opt[apnsKeyBadge]; ok {
		aps[apnsKeyBadge] = v
	}

	message := make(map[string]interface{})
	message["aps"] = aps
	for k, v := range opt {
		switch k {
		case apnsKeySound:
			continue
		case apnsKeyBadge:
			continue
		default:
			message[k] = v
		}
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

// Make sns message for Google Cloud Messaging
func composeMessageGCM(msg string, opt map[string]interface{}) (string, error) {
	data := make(map[string]interface{})
	data[gcmKeyMessage] = msg
	for k, v := range opt {
		data[k] = v
	}

	message := make(map[string]interface{})
	message["data"] = data

	payload, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

// Limit message size to the allowed payload size
func truncateMessage(msg string) string {
	if len(msg) <= MessageBodyLimit {
		return msg
	}
	runes := []rune(msg[:MessageBodyLimit])
	valid := len(runes)
	// traverse runes from last string and detect invalid string
	for i := valid; ; {
		i--
		if runes[i] != utf8.RuneError {
			break
		}
		valid = i
	}
	return string(runes[0:valid])
}
