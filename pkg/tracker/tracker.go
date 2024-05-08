package tracker

import (
	"fmt"
	"net/http"
)

func StartTracker() {

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", hello)

	fmt.Println("Starting tracking server at: http://localhost:7777")
	http.ListenAndServe("localhost:7777", serverMux)

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
