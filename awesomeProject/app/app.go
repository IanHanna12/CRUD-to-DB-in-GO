package app

import (
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Item struct {
	ID       uuid.UUID `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Blogname string    `json:"blogname" gorm:"column:blogname"`
	Author   string    `json:"author" gorm:"column:author"`
	Content  string    `json:"content" gorm:"column:content"`
}

var DB *gorm.DB // Database connection

func InitDB() {
	var err error
	dsn := "root:abcd@tcp(127.0.0.1:3306)/test?allowNativePasswords=true"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = DB.AutoMigrate(&Item{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
}
