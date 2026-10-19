package main

import (
	"log"
	"net/http"

	"rap-lab/cmd/lab-studio/internal/auth"
	"rap-lab/cmd/lab-studio/internal/mail"
	"rap-lab/cmd/lab-studio/internal/server"
	"rap-lab/cmd/lab-studio/internal/store"
)

func main() {
	st, err := store.Open("data/raplab.sqlite")
	if err != nil {
		log.Fatalf("store init failed: %v", err)
	}
	defer st.Close()
	mux := server.NewMux(st, auth.NewCSRF(), auth.NewRateLimiter(), mail.NewFromEnv())

	log.Println("RAP-LAB laeuft auf http://127.0.0.1:8090")
	log.Fatal(http.ListenAndServe("127.0.0.1:8090", mux))
}
