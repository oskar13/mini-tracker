package handlers

import (
	"fmt"
	"net/http"

	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func FriendsPage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		Error          bool
		ErrorText      string
		DisplayedUser  webdata.User
		UserData       webdata.User
		SelfEdit       bool
		TorrentList    []webdata.TorrentWeb
		FriendList     []webdata.User
		SiteName       string
		PageName       string
		FriendRequests webdata.FriendRequests
	}

	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Friends"

	pageStruct.UserData = userData

	pageStruct.FriendRequests = webutils.GetFriendRequests(pageStruct.UserData.UserID)
	pageStruct.FriendList = webutils.GetUserFriends(pageStruct.UserData.UserID)

	if r.Method == "POST" {

		r.ParseForm()
		//TODO Handle posted data
		fmt.Println(r.Form)
	}

	webutils.RenderTemplate(w, []string{"pkg/web/templates/friends.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
