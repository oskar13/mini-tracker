package installer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/oskar13/mini-tracker/pkg/data"
	webutils "github.com/oskar13/mini-tracker/pkg/web/webUtils"
)

var installerToken = data.ReadPassword(os.Getenv("INSTALLER_TOKEN_FILE"))

func Run() error {
	log.Println("Starting Installer")

	serverDone := &sync.WaitGroup{}
	serverDone.Add(1)
	Start(serverDone)
	serverDone.Wait()
	log.Println("Shutting down installer")

	return nil
}

func Start(wg *sync.WaitGroup) {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var pageStruct struct {
			Error       bool
			ErrorText   string
			Message     bool
			MessageText string
			Token       string
			NoToken     bool
		}

		token := r.URL.Query().Get("token")
		if token != "" {
			fmt.Println("Found Token:", token)

			if token == installerToken {
				fmt.Println("YAAY right token")
				pageStruct.Token = token

			} else {
				fmt.Print("BOOO!!!")
				fmt.Println(" Correct token is:")
				fmt.Println(installerToken)
				pageStruct.Error = true
				pageStruct.ErrorText = "Incorrect installer key."
			}

		} else {
			pageStruct.NoToken = true

		}
		webutils.RenderTemplate(w, []string{"pkg/installer/templates/installer.html"}, pageStruct)
	})

	http.HandleFunc("/install-success", func(w http.ResponseWriter, r *http.Request) {
		//Display success message after install, steps to do next
		token := r.URL.Query().Get("token")
		if token != "" {
			if token == installerToken {
				var emptyStruct struct{}
				//Render a welcome page after install
				webutils.RenderTemplate(w, []string{"pkg/installer/templates/installer_success.html"}, emptyStruct)
				// Shut down server here
				err := srv.Shutdown(context.Background())
				if err != nil {
					log.Println("server.Shutdown:", err)
				}
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./pkg/installer/static/"))))

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		fmt.Println("Bye.")
	}()
}
