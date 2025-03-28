package web

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/installer"
	adminHandler "github.com/oskar13/mini-tracker/pkg/web/adminPanel/handlers"
	"github.com/oskar13/mini-tracker/pkg/web/handlers"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
)

func StartWebsite() {

	// Initialize the database
	db.InitDB()
	defer db.Close()

	//Check for database tables,

	if err := checkInitData(); err != nil {
		log.Println("Failed to check initial data!")
		log.Fatal(err)
	}

	if err := webutils.LoadSiteData(); err != nil {
		log.Println("Failed to load site data!")
		log.Fatal(err)
	}

	r := gin.Default()

	serverMux := http.NewServeMux()
	r.GET("/", handlers.MainPage)
	r.GET("/login", handlers.LoginPage)
	r.GET("/logout", handlers.LoginPage)
	r.GET("/signup", handlers.SignupPage)
	r.GET("/profile/{id}/", handlers.ProfilePage)
	r.GET("/profile", handlers.ProfilePage)
	r.GET("/friends", handlers.FriendsPage)
	r.GET("/dms/{id}/", handlers.DirectMessages)
	r.GET("/dms", handlers.DirectMessages)
	r.GET("/new", handlers.NewTorrentPage)
	r.GET("/news/{id}/", handlers.NewsPage)
	r.GET("/news", handlers.NewsPage)
	r.GET("/groups/{id}/", handlers.GroupPage)
	r.GET("/groups", handlers.GroupListPage)
	r.GET("/groups/{groupid}/post/{postid}", handlers.GroupPostPage)
	r.GET("/t/{uuid}/", handlers.TorrentPage)
	r.GET("/t-dl/{uuid}/", handlers.TorrentDownloadPage)
	r.GET("/my-t", handlers.MyTorrentsPage)
	r.GET("/my-groups", handlers.MyGroupsPage)
	r.GET("/cat/{catid}/", handlers.TorrentSearchPage)
	//Admin panel routes
	r.GET("/admin/main", adminHandler.MainPage)
	r.GET("/admin/reports", adminHandler.ReportsPage)
	r.GET("/admin/user-list", adminHandler.UserListPage)
	r.GET("/admin/site-news", adminHandler.SiteNewsPage)
	r.GET("/admin/site-settings", adminHandler.SiteSettingsPage)

	//Static content
	serverMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/web/static/"))))

	fmt.Println("Starting web interface at: http://localhost:8080")
	http.ListenAndServe(":8080", serverMux)

}

func checkInitData() error {
	//Validity of user table
	err := webutils.ValidateSchema()
	if err != nil {
		//Revision number miss match
		if err == sql.ErrNoRows {
			return errors.New("database exists but schema revision number does not match, please see guide to migrate your current schema, quitting")
		} else {
			//Probably need to freshly initialize database
			return installer.Run()
		}

	}

	//Check if any user exists, if not run the installer
	tableEmpty, err := webutils.CheckUsers()
	if err != nil {
		return err
	}

	if tableEmpty {
		return installer.Run()
	}

	return nil
}
