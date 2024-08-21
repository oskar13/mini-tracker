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
		}

		token := r.URL.Query().Get("token")
		if token != "" {
			fmt.Println("Found Token:", token)

			if token == installerToken {
				fmt.Println("YAAY right token")

				fmt.Println("Shuting down ...")
				// Shut down server here
				err := srv.Shutdown(context.Background())
				if err != nil {
					log.Println("server.Shutdown:", err)
				}
			} else {
				fmt.Print("BOOO!!!")
				fmt.Println(" Correct token is:")
				fmt.Println(installerToken)
			}

		} else {
			fmt.Fprintln(w, "Invalid  installer token") // Server HTML page to fetch token and return to server at /callback

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

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		fmt.Println("Bye.")
	}()
}
