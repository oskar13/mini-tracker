package webutils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func ValidateSchema() error {
	var schema_revision string
	q := "SELECT schema_revision FROM sys_info WHERE schema_revision = ?"
	return db.DB.QueryRow(q, db.SchemaRevision).Scan(&schema_revision)
}

// Check if any users exist in database
func CheckUsers() (bool, error) {
	var count int
	q := "select count(1) where exists (select * from minitorrent.users)"
	err := db.DB.QueryRow(q).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return true, nil
	}

	return false, nil
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

		log.Printf("error executing template: %v", err)
		return
		//http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}

func ReturnErrorResponse(w http.ResponseWriter, r *http.Request, errMsg string, statusCode int) {

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/html")

	// Send the error message to the response writer
	fmt.Fprintf(w, "<html><body><h1>Error %d</h1><p>%s</p></body></html>", statusCode, errMsg)

	return
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

// Load user data for profile page or profile cards
func LoadUserProfileData(user_ID int) (webdata.User, error) {

	var user webdata.User

	q := `SELECT users.username, users.email, users.admin_level, users.profile_pic, users.banner_image, users.created, users.disabled, users.tagline, users.bio, users.gender FROM users WHERE user_ID = ?`
	err := db.DB.QueryRow(q, user_ID).Scan(&user.Username, &user.Email, &user.AdminLevel, &user.Cover, &user.Banner, &user.Joined, &user.Disabled, &user.Tagline, &user.Bio, &user.Gender)
	if err != nil {
		if err == sql.ErrNoRows {
			return webdata.User{}, errors.New("no account found")
		}
		return webdata.User{}, err
	}

	user.UserID = user_ID

	return user, nil
}

func GetUserFriends(userID int) []webdata.User {

	var friends []webdata.User

	q := "SELECT users.user_ID, users.username, profile_pic, banner_image, created, disabled, tagline, bio, gender, user_badges_blob FROM friends LEFT JOIN users ON friends.friend_ID = users.user_ID WHERE friends.user_ID = ?"
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
	// fmt.Println("JSON parse result")
	// fmt.Println(result)
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

func GetFriendRequests(userID int) webdata.FriendRequests {
	var friendRequests webdata.FriendRequests
	friendRequests.Incoming = getIncomingFriendRequests(userID)
	friendRequests.Outgoing = getOutgoingFriendRequests(userID)
	return friendRequests
}

func getIncomingFriendRequests(userID int) []webdata.FriendRequest {

	var incomingList []webdata.FriendRequest

	q := "SELECT friend_request_ID, users.user_ID, users.username, users.profile_pic, friend_requests.message, friend_requests.date FROM friend_requests LEFT JOIN users ON sender_user_ID = users.user_ID WHERE receiver_user_ID = ?"

	rows, err := db.DB.Query(q, userID)
	if err != nil {
		return []webdata.FriendRequest{}
	}
	defer rows.Close()
	for rows.Next() {
		var result webdata.FriendRequest

		err = rows.Scan(&result.FriendRequestID, &result.User.UserID, &result.User.Username, &result.User.Cover, &result.Message, &result.Date)
		if err != nil {
			return []webdata.FriendRequest{}
		}

		incomingList = append(incomingList, result)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return []webdata.FriendRequest{}
	}

	return incomingList
}

func getOutgoingFriendRequests(userID int) []webdata.FriendRequest {

	var outgoingList []webdata.FriendRequest

	q := "SELECT friend_request_ID, users.user_ID, users.username, users.profile_pic, friend_requests.message, friend_requests.date FROM friend_requests LEFT JOIN users ON receiver_user_ID = users.user_ID WHERE sender_user_ID = ?"

	rows, err := db.DB.Query(q, userID)
	if err != nil {
		return []webdata.FriendRequest{}
	}
	defer rows.Close()
	for rows.Next() {
		var result webdata.FriendRequest

		err = rows.Scan(&result.FriendRequestID, &result.User.UserID, &result.User.Username, &result.User.Cover, &result.Message, &result.Date)
		if err != nil {
			return []webdata.FriendRequest{}
		}

		outgoingList = append(outgoingList, result)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return []webdata.FriendRequest{}
	}

	return outgoingList

}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
