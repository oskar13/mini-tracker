package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		Error         bool
		ErrorText     string
		DisplayedUser webdata.User
		UserData      webdata.User
		ViewSelf      bool
		CanEdit       bool
		TorrentList   []webdata.TorrentWeb
		SiteName      string
		PageName      string
		Strikes       []webdata.Strike
		UserGroups    []groups.GroupInfo
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Profile"

	idString := r.PathValue("id")

	if idString != "" {
		// Process user ID to load a profile
		userId, err := strconv.Atoi(idString)
		if err != nil {
			webutils.ReturnErrorResponse(w, r, "Bad request", http.StatusBadRequest)
			return
		}

		loadedUserData, err2 := webutils.LoadUserProfileData(userId)

		if err2 != nil {
			webutils.ReturnErrorResponse(w, r, "User not found", http.StatusNotFound)
			return
		}

		pageStruct.DisplayedUser = loadedUserData

		// Load password field for user self password edit
		if pageStruct.DisplayedUser.UserID == pageStruct.UserData.UserID {
			pageStruct.DisplayedUser.Password = pageStruct.UserData.Password
		}

		pageStruct.Strikes = webutils.LoadStrikes(pageStruct.DisplayedUser.UserID)
		pageStruct.UserGroups = groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Public")
		pageStruct.TorrentList = torrentweb.ListTorrents(pageStruct.DisplayedUser.UserID, []string{"Public Listed", "Members Listed", "Group Public"}, 10)

		if pageStruct.DisplayedUser.UserID == pageStruct.UserData.UserID {
			pageStruct.ViewSelf = true
		}

		if pageStruct.UserData.AdminLevel >= 90 && pageStruct.UserData.AdminLevel > pageStruct.DisplayedUser.AdminLevel {
			pageStruct.CanEdit = true

			pageStruct.UserGroups = append(pageStruct.UserGroups, groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Private")...)
		} else if pageStruct.ViewSelf {
			pageStruct.CanEdit = true

			pageStruct.UserGroups = append(pageStruct.UserGroups, groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Private")...)
		}

		if r.Method == "POST" {
			if !pageStruct.CanEdit {
				http.Error(w, "You have no edit rights", http.StatusForbidden)
				return
			}

			newData, err := handleProfilePost(r, pageStruct.DisplayedUser)
			if err != nil {
				http.Error(w, "Error processing request:"+fmt.Sprint(err), http.StatusInternalServerError)
				return
			}
			pageStruct.DisplayedUser = newData
		}

	} else {
		webutils.ReturnErrorResponse(w, r, "User not found", http.StatusNotFound)
		return
	}

	webutils.RenderTemplate(w, []string{"pkg/web/templates/sidebar.html", "pkg/web/templates/profile.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}

func handleProfilePost(r *http.Request, userData webdata.User) (webdata.User, error) {

	if err := r.ParseForm(); err != nil {
		return webdata.User{}, nil
	}

	email := r.FormValue("email")
	if email != "" {
		if webutils.ValidateEmail(email) {
			userData.Email = email
		} else {
			return webdata.User{}, errors.New("invalid email address")
		}
	}

	fmt.Println(userData.Password)

	oldPassword := r.FormValue("password-old")
	newPassword := r.FormValue("password-new")
	newPassword2 := r.FormValue("password-new2")

	if userData.Password == "" && newPassword != "" {
		//Editing someone else's password, don't ask for old password
		err := accounts.SetPassword(userData.UserID, newPassword)
		if err != nil {
			return webdata.User{}, err
		}
	} else if oldPassword != "" && newPassword != "" && newPassword2 != "" {
		err := accounts.UserUpdatePassword(userData.UserID, oldPassword, newPassword, newPassword2)
		if err != nil {
			return webdata.User{}, err
		}
	}

	tagline := r.FormValue("tagline")
	if tagline != "" {
		if len(tagline) > 140 {
			return webdata.User{}, errors.New("tagline too long")
		} else {
			userData.Tagline = &tagline
		}
	}
	bio := r.FormValue("bio")
	if bio != "" {
		if len(bio) > 500 {
			return webdata.User{}, errors.New("bio too long")
		} else {
			userData.Bio = &bio
		}
	}
	gender := r.FormValue("gender")
	if gender != "" {
		if len(gender) > 100 {
			return webdata.User{}, errors.New("gender too long")
		} else {
			userData.Gender = &gender
		}
	}

	err := accounts.UpdateAccountFields(userData)

	if err != nil {
		return webdata.User{}, errors.New(fmt.Sprint(err))
	}

	return userData, nil
}
