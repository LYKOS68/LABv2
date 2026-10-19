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
	mux := newMux(store, newRateLimiter(), newMailerFromEnv())

	log.Println("RAP-LAB laeuft auf http://127.0.0.1:8090")
	log.Fatal(http.ListenAndServe("127.0.0.1:8090", mux))
}

// newMux registriert alle Routen. Server und Tests nutzen dieselbe Liste.
func newMux(store *Store, limiter *RateLimiter, mailer *Mailer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", homeHandler(store))
	mux.HandleFunc("GET /register", pageHandler(store, "content-register"))
	mux.HandleFunc("POST /register", registerHandler(store, limiter, mailer))
	mux.HandleFunc("GET /verify", verifyHandler(store))
	mux.HandleFunc("POST /resend-verification", resendVerificationHandler(store, limiter, mailer))
	mux.HandleFunc("GET /forgot-password", pageHandler(store, "content-forgot"))
	mux.HandleFunc("POST /forgot-password", forgotPasswordHandler(store, limiter, mailer))
	mux.HandleFunc("GET /reset-password", resetPasswordPageHandler(store))
	mux.HandleFunc("POST /reset-password", resetPasswordHandler(store, limiter))
	mux.HandleFunc("GET /login", pageHandler(store, "content-login"))
	mux.HandleFunc("POST /login", loginHandler(store, limiter))
	mux.HandleFunc("POST /logout", logoutHandler(store))
	mux.HandleFunc("POST /delete-account", deleteAccountHandler(store))

	mux.Handle("GET /static/", staticFiles())
	mux.HandleFunc("GET /connect", connectPageHandler())
	mux.HandleFunc("GET /studio", studioPageHandler(store))
	mux.HandleFunc("GET /studio/api/csrf", studioAPI(store, handleStudioCSRF(store)))
	mux.HandleFunc("GET /studio/api/songs", studioAPI(store, handleListSongs(store)))
	mux.HandleFunc("POST /studio/api/songs", studioAPI(store, handleCreateSong(store)))
	mux.HandleFunc("GET /studio/api/songs/{id}", studioAPI(store, handleGetSong(store)))
	mux.HandleFunc("PUT /studio/api/songs/{id}", studioAPI(store, handleSaveSong(store)))
	mux.HandleFunc("GET /studio/api/workspace", studioAPI(store, handleWorkspace(store)))
	mux.HandleFunc("PUT /studio/api/workspace/view", studioAPI(store, handleSaveView(store)))
	mux.HandleFunc("POST /studio/api/fields", studioAPI(store, handleCreateField(store)))
	mux.HandleFunc("PATCH /studio/api/fields/{id}", studioAPI(store, handleUpdateField(store)))
	return mux
}
