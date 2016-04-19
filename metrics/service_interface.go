package metrics

import (
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetAccountsService() accounts.ServiceInterface
	PartitionResponseTime(parentTableName string, now time.Time) error
	RotateSubTables() error
	LogResponseTime(timestamp time.Time, referenceID uint, value int64) error
	PaginatedResponseTimesCount(referenceID int, dateTrunc string, from, to *time.Time) (int, error)
	FindPaginatedResponseTimes(offset, limit int, orderBy string, referenceID int, dateTrunc string, from, to *time.Time) ([]*ResponseTime, error)
}
