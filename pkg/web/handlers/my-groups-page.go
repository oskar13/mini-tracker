package handlers

import (
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func MyGroupsPage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r, db.DB)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		Error         bool
		ErrorText     string
		UserData      webdata.User
		SelfEdit      bool
		SiteName      string
		PageName      string
		PublicGroups  []groups.GroupInfo
		PrivateGroups []groups.GroupInfo
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "My Groups"
	pageStruct.PublicGroups = groups.GetUserGroupsList(userData.UserID, "Public")
	pageStruct.PrivateGroups = groups.GetUserGroupsList(userData.UserID, "Private")

	webutils.RenderTemplate(w, []string{"pkg/web/templates/sidebar.html", "pkg/web/templates/my-groups.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
