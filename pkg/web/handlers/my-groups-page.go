package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func MyGroupsPage(w http.ResponseWriter, r *http.Request) {

	userData := webutils.GetUserData(r, db.DB)

	if !webutils.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	var pageStruct struct {
		Error               bool
		ErrorText           string
		UserData            webdata.User
		SelfEdit            bool
		UnlistedTorrentList []webdata.TorrentWeb
		PublicTorrentList   []webdata.TorrentWeb
		SiteName            string
		PageName            string
		Strikes             []webdata.Strike
		UserGroups          []groups.GroupInfo
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "My Groups"

	pageStruct.PublicTorrentList = torrentweb.LoadUserTorrents(pageStruct.UserData.UserID, []string{"Public", "WWW"})
	pageStruct.UnlistedTorrentList = torrentweb.LoadUserTorrents(pageStruct.UserData.UserID, []string{"Site Public"})

	webutils.RenderTemplate(w, []string{"pkg/web/templates/sidebar.html", "pkg/web/templates/my-groups.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
