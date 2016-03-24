package metrics

import (
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetAccountsService() accounts.ServiceInterface
	PartitionRequestTime(parentTableName string, now time.Time) error
	RotateSubTables() error
	LogRequestTime(timestamp time.Time, referenceID uint, value int64) error
	PaginatedRequestTimesCount(referenceID uint) (int, error)
	FindPaginatedRequestTimes(offset, limit int, orderBy string, referenceID uint) ([]*RequestTime, error)

	// Needed for the newRoutes to be able to register handlers
	// TODO
}
