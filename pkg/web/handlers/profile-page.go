package handlers

import (
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
		pageStruct.Strikes = webutils.LoadStrikes(pageStruct.DisplayedUser.UserID)
		pageStruct.UserGroups = groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Public")
		pageStruct.TorrentList = torrentweb.ListTorrents(pageStruct.DisplayedUser.UserID, []string{"Public Listed", "Members Listed", "Group Public"}, 10)

		if pageStruct.DisplayedUser.UserID == pageStruct.UserData.UserID {
			pageStruct.ViewSelf = true
		}

		if pageStruct.UserData.AdminLevel <= 3 || pageStruct.ViewSelf {
			pageStruct.CanEdit = true

			pageStruct.UserGroups = append(pageStruct.UserGroups, groups.GetUserGroupsList(pageStruct.DisplayedUser.UserID, "Private")...)
		}

		if r.Method == "POST" {
			if err := handleProfilePost(r); err != nil {
				http.Error(w, "Error processing request", http.StatusInternalServerError)
				return
			}
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

func handleProfilePost(r *http.Request) error {
	return nil
}
