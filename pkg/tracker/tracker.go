package tracker

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/oskar13/mini-tracker/pkg/data"
	db "github.com/oskar13/mini-tracker/pkg/db"
)

func StartTracker() {

	// Initialize the database
	db.InitDB()
	defer db.Close()

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/www", HandlePublicTorrents)

	fmt.Println("Starting tracking server at: http://", data.TrackerHostAndPort)
	http.ListenAndServe(data.TrackerHostAndPort, serverMux)

}

// Handle anonymous public requests
func HandlePublicTorrents(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		var newPeer data.Peer

		fmt.Println(r.URL.Query())

		port, err := strconv.Atoi(r.URL.Query().Get("port"))
		if err != nil {
			http.Error(w, "Invalid port number", http.StatusBadRequest)
			return
		}

		left, err := strconv.Atoi(r.URL.Query().Get("left"))
		if err != nil {
			http.Error(w, "Could not parse how much is left to download", http.StatusBadRequest)
			return
		}

		newPeer.InfoHash = fmt.Sprintf("%x", r.URL.Query().Get("info_hash"))
		if len(newPeer.InfoHash) != 40 {
			http.Error(w, "Invalid info hash size", http.StatusBadRequest)
			return
		}

		newPeer.PeerID = r.URL.Query().Get("peer_id")
		if len(newPeer.PeerID) != 20 {
			http.Error(w, "Invalid peer id", http.StatusBadRequest)
			return
		}

		newPeer.Port = port
		newPeer.Left = left
		torrentID, err := GetTorrentIDFromHash(newPeer.InfoHash)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Could not fetch torrent", http.StatusInternalServerError)
			return
		}

		newPeer.TorrentID = torrentID

		peers, err := LoadPeers(newPeer.TorrentID, newPeer.PeerID)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = AddPeer(newPeer)

		if err != nil {
			fmt.Println("Error adding peer")
			http.Error(w, "Error adding peer", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		err = EncodePeerListAndRespond(w, 20, peers)

		if err != nil {
			fmt.Println("Error encoding peer list")
			http.Error(w, "Error encoding peer list", http.StatusInternalServerError)
			fmt.Println(err)
			return
		} else {
			fmt.Println("Successfully updated peerlist")
		}

	} else {
		http.Error(w, "Invalid ", http.StatusMethodNotAllowed)
		return
	}
}
