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

func GroupPostPage(w http.ResponseWriter, r *http.Request) {
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
		Post      groups.GroupPost
	}

	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Group Post"

	groupIdString := r.PathValue("groupid")
	postIdString := r.PathValue("postid")

	if groupIdString != "" {
		groupID, err := strconv.Atoi(groupIdString)
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

				if postIdString != "" {
					postID, err := strconv.Atoi(postIdString)
					if err != nil {
						pageStruct.Error = true
						pageStruct.ErrorText = fmt.Sprint(err)
					} else {
						// Load the replies for the post
						pageStruct.Post = groups.LoadGroupPost(groupID, postID)
					}
				} else {
					// No post ID string found

					pageStruct.Error = true
					pageStruct.ErrorText = "No group specified."

				}

			}
		}
	} else {
		// No ID string found

		pageStruct.Error = true
		pageStruct.ErrorText = "No group specified."

	}

	pageStruct.UserData = userData

	webutils.RenderTemplate(w, []string{"pkg/web/templates/groups/group-post-page.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
