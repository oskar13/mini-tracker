package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func MyTorrentsPage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r, db.DB)

	if !accounts.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	var pageStruct struct {
		Error                  bool
		ErrorText              string
		UserData               webdata.User
		SiteName               string
		PageName               string
		PublicUnlistedTorrents []webdata.TorrentWeb
		PublicListedTorrents   []webdata.TorrentWeb
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "My Torrents"

	pageStruct.PublicListedTorrents = torrentweb.ListTorrents(pageStruct.UserData.UserID, []string{"Public Listed"}, 99)
	pageStruct.PublicUnlistedTorrents = torrentweb.ListTorrents(pageStruct.UserData.UserID, []string{"Public Unlisted"}, 99)

	webutils.RenderTemplate(w, []string{"pkg/web/templates/sidebar.html", "pkg/web/templates/my-torrents.html",
		"pkg/web/templates/torrent-list-item.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
