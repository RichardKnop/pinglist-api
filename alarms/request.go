package alarms

// AlarmRequest ...
type AlarmRequest struct {
	Region                 string `json:"region"`
	EndpointURL            string `json:"endpoint_url"`
	ExpectedHTTPCode       uint   `json:"expected_http_code"`
	MaxResponseTime        uint   `json:"max_response_time"`
	Interval               uint   `json:"interval"`
	EmailAlerts            bool   `json:"email_alerts"`
	PushNotificationAlerts bool   `json:"push_notification_alerts"`
	Active                 bool   `json:"active"`
}
