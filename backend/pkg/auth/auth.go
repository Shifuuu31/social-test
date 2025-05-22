package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	models "social-network/pkg/models"
	utils "social-network/pkg/utils"
	sqlite "social-network/pkg/db/sqlite"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var user models.Logininfo
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	db := sqlite.GetDB()

	var hashedPassword string
	var userID int
	err = db.QueryRow("SELECT password, id FROM users WHERE nickname = ? OR email = ?", user.Email, user.Email).Scan(&hashedPassword, &userID)
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "User not found or incorrect email")
		return
	}

	err = utils.CheckPassword(hashedPassword, user.Password)
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	_, err = db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("Error deleting old sessions:", err)
	}

	cookie := utils.CookieMaker(w)
	err = utils.InsretCookie(db, userID, cookie, time.Now().Add(24*time.Hour))
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "Failed to insert cookie")
		return
	}

	utils.JsonResponse(w, http.StatusOK, "Login successful")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JsonResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newUser models.RegisterInfo
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if newUser.Email == "" || newUser.Password == "" || newUser.Nickname == "" {
		utils.JsonResponse(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	db := sqlite.GetDB()

	// Check if email or nickname already exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR nickname = ?", newUser.Email, newUser.Nickname).Scan(&exists)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "Database error")
		return
	}
	if exists > 0 {
		utils.JsonResponse(w, http.StatusConflict, "Email or nickname already in use")
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "Password hashing failed")
		return
	}

	res, err := db.Exec(
		"INSERT INTO users (email, password, nickname, created_at) VALUES (?, ?, ?, ?)",
		newUser.Email, hashedPassword, newUser.Nickname, time.Now(),
	)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "Error inserting user")
		return
	}

	userID, err := res.LastInsertId()
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "Could not retrieve user ID")
		return
	}

	cookie := utils.CookieMaker(w)
	if err := utils.InsretCookie(db, int(userID), cookie, time.Now().Add(24*time.Hour)); err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "Failed to set session cookie")
		return
	}

	utils.JsonResponse(w, http.StatusOK, "Registration successful")
}
