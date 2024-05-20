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
	userData.LoggedIn = false

	userData, err := ValidateSessionData(r)
	if err != nil || userData.UserID == 0 {
		return userData
	}

	userData.LoggedIn = true

	return userData
}

func ValidateSessionData(r *http.Request) (webdata.User, error) {

	var userData webdata.User

	cookie, err := r.Cookie("session-token")
	if err != nil {
		log.Printf("Error retrieving session cookie: %v", err)
		return webdata.User{}, err
	}

	sessionID := cookie.Value
	fmt.Println("Session-cookie:", sessionID)
	if sessionID == "" {
		log.Println("Session uid is empty in the cookie")
		return webdata.User{}, fmt.Errorf("session uid is empty")
	}

	q := "SELECT user_ID, username, profile_pic, disabled, session_expiry FROM users WHERE session_uid = ?"

	err = db.DB.QueryRow(q, sessionID).Scan(&userData.UserID, &userData.Username, &userData.Cover, &userData.Disabled, &userData.SessionExpiry)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Session uid not found in the database")
			return webdata.User{}, fmt.Errorf("session not found")
		}
		return webdata.User{}, err
	}

	layout := "2006-01-02 15:04:05.000000"
	expiry, _ := time.Parse(layout, userData.SessionExpiry)

	if time.Now().After(expiry) {
		log.Println("Session has expired")
		return webdata.User{}, fmt.Errorf("session expired")
	}

	fmt.Println("UserID validated:", userData.Username)

	return userData, nil
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
