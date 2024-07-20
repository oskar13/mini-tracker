package tracker

import (
	"fmt"
	"net/http"
	"net/url"

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

		fmt.Println(r.URL.Query().Get("info_hash"))

		escaped := url.QueryEscape(r.URL.Query().Get("info_hash"))

		fmt.Println(escaped)

		var response = make(map[string]interface{})

		response["interval"] = 1800

		var peerlist []interface{}

		var peer = make(map[string]interface{})

		peerlist = append(peerlist, peer)

		peer["ip"] = "192.168.189.1"
		peer["port"] = 20111

		response["peers"] = peerlist

		enc := bencode.NewEncoder(w)
		if err := enc.Encode(response); err != nil {
			panic(err)
		}

	}
}
