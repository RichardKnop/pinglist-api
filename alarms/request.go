package alarms

// AlarmRequest ...
type AlarmRequest struct {
	Region           string `json:"region"`
	EndpointURL      string `json:"endpoint_url"`
	ExpectedHTTPCode uint   `json:"expected_http_code"`
	Interval         uint   `json:"interval"`
	Active           bool   `json:"active"`
}
