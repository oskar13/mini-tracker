package handlers

import (
	"fmt"
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	torrenttools "github.com/oskar13/mini-tracker/pkg/torrent-tools"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	"github.com/oskar13/mini-tracker/pkg/web/news"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func NewTorrentPage(w http.ResponseWriter, r *http.Request) {

	userData := webutils.GetUserData(r, db.DB)

	if !webutils.CheckLogin(w, r, userData) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var pageStruct struct {
		UserData   webdata.User
		SiteName   string
		PageName   string
		NewsList   []news.NewsArticle
		Community  []groups.GroupPost
		MyTorrents []webdata.TorrentWeb
	}
	pageStruct.UserData = userData
	pageStruct.SiteName = webdata.SiteName
	pageStruct.PageName = "New Torrent"

	if r.Method == "POST" {

		if r.Body == nil {
			http.Error(w, "Empty request", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()

		file, handler, err := r.FormFile("torrent")

		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)

			http.Error(w, "Problem uploading file", http.StatusMethodNotAllowed)

			return
		}

		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		fileType := handler.Header.Get("Content-type")

		if fileType != "application/x-bittorrent" {
			http.Error(w, "Invalid file uploaded", http.StatusMethodNotAllowed)
			return
		}

		fmt.Println("Reeee")
		fmt.Println(r.Form.Get("category"))

		prettyTitle := r.Form.Get("pretty_title")
		description := r.Form.Get("description")
		visibility := r.Form.Get("visibility")

		if visibility == "" || visibility == "Select" {
			//TODO proper validation
			http.Error(w, "Torrent visibility not set, must be set.", http.StatusMethodNotAllowed)
			return
		}

		torrent, err := torrenttools.DecodeUploadedTorrent(file)
		if err != nil {
			http.Error(w, "Could not parse uploaded torrent.", http.StatusMethodNotAllowed)
			return
		}

		fmt.Println(prettyTitle)
		fmt.Println(torrent.Announce)
		fmt.Println(description)
		fmt.Println("Done")

	} else {
		//If was not a POST request show torrent upload form
		webutils.RenderTemplate(w, []string{"pkg/web/templates/new-torrent.html",
			"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
			"pkg/web/templates/end.html",
			"pkg/web/templates/commandbar.html"}, pageStruct)
	}

	return

}
