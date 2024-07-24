package model

import (
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/permissions/user"
	"github.com/google/uuid"
)

type Item struct {
	ID           uuid.UUID `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Blogname     string    `json:"blogname" gorm:"column:blogname"`
	Author       user.User `json:"author" gorm:"embedded"`
	Content      string    `json:"content" gorm:"column:content"`
	Identifier   string    `json:"identifier" gorm:"column:identifier"`
	IdentifierID uuid.UUID `json:"identifier_id" gorm:"column:identifier_id"`
	Permissions  user.PermissionView
}

func (i Item) GetPermissions() user.PermissionView {
	return i.Permissions
}
