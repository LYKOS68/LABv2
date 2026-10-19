package main

import (
	"database/sql"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"rap-lab/web"
)

// Seiten, Vorlagen und static/ kommen eingebettet aus dem Paket web (Ordner web/ in der Wurzel).
var webFiles = web.Files

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

// Store haelt die Datenbank. Konto- und Studio-Abfragen haengen als Methoden in auth.go und studio.go.
type Store struct {
	db *sql.DB
}

func newStore(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return nil, err
	}
	if err := createAuthTables(db); err != nil {
		return nil, err
	}
	if err := createStudioTables(db); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// ── Seiten ──────────────────────────────────────────────────────────────────────────────────────

func staticFiles() http.Handler {
	static, err := fs.Sub(webFiles, "static")
	if err != nil {
		panic(err)
	}
	return http.StripPrefix("/static/", http.FileServer(http.FS(static)))
}

// Startseite und Konto-Seiten: layout.html ist der Rahmen, home.html und auth.html liefern
// die Inhalte als {{define "content-..."}}-Bloecke.
var tpl = template.Must(template.ParseFS(webFiles, "layout.html", "home.html", "auth.html"))

type pageData struct {
	LoggedIn   bool
	Email      string
	Error      string
	CSRFToken  string
	ShowResend bool
	Token      string
}

func render(w http.ResponseWriter, r *http.Request, store *Store, contentName string, data pageData) {
	data.CSRFToken = csrfTokenFor(w, r, store)
	clone, err := tpl.Clone()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	clone = template.Must(clone.Parse(`{{define "content"}}{{template "` + contentName + `" .}}{{end}}`))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := clone.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("render error: %v", err)
	}
}

func homeHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _, loggedIn := currentUser(r, store)
		data := pageData{LoggedIn: loggedIn}
		if loggedIn {
			data.Email = u.Email
		}
		render(w, r, store, "content-home", data)
	}
}

// LAB-CONNECT ist vorerst ein reiner Oberflaechen-Prototyp:
// die Seite laeuft komplett im Browser mit Beispielanzeigen und schickt nichts an den Server.
func connectPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, webFiles, "connect.html")
	}
}
