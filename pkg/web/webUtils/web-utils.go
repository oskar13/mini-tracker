package webutils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func CheckLogin(w http.ResponseWriter, r *http.Request, userData webdata.User) bool {
	if userData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return false
	}

	return true
}

func GetUserData(r *http.Request, DB *sql.DB) webdata.User {

	var userData webdata.User

	userID, err := GetCurrentUserID(r)
	if err != nil || userID == 0 {
		userData.LoggedIn = false
		return userData
	}

	user, err := FetchUserByUserID(DB, userID)
	if err != nil || !IsSessionActive(DB, userID) {
		userData.LoggedIn = false
		return userData
	}

	userData.LoggedIn = true
	userData = user

	return userData
}

func GetCurrentUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session-token")
	if err != nil {
		log.Printf("Error retrieving session cookie: %v", err)
		return 0, err
	}

	sessionID := cookie.Value
	fmt.Println("Session-cookie:", sessionID)
	if sessionID == "" {
		log.Println("Session ID is empty in the cookie")
		return 0, fmt.Errorf("session ID is empty")
	}

	return ValidateSessionAndGetUserID(db.DB, sessionID)
}

func RenderTemplate(w http.ResponseWriter, templates []string, data interface{}) {

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func ValidateSessionAndGetUserID(DB *sql.DB, sessionID string) (int, error) {
	var userID int
	var expiresAt time.Time
	q := "SELECT userID, expiresAt FROM Sessions WHERE sessionID = ?"

	err := DB.QueryRow(q, sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Session ID not found in the database")
			return 0, fmt.Errorf("session not found")
		}
		return 0, err
	}

	if time.Now().After(expiresAt) {
		log.Println("Session has expired")
		return 0, fmt.Errorf("session expired")
	}

	fmt.Println("UserID validated:", userID)

	return userID, nil

}

func FetchUserByUserID(DB *sql.DB, userID int) (webdata.User, error) {
	var user webdata.User

	q := "SELECT userID, username, firstName, lastName, email, age, gender, profilePic FROM Users WHERE userID = ?"

	err := db.DB.QueryRow(q, userID).Scan(&user.UserID, &user.Username, &user.Email, &user.Cover)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
			return webdata.User{}, nil
		}
	}

	fmt.Println("Fetched user data:", user)
	return user, nil
}

func IsSessionActive(DB *sql.DB, userID int) bool {
	var expiresAt time.Time
	q := "SELECT expiresAT FROM Sessions WHERE userID = ? ORDER BY expiresAT DESC LIMIT 1"

	err := db.DB.QueryRow(q, userID).Scan(&expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("session inactive")
			return false
		}
	}
	return time.Now().Before(expiresAt)
}
