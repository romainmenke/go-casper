package main

import (
	"log"
	"net/http"
	"path/filepath"

	casper "github.com/tcnksm/go-casper"
)

const (
	defaultPort = "3000"
)

func main() {
	// See README.md to generate crt and key for this testing.
	certFile, _ := filepath.Abs("crts/server.crt")
	keyFile, _ := filepath.Abs("crts/server.key")

	// Handle root
	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s %s", r.Method, r.URL.String())

		assets := []string{
			"/static/example.jpg",
		}

		pusher, ok := w.(http.Pusher)
		if !ok {
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(`<img src="/static/example.jpg"/>`))
			return
		}

		for _, asset := range assets {
			// Server push!
			if err := pusher.Push(asset, nil); err != nil {
				log.Fatalf("[ERROR] Failed to push asset %v: %s", asset, err)
			}
		}

		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`<img src="/static/example.jpg"/>`))
	})

	// Initialize casper.
	pusher := casper.Handler(1<<6, 10, root)

	http.Handle("/", pusher)

	// Handle static assets.
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s %s", r.Method, r.URL.String())
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
	})

	// Start listening.
	log.Println("[INFO] Start listening on port", ":"+defaultPort)
	if err := http.ListenAndServeTLS(":"+defaultPort, certFile, keyFile, nil); err != nil {
		log.Fatalf("[ERROR] Failed to listen: %s", err)
	}
}
