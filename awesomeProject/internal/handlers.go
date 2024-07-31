package handlers

import (
	"encoding/json"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/db"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	var items []model.Item
	if err := DB.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	json.NewEncoder(w).Encode(responses)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
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
	json.NewEncoder(w).Encode(response)
}

func DeleteItemByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := db.DeleteItemByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item deleted successfully"))
}

func DeleteAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	if err := DB.Where("1 = 1").Delete(&model.Item{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All items deleted successfully"))
}

func CreateUser(username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Username: username,
		Password: string(hashedPassword),
		IsAdmin:  username == "admin",
	}

	return DB.Create(&user).Error
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := CreateUser(loginRequest.Username, loginRequest.Password); err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
			DB.Where("username = ?", loginRequest.Username).First(&user)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	isAdmin := user.IsAdmin
	redirectURL := "/static/user/user_view.html"
	if isAdmin {
		redirectURL = "/static/admin/admin_view.html"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"isAdmin":     isAdmin,
		"redirectURL": redirectURL,
	})
	log.Printf("User logged in: User=%s, IsAdmin=%v", user.Username, isAdmin)
}

func AuthMiddleware(adminRequired bool) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			isAdmin, _ := strconv.ParseBool(r.Header.Get("isAdmin"))
			if adminRequired && !isAdmin {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}

func PrefetchItemsHandler(w http.ResponseWriter, r *http.Request) {
	var items []model.Item
	if err := DB.Limit(1000).Order("created_at desc").Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	json.NewEncoder(w).Encode(responses)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	}
	http.FileServer(http.Dir("./frontend/static")).ServeHTTP(w, r)
}

func GetItemByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var item model.Item
	if err := DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
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
