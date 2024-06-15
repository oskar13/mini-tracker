package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/news"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func NewsPage(w http.ResponseWriter, r *http.Request) {

	userData := webutils.GetUserData(r, db.DB)

	if !webutils.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	var pageStruct struct {
		Error       bool
		ErrorText   string
		UserData    webdata.User
		NewsArticle news.NewsArticle
		NewsList    []news.NewsArticle
		PageName    string
	}

	pageStruct.UserData = userData
	pageStruct.PageName = "news"

	idString := r.PathValue("id")

	if idString != "" {
		articleID, err := strconv.Atoi(idString)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err)
		} else {
			loadedNewsArticle, err2 := news.LoadNewsArticle(articleID)

			if err2 != nil {
				pageStruct.Error = true
				pageStruct.ErrorText = fmt.Sprint(err2)

			} else {
				pageStruct.NewsArticle = loadedNewsArticle
			}
		}
	} else {
		// Show list of news
		loadedNewsList, err := news.LoadNewsList()

		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprintf("%v", err)
		} else {
			pageStruct.NewsList = loadedNewsList
		}

	}

	webutils.RenderTemplate(w, []string{"pkg/web/templates/sidebar.html", "pkg/web/templates/news.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
