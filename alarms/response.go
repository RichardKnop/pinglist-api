package alarms

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/pinglist-api/metrics"
)

// RegionResponse ...
type RegionResponse struct {
	jsonhal.Hal
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListRegionsResponse ...
type ListRegionsResponse struct {
	jsonhal.Hal
	Count uint `json:"count"`
	Page  uint `json:"page"`
}

// AlarmResponse ...
type AlarmResponse struct {
	jsonhal.Hal
	ID                     uint   `json:"id"`
	UserID                 uint   `json:"user_id"`
	Region                 string `json:"region"`
	EndpointURL            string `json:"endpoint_url"`
	ExpectedHTTPCode       uint   `json:"expected_http_code"`
	MaxResponseTime        uint   `json:"max_response_time"`
	Interval               uint   `json:"interval"`
	EmailAlerts            bool   `json:"email_alerts"`
	PushNotificationAlerts bool   `json:"push_notification_alerts"`
	Active                 bool   `json:"active"`
	State                  string `json:"state"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
}

// ListAlarmsResponse ...
type ListAlarmsResponse struct {
	jsonhal.Hal
	Count uint `json:"count"`
	Page  uint `json:"page"`
}

// IncidentResponse ...
type IncidentResponse struct {
	jsonhal.Hal
	ID           uint    `json:"id"`
	AlarmID      uint    `json:"alarm_id"`
	Type         string  `json:"type"`
	HTTPCode     *uint   `json:"http_code"`
	Response     *string `json:"response"`
	ResponseTime *uint   `json:"response_time"`
	ErrorMessage *string `json:"error_message"`
	ResolvedAt   *string `json:"resolved_at"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// ListIncidentsResponse ...
type ListIncidentsResponse struct {
	jsonhal.Hal
	Count uint `json:"count"`
	Page  uint `json:"page"`
}

// ListResponseTimesResponse ...
type ListResponseTimesResponse struct {
	jsonhal.Hal
	Uptime             float64        `json:"uptime"`
	Average            float64        `json:"average"`
	IncidentTypeCounts map[string]int `json:"incident_type_counts"`
	Count              uint           `json:"count"`
	Page               uint           `json:"page"`
}

// NewRegionResponse creates new ResultResponse instance
func NewRegionResponse(region *Region) (*RegionResponse, error) {
	response := &RegionResponse{
		ID:   region.ID,
		Name: region.Name,
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/alarms/regions/%s", region.ID), // href
		"", // title
	)

	return response, nil
}

// NewListRegionsResponse creates new ListRegionsResponse instance
func NewListRegionsResponse(count, page int, self, first, last, previous, next string, regions []*Region) (*ListRegionsResponse, error) {
	response := &ListRegionsResponse{
		Count: uint(count),
		Page:  uint(page),
	}

	// Set the self link
	response.SetLink("self", self, "")

	// Set the first link
	response.SetLink("first", first, "")

	// Set the last link
	response.SetLink("last", last, "")

	// Set the previous link
	response.SetLink("prev", previous, "")

	// Set the next link
	response.SetLink("next", next, "")

	// Create slice of region responses
	regionResponses := make([]*RegionResponse, len(regions))
	for i, region := range regions {
		regionResponse, err := NewRegionResponse(region)
		if err != nil {
			return nil, err
		}
		regionResponses[i] = regionResponse
	}

	// Set embedded regions
	response.SetEmbedded(
		"regions",
		jsonhal.Embedded(regionResponses),
	)

	return response, nil
}

// NewAlarmResponse creates new AlarmResponse instance
func NewAlarmResponse(alarm *Alarm) (*AlarmResponse, error) {
	response := &AlarmResponse{
		ID:                     alarm.ID,
		UserID:                 uint(alarm.UserID.Int64),
		Region:                 alarm.RegionID.String,
		EndpointURL:            alarm.EndpointURL,
		ExpectedHTTPCode:       alarm.ExpectedHTTPCode,
		MaxResponseTime:        alarm.MaxResponseTime,
		Interval:               alarm.Interval,
		EmailAlerts:            alarm.EmailAlerts,
		PushNotificationAlerts: alarm.PushNotificationAlerts,
		Active:                 alarm.Active,
		State:                  alarm.AlarmStateID.String,
		CreatedAt:              alarm.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:              alarm.UpdatedAt.UTC().Format(time.RFC3339),
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/alarms/%d", alarm.ID), // href
		"", // title
	)

	return response, nil
}

// NewListAlarmsResponse creates new ListAlarmsResponse instance
func NewListAlarmsResponse(count, page int, self, first, last, previous, next string, alarms []*Alarm) (*ListAlarmsResponse, error) {
	response := &ListAlarmsResponse{
		Count: uint(count),
		Page:  uint(page),
	}

	// Set the self link
	response.SetLink("self", self, "")

	// Set the first link
	response.SetLink("first", first, "")

	// Set the last link
	response.SetLink("last", last, "")

	// Set the previous link
	response.SetLink("prev", previous, "")

	// Set the next link
	response.SetLink("next", next, "")

	// Create slice of alarm responses
	alarmResponses := make([]*AlarmResponse, len(alarms))
	for i, alarm := range alarms {
		alarmResponse, err := NewAlarmResponse(alarm)
		if err != nil {
			return nil, err
		}
		alarmResponses[i] = alarmResponse
	}

	// Set embedded alarms
	response.SetEmbedded(
		"alarms",
		jsonhal.Embedded(alarmResponses),
	)

	return response, nil
}

// NewIncidentResponse creates new IncidentResponse instance
func NewIncidentResponse(incident *Incident) (*IncidentResponse, error) {
	response := &IncidentResponse{
		ID:        incident.ID,
		AlarmID:   uint(incident.AlarmID.Int64),
		Type:      incident.IncidentTypeID.String,
		CreatedAt: incident.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: incident.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if incident.HTTPCode.Valid {
		httpCode := uint(incident.HTTPCode.Int64)
		response.HTTPCode = &httpCode
	}
	if incident.Response.Valid {
		r := incident.Response.String
		response.Response = &r
	}
	if incident.ResponseTime.Valid {
		rt := uint(incident.ResponseTime.Int64)
		response.ResponseTime = &rt
	}
	if incident.ErrorMessage.Valid {
		e := incident.ErrorMessage.String
		response.ErrorMessage = &e
	}
	if incident.ResolvedAt.Valid {
		resolvedAt := incident.ResolvedAt.Time.UTC().Format(time.RFC3339)
		response.ResolvedAt = &resolvedAt
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf(
			"/v1/alarms/%d/incidents/%d",
			incident.AlarmID.Int64,
			incident.ID,
		), // href
		"", // title
	)

	return response, nil
}

// NewListIncidentsResponse creates new ListIncidentsResponse instance
func NewListIncidentsResponse(count, page int, self, first, last, previous, next string, incidents []*Incident) (*ListIncidentsResponse, error) {
	response := &ListIncidentsResponse{
		Count: uint(count),
		Page:  uint(page),
	}

	// Set the self link
	response.SetLink("self", self, "")

	// Set the first link
	response.SetLink("first", first, "")

	// Set the last link
	response.SetLink("last", last, "")

	// Set the previous link
	response.SetLink("prev", previous, "")

	// Set the next link
	response.SetLink("next", next, "")

	// Create slice of incident responses
	incidentResponses := make([]*IncidentResponse, len(incidents))
	for i, incident := range incidents {
		incidentResponse, err := NewIncidentResponse(incident)
		if err != nil {
			return nil, err
		}
		incidentResponses[i] = incidentResponse
	}

	// Set embedded incidents
	response.SetEmbedded(
		"incidents",
		jsonhal.Embedded(incidentResponses),
	)

	return response, nil
}

// NewListResponseTimesResponse creates new ListResponseTimesResponse instance
func NewListResponseTimesResponse(count, page int, self, first, last, previous, next string, responseTimes []*metrics.ResponseTime, uptime float64, incidentTypeCount map[string]int) (*ListResponseTimesResponse, error) {
	// Format uptime to 4 decimal numbers
	uptime, err := strconv.ParseFloat(fmt.Sprintf("%.4f", uptime), 64)
	if err != nil {
		return nil, err
	}

	response := &ListResponseTimesResponse{
		Uptime: uptime,
		Count:  uint(count),
		Page:   uint(page),
	}

	// Set the self link
	response.SetLink("self", self, "")

	// Set the first link
	response.SetLink("first", first, "")

	// Set the last link
	response.SetLink("last", last, "")

	// Set the previous link
	response.SetLink("prev", previous, "")

	// Set the next link
	response.SetLink("next", next, "")

	// Create slice of metrics responses
	metricResponses := make([]*metrics.MetricResponse, len(responseTimes))
	for i, responseTime := range responseTimes {
		metricResponse, err := metrics.NewMetricResponse(
			responseTime.Timestamp,
			responseTime.Value,
		)
		if err != nil {
			return nil, err
		}
		metricResponses[i] = metricResponse
	}

	// Average response time
	if len(responseTimes) > 0 {
		var sum float64
		for _, responseTime := range responseTimes {
			sum += float64(responseTime.Value)
		}
		response.Average = sum / float64(len(responseTimes))
	}

	// Aggregate incident counts based on type
	response.IncidentTypeCounts = incidentTypeCount

	// Set embedded response times
	response.SetEmbedded(
		"response_times",
		jsonhal.Embedded(metricResponses),
	)

	return response, nil
}
