package alarms

import (
	"database/sql"
	"net/http"
	"net/http/httputil"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/database"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Region is a region from where alarm checks will be run
type Region struct {
	database.TimestampModel
	ID       string         `gorm:"primary_key" sql:"type:varchar(20)"`
	Name     string         `sql:"type:varchar(50);unique;not null"`
	ProxyURL sql.NullString `sql:"type:varchar(200)"`
}

// TableName specifies table name
func (r *Region) TableName() string {
	return "alarm_regions"
}

// AlarmState is a state that an alarm can be in
type AlarmState struct {
	database.TimestampModel
	ID   string `gorm:"primary_key" sql:"type:varchar(20)"`
	Name string `sql:"type:varchar(50);unique;not null"`
}

// TableName specifies table name
func (s *AlarmState) TableName() string {
	return "alarm_states"
}

// Alarm ...
type Alarm struct {
	gorm.Model
	UserID                 sql.NullInt64 `sql:"index;not null"`
	User                   *accounts.User
	RegionID               sql.NullString `sql:"type:varchar(20);index;not null"`
	Region                 *Region
	AlarmStateID           sql.NullString `sql:"type:varchar(20);index;not null"`
	AlarmState             *AlarmState
	Incidents              []*Incident
	EndpointURL            string      `sql:"type:varchar(254);not null"`
	ExpectedHTTPCode       uint        `sql:"default:200;not null"`
	MaxResponseTime        uint        `sql:"default:60;not null"` // miliseconds
	Interval               uint        `sql:"default:60;not null"` // seconds
	EmailAlerts            bool        `sql:"default:false;index;not null"`
	PushNotificationAlerts bool        `sql:"default:false;index;not null"`
	SlackAlerts            bool        `sql:"default:false;index;not null"`
	Active                 bool        `sql:"index;not null"`
	Watermark              pq.NullTime `sql:"index"`
	LastDowntimeStartedAt  pq.NullTime `sql:"index"`
	LastUptimeStartedAt    pq.NullTime `sql:"index"`
}

// TableName specifies table name
func (a *Alarm) TableName() string {
	return "alarm_alarms"
}

// IncidentType ...
type IncidentType struct {
	database.TimestampModel
	ID   string `gorm:"primary_key"`
	Name string `sql:"type:varchar(50);unique;not null"`
}

// TableName specifies table name
func (t *IncidentType) TableName() string {
	return "alarm_incident_types"
}

// Incident ...
type Incident struct {
	gorm.Model
	AlarmID        sql.NullInt64  `sql:"index;not null"`
	IncidentTypeID sql.NullString `sql:"index;not null"`
	Alarm          *Alarm
	IncidentType   *IncidentType
	HTTPCode       sql.NullInt64
	ResponseTime   sql.NullInt64  // nanoseconds
	Response       sql.NullString `sql:"type:text"`
	ErrorMessage   sql.NullString
	ResolvedAt     pq.NullTime `sql:"index"`
}

// TableName specifies table name
func (i *Incident) TableName() string {
	return "alarm_incidents"
}

// NotificationCounter ...
type NotificationCounter struct {
	gorm.Model
	UserID sql.NullInt64 `sql:"index;not null"`
	User   *accounts.User
	Year   uint `sql:"index;not null"`
	Month  uint `sql:"index;not null"`
	Email  uint `sql:"default:0;not null"`
	Push   uint `sql:"default:0;not null"`
	Slack  uint `sql:"default:0;not null"`
}

// TableName specifies table name
func (i *NotificationCounter) TableName() string {
	return "alarm_notification_counters"
}

// NewAlarm creates new Alarm instance
func NewAlarm(user *accounts.User, region *Region, alarmState *AlarmState, alarmRequest *AlarmRequest) *Alarm {
	userID := util.PositiveIntOrNull(int64(user.ID))
	regionID := util.StringOrNull(region.ID)
	alarmStateID := util.StringOrNull(alarmState.ID)
	alarm := &Alarm{
		UserID:                 userID,
		RegionID:               regionID,
		AlarmStateID:           alarmStateID,
		EndpointURL:            alarmRequest.EndpointURL,
		ExpectedHTTPCode:       alarmRequest.ExpectedHTTPCode,
		MaxResponseTime:        alarmRequest.MaxResponseTime,
		Interval:               alarmRequest.Interval,
		EmailAlerts:            alarmRequest.EmailAlerts,
		PushNotificationAlerts: alarmRequest.PushNotificationAlerts,
		SlackAlerts:            alarmRequest.SlackAlerts,
		Active:                 alarmRequest.Active,
	}
	return alarm
}

// NewIncident creates new Incident instance
func NewIncident(alarm *Alarm, incidentType *IncidentType, resp *http.Response, responseTime int64, errMsg string) *Incident {
	alarmID := util.PositiveIntOrNull(int64(alarm.ID))
	incidentTypeID := util.StringOrNull(incidentType.ID)
	incident := &Incident{
		AlarmID:        alarmID,
		IncidentTypeID: incidentTypeID,
		ErrorMessage:   util.StringOrNull(errMsg),
	}

	// If the response is not nil
	if resp != nil {
		// Save the response status code
		incident.HTTPCode = util.IntOrNull(int64(resp.StatusCode))

		// Save the response time
		incident.ResponseTime = util.IntOrNull(responseTime)

		// Save the response dump
		var respDump string
		respBytes, err := httputil.DumpResponse(
			resp,
			false, // body
		)
		if err == nil {
			respDump = string(respBytes)
		}
		incident.Response = util.StringOrNull(respDump)
	}

	return incident
}
