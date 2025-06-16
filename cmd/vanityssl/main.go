package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/example/vanityssl/internal/api"
	"github.com/example/vanityssl/internal/proxy"
	"github.com/example/vanityssl/internal/store"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		log.Fatal("BACKEND_URL is required")
	}
	bURL, err := url.Parse(backendURL)
	if err != nil {
		log.Fatalf("invalid BACKEND_URL: %v", err)
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data"
	}
	backendStore, err := store.NewBadgerStore(dbPath)
	if err != nil {
		log.Fatalf("error opening store: %v", err)
	}
	defer backendStore.Close()

	cacheStore, err := store.NewCachedStore(backendStore, 1000)
	if err != nil {
		log.Fatalf("error creating cache: %v", err)
	}

	email := os.Getenv("ACME_EMAIL")
	cm := proxy.NewAutoCertManager(cacheStore, email, autocert.HostWhitelist())

	apiToken := os.Getenv("API_TOKEN")
	apiServer := api.New(cacheStore, apiToken)

	p := proxy.New(bURL, cacheStore, cm)

	r := mux.NewRouter()

	r.PathPrefix("/").Handler(p)

	go func() {
		log.Println("Internal API running on :8081")
		http.ListenAndServe(":8081", apiServer.Router())
	}()

	go func() {
		log.Println("ACME HTTP challenge on :80")
		http.ListenAndServe(":80", cm.HTTPHandler(nil))
	}()

  log.Println("Proxy running on :https (port 443)")
	server := &http.Server{
		Addr:      ":443",
		Handler:   cm.HTTPHandler(r),
		TLSConfig: cm.TLSConfig(),
	}
	log.Fatal(server.ListenAndServeTLS("", ""))
}
