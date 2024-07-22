package app

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"sync"
)

type Item struct {
	ID       uuid.UUID `json:"id,omitempty" gorm:"type:char(36);primaryKey"`
	Blogname string    `json:"blogname" gorm:"column:blogname"`
	Author   string    `json:"author" gorm:"column:author"`
}

var DB *gorm.DB // Database connection
var mu sync.Mutex

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

func CreateItem(item *Item) error {
	if item.Blogname == "" || item.Author == "" {
		return errors.New("blogname and author are required")
	}

	if item.ID == uuid.Nil {
		item.ID = uuid.New() // Generate a new UUID if the ID is empty
	}
	if err := DB.Create(item).Error; err != nil {
		return err
	}
	return nil
}

func GetAllItems() ([]Item, error) {
	var items []Item
	if err := DB.Find(&items).Error; err != nil {
		log.Printf("Error fetching items: %v", err)
		return nil, err
	}
	return items, nil
}

func GetItemByID(id uuid.UUID) (Item, error) {
	var item Item
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

func UpdateItem(item Item) error {
	if item.Blogname == "" || item.Author == "" {
		return errors.New("blogname and author are required")
	}

	if err := DB.Save(&item).Error; err != nil {
		return err
	}
	return nil
}

func DeleteItemByID(id uuid.UUID) error {
	if err := DB.Where("id = ?", id).Delete(&Item{}).Error; err != nil {
		log.Printf("Error deleting item by ID: %v", err)
		return err
	}
	log.Printf("Item with ID %v successfully deleted", id)
	return nil
}

func DeleteAllItems() error {
	if err := DB.Exec("DELETE FROM items").Error; err != nil { // Delete all items from the database
		log.Printf("Error deleting all items: %v", err)
		return err
	}
	return nil
}

func HandleCreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		httpError(w, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.Blogname == "" || item.Author == "" {
		httpError(w, "blogname and author are required", http.StatusBadRequest)
		return
	}

	if err := CreateItem(&item); err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		log.Printf("Error encoding response JSON: %v", err)
	}
}

func HandleGetAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	items, err := GetAllItems()
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("Error encoding response JSON: %v", err)
	}
}

func HandleGetItemByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, err := GetItemByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpError(w, "Item not found", http.StatusNotFound)
		} else {
			httpError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(item); err != nil {
		log.Printf("Error encoding response JSON: %v", err)
	}
}

func HandleUpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		httpError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if item.Blogname == "" || item.Author == "" {
		httpError(w, "blogname and author are required", http.StatusBadRequest)
		return
	}

	// Generate a new UUID
	newUUID := uuid.New()
	item.ID = newUUID

	if err := UpdateItem(item); err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		log.Printf("Error encoding response JSON: %v", err)
	}
}

func HandleDeleteItemByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := DeleteItemByID(id); err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleDeleteAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	if err := DeleteAllItems(); err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func httpError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
	log.Printf("HTTP error: %v, Code: %d", message, code)
}
