package webutils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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

	q := "SELECT user_ID, username, profile_pic, created, disabled, session_expiry, gender FROM users WHERE session_uid = ?"

	err = db.DB.QueryRow(q, sessionID).Scan(&userData.UserID, &userData.Username, &userData.Cover, &userData.Joined, &userData.Disabled, &userData.SessionExpiry, &userData.Gender)
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

func LoginUser(w http.ResponseWriter, r *http.Request, username string, password string) (error, webdata.User) {
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

		UpdateSessionInDB(db.DB, sessionToken, user.UserID, true, ReadUserIP(r))

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

	UpdateSessionInDB(db.DB, "null", user.UserID, false, "0")

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

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func CreateUser(username string, password string, password2 string, ref string) error {

	//

	_, err := CheckRefCode(ref)

	if err != nil {
		return err
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

	q := `INSERT INTO users (users.username, users.salt, users.password) VALUES (?,?,?)`

	salt := password

	res, err := db.DB.Exec(q, username, salt, password)

	if err != nil {
		return err
	}

	inserteId, err := res.LastInsertId()

	if err != nil {
		return err
	}

	err = UseRefCode(int(inserteId), ref)

	if err != nil {
		return err
	}

	return nil
}

func UsernameExists(username string) (bool, error) {
	var user_ID int
	q := "SELECT users.user_ID FROM users WHERE users.username = ?"

	err := db.DB.QueryRow(q, username).Scan(&user_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Returns true if refcode is valid
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

// Load user data for profile page or profile cards
func LoadUserProfileData(user_ID int) (webdata.User, error) {

	var user webdata.User

	q := `SELECT users.username, users.profile_pic, users.banner_image, users.created, users.disabled, users.tagline, users.bio, users.gender FROM users WHERE user_ID = ?`
	err := db.DB.QueryRow(q, user_ID).Scan(&user.Username, &user.Cover, &user.Banner, &user.Joined, &user.Disabled, &user.Tagline, &user.Bio, &user.Gender)
	if err != nil {
		if err == sql.ErrNoRows {
			return webdata.User{}, errors.New("no account found")
		}
		return webdata.User{}, err
	}

	user.UserID = user_ID

	return user, nil
}

// Load public torrents in user profile view or elsewhere, set flag to true to
func LoadUserTorrents(user_ID int, access_type []string) []webdata.TorrentWeb {

	var resultTorrents []webdata.TorrentWeb

	if len(access_type) == 0 {
		return []webdata.TorrentWeb{}
	}

	q := `SELECT torrents.torrent_ID, torrents.created, torrents.name, torrents.upvotes, torrents.downvotes
	      FROM torrents
	      WHERE torrents.users_user_ID = ? AND torrents.access_type IN (` + strings.Repeat("?,", len(access_type)-1) + `?)`

	args := make([]interface{}, 0, len(access_type)+1)
	args = append(args, user_ID)
	for _, at := range access_type {
		args = append(args, at)
	}

	rows, err := db.DB.Query(q, args...)

	fmt.Println(err)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return []webdata.TorrentWeb{}
	}
	defer rows.Close()

	for rows.Next() {
		var row webdata.TorrentWeb
		if err := rows.Scan(&row.TorrentID, &row.Created, &row.Name, &row.UpVotes, &row.DownVotes); err != nil {
			// do something with error
		} else {
			resultTorrents = append(resultTorrents, row)
		}
	}

	return resultTorrents
}

func GetUserFriends(userID int) []webdata.User {

	var friends []webdata.User

	q := "SELECT users.user_ID, users.username, profile_pic, banner_image, created, disabled, tagline, bio, gender, user_badges_blob FROM friends LEFT JOIN users ON friends.friend_ID = users.user_ID WHERE friends.users_user_ID = ?"
	rows, err := db.DB.Query(q, userID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user webdata.User
		err = rows.Scan(&user.UserID, &user.Username, &user.Cover, &user.Banner, &user.Joined, &user.Disabled, &user.Tagline, &user.Bio, &user.Gender, &user.UserBadgesBlob)
		if err != nil {
			// handle this error
			panic(err)
		}

		if user.UserBadgesBlob != nil {
			badges := ParseBadgeBlob(user.UserBadgesBlob)

			user.UserBadges = &badges
		}

		friends = append(friends, user)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return friends
}

func ParseBadgeBlob(blob *string) []webdata.Badges {
	var result []webdata.Badges
	err := json.Unmarshal([]byte(*blob), &result)
	if err != nil {
		fmt.Println("Error parsing badges")
		return []webdata.Badges{}

	}
	fmt.Println("JSON parse result")
	fmt.Println(result)
	return result
}

func LoadStrikes(userID int) []webdata.Strike {
	var strikes []webdata.Strike

	q := "SELECT strike_ID, user_ID, heading, description, date FROM strikes WHERE user_ID = ?"
	rows, err := db.DB.Query(q, userID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var strike webdata.Strike
		err = rows.Scan(&strike.StrikeID, &strike.UserID, &strike.Heading, &strike.Description, &strike.Date)
		if err != nil {
			// handle this error
			panic(err)
		}

		strikes = append(strikes, strike)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return strikes
}

func GetUserGroupsList(userID int, visibility string) []webdata.UserGroupListObject {
	var groupList []webdata.UserGroupListObject
	q := " SELECT group_members.group_ID, groups.group_name, groups.group_icon, group_roles.role_type, groups.group_visibility FROM group_members LEFT JOIN groups ON group_members.group_ID = groups.group_ID LEFT JOIN group_roles on group_roles.group_ID = group_members.group_ID WHERE group_members.user_ID = ? AND groups.group_visibility = ?"

	rows, err := db.DB.Query(q, userID, visibility)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var listItem webdata.UserGroupListObject
		err = rows.Scan(&listItem.GroupID, &listItem.GroupName, &listItem.GroupIcon, &listItem.GroupRole, &listItem.GroupVisibility)
		if err != nil {
			// handle this error
			panic(err)
		}

		groupList = append(groupList, listItem)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return groupList

}
