package handlers

import (
	"net/http"

	"github.com/google/uuid"
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

			newTorrentFile, err := torrentweb.CreatePublicTorrentFile(torrent)

			if err != nil {
				http.Error(w, "Error creating new torrent file", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Disposition", "attachment; filename="+torrent.Name+".torrent")
			w.Header().Set("Content-Type", "application/x-bittorrent")

			w.Write(newTorrentFile)
			return
		}
	} else {
		// No ID string found
		http.Error(w, "Bad request", http.StatusBadRequest)
		return

	}
}
