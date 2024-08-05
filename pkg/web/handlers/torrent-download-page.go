package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/oskar13/mini-tracker/pkg/data"
	torrenttools "github.com/oskar13/mini-tracker/pkg/torrent-tools"
	torrentweb "github.com/oskar13/mini-tracker/pkg/web/torrent-web"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func TorrentDownloadPage(w http.ResponseWriter, r *http.Request) {
	torrentUuidString := r.PathValue("uuid")

	var torrent webdata.TorrentWeb

	if torrentUuidString != "" {
		torrentUuid, err := uuid.Parse(torrentUuidString)
		if err != nil {
			http.Error(w, "Invalid uuid", http.StatusBadRequest)
			return
		} else {
			// Try torrent info
			torrent, err = torrentweb.LoadTorrentData(torrentUuid.String(), 0)
			if err != nil {

				http.Error(w, "Error loading torrent data", http.StatusBadRequest)
				return
			}

			if torrent.AccessType == "Public Listed" || torrent.AccessType == "Public Unlisted" {
				torrent.Announce = "http://" + data.TrackerHost + data.TrackerPort + "/www"
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
			newTorrentFile, err := torrenttools.TorrentFromDatabase(torrent)

			if err != nil {
				http.Error(w, "Error creating new torrent file", http.StatusBadRequest)
				return
			}

			fmt.Println(string(newTorrentFile)[:])

			fmt.Println(torrent.Announce)

			fmt.Println("Done")

			w.Header().Set("Content-Disposition", "attachment; filename="+torrent.Name+".torrent")
			w.Header().Set("Content-Type", "application/x-bittorrent")

			w.Write(newTorrentFile)

		}
	} else {
		// No ID string found
		http.Error(w, "Bad request", http.StatusBadRequest)
		return

	}
}
