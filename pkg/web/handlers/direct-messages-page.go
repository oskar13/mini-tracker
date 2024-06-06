package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func DirectMessages(w http.ResponseWriter, r *http.Request) {

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
		FriendList    []webdata.User
	}

	pageStruct.UserData = userData

	pageStruct.FriendList = webutils.GetUserFriends(pageStruct.UserData.UserID)

	webutils.RenderTemplate(w, []string{"pkg/web/templates/directmessages.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
