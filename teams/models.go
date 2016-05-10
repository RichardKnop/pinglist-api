package teams

import (
	"database/sql"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

// Team ...
type Team struct {
	gorm.Model
	OwnerID sql.NullInt64 `sql:"index;not null"`
	Owner   *accounts.User
	Name    string           `sql:"type:varchar(40);not null"`
	Members []*accounts.User `gorm:"many2many:team_team_members"`
}

// TableName specifies table name
func (t *Team) TableName() string {
	return "team_teams"
}

// NewTeam creates new Team instance
func NewTeam(owner *accounts.User, members []*accounts.User, name string) *Team {
	ownerID := util.PositiveIntOrNull(int64(owner.ID))
	team := &Team{
		OwnerID: ownerID,
		Name:    name,
		Members: members,
	}
	if ownerID.Valid {
		team.Owner = owner
	}
	return team
}
