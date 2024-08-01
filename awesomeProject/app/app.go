package app

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID         uuid.UUID `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Blogname   string    `json:"blogname" gorm:"column:blogname"`
	Author     string    `json:"author" gorm:"column:author"`
	Content    string    `json:"content" gorm:"column:content"`
	Permission string    `json:"permission" gorm:"column:permission"`
}

var DB *gorm.DB
