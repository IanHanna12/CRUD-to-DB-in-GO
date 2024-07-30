package handlers

import (
	"context"
	"encoding/json"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

type Session struct {
	UserID    string
	Username  string
	IsAdmin   bool
	ExpiresAt time.Time
}

var sessions = make(map[string]Session)

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

	w.Header().Set("Content-Type", "application/json")
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
	sessionID := r.Header.Get("sessionID")
	session, exists := sessions[sessionID]
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var item model.Item
	if err := DB.First(&item, id).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if !session.IsAdmin && item.Author != session.Username {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	response := ItemResponse{
		ID:       item.ID,
		Blogname: item.Blogname,
		Author:   item.Author,
		Content:  item.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateitemHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("sessionID")
	session, exists := sessions[sessionID]
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var existingItem model.Item
	if err := DB.First(&existingItem, id).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if !session.IsAdmin && existingItem.Author != session.Username {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var itemRequest struct {
		Blogname string `json:"blogname"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&itemRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingItem.Blogname = itemRequest.Blogname
	existingItem.Content = itemRequest.Content

	if err := DB.Save(&existingItem).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ItemResponse{
		ID:       existingItem.ID,
		Blogname: existingItem.Blogname,
		Author:   existingItem.Author,
		Content:  existingItem.Content,
	}

	w.Header().Set("Content-Type", "application/json")
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
	w.Write([]byte("Item deleted successfully"))
}

func DeleteallitemsHandler(w http.ResponseWriter, r *http.Request) {
	if err := DB.Where("1 = 1").Delete(&model.Item{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All items deleted successfully"))
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
	log.Println("LoginHandler: Started")
	w.Header().Set("Content-Type", "application/json")

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Printf("LoginHandler: Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	log.Printf("LoginHandler: Received login request for username: %s", loginRequest.Username)

	isAdmin := loginRequest.Username == "admin"
	var hashedPassword []byte
	if isAdmin {
		hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	} else {
		hashedPassword = []byte(loginRequest.Password)
	}

	if (isAdmin && bcrypt.CompareHashAndPassword(hashedPassword, []byte(loginRequest.Password)) == nil) ||
		(!isAdmin && loginRequest.Password != "") {

		sessionID := uuid.New().String()
		sessions[sessionID] = Session{
			UserID:    uuid.New().String(),
			Username:  loginRequest.Username,
			IsAdmin:   isAdmin,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		redirectURL := "/static/main_page/mainpage.html"
		if isAdmin {
			redirectURL = "/static/admin/admin_view.html"
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"isAdmin":     isAdmin,
			"sessionID":   sessionID,
			"message":     "Login successful",
			"redirectURL": redirectURL,
		})
		log.Printf("LoginHandler: Login successful for user: %s", loginRequest.Username)
	} else {
		log.Println("LoginHandler: Invalid credentials")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid credentials",
		})
	}
}

func WithAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("sessionID")
		session, exists := sessions[sessionID]
		if !exists || time.Now().After(session.ExpiresAt) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func WithAdminAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("sessionID")
		session, exists := sessions[sessionID]
		if !exists || time.Now().After(session.ExpiresAt) || !session.IsAdmin {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func CleanupSessions() {
	for id, session := range sessions {
		if time.Now().After(session.ExpiresAt) {
			delete(sessions, id)
		}
	}
}
