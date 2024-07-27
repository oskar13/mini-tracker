package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func MyTorrentsPage(w http.ResponseWriter, r *http.Request) {

	userData := webutils.GetUserData(r, db.DB)

	if !webutils.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	var pageStruct struct {
		Error               bool
		ErrorText           string
		UserData            webdata.User
		SiteName            string
		PageName            string
		UnlistedTorrentList []webdata.TorrentWeb
		PublicTorrentList   []webdata.TorrentWeb
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "My Torrents"

	pageStruct.PublicTorrentList = torrentweb.LoadUserTorrents(pageStruct.UserData.UserID, []string{"Public", "WWW"})
	pageStruct.UnlistedTorrentList = torrentweb.LoadUserTorrents(pageStruct.UserData.UserID, []string{"Site Public"})

	webutils.RenderTemplate(w, []string{"pkg/web/templates/sidebar.html", "pkg/web/templates/my-torrents.html",
		"pkg/web/templates/torrent-list-item.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
