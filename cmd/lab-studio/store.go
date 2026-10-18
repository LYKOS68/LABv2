package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type User struct {
	Email        string
	PasswordHash string
}

type Store struct {
	db       *sql.DB
	mu       sync.Mutex
	sessions map[string]string    // token -> lowercase email, in-memory (nicht persistiert)
	csrf     map[string]csrfEntry // cookie-id -> Token+Ablauf, in-memory
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
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		email TEXT NOT NULL UNIQUE COLLATE NOCASE,
		password_hash TEXT NOT NULL
	)`); err != nil {
		return nil, err
	}
	if err := migrateVerifiedColumn(db); err != nil {
		return nil, err
	}
	if err := createStudioTables(db); err != nil {
		return nil, err
	}
	for _, table := range []string{verificationTable, resetTable} {
		if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS ` + table + ` (
			token_hash TEXT PRIMARY KEY,
			email TEXT NOT NULL COLLATE NOCASE,
			expires_at INTEGER NOT NULL
		)`); err != nil {
			return nil, err
		}
	}
	return &Store{db: db, sessions: make(map[string]string), csrf: make(map[string]csrfEntry)}, nil
}

// migrateVerifiedColumn fuegt users.verified hinzu, falls es fehlt. Konten, die schon vor
// der E-Mail-Bestaetigung existierten, werden dabei einmalig als bestaetigt markiert.
func migrateVerifiedColumn(db *sql.DB) error {
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name = 'verified'").Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	if _, err := db.Exec("ALTER TABLE users ADD COLUMN verified INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	_, err := db.Exec("UPDATE users SET verified = 1")
	return err
}

// Tabellen fuer Einmal-Links per Mail. Nur diese festen Namen gehen in SQL-Strings ein.
const (
	verificationTable = "email_verifications"
	resetTable        = "password_resets"
	verificationTTL   = 24 * time.Hour
	resetTTL          = 1 * time.Hour
)

var (
	errAccountExists      = errors.New("Konto mit dieser E-Mail existiert bereits.")
	errInvalidCredentials = errors.New("E-Mail oder Passwort ist falsch.")
	errNotVerified        = errors.New("Bitte bestaetige zuerst deine E-Mail-Adresse.")
	errInvalidToken       = errors.New("Der Link ist ungueltig oder abgelaufen.")
)

func (s *Store) CreateUser(email, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, hash)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errAccountExists
		}
		return err
	}
	return nil
}

func (s *Store) Authenticate(email, password string) (*User, error) {
	var storedEmail, hash string
	var verified bool
	err := s.db.QueryRow("SELECT email, password_hash, verified FROM users WHERE email = ? COLLATE NOCASE", email).Scan(&storedEmail, &hash, &verified)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	ok, err := verifyPassword(hash, password)
	if err != nil {
		log.Printf("verifyPassword error for %s: %v", storedEmail, err)
		return nil, errInvalidCredentials
	}
	if !ok {
		return nil, errInvalidCredentials
	}
	// Erst nach korrektem Passwort pruefen, damit der Hinweis nichts ueber fremde Konten verraet.
	if !verified {
		return nil, errNotVerified
	}
	return &User{Email: storedEmail, PasswordHash: hash}, nil
}

// createToken erzeugt einen Einmal-Token fuer email in table und ersetzt aeltere derselben Adresse.
// In der DB liegt nur der SHA-256-Hash; der Klartext-Token geht ausschliesslich in den Link.
func (s *Store) createToken(table, email string, ttl time.Duration) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}
	now := time.Now().Unix()
	if _, err := s.db.Exec("DELETE FROM "+table+" WHERE email = ? COLLATE NOCASE OR expires_at <= ?", email, now); err != nil {
		return "", err
	}
	_, err = s.db.Exec("INSERT INTO "+table+" (token_hash, email, expires_at) VALUES (?, ?, ?)",
		hashToken(token), email, time.Now().Add(ttl).Unix())
	if err != nil {
		return "", err
	}
	return token, nil
}

