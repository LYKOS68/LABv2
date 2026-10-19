package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

var studioTpl = template.Must(template.ParseFS(webFiles, "web/studio.html"))

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

// ── Datenbank ────────────────────────────────────────────────────────────────────────────────
// Jede Abfrage enthaelt "owner = ?": ein Konto sieht und aendert nur seine eigenen Daten.

var errNotFound = errors.New("nicht gefunden")

type Song struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type View struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Scale float64 `json:"scale"`
}

type Field struct {
	ID     int64   `json:"id"`
	Type   string  `json:"type"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	W      float64 `json:"w"`
	H      float64 `json:"h"`
	Text   string  `json:"text"`
	SongID int64   `json:"song_id"`
}

func createStudioTables(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS songs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			owner TEXT NOT NULL COLLATE NOCASE,
			title TEXT NOT NULL,
			text TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS studio_fields (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			owner TEXT NOT NULL COLLATE NOCASE,
			type TEXT NOT NULL,
			x REAL NOT NULL,
			y REAL NOT NULL,
			w REAL NOT NULL,
			h REAL NOT NULL,
			text TEXT NOT NULL DEFAULT '',
			song_id INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS studio_view (
			owner TEXT PRIMARY KEY COLLATE NOCASE,
			x REAL NOT NULL,
			y REAL NOT NULL,
			scale REAL NOT NULL
		)`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListSongs(owner string) ([]Song, error) {
	rows, err := s.db.Query("SELECT id, title FROM songs WHERE owner = ? ORDER BY title COLLATE NOCASE", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	songs := []Song{}
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.ID, &song.Title); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, rows.Err()
}

func (s *Store) CreateSong(owner, title string) (Song, error) {
	result, err := s.db.Exec("INSERT INTO songs (owner, title) VALUES (?, ?)", owner, title)
	if err != nil {
		return Song{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Song{}, err
	}
	return Song{ID: id, Title: title}, nil
}

func (s *Store) GetSong(owner string, id int64) (Song, error) {
	var song Song
	err := s.db.QueryRow("SELECT id, title, text FROM songs WHERE id = ? AND owner = ?", id, owner).
		Scan(&song.ID, &song.Title, &song.Text)
	if errors.Is(err, sql.ErrNoRows) {
		return Song{}, errNotFound
	}
	return song, err
}

func (s *Store) SaveSongText(owner string, id int64, text string) error {
	result, err := s.db.Exec("UPDATE songs SET text = ? WHERE id = ? AND owner = ?", text, id, owner)
	return requireOneRow(result, err)
}

// LoadView liefert die gespeicherte Ansicht oder die Startansicht (0, 0, 100 %).
func (s *Store) LoadView(owner string) (View, error) {
	view := View{Scale: 1}
	err := s.db.QueryRow("SELECT x, y, scale FROM studio_view WHERE owner = ?", owner).
		Scan(&view.X, &view.Y, &view.Scale)
	if errors.Is(err, sql.ErrNoRows) {
		return View{Scale: 1}, nil
	}
	return view, err
}

func (s *Store) SaveView(owner string, view View) error {
	_, err := s.db.Exec(`INSERT INTO studio_view (owner, x, y, scale) VALUES (?, ?, ?, ?)
		ON CONFLICT(owner) DO UPDATE SET x = excluded.x, y = excluded.y, scale = excluded.scale`,
		owner, view.X, view.Y, view.Scale)
	return err
}

const fieldColumns = "id, type, x, y, w, h, text, song_id"

func scanField(row interface{ Scan(...any) error }) (Field, error) {
	var f Field
	err := row.Scan(&f.ID, &f.Type, &f.X, &f.Y, &f.W, &f.H, &f.Text, &f.SongID)
	return f, err
}

func (s *Store) ListFields(owner string) ([]Field, error) {
	rows, err := s.db.Query("SELECT "+fieldColumns+" FROM studio_fields WHERE owner = ? ORDER BY id", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fields := []Field{}
	for rows.Next() {
		f, err := scanField(rows)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, rows.Err()
}

func (s *Store) GetField(owner string, id int64) (Field, error) {
	row := s.db.QueryRow("SELECT "+fieldColumns+" FROM studio_fields WHERE id = ? AND owner = ?", id, owner)
	f, err := scanField(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Field{}, errNotFound
	}
	return f, err
}

func (s *Store) CreateField(owner string, f Field) (Field, error) {
	result, err := s.db.Exec(`INSERT INTO studio_fields (owner, type, x, y, w, h, text, song_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, owner, f.Type, f.X, f.Y, f.W, f.H, f.Text, f.SongID)
	if err != nil {
		return Field{}, err
	}
	f.ID, err = result.LastInsertId()
	return f, err
}

func (s *Store) UpdateField(owner string, f Field) error {
	result, err := s.db.Exec(`UPDATE studio_fields SET x = ?, y = ?, w = ?, h = ?, text = ?, song_id = ?
		WHERE id = ? AND owner = ?`, f.X, f.Y, f.W, f.H, f.Text, f.SongID, f.ID, owner)
	return requireOneRow(result, err)
}

// requireOneRow meldet errNotFound, wenn ein UPDATE keine Zeile getroffen hat
// (unbekannte ID oder fremdes Konto).
func requireOneRow(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errNotFound
	}
	return nil
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
