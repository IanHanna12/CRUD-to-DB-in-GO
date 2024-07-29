package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

var adminSessions = make(map[string]bool)

// unhashed pwd for admin = admin
const adminpwdHash = "$2a$10$XWN1bGzK5Y5JZw.Qx9Yl6O5tLtq5jZ1bJ1tQ5Yl6O5tLtq5jZ1bJ1"

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
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

func GetallitemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getallitemsHandler: Started")

	var items []model.Item
	if err := DB.Find(&items).Error; err != nil {
		log.Printf("getallitemsHandler: Database error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("getallitemsHandler: Found %d items", len(items))

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
		log.Printf("getallitemsHandler: JSON encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("getallitemsHandler: Completed successfully")
}

func GetitembyIDHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("getitembyIDHandler: Started")

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("getitembyIDHandler: Invalid UUID: %s", idStr)
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	log.Printf("getitembyIDHandler: Parsed UUID: %s", id)
	cacheKey := "item:" + id.String()
	cachedItem, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Println("getitembyIDHandler: Cache hit, returning cached item")
		w.Write([]byte(cachedItem))
		return
	}
	log.Println("getitembyIDHandler: Cache miss, querying database")

	var item model.Item
	result := DB.First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("getitembyIDHandler: Item not found for ID: %s", id)
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			log.Printf("getitembyIDHandler: Database error: %v", result.Error)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	log.Printf("getitembyIDHandler: Item found: %+v", item)

	response := ItemResponse{
		ID:       item.ID,
		Blogname: item.Blogname,
		Author:   item.Author,
		Content:  item.Content,
	}

	responseJSON, err := json.Marshal(response)
	if err == nil {
		log.Println("getitembyIDHandler: Caching item")
		rdb.Set(ctx, cacheKey, responseJSON, 10*time.Minute)
	}

	log.Println("getitembyIDHandler: Sending response")
	w.Write(responseJSON)
	log.Println("getitembyIDHandler: Response sent")
}

func UpdateitemHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteitembyIDHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteallitemsHandler(w http.ResponseWriter, r *http.Request) {
	if err := DB.Where("1 = 1").Delete(&model.Item{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func MainpageHandler(w http.ResponseWriter, r *http.Request) {
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

func ServewithproperMIME(w http.ResponseWriter, r *http.Request) {
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Vary", "Origin")

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	isAdmin := loginRequest.Username == "admin" && loginRequest.Password == adminpwdHash

	if loginRequest.Username != "" && loginRequest.Password != "" {
		sessionID := ""
		if isAdmin {
			sessionID = uuid.New().String()
			adminSessions[sessionID] = true
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"isAdmin":   isAdmin,
			"sessionID": sessionID,
			"message":   "Login successful",
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid credentials",
		})
	}
}

func WithadminAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("Admin-Session-ID")
		if sessionID == "" {
			http.Error(w, "Unauthorized: No session ID provided", http.StatusUnauthorized)
			return
		}
		//if sessionID not in map {
		if !adminSessions[sessionID] {
			http.Error(w, "Unauthorized: Invalid session ID", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
