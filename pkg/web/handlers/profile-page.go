package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	db "github.com/oskar13/mini-tracker/pkg/db"
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
	}

	pageStruct.UserData = userData

	idString := r.PathValue("id")
	fmt.Printf("ID: %v", idString)

	if idString != "" {
		userId, err := strconv.Atoi(idString)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err)
		} else {
			loadedUserData, err2 := webutils.LoadUserProfileData(userId)

			if err2 != nil {
				pageStruct.Error = true
				pageStruct.ErrorText = fmt.Sprint(err2)

			} else {
				pageStruct.DisplayedUser = loadedUserData
			}
		}
	} else {
		// Display data for self
		pageStruct.DisplayedUser = pageStruct.UserData
		pageStruct.SelfEdit = true
	}

	pageStruct.TorrentList = webutils.LoadUserTorrents(pageStruct.DisplayedUser.UserID, []string{"Public"})

	fmt.Println(pageStruct.TorrentList)
	fmt.Println(pageStruct.TorrentList)
	webutils.RenderTemplate(w, []string{"pkg/web/templates/profile.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