// consumeToken loest einen Token aus table genau einmal ein und liefert die zugehoerige E-Mail.
func (s *Store) consumeToken(table, token string) (string, error) {
	if token == "" {
		return "", errInvalidToken
	}
	var email string
	err := s.db.QueryRow("DELETE FROM "+table+" WHERE token_hash = ? AND expires_at > ? RETURNING email",
		hashToken(token), time.Now().Unix()).Scan(&email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errInvalidToken
	}
	return email, err
}

func (s *Store) CreateVerification(email string) (string, error) {
	return s.createToken(verificationTable, email, verificationTTL)
}

// VerifyEmail loest einen Token genau einmal ein und markiert das Konto als bestaetigt.
func (s *Store) VerifyEmail(token string) error {
	email, err := s.consumeToken(verificationTable, token)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE users SET verified = 1 WHERE email = ? COLLATE NOCASE", email)
	return err
}

func (s *Store) CreatePasswordReset(email string) (string, error) {
	return s.createToken(resetTable, email, resetTTL)
}

// ResetPassword setzt per Einmal-Token ein neues Passwort, markiert das Konto als bestaetigt
// (der Link beweist den Postfach-Besitz) und beendet alle offenen Sitzungen des Kontos.
func (s *Store) ResetPassword(token, newPassword string) error {
	hash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	email, err := s.consumeToken(resetTable, token)
	if err != nil {
		return err
	}
	if _, err := s.db.Exec("UPDATE users SET password_hash = ?, verified = 1 WHERE email = ? COLLATE NOCASE", hash, email); err != nil {
		return err
	}
	s.deleteSessionsOf(email)
	return nil
}

// AccountExists meldet, ob ein Konto mit dieser E-Mail existiert.
func (s *Store) AccountExists(email string) bool {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? COLLATE NOCASE", email).Scan(&n)
	return err == nil && n > 0
}

// IsUnverified meldet, ob ein Konto mit dieser E-Mail existiert und noch unbestaetigt ist.
func (s *Store) IsUnverified(email string) bool {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? COLLATE NOCASE AND verified = 0", email).Scan(&n)
	return err == nil && n > 0
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Store) CreateSession(email string) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.sessions[token] = strings.ToLower(email)
	s.mu.Unlock()
	return token, nil
}

func (s *Store) UserByToken(token string) (*User, bool) {
	s.mu.Lock()
	email, ok := s.sessions[token]
	s.mu.Unlock()
	if !ok {
		return nil, false
	}
	var storedEmail, hash string
	err := s.db.QueryRow("SELECT email, password_hash FROM users WHERE email = ? COLLATE NOCASE", email).Scan(&storedEmail, &hash)
	if err != nil {
		return nil, false
	}
	return &User{Email: storedEmail, PasswordHash: hash}, true
}

func (s *Store) DeleteSession(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func (s *Store) DeleteUser(email string) error {
	key := strings.ToLower(email)
	if _, err := s.db.Exec("DELETE FROM users WHERE email = ? COLLATE NOCASE", key); err != nil {
		return err
	}
	for _, table := range []string{verificationTable, resetTable} {
		if _, err := s.db.Exec("DELETE FROM "+table+" WHERE email = ? COLLATE NOCASE", key); err != nil {
			return err
		}
	}
	for _, table := range []string{"songs", "studio_fields", "studio_view"} {
		if _, err := s.db.Exec("DELETE FROM "+table+" WHERE owner = ?", key); err != nil {
			return err
		}
	}
	s.deleteSessionsOf(key)
	return nil
}

func (s *Store) deleteSessionsOf(email string) {
	key := strings.ToLower(email)
	s.mu.Lock()
	for tok, e := range s.sessions {
		if e == key {
			delete(s.sessions, tok)
		}
	}
	s.mu.Unlock()
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ── LAB-STUDIO ────────────────────────────────────────────────────────────────────────────
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
