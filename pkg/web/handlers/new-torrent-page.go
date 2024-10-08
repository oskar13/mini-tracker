package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/oskar13/mini-tracker/pkg/data"
	torrenttools "github.com/oskar13/mini-tracker/pkg/torrent-tools"
	"github.com/oskar13/mini-tracker/pkg/web/accounts"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func NewTorrentPage(w http.ResponseWriter, r *http.Request) {

	userData := accounts.GetUserData(r)

	if !accounts.CheckLogin(w, r, userData) {
		return
	}

	var pageStruct struct {
		UserData   webdata.User
		SiteName   string
		PageName   string
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

		fmt.Println(r.Form.Get("category"))

		description := r.Form.Get("description")
		visibility := r.Form.Get("visibility")
		category := r.Form.Get("category")
		keepTrackersString := r.Form.Get("keepTrackers")
		keepTracker := false

		var categoryID int

		if visibility == "" || visibility == "Select" {
			//TODO proper validation
			http.Error(w, "Torrent visibility not set, must be set.", http.StatusMethodNotAllowed)
			return
		}

		if keepTrackersString != "yes" && keepTrackersString != "no" {
			http.Error(w, "Invalid form data.", http.StatusMethodNotAllowed)
			return
		} else {
			if keepTrackersString == "yes" {
				keepTracker = true
			}
		}

		if category == "" {
			//TODO proper validation
			http.Error(w, "Torrent category not set, must be set.", http.StatusMethodNotAllowed)
			return
		} else {
			categoryID, err = strconv.Atoi(category)
			if err != nil {
				http.Error(w, "Invalid", http.StatusBadRequest)
				return
			}
			if categoryID < 0 && categoryID > 999 {
				http.Error(w, "Invalid", http.StatusBadRequest)
				return
			}
		}

		torrent, err := torrenttools.DecodeUploadedTorrent(file)
		if err != nil {
			http.Error(w, "Could not parse uploaded torrent.", http.StatusMethodNotAllowed)
			return
		}

		if visibility == "Public Listed" || visibility == "Public Unlisted" {
			torrent.Announce = "http://" + data.TrackerHost + data.TrackerPort + "/www"
		}

		if visibility != "Public Listed" && keepTracker {
			http.Error(w, "Only public listed torrents are allowed to have other trackers included.", http.StatusBadRequest)
			return
		}

		torrent.Description = description
		torrent.AccessType = visibility
		torrent.CategoryID = categoryID

		uploadedTorrentUuid, err := torrentweb.CreateTorrentEntry(torrent, userData.UserID, keepTracker)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Could not save uploaded torrent.", http.StatusInternalServerError)
			return
		}

		fmt.Println("Uploaded uuid ", uploadedTorrentUuid)

		http.Redirect(w, r, "/t/"+uploadedTorrentUuid, http.StatusSeeOther)

		/*
			//Send client a generated torrent file
			newTorrentFile, err := torrenttools.TorrentFromWebTorrent(torrent)

			if err != nil {
				fmt.Println("Error creating a torrent file")
				fmt.Println(err)
				return
			}

			fmt.Println(string(newTorrentFile)[:])

			fmt.Println(torrent.Announce)
			fmt.Println(description)
			fmt.Println("Done")

			w.Header().Set("Content-Disposition", "attachment; filename=foo.torrent")
			w.Header().Set("Content-Type", "application/x-bittorrent")

			w.Write(newTorrentFile)

		*/

	} else {

		//If was not a POST request show torrent upload form
		webutils.RenderTemplate(w, []string{"pkg/web/templates/new-torrent.html",
			"pkg/web/templates/sidebar.html", "pkg/web/templates/head.html",
			"pkg/web/templates/end.html",
			"pkg/web/templates/commandbar.html"}, pageStruct)
	}

}
