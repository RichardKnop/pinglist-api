package metrics

import (
	"time"

	"github.com/jinzhu/gorm"
)

// RequestTimeParentTableName defines request time parent table name
const RequestTimeParentTableName = "metrics_request_times"

// SubTable keeps track of all result sub tables
type SubTable struct {
	gorm.Model
	ParentTable string `sql:"type:varchar(254);index;not null"`
	Name        string `sql:"type:varchar(254);unique;not null"`
}

// TableName specifies table name
func (t *SubTable) TableName() string {
	return "metrics_sub_tables"
}

// RequestTime represents a parent table used to vertically partition request times,
// sub tables will inherit from this table and split data by day
type RequestTime struct {
	ReferenceID uint      `sql:"index;not null"`
	Timestamp   time.Time `sql:"index;not null"`
	Value       int64     // request time in nanoseconds
	Table       string    `sql:"-"` // ignore this field
}

// TableName specifies table name
func (r *RequestTime) TableName() string {
	if r.Table == "" {
		return RequestTimeParentTableName
	}
	return r.Table
}

// NewRequestTime creates new RequestTime instance
func NewRequestTime(table string, referenceID uint, timestamp time.Time, value int64) *RequestTime {
	return &RequestTime{
		ReferenceID: referenceID,
		Timestamp:   timestamp,
		Value:       value,
		Table:       table,
	}
}
