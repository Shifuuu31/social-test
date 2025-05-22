package utils

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func InsretCookie(db *sql.DB, userID int, sessionID string, expiresAt time.Time) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, sessionID, userID, expiresAt.Format(time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}

func CookieMaker(w http.ResponseWriter) string {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}

	sessionID := u.String()
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return sessionID
}

func ValidateCookie(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, err
	}

	sessionID := cookie.Value
	query := `SELECT user_id FROM sessions WHERE id = ? AND expires_at > datetime('now')`
	
	var userID int
	err = db.QueryRow(query, sessionID).Scan(&userID)
	if err != nil {
		log.Printf("Failed to validate session: %v", err)
		return 0, errors.New("invalid session")
	}
	
	return userID, nil
}

func IsLoggedIn(db *sql.DB, r *http.Request) int {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0
	}
	
	sessionID := cookie.Value
	query := `SELECT user_id FROM sessions WHERE id = ? AND expires_at > datetime('now')`
	
	var userID int
	err = db.QueryRow(query, sessionID).Scan(&userID)
	if err != nil {
		return 0
	}
	
	return userID
}

func Logout(db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil // Already logged out
	}
	
	sessionID := cookie.Value
	_, err = db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return err
	}
	
	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	
	return nil
}