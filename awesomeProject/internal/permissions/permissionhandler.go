package permissions

import (
	"errors"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/model"
	"github.com/IanHanna/CRUD-to-DB-in-GO/internal/permissions/user"
	"net/http"
)

func CanView(u *user.User, isAdminContent bool) bool {
	if u.Permissions == nil {
		u.Permissions = user.GetView(u.Identifier)
	}
	return u.Permissions.CanView(isAdminContent)
}

func CanCreate(u *user.User) bool {
	if u.Permissions == nil {
		u.Permissions = user.GetView(u.Identifier)
	}
	return u.Permissions.CanCreate()
}

func CanUpdate(u *user.User) bool {
	if u.Permissions == nil {
		u.Permissions = user.GetView(u.Identifier)
	}
	return u.Permissions.CanUpdate()
}

func CanDelete(u *user.User) bool {
	if u.Permissions == nil {
		u.Permissions = user.GetView(u.Identifier)
	}
	return u.Permissions.CanDelete()
}

func CanViewAll(u *user.User) bool {
	if u.Permissions == nil {
		u.Permissions = user.GetView(u.Identifier)
	}
	return u.Permissions.CanViewAll()
}

func FilterItemsForUser(items []model.Item, user *user.User) []model.Item {
	var filteredItems []model.Item
	for _, item := range items {
		if CanView(user, false) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func GetAdminOnlyContent(user *user.User) (string, error) {
	if CanView(user, true) {
		return "This is admin-only content", nil
	}
	return "", errors.New("not authorized")
}

func UserView(user *user.User) (string, error) {
	if CanView(user, false) {
		return "This is user-level content", nil
	}
	return "", errors.New("not authorized")
}

func AdminView(user *user.User) (string, error) {
	if CanView(user, true) {
		return "This is admin-level content", nil
	}
	return "", errors.New("not authorized")
}

func CanViewHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item := r.Context().Value("item").(model.Item)
		permissionView := user.GetView(item.Author.Username)
		if !permissionView.CanView(false) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func CanCreateHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value("username").(string)
		permissionView := user.GetView(username)
		if !permissionView.CanCreate() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func CanUpdateHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item := r.Context().Value("item").(model.Item)
		permissionView := user.GetView(item.Author.Username)
		if !permissionView.CanUpdate() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func CanDeleteHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item := r.Context().Value("item").(model.Item)
		permissionView := user.GetView(item.Author.Username)
		if !permissionView.CanDelete() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func CanViewAllHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value("username").(string)
		permissionView := user.GetView(username)
		if !permissionView.CanViewAll() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
