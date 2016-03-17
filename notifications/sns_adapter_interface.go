package notifications

// SNSAdapterInterface defines exported methods
type SNSAdapterInterface interface {
	// Exported methods
	CreateEndpoint(applicationARN, customUserData, deviceToken string) (string, error)
	GetEndpointAttributes(endpointARN string) (*EndpointAttributes, error)
	SetEndpointAttributes(endpointARN string, endpointAttributes *EndpointAttributes) error
	PublishMessage(endpointARN, message string, opt map[string]interface{}) (string, error)
}
