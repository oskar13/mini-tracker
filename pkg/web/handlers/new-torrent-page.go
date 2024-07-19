package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	"github.com/oskar13/mini-tracker/pkg/web/news"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func NewTorrentPage(w http.ResponseWriter, r *http.Request) {

	userData := webutils.GetUserData(r, db.DB)

	if !webutils.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var pageStruct struct {
		UserData   webdata.User
		SiteName   string
		PageName   string
		NewsList   []news.NewsArticle
		Community  []groups.GroupPost
		MyTorrents []webdata.TorrentWeb
	}
	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "New Torrent"

	webutils.RenderTemplate(w, []string{"pkg/web/templates/new-torrent.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)

}
