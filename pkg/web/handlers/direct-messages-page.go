package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	directmessages "github.com/oskar13/mini-tracker/pkg/web/directMessages"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func DirectMessages(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		Error     bool
		ErrorText string
		UserData  webdata.User
		SiteName  string
		PageName  string
		Thread    webdata.DMThread
		Threads   []webdata.DMThreadListItem
	}

	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Direct Messages"

	threadIdString := r.PathValue("id")

	if threadIdString != "" {
		threadID, err := strconv.Atoi(threadIdString)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err)
		} else {
			// Load thread contents
			pageStruct.Thread = directmessages.LoadDMThread(threadID)
		}
	} else {
		// Display list of threads
		pageStruct.Threads = directmessages.LoadDMThreadList(pageStruct.UserData.UserID)

	}

	webutils.RenderTemplate(w, []string{"pkg/web/templates/directmessages.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
