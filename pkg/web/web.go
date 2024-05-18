package web

import (
	"fmt"
	"net/http"
)

func StartWebsite() {

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", hello)
	serverMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Starting web interface at: http://localhost:8080")
	http.ListenAndServe("localhost:8080", serverMux)

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
