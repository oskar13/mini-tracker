package accounts

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/oskar13/mini-tracker/pkg/db"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
	"golang.org/x/crypto/bcrypt"
)

func CheckLogin(w http.ResponseWriter, r *http.Request, userData webdata.User) bool {
	if userData.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}

	return true
}

func LoginUser(w http.ResponseWriter, r *http.Request, username string, givenPassword string) (webdata.User, error) {
	var user webdata.User

	q := "SELECT user_ID, username, profile_pic, disabled, password FROM users WHERE username = ?"

	err := db.DB.QueryRow(q, username).Scan(&user.UserID, &user.Username, &user.Cover, &user.Disabled, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return webdata.User{}, errors.New("no account found")
		}
		return webdata.User{}, err
	}

	if ComparePasswords(user.Password, givenPassword) {

		//set up session

		sessionToken := GenerateSession()

		UpdateSessionInDB(db.DB, sessionToken, user.UserID, true, webutils.ReadUserIP(r))

		http.SetCookie(w, &http.Cookie{
			Name:     "session-token",
			Value:    sessionToken,
			MaxAge:   60 * 30,
			HttpOnly: true,
			Secure:   false, //fix this in production
		})

		return user, nil
	}

	return webdata.User{}, errors.New("wrong password")

}

func LogOutUser(w http.ResponseWriter, user webdata.User) webdata.User {

	UpdateSessionInDB(db.DB, "null", user.UserID, false, "0")

	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	return webdata.User{}

}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	bytePlainPwd := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// save session in db
func UpdateSessionInDB(DB *sql.DB, Token string, userID int, login bool, ipAddr string) {
	expireTime := time.Time{}
	if login {
		expireTime = time.Now().Add(30 * time.Minute)
	}

	stmt, err := db.DB.Prepare(`UPDATE users SET session_uid = ?, session_expiry = ?, session_ip = ? WHERE user_ID = ?`)

	if err != nil {
		log.Printf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(Token, expireTime, ipAddr, userID)
	if err != nil {
		log.Printf("error storing session in database: %v", err)
	}
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

func ValidateSessionData(r *http.Request) (webdata.User, error) {

	var userData webdata.User

	cookie, err := r.Cookie("session-token")
	if err != nil {
		//log.Printf("Error retrieving session cookie: %v", err)
		//No cookie found, return empty
		return webdata.User{}, err
	}

	sessionID := cookie.Value
	//fmt.Println("Session-cookie:", sessionID)
	if sessionID == "" {
		log.Println("Session uid is empty in the cookie")
		return webdata.User{}, fmt.Errorf("session uid is empty")
	}

	q := "SELECT user_ID, admin_level, username, password, profile_pic, created, disabled, session_expiry, gender FROM users WHERE session_uid = ?"

	err = db.DB.QueryRow(q, sessionID).Scan(&userData.UserID, &userData.AdminLevel, &userData.Username, &userData.Password, &userData.Cover, &userData.Joined, &userData.Disabled, &userData.SessionExpiry, &userData.Gender)
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

	//fmt.Println("UserID validated:", userData.Username)

	return userData, nil
}

// Validated and returns user data
func GetUserData(r *http.Request) webdata.User {

	var userData webdata.User
	userData.LoggedIn = false

	userData, err := ValidateSessionData(r)
	if err != nil || userData.UserID == 0 {
		return userData
	}

	userData.LoggedIn = true

	return userData
}

func HashAndSalt(pwd string) (string, error) {
	pwd_byte := []byte(pwd)

	hash, err := bcrypt.GenerateFromPassword(pwd_byte, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(hash), nil
}

func CreateUser(username string, password string, password2 string, ref string, adminLevel int, ignoreRefCode bool) error {

	//

	if !ignoreRefCode {
		_, err := CheckRefCode(ref)

		if err != nil {
			return err
		}
	}

	if password != password2 {
		return errors.New("password mismatch")
	}

	usernameTaken, err := UsernameExists(username)

	if err != nil {
		return err
	}

	if usernameTaken {
		return errors.New("username taken")
	}

	hashedPwd, err := HashAndSalt(password)

	if err != nil {
		return err
	}

	q := `INSERT INTO users (users.username, users.password, users.admin_level) VALUES (?,?,?)`

	res, err := db.DB.Exec(q, username, hashedPwd, adminLevel)

	if err != nil {
		return err
	}

	insertedId, err := res.LastInsertId()

	if err != nil {
		return err
	}
	if !ignoreRefCode {
		err = UseRefCode(int(insertedId), ref)

		if err != nil {
			return err
		}
	}

	return nil
}

// Lets user update their password, checks if their old password is correct
func UserUpdatePassword(userID int, oldPassword string, newPassword string, newPassword2 string) error {

	if newPassword != newPassword2 {
		return errors.New("new passwords do not match")
	}

	var storedHash string
	q := "SELECT users.password FROM users WHERE LOWER(users.user_ID) = ?"

	err := db.DB.QueryRow(q, userID).Scan(&storedHash)
	if err != nil {
		return err
	}

	if ComparePasswords(storedHash, oldPassword) {
		SetPassword(userID, newPassword)
	} else {
		return errors.New("invalid password")
	}

	return nil
}

// Sets password with not checking the previous password
func SetPassword(userID int, password string) error {
	hashedPwd, err := HashAndSalt(password)
	if err != nil {
		return err
	}

	q := `UPDATE users SET users.password = ? WHERE users.user_id = ?`

	_, err = db.DB.Exec(q, hashedPwd, userID)

	if err != nil {
		return err
	}

	return nil
}

func UsernameExists(username string) (bool, error) {
	var user_ID int
	q := "SELECT users.user_ID FROM users WHERE LOWER(users.username) = ?"

	usernameLower := strings.ToLower(username)

	err := db.DB.QueryRow(q, usernameLower).Scan(&user_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Returns true if reference code is valid
func CheckRefCode(ref string) (int, error) {
	var user_ID int
	var invited_user *int
	q := "SELECT invites.inviting_user_ID, invites.invited_user_ID FROM invites WHERE invite_code = ?"

	err := db.DB.QueryRow(q, ref).Scan(&user_ID, &invited_user)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no ref code found in db -", ref)
			return 0, errors.New("no ref code found")
		}
		return 0, err
	}

	if invited_user != nil {
		return 0, errors.New("invite code already used")
	}

	return user_ID, nil
}

// Ties user ID to ref code marking it used
func UseRefCode(user_ID int, ref string) error {
	// UPDATE invites SET invites.invited_user_ID = 1 WHERE invites.invite_code = "asdf";
	q := `UPDATE invites SET invites.invited_user_ID = ? WHERE invites.invite_code = ?`

	_, err := db.DB.Exec(q, user_ID, ref)

	if err != nil {
		return err
	}

	return nil

}

// Creates invite code with inviting user ID
func CreateRefCode(user_ID int, ref string) error {

	q := `INSERT INTO minitorrent.invites (inviting_user_ID, invite_code,  invite_creation_date) VALUES (?,?, CURRENT_TIMESTAMP())`
	id := 0
	err := db.DB.QueryRow(q, user_ID, ref).Scan(&id)

	if err != nil {
		return err
	}
	return nil
}

func UpdateAccountFields(userData webdata.User) error {
	q := `UPDATE users SET users.email = ? ,users.tagline = ?, users.bio = ?,  users.gender = ? WHERE users.user_id = ?`

	fmt.Println(userData.Password)
	_, err := db.DB.Exec(q, userData.Email, userData.Tagline, userData.Bio, userData.Gender, userData.UserID)

	if err != nil {
		return err
	}

	return nil
}
