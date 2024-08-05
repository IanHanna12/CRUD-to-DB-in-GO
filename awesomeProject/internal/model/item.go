package model

import (
	"github.com/google/uuid"
)

type Item struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Blogname string    `json:"blogname" gorm:"column:blogname"`
	Author   string    `json:"author" gorm:"column:author"`
	Content  string    `json:"content" gorm:"column:content"`
	UserID   uuid.UUID `json:"userID" gorm:"column:user_id;type:char(36)"`
}

type User struct {
	ID        uuid.UUID `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Username  string    `json:"username" gorm:"column:username"`
	Password  string    `json:"password" gorm:"column:password"`
	IsAdmin   bool      `json:"isAdmin" gorm:"column:isAdmin"`
	Session   string    `json:"session" gorm:"column:session"`
	SessionID string    `json:"session_id" gorm:"column:session_id"`
}
type ItemResponse struct {
	ID       uuid.UUID `json:"id"`
	Blogname string    `json:"blogname"`
	Author   string    `json:"author"`
	Content  string    `json:"content"`
}
