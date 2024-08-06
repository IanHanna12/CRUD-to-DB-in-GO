package handlers

import (
	"encoding/json"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"net/http"
)

var DB *gorm.DB

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
			// Create a new user
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
				// if username is admin --> isAdmin = true
				IsAdmin: credentials.Username == "admin",
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
		// Verify password for existing user
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			http.Error(responseWriter, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	sessionID, sessionToken := SetSessionCookie(responseWriter, user.ID)

	user.SessionID = sessionID
	if err := DB.Save(&user).Error; err != nil {
		http.Error(responseWriter, "Error saving session", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"success":       true,
		"message":       "Logged in successfully",
		"isAdmin":       user.IsAdmin,
		"session_id":    sessionID,
		"session_token": sessionToken,
	})
}

func GetAllItemsHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	var items []model.Item
	if err := DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		http.Error(responseWriter, "Error fetching items", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(items)
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

	var item model.Item
	if err := DB.Where("id = ? AND user_id = ?", itemID, userID).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(responseWriter, "Item not found", http.StatusNotFound)
		} else {
			http.Error(responseWriter, "Database error", http.StatusInternalServerError)
		}
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(item)
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

	responseWriter.WriteHeader(http.StatusNoContent)
}

func DeleteAllItemsHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	if err := DB.Where("user_id = ?", userID).Delete(&model.Item{}).Error; err != nil {
		http.Error(responseWriter, "Error deleting items", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}

func PrefetchItemsHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	authReq := request.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	var items []model.Item
	if err := DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		http.Error(responseWriter, "Error fetching items", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(map[string]interface{}{
		"prefetchedItems": items,
	})
}

func ValidateSessionHandler(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	userID, err := GetUserIDFromSessionCookie(request)
	if err != nil {
		http.Error(responseWriter, "Invalid session", http.StatusUnauthorized)
		return
	}

	var user model.User
	if err := DB.First(&user, userID).Error; err != nil {
		http.Error(responseWriter, "Invalid session", http.StatusUnauthorized)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(map[string]bool{"valid": true})
}

func AuthMiddleware(adminRequired bool) func(httprouter.Handle) httprouter.Handle {
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
			// Add the auth request to the context, pass it to the next handler, and then move it along
			authenticatedContext := context.WithValue(request.Context(), "authRequest", authReq)
			next(responseWriter, request.WithContext(authenticatedContext), params)
		}
	}
}
