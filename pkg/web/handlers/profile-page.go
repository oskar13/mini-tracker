package handlers

import (
	"net/http"
	"strconv"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r, db.DB)

	if !accounts.CheckLogin(w, r, userData) {
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
		pageStruct.Strikes = webutils.LoadStrikes(pageStruct.DisplayedUser.UserID)
		pageStruct.UserGroups = groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Public")

	} else {
		// Display data for self
		loadedUserData, err2 := webutils.LoadUserProfileData(pageStruct.UserData.UserID)
		if err2 != nil {
			webutils.ReturnErrorResponse(w, r, "User not found / Internal server error", http.StatusNotFound)
			return
		}

		pageStruct.DisplayedUser = loadedUserData
		pageStruct.SelfEdit = true
		pageStruct.Strikes = webutils.LoadStrikes(pageStruct.DisplayedUser.UserID)
		pageStruct.UserGroups = groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Public")
		pageStruct.UserGroups = append(pageStruct.UserGroups, groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Private")...)

	}

	pageStruct.TorrentList = torrentweb.ListTorrents(pageStruct.DisplayedUser.UserID, []string{"Public Listed", "Members Listed", "Group Public"}, 10)

	webutils.RenderTemplate(w, []string{"pkg/web/templates/sidebar.html", "pkg/web/templates/profile.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
