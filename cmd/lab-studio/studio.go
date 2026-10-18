package main

import (
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// web/ wird ins Programm eingebaut: studio.html als Vorlage, static/ als Dateien unter /static/.
//
//go:embed web
var webFiles embed.FS

var studioTpl = template.Must(template.ParseFS(webFiles, "web/studio.html"))

func staticFiles() http.Handler {
	static, err := fs.Sub(webFiles, "web/static")
	if err != nil {
		panic(err)
	}
	return http.StripPrefix("/static/", http.FileServer(http.FS(static)))
}

const (
	minScale      = 0.25
	maxScale      = 4.0
	maxTitleLen   = 100
	maxTextBytes  = 1 << 20
	maxBodyBytes  = maxTextBytes + 4096 // Text plus JSON-Huelle
	csrfHeaderKey = "X-CSRF-Token"
)

// Feldtypen mit Standardgroesse; unbekannter Typ = 400. Gleiche Werte in web/static/studio.js.
var fieldSizes = map[string][2]float64{
	"note":  {240, 160},
	"songs": {640, 420},
}

func studioPageHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, _, loggedIn := currentUser(r, store); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		data := pageData{CSRFToken: csrfTokenFor(w, r, store)}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := studioTpl.Execute(w, data); err != nil {
			log.Printf("studio render error: %v", err)
		}
	}
}

// studioHandler bekommt das angemeldete Konto bereits geprueft uebergeben.
type studioHandler func(w http.ResponseWriter, r *http.Request, owner string)

// studioAPI prueft vor jedem Studio-Zugriff die Sitzung und bei Schreibzugriffen das
// CSRF-Token aus dem Header. Erst danach laeuft next.
func studioAPI(store *Store, next studioHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _, loggedIn := currentUser(r, store)
		if !loggedIn {
			writeJSON(w, http.StatusUnauthorized, errMsg("nicht angemeldet"))
			return
		}
		if r.Method != http.MethodGet && !validCSRFHeader(r, store) {
			writeJSON(w, http.StatusForbidden, errMsg("ungueltiges csrf-token"))
			return
		}
		next(w, r, u.Email)
	}
}

func validCSRFHeader(r *http.Request, store *Store) bool {
	cookieID := ""
	if c, err := r.Cookie(csrfCookieName); err == nil {
		cookieID = c.Value
	}
	return store.ValidateCSRF(cookieID, r.Header.Get(csrfHeaderKey))
}

// Das CSRF-Token laeuft nach 30 Minuten ab; die offene Studio-Seite holt sich hier ein neues.
func handleStudioCSRF(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		writeJSON(w, http.StatusOK, map[string]string{"token": csrfTokenFor(w, r, store)})
	}
}

// ── Songs ──────────────────────────────────────────────────────────────────────────────────

func handleListSongs(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		songs, err := store.ListSongs(owner)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, songs)
	}
}

func handleCreateSong(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		var body struct {
			Title string `json:"title"`
		}
		if !readJSON(w, r, &body) {
			return
		}
		title := strings.TrimSpace(body.Title)
		if title == "" || len(title) > maxTitleLen {
			writeJSON(w, http.StatusBadRequest, errMsg("titel fehlt oder ist zu lang"))
			return
		}
		song, err := store.CreateSong(owner, title)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, song)
	}
}

func handleGetSong(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}
		song, err := store.GetSong(owner, id)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, song)
	}
}

func handleSaveSong(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}
		var body struct {
			Text string `json:"text"`
		}
		if !readJSON(w, r, &body) {
			return
		}
		if len(body.Text) > maxTextBytes {
			writeJSON(w, http.StatusBadRequest, errMsg("text zu lang"))
			return
		}
		if err := store.SaveSongText(owner, id, body.Text); err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// ── Arbeitsflaeche ─────────────────────────────────────────────────────────────────────────

func handleWorkspace(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		view, err := store.LoadView(owner)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		fields, err := store.ListFields(owner)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"view": view, "fields": fields})
	}
}

func handleSaveView(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		view, err := store.LoadView(owner)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		if !readJSON(w, r, &view) {
			return
		}
		view.Scale = clampScale(view.Scale)
		if err := store.SaveView(owner, view); err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, view)
	}
}

func handleCreateField(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		var f Field
		if !readJSON(w, r, &f) {
			return
		}
		size, known := fieldSizes[f.Type]
		if !known {
			writeJSON(w, http.StatusBadRequest, errMsg("unbekannter type"))
			return
		}
		if f.W <= 0 || f.H <= 0 {
			f.W, f.H = size[0], size[1]
		}
		if !validFieldSong(w, store, owner, &f) {
			return
		}
		created, err := store.CreateField(owner, f)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	}
}

// handleUpdateField uebernimmt nur die mitgeschickten Werte; ID und Typ bleiben unveraendert.
func handleUpdateField(store *Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}
		f, err := store.GetField(owner, id)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		fieldType := f.Type
		if !readJSON(w, r, &f) {
			return
		}
		f.ID = id
		f.Type = fieldType
		if f.W <= 0 || f.H <= 0 {
			writeJSON(w, http.StatusBadRequest, errMsg("ungueltige groesse"))
			return
		}
		if !validFieldSong(w, store, owner, &f) {
			return
		}
		if err := store.UpdateField(owner, f); err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, f)
	}
}

// validFieldSong: Nur ein SONGS-Feld darf auf einen Song zeigen, und nur auf einen eigenen.
func validFieldSong(w http.ResponseWriter, store *Store, owner string, f *Field) bool {
	if f.Type != "songs" {
		f.SongID = 0
		return true
	}
	if f.SongID == 0 {
		return true
	}
	if _, err := store.GetSong(owner, f.SongID); err != nil {
		writeJSON(w, http.StatusBadRequest, errMsg("unbekannte song_id"))
		return false
	}
	return true
}

// ── Hilfsfunktionen ────────────────────────────────────────────────────────────────────────

func clampScale(s float64) float64 {
	if math.IsNaN(s) {
		return 1
	}
	return math.Max(minScale, math.Min(maxScale, s))
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errMsg("unbekannte id"))
		return 0, false
	}
	return id, true
}

func readJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(v)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errMsg("ungueltiger body"))
		return false
	}
	return true
}

func writeStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, errNotFound) {
		writeJSON(w, http.StatusNotFound, errMsg("nicht gefunden"))
		return
	}
	log.Printf("studio store error: %v", err)
	writeJSON(w, http.StatusInternalServerError, errMsg("interner fehler"))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func errMsg(text string) map[string]string {
	return map[string]string{"error": text}
}
