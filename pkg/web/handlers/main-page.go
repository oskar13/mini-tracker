package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/news"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func MainPage(w http.ResponseWriter, r *http.Request) {

	userData := webutils.GetUserData(r, db.DB)

	if !webutils.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var pageStruct struct {
		UserData webdata.User
		SiteName string
		PageName string
		NewsList []news.NewsArticle
	}
	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Main"
	pageStruct.NewsList, _ = news.LoadNewsList(3)

	webutils.RenderTemplate(w, []string{"pkg/web/templates/home.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html",
		"pkg/web/templates/newslist-item.html"}, pageStruct)

}
