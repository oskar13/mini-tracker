package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	db "github.com/oskar13/mini-tracker/pkg/db"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func GroupPage(w http.ResponseWriter, r *http.Request) {

	userData := webutils.GetUserData(r, db.DB)

	if !webutils.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	var pageStruct struct {
		Error     bool
		ErrorText string
		UserData  webdata.User
		SiteName  string
		PageName  string
	}

	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Group page title"

	idString := r.PathValue("id")

	if idString != "" {
		groupId, err := strconv.Atoi(idString)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err)
		} else {
			// Try loading group info
			fmt.Println("Loading group id: ", groupId)
		}
	} else {
		// No ID string found

		pageStruct.Error = true
		pageStruct.ErrorText = "No group specified."

	}

	pageStruct.UserData = userData

	webutils.RenderTemplate(w, []string{"pkg/web/templates/groups/group-hub.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
