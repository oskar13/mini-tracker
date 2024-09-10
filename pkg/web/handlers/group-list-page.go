package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func GroupListPage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r, db.DB)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		Error     bool
		ErrorText string
		UserData  webdata.User
		SiteName  string
		PageName  string
		Groups    []groups.GroupInfo
	}

	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Groups"
	pageStruct.Groups = groups.ListPublicGroups()

	pageStruct.UserData = userData

	webutils.RenderTemplate(w, []string{"pkg/web/templates/groups/group-list.html", "pkg/web/templates/groups/group-list-item.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
