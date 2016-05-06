package notifications

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SNSAdapter struct keeps objects to avoid passing them around
type SNSAdapter struct {
	svc *sns.SNS
}

// NewSNSAdapter starts a new SNSAdapter instance
func NewSNSAdapter(awsRegion string) *SNSAdapter {
	session := session.New(&aws.Config{
		Region: aws.String(awsRegion),
	})
	return &SNSAdapter{svc: sns.New(session)}
}

// CreateEndpoint creates a new endpoint for this device
func (a *SNSAdapter) CreateEndpoint(applicationARN, customUserData, deviceToken string) (string, error) {
	// Call AWS to create a new endpoint
	params := &sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(applicationARN),
		Token:          aws.String(deviceToken),
		CustomUserData: aws.String(customUserData),
	}
	resp, err := a.svc.CreatePlatformEndpoint(params)
	if err != nil {
		return "", err
	}

	return *resp.EndpointArn, nil
}

// EndpointAttributes is a wrapper around important endpoint data
type EndpointAttributes struct {
	CustomUserData string
	Enabled        bool
	Token          string
}

// GetEndpointAttributes returns endpoint attributes (customUserData, enabled, token)
func (a *SNSAdapter) GetEndpointAttributes(endpointARN string) (*EndpointAttributes, error) {
	// Call AWS to get the endpoint attributes data
	params := &sns.GetEndpointAttributesInput{
		EndpointArn: aws.String(endpointARN),
	}
	resp, err := a.svc.GetEndpointAttributes(params)
	if err != nil {
		return nil, err
	}

	// Prepare variables to extract data from the attributes map (map[string]*string)
	var (
		customUserDataStr *string
		enabledStr        *string
		tokenStr          *string
		ok                bool
	)

	// Custom user data
	customUserDataStr, ok = resp.Attributes["CustomUserData"]
	if !ok {
		logger.Info(resp.Attributes)
		return nil, errors.New("CustomUserData key not found in attributes")
	}

	// Enabled
	enabledStr, ok = resp.Attributes["Enabled"]
	if !ok {
		logger.Info(resp.Attributes)
		return nil, errors.New("Enabled key not found in attributes")
	}

	// Token
	tokenStr, ok = resp.Attributes["Token"]
	if !ok {
		logger.Info(resp.Attributes)
		return nil, errors.New("Token key not found in attributes")
	}

	// Parse the enabled key from string to a boolean
	enabled, err := strconv.ParseBool(*enabledStr)
	if err != nil {
		logger.Info(resp.Attributes)
		return nil, errors.New("Enabled key could not be parsed into a boolean")
	}

	return &EndpointAttributes{
		CustomUserData: *customUserDataStr,
		Enabled:        enabled,
		Token:          *tokenStr,
	}, nil
}

// SetEndpointAttributes sets endpoint attributes (customUserData, enabled, token)
func (a *SNSAdapter) SetEndpointAttributes(endpointARN string, endpointAttributes *EndpointAttributes) error {
	// Call AWS to set the endpoint attributes data
	params := &sns.SetEndpointAttributesInput{
		Attributes: map[string]*string{
			"CustomUserData": aws.String(endpointAttributes.CustomUserData),
			"Enabled":        aws.String(fmt.Sprintf("%v", endpointAttributes.Enabled)),
			"Token":          aws.String(endpointAttributes.Token),
		},
		EndpointArn: aws.String(endpointARN),
	}
	_, err := a.svc.SetEndpointAttributes(params)
	if err != nil {
		return err
	}

	return nil
}

// PublishMessage sends a message to an endpoint and returns message ID
func (a *SNSAdapter) PublishMessage(endpointARN, msg string, opt map[string]interface{}) (string, error) {
	// Truncate the message
	msg = truncateMessage(msg)

	// Build a message string in an agnostic way to support both APNS and GCM
	m := make(map[string]string)
	m["default"] = msg

	// GCM
	gcmPayload, err := composeMessageGCM(msg, opt)
	if err != nil {
		return "", err
	}
	m["GCM"] = gcmPayload

	// APNS
	apnsPayload, err := composeMessageAPNS(msg, opt)
	if err != nil {
		return "", err
	}
	m["APNS"] = apnsPayload
	m["APNS_SANDBOX"] = apnsPayload

	// And put it all together
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	// Publish the message to SNS
	params := &sns.PublishInput{
		Message:          aws.String(string(msgBytes)), // Required
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(endpointARN),
	}
	resp, err := a.svc.Publish(params)
	if err != nil {
		return "", err
	}

	return *resp.MessageId, nil
}
