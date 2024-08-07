package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var DB *gorm.DB
var RedisClient *redis.Client

type ItemResponse struct {
	ID       uuid.UUID `json:"id"`
	Blogname string    `json:"blogname"`
	Author   string    `json:"author"`
	Content  string    `json:"content"`
}

type AuthenticatedRequest struct {
	*http.Request
	User model.User
}

func InitHandlers(db *gorm.DB) {
	DB = db
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func LoginHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&credentials); err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(responseWriter, "Error creating user", http.StatusInternalServerError)
				return
			}
			user = model.User{
				ID:        uuid.New(),
				Username:  credentials.Username,
				Password:  string(hashedPassword),
				SessionID: "",
				IsAdmin:   credentials.Username == "admin",
			}
			if err := DB.Create(&user).Error; err != nil {
				http.Error(responseWriter, "Error creating user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(responseWriter, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			http.Error(responseWriter, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	sessionID, AuthToken := SetSessionCookie(responseWriter, user.ID)

	user.SessionID = sessionID
	if err := DB.Save(&user).Error; err != nil {
		http.Error(responseWriter, "Error saving session", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"success":    true,
		"message":    "Logged in successfully",
		"isAdmin":    user.IsAdmin,
		"session_id": sessionID,
		"authtoken":  AuthToken,
	})
}

// validate session by comparing sessionID with userID retrieved from cookie
func ValidateSessionHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userID, err := GetUserIDFromSessionCookie(request)
	if err != nil {
		http.Error(responseWriter, "Invalid session", http.StatusUnauthorized)
		return
	}

	cacheKey := fmt.Sprintf("user:%s", userID)
	cachedUser, err := RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Write([]byte(cachedUser))
		return
	}

	var user model.User
	if err := DB.First(&user, userID).Error; err != nil {
		http.Error(responseWriter, "Invalid session", http.StatusUnauthorized)
		return
	}

	userJSON, _ := json.Marshal(map[string]bool{"valid": true})
	RedisClient.Set(context.Background(), cacheKey, userJSON, 30*time.Minute)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(userJSON)
}

func AuthMiddleware(adminRequired bool) func(httprouter.Handle) httprouter.Handle {
	//use authenticated request (user or admin logged in) as a wrapper around actual request
	return func(next httprouter.Handle) httprouter.Handle {
		return func(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
			sessionID, err := GetSessionIDFromCookie(request)
			if err != nil {
				http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var user model.User
			if err := DB.Where("session_id = ?", sessionID).First(&user).Error; err != nil {
				http.Error(responseWriter, "Invalid session", http.StatusUnauthorized)
				return
			}

			if adminRequired && !user.IsAdmin {
				http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
				return
			}

			authReq := &AuthenticatedRequest{
				Request: request,
				User:    user,
			}
			authenticatedContext := context.WithValue(request.Context(), "authRequest", authReq)
			next(responseWriter, request.WithContext(authenticatedContext), params)
		}
	}
}

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authReq := r.Context().Value("authRequest").(*AuthenticatedRequest)

	if !authReq.User.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cacheKey := "all_items"
	cachedItems, err := RedisClient.Get(context.Background(), cacheKey).Result()

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedItems))
		return
	}

	var items []model.Item
	if err := DB.Find(&items).Error; err != nil {
		http.Error(w, "Error fetching items", http.StatusInternalServerError)
		return
	}

	itemsJSON, _ := json.Marshal(items)
	RedisClient.Set(context.Background(), cacheKey, itemsJSON, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(itemsJSON)
}

func CreateItemHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	var item model.Item
	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	item.ID = uuid.New()
	item.UserID = userID

	if err := DB.Create(&item).Error; err != nil {
		http.Error(responseWriter, "Error creating item", http.StatusInternalServerError)
		return
	}

	RedisClient.Del(context.Background(), fmt.Sprintf("user:%s:items", userID))
	RedisClient.Del(context.Background(), "all_items")

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseWriter).Encode(item)
}

