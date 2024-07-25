package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/oskar13/mini-tracker/pkg/data"
	torrenttools "github.com/oskar13/mini-tracker/pkg/torrent-tools"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func TorrentDownloadPage(w http.ResponseWriter, r *http.Request) {
	torrentIdString := r.PathValue("id")

	var torrent webdata.TorrentWeb

	if torrentIdString != "" {
		torrentID, err := strconv.Atoi(torrentIdString)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		} else {
			// Try torrent info
			torrent, err = torrentweb.LoadTorrentData(torrentID, 0)
			if err != nil {

				http.Error(w, "Error loading torrent data", http.StatusBadRequest)
				return
			}

			if torrent.AccessType == "Public" || torrent.AccessType == "WWW" {
				torrent.Announce = "http://" + data.TrackerHostAndPort + "/www"
				torrent.InfoField, err = torrentweb.LoadTorrentInfoField(torrent.TorrentID)
				if err != nil {
					http.Error(w, "Error loading info field", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}

			//Send client a generated torrent file
			newTorrentFile, err := torrenttools.TorrentFromWebTorrent(torrent)

			if err != nil {
				http.Error(w, "Error creating new torrent file", http.StatusBadRequest)
				return
			}

			fmt.Println(string(newTorrentFile)[:])

			fmt.Println(torrent.Announce)

			fmt.Println("Done")

			w.Header().Set("Content-Disposition", "attachment; filename=foo.torrent")
			w.Header().Set("Content-Type", "application/x-bittorrent")

			w.Write(newTorrentFile)

		}
	} else {
		// No ID string found
		http.Error(w, "Bad request", http.StatusBadRequest)
		return

	}
}
