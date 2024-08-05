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

func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new user
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return
			}
			user = model.User{
				ID:        uuid.New(),
				Username:  credentials.Username,
				Password:  string(hashedPassword),
				SessionID: "",
			}
			if err := DB.Create(&user).Error; err != nil {
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		// Verify password for existing user
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	sessionID, sessionToken := SetSessionCookie(w, user.ID)

	user.SessionID = sessionID
	if err := DB.Save(&user).Error; err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"message":       "Logged in successfully",
		"isAdmin":       user.IsAdmin,
		"session_id":    sessionID,
		"session_token": sessionToken,
	})
}

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authReq := r.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	var items []model.Item
	if err := DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		http.Error(w, "Error fetching items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authReq := r.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item.ID = uuid.New()
	item.UserID = userID

	if err := DB.Create(&item).Error; err != nil {
		http.Error(w, "Error creating item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func GetItemByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authReq := r.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID
	itemID, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var item model.Item
	if err := DB.Where("id = ? AND user_id = ?", itemID, userID).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authReq := r.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID
	itemID, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var updatedItem model.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingItem model.Item
	if err := DB.Where("id = ? AND user_id = ?", itemID, userID).First(&existingItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	existingItem.Blogname = updatedItem.Blogname
	existingItem.Author = updatedItem.Author
	existingItem.Content = updatedItem.Content

	if err := DB.Save(&existingItem).Error; err != nil {
		http.Error(w, "Error updating item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingItem)
}

func DeleteItemByIDHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authReq := r.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID
	itemID, err := uuid.Parse(ps.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	if err := DB.Where("id = ? AND user_id = ?", itemID, userID).Delete(&model.Item{}).Error; err != nil {
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteAllItemsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authReq := r.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	if err := DB.Where("user_id = ?", userID).Delete(&model.Item{}).Error; err != nil {
		http.Error(w, "Error deleting items", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func PrefetchItemsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authReq := r.Context().Value("authRequest").(*AuthenticatedRequest)
	userID := authReq.User.ID

	var items []model.Item
	if err := DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		http.Error(w, "Error fetching items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"prefetchedItems": items,
	})
}

func ValidateSessionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, err := GetUserIDFromSessionCookie(r)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var user model.User
	if err := DB.First(&user, userID).Error; err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"valid": true})
}

func AuthMiddleware(adminRequired bool) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			sessionID, err := GetSessionIDFromCookie(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var user model.User
			if err := DB.Where("session_id = ?", sessionID).First(&user).Error; err != nil {
				http.Error(w, "Invalid session", http.StatusUnauthorized)
				return
			}

			if adminRequired && !user.IsAdmin {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			authReq := &AuthenticatedRequest{
				Request: r,
				User:    user,
			}
			ctx := context.WithValue(r.Context(), "authRequest", authReq)
			next(w, r.WithContext(ctx), ps)
		}
	}
}
