package metrics

import (
	"time"

	"github.com/jinzhu/gorm"
)

// ResponseTimeParentTableName defines request time parent table name
const ResponseTimeParentTableName = "metrics_response_times"

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

// ResponseTime represents a parent table used to vertically partition request times,
// sub tables will inherit from this table and split data by day
type ResponseTime struct {
	ReferenceID uint      `sql:"index;not null"`
	Timestamp   time.Time `sql:"index;not null"`
	Value       int64     // request time in nanoseconds
	Table       string    `sql:"-"` // ignore this field
}

// TableName specifies table name
func (r *ResponseTime) TableName() string {
	if r.Table == "" {
		return ResponseTimeParentTableName
	}
	return r.Table
}

// NewResponseTime creates new ResponseTime instance
func NewResponseTime(table string, referenceID uint, timestamp time.Time, value int64) *ResponseTime {
	return &ResponseTime{
		ReferenceID: referenceID,
		Timestamp:   timestamp,
		Value:       value,
		Table:       table,
	}
}
