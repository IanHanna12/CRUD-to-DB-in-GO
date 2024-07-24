package user

import (
	"github.com/google/uuid"
	"strings"
)

type User struct {
	ID           uuid.UUID `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Username     string    `json:"username" gorm:"column:username"`
	Identifier   string    `json:"identifier" gorm:"column:identifier"`
	IdentifierID uuid.UUID `json:"identifier_id" gorm:"column:identifier_id"`
	Permissions  PermissionView
}

type PermissionView interface {
	CanView(isAdminContent bool) bool
	CanCreate() bool
	CanUpdate() bool
	CanDelete() bool
	CanViewAll() bool
}

// AdminView with full permissions
type AdminView struct{}

func (v AdminView) CanView(bool) bool { return true }
func (v AdminView) CanCreate() bool   { return true }
func (v AdminView) CanUpdate() bool   { return true }
func (v AdminView) CanDelete() bool   { return true }
func (v AdminView) CanViewAll() bool  { return true }

// UserView  with limited permissions
type UserView struct{}

func (v UserView) CanView(isAdminContent bool) bool { return !isAdminContent }
func (v UserView) CanCreate() bool                  { return true }
func (v UserView) CanUpdate() bool                  { return false }
func (v UserView) CanDelete() bool                  { return false }
func (v UserView) CanViewAll() bool                 { return false }

// DefaultView  no permissions
type DefaultView struct{}

func (v DefaultView) CanView(isAdminContent bool) bool { return false }
func (v DefaultView) CanCreate() bool                  { return false }
func (v DefaultView) CanUpdate() bool                  { return false }
func (v DefaultView) CanDelete() bool                  { return false }
func (v DefaultView) CanViewAll() bool                 { return false }

func GetView(identifier string) PermissionView {
	switch strings.ToLower(identifier) {
	case "admin":
		return AdminView{}
	case "user":
		return UserView{}
	default:
		return DefaultView{}
	}
}
