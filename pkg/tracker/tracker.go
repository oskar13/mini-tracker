package tracker

import (
	"fmt"
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
)

func StartTracker() {

	// Initialize the database
	db.InitDB()
	defer db.Close()

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/www", hello)

	fmt.Println("Starting tracking server at: http://localhost:7777")
	http.ListenAndServe("localhost:7777", serverMux)

}

func hello(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hello")

	if r.Method == "GET" {
		//message := r.URL.Query().Get("message")

		fmt.Println(r.URL.Query())

		fmt.Println("Sprintf x")
		fmt.Println(fmt.Sprintf("%x", r.URL.Query().Get("info_hash")))
		fmt.Println(r.URL.Query().Get("peer_id"))
		fmt.Println(r.URL.Query().Get("port"))

		fmt.Println(GetHTTPRequestIP(r))

		peers, err := LoadPeers(1)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = EncodePeerListAndRespond(w, 18, peers)

		if err != nil {
			fmt.Println("Error encoding peer list")
		} else {
			fmt.Println("Successfully updated peerlist")
		}

	}
}
