package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func TorrentPage(w http.ResponseWriter, r *http.Request) {
	userData := accounts.GetUserData(r, db.DB)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		Error      bool
		ErrorText  string
		UserData   webdata.User
		SiteName   string
		PageName   string
		TheTorrent webdata.TorrentWeb
	}

	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "Torrent"

	torrentUuidString := r.PathValue("uuid")

	if torrentUuidString != "" {
		torrentUuid, err := uuid.Parse(torrentUuidString)
		if err != nil {
			pageStruct.Error = true
			pageStruct.ErrorText = fmt.Sprint(err)
		} else {
			// Try torrent info
			pageStruct.TheTorrent, err = torrentweb.LoadTorrentData(torrentUuid.String(), userData.UserID)
			if err != nil {
				pageStruct.Error = true
				pageStruct.ErrorText = fmt.Sprintf("%v", err)
			}

		}
	} else {
		// No ID string found

		pageStruct.Error = true
		pageStruct.ErrorText = "No torrent specified."

	}

	pageStruct.UserData = userData

	webutils.RenderTemplate(w, []string{"pkg/web/templates/torrent.html",
		"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
		"pkg/web/templates/end.html",
		"pkg/web/templates/commandbar.html"}, pageStruct)
}
