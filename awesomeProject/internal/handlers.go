package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

var isAdminMode bool

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:8080",
	DB:   0,
})

var DB *gorm.DB

type ItemResponse struct {
	ID       uuid.UUID `json:"id"`
	Blogname string    `json:"blogname"`
	Author   string    `json:"author"`
	Content  string    `json:"content"`
}

func InitHandlers(db *gorm.DB) {
	DB = db
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var itemRequest struct {
		Blogname string `json:"blogname"`
		Author   string `json:"author"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&itemRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := model.Item{
		ID:       uuid.New(),
		Blogname: itemRequest.Blogname,
		Author:   itemRequest.Author,
		Content:  itemRequest.Content,
	}

	if err := DB.Create(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ItemResponse{
		ID:       item.ID,
		Blogname: item.Blogname,
		Author:   item.Author,
		Content:  item.Content,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAllItemsHandler: Started")

	var items []model.Item
	if err := DB.Find(&items).Error; err != nil {
		log.Printf("GetAllItemsHandler: Database error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GetAllItemsHandler: Found %d items", len(items))

	var responses []ItemResponse
	for _, item := range items {
		responses = append(responses, ItemResponse{
			ID:       item.ID,
			Blogname: item.Blogname,
			Author:   item.Author,
			Content:  item.Content,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responses); err != nil {
		log.Printf("GetAllItemsHandler: JSON encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("GetAllItemsHandler: Completed successfully")
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

	response := ItemResponse{
		ID:       item.ID,
		Blogname: item.Blogname,
		Author:   item.Author,
		Content:  item.Content,
	}

	responseJSON, err := json.Marshal(response)
	if err == nil {
		log.Println("GetItemByIDHandler: Caching item")
		rdb.Set(ctx, cacheKey, responseJSON, 10*time.Minute)
	}

	log.Println("GetItemByIDHandler: Sending response")
	w.Write(responseJSON)
	log.Println("GetItemByIDHandler: Response sent")
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var itemRequest struct {
		Blogname string `json:"blogname"`
		Author   string `json:"author"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&itemRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := model.Item{
		ID:       id,
		Blogname: itemRequest.Blogname,
		Author:   itemRequest.Author,
		Content:  itemRequest.Content,
	}

	if err := DB.Save(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cacheKey := "item:" + id.String()
	rdb.Del(ctx, cacheKey)

	response := ItemResponse{
		ID:       item.ID,
		Blogname: item.Blogname,
		Author:   item.Author,
		Content:  item.Content,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Vary", "Origin")

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if loginRequest.Username != "" && loginRequest.Password != "" {
		isAdminMode = loginRequest.IsAdmin
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"isAdmin": isAdminMode,
			"message": "Login successful",
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid credentials",
		})
	}
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "frontend/static/main_page/mainPage.js")
	} else if strings.HasSuffix(r.URL.Path, "user_view.html") {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		http.ServeFile(w, r, "frontend/static/user/user_view.html")
	} else {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		http.ServeFile(w, r, "awesomeProject/frontend/static/main_page/mainPage.html")
	}
}

func ServeWithProperMIME(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, ".html") {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	} else if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	}

	http.ServeFile(w, r, "frontend/static"+path)
}
