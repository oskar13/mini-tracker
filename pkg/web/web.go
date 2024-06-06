package web

import (
	"fmt"
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/handlers"
)

func StartWebsite() {

	// Initialize the database
	db.InitDB()
	defer db.Close()

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handlers.MainPage)
	serverMux.HandleFunc("/login", handlers.LoginPage)
	serverMux.HandleFunc("/logout", handlers.LoginPage)
	serverMux.HandleFunc("/signup", handlers.SignupPage)
	serverMux.HandleFunc("/profile/{id}/", handlers.ProfilePage)
	serverMux.HandleFunc("/profile", handlers.ProfilePage)
	serverMux.HandleFunc("/friends", handlers.FriendsPage)
	serverMux.HandleFunc("/dms", handlers.DirectMessages)

	serverMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/web/static/"))))

	fmt.Println("Starting web interface at: http://localhost:8080")
	http.ListenAndServe("localhost:8080", serverMux)

}
