package web

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/handlers"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
)

func StartWebsite() {

	// Initialize the database
	db.InitDB()
	defer db.Close()

	//Check for database tables,
	err := webutils.CheckInitData()

	if err != nil {
		log.Println("Failed to check initial data!")
		log.Fatal(err)
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", handlers.MainPage)
	serverMux.HandleFunc("/login", handlers.LoginPage)
	serverMux.HandleFunc("/logout", handlers.LoginPage)
	serverMux.HandleFunc("/signup", handlers.SignupPage)
	serverMux.HandleFunc("/profile/{id}/", handlers.ProfilePage)
	serverMux.HandleFunc("/profile", handlers.ProfilePage)
	serverMux.HandleFunc("/friends", handlers.FriendsPage)
	serverMux.HandleFunc("/dms/{id}/", handlers.DirectMessages)
	serverMux.HandleFunc("/dms", handlers.DirectMessages)
	serverMux.HandleFunc("/new", handlers.NewTorrentPage)
	serverMux.HandleFunc("/news/{id}/", handlers.NewsPage)
	serverMux.HandleFunc("/news", handlers.NewsPage)
	serverMux.HandleFunc("/groups/{id}/", handlers.GroupPage)
	serverMux.HandleFunc("/groups", handlers.GroupListPage)
	serverMux.HandleFunc("/groups/{groupid}/post/{postid}", handlers.GroupPostPage)
	serverMux.HandleFunc("/t/{uuid}/", handlers.TorrentPage)
	serverMux.HandleFunc("/t-dl/{uuid}/", handlers.TorrentDownloadPage)
	serverMux.HandleFunc("/my-t", handlers.MyTorrentsPage)
	serverMux.HandleFunc("/my-groups", handlers.MyGroupsPage)
	serverMux.HandleFunc("/cat/{catid}/", handlers.TorrentSearchPage)

	serverMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/web/static/"))))

	fmt.Println("Starting web interface at: http://localhost:8080")
	http.ListenAndServe(":8080", serverMux)

}
