package user

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Username     string    `json:"username" gorm:"column:username"`
	Identifier   string    `json:"identifier" gorm:"column:identifier"`
	IdentifierID uuid.UUID `json:"identifier_id" gorm:"column:identifier_id"`
	Role         string    `json:"role" gorm:"column:role"`
}

func (u *User) canView(isAdminContent bool) bool {
	return u.Role == "admin" || !isAdminContent
}

func (u *User) canCreate() bool {
	return u.Role == "admin" || u.Role == "user"
}

func (u *User) canUpdate() bool {
	return u.Role == "admin" || u.Role == "user"
}

func (u *User) canDelete() bool {
	return u.Role == "admin"
}

func GetView(identifier string) *User {
	role := "guest"
	if identifier == "admin" {
		role = "admin"
	} else if identifier == "user" {
		role = "user"
	}
	return &User{
		Identifier: identifier,
		Role:       role,
	}
}
