package server

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"rap-lab/cmd/lab-studio/internal/auth"
	"rap-lab/cmd/lab-studio/internal/store"
	"rap-lab/cmd/lab-studio/web"
)

var studioTpl = template.Must(template.ParseFS(web.Files, "studio.html"))

func staticFiles() http.Handler {
	static, err := fs.Sub(web.Files, "static")
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

// Feldtypen mit Standardgroesse; unbekannter Typ = 400. Gleiche Werte in cmd/lab-studio/web/static/studio.js.
var fieldSizes = map[string][2]float64{
	"note":  {240, 160},
	"songs": {640, 420},
}

func studioPageHandler(st *store.Store, csrf *auth.CSRF) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, _, loggedIn := currentUser(r, st); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		data := pageData{CSRFToken: csrf.Token(w, r)}
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
func studioAPI(st *store.Store, csrf *auth.CSRF, next studioHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _, loggedIn := currentUser(r, st)
		if !loggedIn {
			writeJSON(w, http.StatusUnauthorized, errMsg("nicht angemeldet"))
			return
		}
		if r.Method != http.MethodGet && !csrf.ValidHeader(r, csrfHeaderKey) {
			writeJSON(w, http.StatusForbidden, errMsg("ungueltiges csrf-token"))
			return
		}
		next(w, r, u.Email)
	}
}

// Das CSRF-Token laeuft nach 30 Minuten ab; die offene Studio-Seite holt sich hier ein neues.
func handleStudioCSRF(csrf *auth.CSRF) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		writeJSON(w, http.StatusOK, map[string]string{"token": csrf.Token(w, r)})
	}
}

// ── Songs ──────────────────────────────────────────────────────────────────────────────────

func handleListSongs(st *store.Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		songs, err := st.ListSongs(owner)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, songs)
	}
}

func handleCreateSong(st *store.Store) studioHandler {
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
		song, err := st.CreateSong(owner, title)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, song)
	}
}

func handleGetSong(st *store.Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}
		song, err := st.GetSong(owner, id)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, song)
	}
}

func handleSaveSong(st *store.Store) studioHandler {
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
		if err := st.SaveSongText(owner, id, body.Text); err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// ── Arbeitsflaeche ─────────────────────────────────────────────────────────────────────────

func handleWorkspace(st *store.Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		view, err := st.LoadView(owner)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		fields, err := st.ListFields(owner)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"view": view, "fields": fields})
	}
}

func handleSaveView(st *store.Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		view, err := st.LoadView(owner)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		if !readJSON(w, r, &view) {
			return
		}
		view.Scale = clampScale(view.Scale)
		if err := st.SaveView(owner, view); err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, view)
	}
}

func handleCreateField(st *store.Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		var f store.Field
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
		if !validFieldSong(w, st, owner, &f) {
			return
		}
		created, err := st.CreateField(owner, f)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	}
}

// handleUpdateField uebernimmt nur die mitgeschickten Werte; ID und Typ bleiben unveraendert.
func handleUpdateField(st *store.Store) studioHandler {
	return func(w http.ResponseWriter, r *http.Request, owner string) {
		id, ok := pathID(w, r)
		if !ok {
			return
		}
		f, err := st.GetField(owner, id)
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
		if !validFieldSong(w, st, owner, &f) {
			return
		}
		if err := st.UpdateField(owner, f); err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, f)
	}
}

// validFieldSong: Nur ein SONGS-Feld darf auf einen Song zeigen, und nur auf einen eigenen.
func validFieldSong(w http.ResponseWriter, st *store.Store, owner string, f *store.Field) bool {
	if f.Type != "songs" {
		f.SongID = 0
		return true
	}
	if f.SongID == 0 {
		return true
	}
	if _, err := st.GetSong(owner, f.SongID); err != nil {
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
	if errors.Is(err, store.ErrNotFound) {
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
