package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func TorrentSearchPage(c *gin.Context) {
	userData := accounts.GetUserData(c.Request)

	// Check login
	if !accounts.CheckLogin(c.Writer, c.Request, userData) {
		return
	}

	var pageStruct struct {
		Error     bool
		ErrorText string
		UserData  webdata.User
		SiteName  string
		PageName  string
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Search Results"

	webutils.RenderTemplate(c.Writer, []string{
		"pkg/web/templates/sidebar.html",
		"pkg/web/templates/search.html",
		"pkg/web/templates/torrent-list-item.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html",
	}, pageStruct)
}
