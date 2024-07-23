package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/db"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.CreateItem(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items, err := db.GetAllItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func GetItemByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	// Try to get the item from Redis cache
	cacheKey := "item:" + id.String()
	cachedItem, err := rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		// Cache miss, get the item from the database
		item, err := db.GetItemByID(id)
		if err != nil {
			if err.Error() == "record not found" {
				http.Error(w, "Item not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Cache the item in Redis
		itemJSON, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rdb.Set(ctx, cacheKey, itemJSON, 10*time.Minute)

		w.WriteHeader(http.StatusOK)
		w.Write(itemJSON)
	} else if err != nil {
		// Some other error occurred
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// Cache hit, return the cached item
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedItem))
	}
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
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
	if err := db.UpdateItem(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Invalidate the cache
	cacheKey := "item:" + id.String()
	rdb.Del(ctx, cacheKey)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func DeleteItemByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := db.DeleteItemByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Invalidate the cache
	cacheKey := "item:" + id.String()
	rdb.Del(ctx, cacheKey)

	w.WriteHeader(http.StatusOK)
}

func DeleteAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := db.DeleteAllItems(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
