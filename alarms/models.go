package alarms

import (
	"database/sql"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/database"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// ResultParentTableName defines alarm parent results table name
const ResultParentTableName = "alarm_results"

// Region is a region from where alarm checks will be run
type Region struct {
	database.TimestampModel
	ID   string `gorm:"primary_key" sql:"type:varchar(20)"`
	Name string `sql:"type:varchar(50);unique;not null"`
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
	Results                []*Result
	EndpointURL            string      `sql:"type:varchar(254);not null"`
	ExpectedHTTPCode       uint        `sql:"default:200;not null"`
	Interval               uint        `sql:"default:60;not null"` // seconds
	EmailAlerts            bool        `sql:"default:false;index;not null"`
	PushNotificationAlerts bool        `sql:"default:false;index;not null"`
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
	Response       sql.NullString `sql:"type:text"`
	ErrorMessage   sql.NullString
	ResolvedAt     pq.NullTime `sql:"index"`
}

// TableName specifies table name
func (i *Incident) TableName() string {
	return "alarm_incidents"
}

// ResultSubTable keeps track of all result sub tables
type ResultSubTable struct {
	gorm.Model
	Name string `sql:"type:varchar(254);unique;not null"`
}

// TableName specifies table name
func (t *ResultSubTable) TableName() string {
	return "alarm_result_sub_tables"
}

// Result represents a parent table used to vertically partition results,
// sub tables will inherit from this table and split data by day
type Result struct {
	AlarmID     sql.NullInt64 `sql:"index;not null"`
	Alarm       *Alarm
	Timestamp   time.Time `sql:"index;not null"`
	RequestTime int64     // request time in nanoseconds
	Table       string    `sql:"-"` // ignore this field
}

// TableName specifies table name
func (r *Result) TableName() string {
	if r.Table == "" {
		return ResultParentTableName
	}
	return r.Table
}

// newAlarm creates new Alarm instance
func newAlarm(user *accounts.User, region *Region, alarmState *AlarmState, alarmRequest *AlarmRequest) *Alarm {
	userID := util.PositiveIntOrNull(int64(user.ID))
	regionID := util.StringOrNull(region.ID)
	alarmStateID := util.StringOrNull(alarmState.ID)
	alarm := &Alarm{
		UserID:                 userID,
		RegionID:               regionID,
		AlarmStateID:           alarmStateID,
		EndpointURL:            alarmRequest.EndpointURL,
		ExpectedHTTPCode:       alarmRequest.ExpectedHTTPCode,
		Interval:               alarmRequest.Interval,
		EmailAlerts:            alarmRequest.EmailAlerts,
		PushNotificationAlerts: alarmRequest.PushNotificationAlerts,
		Active:                 alarmRequest.Active,
	}
	if userID.Valid {
		alarm.User = user
	}
	if regionID.Valid {
		alarm.Region = region
	}
	if alarmStateID.Valid {
		alarm.AlarmState = alarmState
	}
	return alarm
}

// newIncident creates new Incident instance
func newIncident(alarm *Alarm, incidentType *IncidentType, resp *http.Response, errMsg string) *Incident {
	alarmID := util.PositiveIntOrNull(int64(alarm.ID))
	incidentTypeID := util.StringOrNull(incidentType.ID)
	incident := &Incident{
		AlarmID:        alarmID,
		IncidentTypeID: incidentTypeID,
		ErrorMessage:   util.StringOrNull(errMsg),
	}
	if alarmID.Valid {
		incident.Alarm = alarm
	}
	if incidentTypeID.Valid {
		incident.IncidentType = incidentType
	}

	// If the response is not nil
	if resp != nil {
		// Save the response status code
		incident.HTTPCode = util.IntOrNull(int64(resp.StatusCode))

		// Save the respnse dump
		var respDump string
		respBytes, err := httputil.DumpResponse(resp, true) // body = true
		if err != nil {
			respDump = string(respBytes)
		}
		incident.Response = util.StringOrNull(respDump)
	}

	return incident
}

// newResult creates new Result instance
func newResult(table string, alarm *Alarm, timestamp time.Time, requestTime int64) *Result {
	alarmID := util.PositiveIntOrNull(int64(alarm.ID))
	result := &Result{
		AlarmID:     alarmID,
		Timestamp:   timestamp,
		RequestTime: requestTime,
		Table:       table,
	}
	if alarmID.Valid {
		result.Alarm = alarm
	}
	return result
}
