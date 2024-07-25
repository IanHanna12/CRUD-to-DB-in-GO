package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Username     string    `json:"username" gorm:"column:username"`
	Identifier   string    `json:"identifier" gorm:"column:identifier"`
	IdentifierID uuid.UUID `json:"identifier_id" gorm:"column:identifier_id"`
	IsAdmin      bool      `json:"is_admin" gorm:"column:is_admin"`
}

func (u *User) CanView(isAdminContent bool) bool {
	return u.IsAdmin || !isAdminContent
}

func (u *User) CanCreate() bool {
	return true // All users can create
}

func (u *User) CanUpdate() bool {
	return true // All users can update
}

func GetView(identifier string) *User {
	isAdmin := identifier == "admin"
	return &User{
		Identifier: identifier,
		IsAdmin:    isAdmin,
	}
}
