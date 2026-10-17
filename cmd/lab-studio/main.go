package main

import (
	"log"
	"net/http"
)

func main() {
	store, err := newStore("data/raplab.sqlite")
	if err != nil {
		log.Fatalf("store init failed: %v", err)
	}
	defer store.db.Close()
	limiter := newRateLimiter()
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler(store))
	mux.HandleFunc("/register", registerHandler(store, limiter))
	mux.HandleFunc("/login", loginHandler(store, limiter))
	mux.HandleFunc("/logout", logoutHandler(store))
	mux.HandleFunc("/delete-account", deleteAccountHandler(store))

	log.Println("RAP-LAB laeuft auf http://127.0.0.1:8090")
	log.Fatal(http.ListenAndServe("127.0.0.1:8090", mux))
}
