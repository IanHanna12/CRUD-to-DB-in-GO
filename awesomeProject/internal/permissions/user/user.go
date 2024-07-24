package user

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID        `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Username     string           `json:"username" gorm:"column:username"`
	Identifier   string           `json:"identifier" gorm:"column:identifier"`
	IdentifierID uuid.UUID        `json:"identifier_id" gorm:"column:identifier_id"`
	Permissions  DBPermissionView `gorm:"type:int"`
}

type PermissionView interface {
	CanView(isAdminContent bool) bool
	CanCreate() bool
	CanUpdate() bool
	CanDelete() bool
	CanViewAll() bool
}

type DBPermissionView struct {
	Permissions int
}

func (pv DBPermissionView) CanView(isAdminContent bool) bool {
	return pv.Permissions&1 != 0 && (!isAdminContent || pv.Permissions&16 != 0)
}
func (pv DBPermissionView) CanCreate() bool  { return pv.Permissions&2 != 0 }
func (pv DBPermissionView) CanUpdate() bool  { return pv.Permissions&4 != 0 }
func (pv DBPermissionView) CanDelete() bool  { return pv.Permissions&8 != 0 }
func (pv DBPermissionView) CanViewAll() bool { return pv.Permissions&16 != 0 }

func (pv DBPermissionView) Value() (driver.Value, error) {
	return pv.Permissions, nil
}

func (pv *DBPermissionView) Scan(value interface{}) error {
	if v, ok := value.(int64); ok {
		pv.Permissions = int(v)
		return nil
	}
	return fmt.Errorf("invalid type for DBPermissionView")
}

func GetView(identifier string) DBPermissionView {
	switch identifier {
	case "admin":
		return DBPermissionView{Permissions: 31}
	case "user":
		return DBPermissionView{Permissions: 7}
	default:
		return DBPermissionView{Permissions: 1}
	}
}
