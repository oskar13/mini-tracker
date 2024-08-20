package installer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func Run() error {
	log.Println("Starting Installer")

	serverDone := &sync.WaitGroup{}
	serverDone.Add(1)
	Start(serverDone)
	serverDone.Wait()
	log.Println("Shutting down installer")

	return nil
}

var ctxShutdown, cancel = context.WithCancel(context.Background())

func Start(wg *sync.WaitGroup) {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctxShutdown.Done():
			fmt.Println("Sorry: Shuting down ...")
			return
		default:
		}
		token := r.URL.Query().Get("token")
		if token != "" {
			fmt.Println("Found Token:", token)
			fmt.Println("Shuting down ...")
			// Shut down server here
			cancel() // to say sorry, above.
			// graceful-shutdown
			err := srv.Shutdown(context.Background())
			if err != nil {
				log.Println("server.Shutdown:", err)
			}

		} else {
			fmt.Fprintln(w, "Hi") // Server HTML page to fetch token and return to server at /callback
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
