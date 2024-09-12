package handlers

import (
	"net/http"

	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	"github.com/oskar13/mini-tracker/pkg/web/news"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func MainPage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		UserData       webdata.User
		SiteName       string
		PageName       string
		NewsList       []news.NewsArticle
		Community      []groups.GroupPost
		MyTorrents     []webdata.TorrentWeb
		LatestTorrents []webdata.TorrentWeb
	}
	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Main"
	pageStruct.NewsList, _ = news.LoadNewsList(3)
	pageStruct.Community = groups.GetCommunityUpdates(userData.UserID, 3)
	pageStruct.MyTorrents = torrentweb.ListTorrents(userData.UserID, []string{"Public Listed", "Public Unlisted", "Members Listed", "Members Unlisted", "Members Access List", "Group Public", "Group Private"}, 5)
	pageStruct.LatestTorrents = torrentweb.ListTorrents(0, []string{"Public Listed", "Members Listed", "Group Public"}, 5)

	webutils.RenderTemplate(w, []string{"pkg/web/templates/main.html",
		"pkg/web/templates/torrent-list-item.html",
		"pkg/web/templates/sidebar.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html",
		"pkg/web/templates/newslist-item.html"}, pageStruct)

}
