package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/permissions/user"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:8080",
	DB:   0,
})

var DB *gorm.DB

func InitHandlers(db *gorm.DB) {
	DB = db
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := DB.Create(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	var items []model.Item
	if err := DB.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func GetItemByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetItemByIDHandler: Started")

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("GetItemByIDHandler: Invalid UUID: %s", idStr)
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	log.Printf("GetItemByIDHandler: Parsed UUID: %s", id)

	cacheKey := "item:" + id.String()

	cachedItem, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Println("GetItemByIDHandler: Cache hit, returning cached item")
		w.Write([]byte(cachedItem))
		return
	}
	log.Println("GetItemByIDHandler: Cache miss, querying database")

	var item model.Item
	result := DB.First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("GetItemByIDHandler: Item not found for ID: %s", id)
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			log.Printf("GetItemByIDHandler: Database error: %v", result.Error)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	log.Printf("GetItemByIDHandler: Item found: %+v", item)

	permissionView := user.GetView(item.Author.Username)
	if !permissionView.CanView(false) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	itemJSON, err := json.Marshal(item)
	if err == nil {
		log.Println("GetItemByIDHandler: Caching item")
		rdb.Set(ctx, cacheKey, itemJSON, 10*time.Minute)
	}

	log.Println("GetItemByIDHandler: Sending response")
	w.Write(itemJSON)
	log.Println("GetItemByIDHandler: Response sent")
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.ID = id
	if err := DB.Save(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cacheKey := "item:" + id.String()
	rdb.Del(ctx, cacheKey)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func DeleteItemByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := DB.Delete(&model.Item{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cacheKey := "item:" + id.String()
	rdb.Del(ctx, cacheKey)

	w.WriteHeader(http.StatusOK)
}

func DeleteAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	if err := DB.Where("1 = 1").Delete(&model.Item{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
