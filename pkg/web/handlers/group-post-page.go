package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func GroupPostPage(c *gin.Context) {
	userData := accounts.GetUserData(c.Request)

	if !accounts.CheckLogin(c.Writer, c.Request, userData) {
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

	groupIdString := c.Param("groupid")
	postIdString := c.Param("postid")

	if groupIdString != "" {
		groupID, err := strconv.Atoi(groupIdString)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err)
		} else {
			pageStruct.UserRole = groups.LoadGroupAccess(userData.UserID, groupID)

			if pageStruct.UserRole == "" {
				pageStruct.Error = true
				pageStruct.ErrorText = "Access denied to group"
			} else {
				pageStruct.Group = groups.LoadGroupInfo(groupID)

				if postIdString != "" {
					postID, err := strconv.Atoi(postIdString)
					if err != nil {
						pageStruct.Error = true
						pageStruct.ErrorText = fmt.Sprint(err)
					} else {
						pageStruct.Post = groups.LoadGroupPost(groupID, postID)
					}
				} else {
					pageStruct.Error = true
					pageStruct.ErrorText = "No post specified."
				}
			}
		}
	} else {
		pageStruct.Error = true
		pageStruct.ErrorText = "No group specified."
	}

	pageStruct.UserData = userData

	webutils.RenderTemplate(c.Writer, []string{
		"pkg/web/templates/groups/group-post-page.html",
		"pkg/web/templates/sidebar.html",
		"pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html",
	}, pageStruct)
}
