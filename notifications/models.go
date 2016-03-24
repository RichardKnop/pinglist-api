package notifications

import (
	"database/sql"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

// Endpoint represents an AWS SNS platform endpoint
type Endpoint struct {
	gorm.Model
	UserID         sql.NullInt64 `sql:"index;not null"`
	User           *accounts.User
	ApplicationARN string `sql:"type:varchar(200);index"`
	ARN            string `sql:"type:varchar(200);unique;index"`
	DeviceToken    string `sql:"type:varchar(200)"`
	CustomUserData string `sql:"type:varchar(200)"`
	Enabled        bool
}

// TableName specifies table name
func (e *Endpoint) TableName() string {
	return "notification_endpoints"
}

// NewEndpoint creates new Endpoint instance
func NewEndpoint(user *accounts.User, applicationARN, arn, deviceToken, customUserData string, enabled bool) *Endpoint {
	userID := util.PositiveIntOrNull(int64(user.ID))
	endpoint := &Endpoint{
		UserID:         userID,
		ApplicationARN: applicationARN,
		ARN:            arn,
		DeviceToken:    deviceToken,
		CustomUserData: customUserData,
		Enabled:        enabled,
	}
	if userID.Valid {
		endpoint.User = user
	}
	return endpoint
}
