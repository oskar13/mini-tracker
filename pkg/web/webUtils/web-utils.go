package webutils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
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

	layout := "2006-01-02 15:04:05"
	expiry, _ := time.Parse(layout, userData.SessionExpiry)

	if time.Now().After(expiry) {
		log.Println("Session has expired")
		log.Println("Time: ", userData.SessionExpiry)
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

func LoginUser(w http.ResponseWriter, username string, password string) (error, webdata.User) {
	var user webdata.User
	var salt string

	q := "SELECT user_ID, username, profile_pic, disabled, password, salt FROM users WHERE username = ?"

	err := db.DB.QueryRow(q, username).Scan(&user.UserID, &user.Username, &user.Cover, &user.Disabled, &user.Password, &salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no account found"), webdata.User{}
		}
		return err, webdata.User{}
	}

	if user.Password == password {

		//set up session

		sessionToken := GenerateSession()

		UpdateSessionInDB(db.DB, sessionToken, user.UserID, true)

		http.SetCookie(w, &http.Cookie{
			Name:     "session-token",
			Value:    sessionToken,
			MaxAge:   60 * 30,
			HttpOnly: true,
			Secure:   true,
		})

		return nil, user
	}

	return errors.New("wrong password"), webdata.User{}

}

func LogOutUser(w http.ResponseWriter, user webdata.User) webdata.User {

	UpdateSessionInDB(db.DB, "null", user.UserID, false)

	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	return webdata.User{}

}

func GenerateSession() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		log.Fatalf("Error generating session token: %v", err)
	}
	token := hex.EncodeToString(b)
	return token
}

// save session in db
func UpdateSessionInDB(DB *sql.DB, Token string, userID int, login bool) {
	expireTime := time.Time{}
	if login {
		expireTime = time.Now().Add(30 * time.Minute)
	}

	stmt, err := db.DB.Prepare(`UPDATE users SET session_uid = ?, session_expiry = ? WHERE user_ID = ?`)

	if err != nil {
		log.Printf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(Token, expireTime, userID)
	if err != nil {
		log.Printf("error storing session in database: %v", err)
		log.Printf(Token)
	}
}
