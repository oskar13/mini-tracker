package handlers

import (
	"net/http"

	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func SiteSettingsPage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	if userData.AdminLevel < webdata.MainAdminLevel {
		webutils.ReturnErrorResponse(w, r, "You have no access rights to browse this page", http.StatusForbidden)
		return
	}

	var pageStruct struct {
		UserData webdata.User
		SiteName string
		PageName string
	}
	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Site Settings"

	webutils.RenderTemplate(w, []string{"pkg/web/adminPanel/templates/site-settings.html", "pkg/web/adminPanel/templates/head.html", "pkg/web/adminPanel/templates/sidebar.html", "pkg/web/adminPanel/templates/footer.html"}, pageStruct)

}
