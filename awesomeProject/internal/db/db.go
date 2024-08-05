package db

import (
	"errors"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "root:abcd@tcp(127.0.0.1:3306)/test?allowNativePasswords=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.AutoMigrate(&model.Item{}, &model.User{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	DB = db
	return db
}

func CreateItem(item *model.Item) error {
	if item.Blogname == "" || item.Author == "" {
		return errors.New("blogname and author are required")
	}

	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	if err := DB.Create(item).Error; err != nil {
		return err
	}
	return nil
}

func GetAllItems() ([]model.Item, error) {
	var items []model.Item
	if err := DB.Find(&items).Error; err != nil {
		log.Printf("Error fetching items: %v", err)
		return nil, err
	}
	return items, nil
}

func GetItemByID(id uuid.UUID) (model.Item, error) {
	var item model.Item
	err := DB.First(&item, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No item found with id: %s", id)
		} else {
			log.Printf("Error fetching item by ID: %v", err)
		}
	}
	return item, err
}

func UpdateItem(item model.Item) error {
	if item.Blogname == "" || item.Author == "" {
		return errors.New("blogname and author are required")
	}

	var existingItem model.Item
	if err := DB.First(&existingItem, "id = ?", item.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item not found")
		}
		return err
	}

	existingItem.Blogname = item.Blogname
	existingItem.Author = item.Author
	existingItem.Content = item.Content

	if err := DB.Save(&existingItem).Error; err != nil {
		return err
	}
	return nil
}

func DeleteItemByID(id uuid.UUID) error {
	if err := DB.Where("id = ?", id).Delete(&model.Item{}).Error; err != nil {
		log.Printf("Error deleting item by ID: %v", err)
		return err
	}
	log.Printf("Item with ID %v successfully deleted", id)
	return nil
}

func DeleteAllItems() error {
	if err := DB.Exec("DELETE FROM items").Error; err != nil {
		log.Printf("Error deleting all items: %v", err)
		return err
	}
	return nil
}
