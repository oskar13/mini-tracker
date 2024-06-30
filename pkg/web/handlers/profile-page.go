package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {

	userData := webutils.GetUserData(r, db.DB)

	if !webutils.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	var pageStruct struct {
		Error         bool
		ErrorText     string
		DisplayedUser webdata.User
		UserData      webdata.User
		SelfEdit      bool
		TorrentList   []webdata.TorrentWeb
		SiteName      string
		PageName      string
		Strikes       []webdata.Strike
		UserGroups    []webdata.GroupListObject
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Profile"

	idString := r.PathValue("id")

	if idString != "" {
		userId, err := strconv.Atoi(idString)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err)
		} else {
			loadedUserData, err2 := webutils.LoadUserProfileData(userId)
			fmt.Println("LOADED USER DATA", loadedUserData)

			fmt.Println("ERROR", err2)
			if err2 != nil {
				pageStruct.Error = true
				pageStruct.ErrorText = fmt.Sprint(err2)

			} else {
				pageStruct.DisplayedUser = loadedUserData
				pageStruct.Strikes = webutils.LoadStrikes(pageStruct.DisplayedUser.UserID)
				pageStruct.UserGroups = groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Public")
			}
		}
	} else {
		// Display data for self
		loadedUserData, err2 := webutils.LoadUserProfileData(pageStruct.UserData.UserID)
		if err2 != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err2)

		} else {
			pageStruct.DisplayedUser = loadedUserData
			pageStruct.SelfEdit = true
			pageStruct.Strikes = webutils.LoadStrikes(pageStruct.DisplayedUser.UserID)
			pageStruct.UserGroups = groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Public")
			pageStruct.UserGroups = append(pageStruct.UserGroups, groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Private")...)

		}

	}

	pageStruct.TorrentList = webutils.LoadUserTorrents(pageStruct.DisplayedUser.UserID, []string{"Public"})

	webutils.RenderTemplate(w, []string{"pkg/web/templates/sidebar.html", "pkg/web/templates/profile.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
