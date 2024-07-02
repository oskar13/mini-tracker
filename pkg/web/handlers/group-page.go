package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
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
		UserRole  string
		Group     groups.GroupInfo
		Posts     []groups.GroupPost
	}

	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Group page title"

	idString := r.PathValue("id")

	if idString != "" {
		groupID, err := strconv.Atoi(idString)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err)
		} else {
			// Try loading group info
			pageStruct.UserRole = groups.LoadGroupAccess(userData.UserID, groupID)
			if pageStruct.UserRole == "" {
				//User has no right to view the page
				pageStruct.Error = true
				pageStruct.ErrorText = "Access denied to group"
			} else {
				//Continue loading data for page

				pageStruct.Group = groups.LoadGroupInfo(groupID)
				pageStruct.Posts = groups.LoadGroupPostsList(groupID)
			}
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
