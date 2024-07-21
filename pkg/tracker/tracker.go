package tracker

import (
	"fmt"
	"net/http"

	"github.com/zeebo/bencode"
)

func StartTracker() {

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", hello)

	fmt.Println("Starting tracking server at: http://localhost:7777")
	http.ListenAndServe("localhost:7777", serverMux)

}

func hello(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hello")

	if r.Method == "GET" {
		//message := r.URL.Query().Get("message")

		fmt.Println("Sprintf x")
		fmt.Println(fmt.Sprintf("%x", r.URL.Query().Get("info_hash")))

		var response = make(map[string]interface{})

		response["interval"] = 18

		var peerlist []interface{}

		var peer = make(map[string]interface{})

		peer["ip"] = "192.168.101.111"
		peer["port"] = 20111

		peerlist = append(peerlist, peer)

		response["peers"] = peerlist

		enc := bencode.NewEncoder(w)
		if err := enc.Encode(response); err != nil {
			panic(err)
		}

	}
}
