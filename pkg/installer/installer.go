package installer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/oskar13/mini-tracker/pkg/data"
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
	})

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		fmt.Println("Bye.")
	}()
}
