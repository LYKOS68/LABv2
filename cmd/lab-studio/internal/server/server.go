// Package server enthaelt die HTTP-Seite von RAP-LAB: Routen, Seiten und Studio-API.
package server

import (
	"net/http"

	"rap-lab/cmd/lab-studio/internal/auth"
	"rap-lab/cmd/lab-studio/internal/mail"
	"rap-lab/cmd/lab-studio/internal/store"
)

// NewMux registriert alle Routen. Server und Tests nutzen dieselbe Liste.
func NewMux(st *store.Store, csrf *auth.CSRF, limiter *auth.RateLimiter, mailer *mail.Mailer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler(st, csrf))
	mux.HandleFunc("/register", registerHandler(st, csrf, limiter, mailer))
	mux.HandleFunc("/verify", verifyHandler(st, csrf))
	mux.HandleFunc("/resend-verification", resendVerificationHandler(st, csrf, limiter, mailer))
	mux.HandleFunc("/forgot-password", forgotPasswordHandler(st, csrf, limiter, mailer))
	mux.HandleFunc("/reset-password", resetPasswordHandler(st, csrf, limiter))
	mux.HandleFunc("/login", loginHandler(st, csrf, limiter))
	mux.HandleFunc("/logout", logoutHandler(st, csrf))
	mux.HandleFunc("/delete-account", deleteAccountHandler(st, csrf))

	mux.Handle("GET /static/", staticFiles())
	mux.HandleFunc("GET /connect", connectPageHandler())
	mux.HandleFunc("GET /studio", studioPageHandler(st, csrf))
	mux.HandleFunc("GET /studio/api/csrf", studioAPI(st, csrf, handleStudioCSRF(csrf)))
	mux.HandleFunc("GET /studio/api/songs", studioAPI(st, csrf, handleListSongs(st)))
	mux.HandleFunc("POST /studio/api/songs", studioAPI(st, csrf, handleCreateSong(st)))
	mux.HandleFunc("GET /studio/api/songs/{id}", studioAPI(st, csrf, handleGetSong(st)))
	mux.HandleFunc("PUT /studio/api/songs/{id}", studioAPI(st, csrf, handleSaveSong(st)))
	mux.HandleFunc("GET /studio/api/workspace", studioAPI(st, csrf, handleWorkspace(st)))
	mux.HandleFunc("PUT /studio/api/workspace/view", studioAPI(st, csrf, handleSaveView(st)))
	mux.HandleFunc("POST /studio/api/fields", studioAPI(st, csrf, handleCreateField(st)))
	mux.HandleFunc("PATCH /studio/api/fields/{id}", studioAPI(st, csrf, handleUpdateField(st)))
	return mux
}
