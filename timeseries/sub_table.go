package timeseries

import (
	"fmt"
	"time"
)

// getSubTableName generates a new sub table name for a given date
func getSubTableName(parentName string, now time.Time) string {
	return fmt.Sprintf("%s_%s", parentName, now.UTC().Format("2006_01_02"))
}
