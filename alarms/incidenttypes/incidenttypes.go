package incidenttypes

const (
	// SlowResponse - The request took unusually long time
	SlowResponse = "slow_response"
	// Timeout - The request timed out
	Timeout = "timeout"
	// BadCode - The request returned a response with a bad status code
	BadCode = "bad_code"
	// Other - Any other request error
	Other = "other"
)
