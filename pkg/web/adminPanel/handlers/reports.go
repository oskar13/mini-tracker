package handlers

import (
	"net/http"

	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func ReportsPage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		UserData webdata.User
		SiteName string
		PageName string
	}
	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Reports Admin Page"

	webutils.RenderTemplate(w, []string{"pkg/web/adminPanel/templates/reports.html", "pkg/web/adminPanel/templates/head.html", "pkg/web/adminPanel/templates/sidebar.html", "pkg/web/adminPanel/templates/footer.html"}, pageStruct)

}
