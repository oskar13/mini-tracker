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
	serverMux.HandleFunc("/www", hello)

	fmt.Println("Starting tracking server at: http://", data.TrackerHostAndPort)
	http.ListenAndServe(data.TrackerHostAndPort, serverMux)

}

func hello(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hello")

	if r.Method == "GET" {
		//message := r.URL.Query().Get("message")

		var newPeer data.Peer

		fmt.Println(r.URL.Query())

		fmt.Println("Sprintf x")

		fmt.Println(fmt.Sprintf("%x", r.URL.Query().Get("info_hash")))
		fmt.Println(r.URL.Query().Get("peer_id"))
		fmt.Println(r.URL.Query().Get("port"))

		port, err := strconv.Atoi(r.URL.Query().Get("port"))
		if err != nil {
			panic(err)
		}

		left, err := strconv.Atoi(r.URL.Query().Get("left"))
		if err != nil {
			panic(err)
		}

		newPeer.InfoHash = fmt.Sprintf("%x", r.URL.Query().Get("info_hash"))
		newPeer.PeerID = r.URL.Query().Get("peer_id")
		newPeer.Port = port
		newPeer.Left = left
		torrentID, err := GetTorrentIDFromHash(newPeer.InfoHash)

		if err != nil {
			fmt.Println(err)
			return
		}

		newPeer.TorrentID = torrentID

		peers, err := LoadPeers(torrentID, newPeer.PeerID)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = AddPeer(newPeer)

		if err != nil {
			fmt.Println("Error adding peer")
			fmt.Println(err)
		}

		err = EncodePeerListAndRespond(w, torrentID, peers)

		if err != nil {
			fmt.Println("Error encoding peer list")
		} else {
			fmt.Println("Successfully updated peerlist")
		}

	}
}