func GetItemByIDHandler(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID
	itemID, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		http.Error(responseWriter, "Invalid item ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("item:%s", itemID)
	cachedItem, err := RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Write([]byte(cachedItem))
		return
	}

	var item model.Item
	if err := DB.Where("id = ? AND user_id = ?", itemID, userID).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(responseWriter, "Item not found", http.StatusNotFound)
		} else {
			http.Error(responseWriter, "Database error", http.StatusInternalServerError)
		}
		return
	}

	itemJSON, _ := json.Marshal(item)
	RedisClient.Set(context.Background(), cacheKey, itemJSON, 5*time.Minute)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(itemJSON)
}

func UpdateItemHandler(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID
	itemID, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		http.Error(responseWriter, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var updatedItem model.Item
	if err := json.NewDecoder(request.Body).Decode(&updatedItem); err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingItem model.Item
	if err := DB.Where("id = ? AND user_id = ?", itemID, userID).First(&existingItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(responseWriter, "Item not found", http.StatusNotFound)
		} else {
			http.Error(responseWriter, "Database error", http.StatusInternalServerError)
		}
		return
	}

	existingItem.Blogname = updatedItem.Blogname
	existingItem.Author = updatedItem.Author
	existingItem.Content = updatedItem.Content

	if err := DB.Save(&existingItem).Error; err != nil {
		http.Error(responseWriter, "Error updating item", http.StatusInternalServerError)
		return
	}

	RedisClient.Del(context.Background(), fmt.Sprintf("item:%s", itemID))
	RedisClient.Del(context.Background(), fmt.Sprintf("user:%s:items", userID))
	RedisClient.Del(context.Background(), "all_items")

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(existingItem)
}

func UpdateItemHandlerForAdmin(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	itemID, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		http.Error(responseWriter, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var updatedItem model.Item
	if err := json.NewDecoder(request.Body).Decode(&updatedItem); err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingItem model.Item
	if err := DB.First(&existingItem, itemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(responseWriter, "Item not found", http.StatusNotFound)
		} else {
			http.Error(responseWriter, "Database error", http.StatusInternalServerError)
		}
		return
	}

	existingItem.Blogname = updatedItem.Blogname
	existingItem.Author = updatedItem.Author
	existingItem.Content = updatedItem.Content

	if err := DB.Save(&existingItem).Error; err != nil {
		http.Error(responseWriter, "Error updating item", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(existingItem)
}

func DeleteItemByIDHandler(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID
	itemID, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		http.Error(responseWriter, "Invalid item ID", http.StatusBadRequest)
		return
	}

	if err := DB.Where("id = ? AND user_id = ?", itemID, userID).Delete(&model.Item{}).Error; err != nil {
		http.Error(responseWriter, "Error deleting item", http.StatusInternalServerError)
		return
	}

	RedisClient.Del(context.Background(), fmt.Sprintf("item:%s", itemID))
	RedisClient.Del(context.Background(), fmt.Sprintf("user:%s:items", userID))
	RedisClient.Del(context.Background(), "all_items")

	responseWriter.WriteHeader(http.StatusNoContent)
}

func DeleteAllItemsHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	if err := DB.Where("user_id = ?", userID).Delete(&model.Item{}).Error; err != nil {
		http.Error(responseWriter, "Error deleting items", http.StatusInternalServerError)
		return
	}

	RedisClient.Del(context.Background(), fmt.Sprintf("user:%s:items", userID))
	RedisClient.Del(context.Background(), "all_items")

	responseWriter.WriteHeader(http.StatusNoContent)
}

func PrefetchItemsHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	cacheKey := fmt.Sprintf("user:%s:items", userID)
	cachedItems, err := RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Write([]byte(cachedItems))
		return
	}

	var items []model.Item
	if err := DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		http.Error(responseWriter, "Error fetching items", http.StatusInternalServerError)
		return
	}

	itemsJSON, _ := json.Marshal(map[string]interface{}{"prefetchedItems": items})
	RedisClient.Set(context.Background(), cacheKey, itemsJSON, 5*time.Minute)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(itemsJSON)
	return
}

func PrefetchAllItemsHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	var items []model.Item
	if err := DB.Find(&items).Error; err != nil {
		http.Error(responseWriter, "Error fetching items", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"prefetchedItems": items,
	})
}
