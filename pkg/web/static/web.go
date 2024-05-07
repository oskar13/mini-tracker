package web

import (
	"fmt"
	"net/http"
)

func StartWebsite() {

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", hello)

	fmt.Println("Starting web interface at: http://localhost:8080")
	http.ListenAndServe("localhost:8080", serverMux)

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
