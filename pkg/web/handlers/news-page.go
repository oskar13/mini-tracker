package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	"github.com/oskar13/mini-tracker/pkg/web/news"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func NewsPage(c *gin.Context) {
	userData := accounts.GetUserData(c.Request)

	// Check login
	if !accounts.CheckLogin(c.Writer, c.Request, userData) {
		return
	}

	var pageStruct struct {
		Error       bool
		ErrorText   string
		UserData    webdata.User
		NewsArticle news.NewsArticle
		NewsList    []news.NewsArticle
		SiteName    string
		PageName    string
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "News"

	idString := c.Param("id") // Gin grabs from route param like /news/:id

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
		loadedNewsList, err := news.LoadNewsList(25)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprintf("%v", err)
		} else {
			pageStruct.NewsList = loadedNewsList
		}
	}

	// Render the same template stack
	webutils.RenderTemplate(c.Writer, []string{
		"pkg/web/templates/sidebar.html",
		"pkg/web/templates/news.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html",
		"pkg/web/templates/newslist-item.html",
	}, pageStruct)
}
