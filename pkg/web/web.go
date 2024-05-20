package web

import (
	"fmt"
	"net/http"

	"github.com/oskar13/mini-tracker/pkg/web/handlers"
)

func StartWebsite() {

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handlers.MainPage)
	serverMux.HandleFunc("/login", handlers.LoginPage)
	serverMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Starting web interface at: http://localhost:8080")
	http.ListenAndServe("localhost:8080", serverMux)

}
