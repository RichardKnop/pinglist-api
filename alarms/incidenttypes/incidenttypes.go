package incidenttypes

const (
	// Slow - The request took unusually long time
	Slow = "slow"
	// Timeout - The request timed out
	Timeout = "timeout"
	// BadCode - The request returned a response with a bad status code
	BadCode = "bad code"
	// Other - Any other request error
	Other = "other"
)
