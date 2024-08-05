package handlers

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

func SetSessionCookie(w http.ResponseWriter, userID uuid.UUID) (string, string) {
	sessionID := uuid.New().String()
	sessionToken := uuid.New().String()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	return sessionID, sessionToken
}

func GetUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return uuid.Nil, err
	}
	return retrieveUserIDFromSessionStorage(cookie.Value)
}

func retrieveUserIDFromSessionStorage(sessionID string) (uuid.UUID, error) {

	userID, err := uuid.Parse(sessionID)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func GetSessionIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
