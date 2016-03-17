package alarms

import (
	"database/sql"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// ResultParentTableName defines parent results table name
const ResultParentTableName = "alarm_results"

// Alarm ...
type Alarm struct {
	gorm.Model
	UserID           sql.NullInt64 `sql:"index;not null"`
	User             *accounts.User
	Incidents        []*Incident
	Results          []*Result
	Watermark        pq.NullTime `sql:"index"`
	EndpointURL      string      `sql:"type:varchar(254);not null"`
	ExpectedHTTPCode uint        `sql:"default:200;not null"`
	Interval         uint        `sql:"default:60;not null"` // seconds
	Active           bool        `sql:"index;not null"`
	State            string      `sql:"type:varchar(20);index;not null"`
}

// TableName specifies table name
func (a *Alarm) TableName() string {
	return "alarm_alarms"
}

// Incident ...
type Incident struct {
	gorm.Model
	AlarmID    sql.NullInt64 `sql:"index;not null"`
	Alarm      *Alarm
	Type       string `sql:"type:varchar(20);index;not null"`
	HTTPCode   sql.NullInt64
	Response   sql.NullString `sql:"type:text"`
	ResolvedAt pq.NullTime    `sql:"index"`
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
func newAlarm(user *accounts.User, alarmRequest *AlarmRequest) *Alarm {
	userID := util.PositiveIntOrNull(int64(user.ID))
	alarm := &Alarm{
		UserID:           userID,
		EndpointURL:      alarmRequest.EndpointURL,
		ExpectedHTTPCode: alarmRequest.ExpectedHTTPCode,
		Interval:         alarmRequest.Interval,
		Active:           alarmRequest.Active,
		State:            alarmstates.InsufficientData,
	}
	if userID.Valid {
		alarm.User = user
	}
	return alarm
}

// newIncident creates new Incident instance
func newIncident(alarm *Alarm, theType string, resp *http.Response) *Incident {
	alarmID := util.PositiveIntOrNull(int64(alarm.ID))
	incident := &Incident{
		AlarmID: alarmID,
		Type:    theType,
	}
	if alarmID.Valid {
		incident.Alarm = alarm
	}

	// If the response is not nil
	if resp != nil {
		// Save the response status code
		incident.HTTPCode = util.IntOrNull(int64(resp.StatusCode))

		// Save the respnse dump
		var respDump string
		respBytes, err := httputil.DumpResponse(resp, false) // body = false
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
