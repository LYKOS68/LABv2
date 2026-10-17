package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"
	_ "modernc.org/sqlite"
)

const csrfTokenTTL = 30 * time.Minute

const (
	rateLimitAttempts = 5
	rateLimitWindow   = 10 * time.Minute
	rateLimitMaxKeys  = 10000 // Obergrenze gegen unbegrenztes Map-Wachstum
)

type rateLimitEntry struct {
	count       int
	windowStart time.Time
}

// RateLimiter begrenzt Versuche pro Schluessel (z.B. "login:"+IP) auf N pro Zeitfenster.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
}

func newRateLimiter() *RateLimiter {
	return &RateLimiter{entries: make(map[string]*rateLimitEntry)}
}

// Allow prueft und zaehlt einen Versuch fuer key. allowed=false liefert zusaetzlich die
// verbleibende Wartezeit bis zum naechsten erlaubten Versuch.
func (rl *RateLimiter) Allow(key string) (allowed bool, retryAfter time.Duration) {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for k, e := range rl.entries {
		if now.Sub(e.windowStart) >= rateLimitWindow {
			delete(rl.entries, k)
		}
	}

	e, ok := rl.entries[key]
	if !ok {
		if len(rl.entries) >= rateLimitMaxKeys {
			return false, rateLimitWindow
		}
		rl.entries[key] = &rateLimitEntry{count: 1, windowStart: now}
		return true, 0
	}
	if e.count >= rateLimitAttempts {
		return false, rateLimitWindow - now.Sub(e.windowStart)
	}
	e.count++
	return true, 0
}

// allowOrTooManyRequests prueft das Rate-Limit fuer key. Bei Ueberschreitung wird 429 mit
// Retry-After geschrieben und false zurueckgegeben; der Aufrufer darf dann nicht fortfahren.
func allowOrTooManyRequests(w http.ResponseWriter, limiter *RateLimiter, key string) bool {
	allowed, retryAfter := limiter.Allow(key)
	if !allowed {
		w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
		http.Error(w, "too many attempts, try again later", http.StatusTooManyRequests)
		return false
	}
	return true
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

const (
	argonMemoryKiB    = 19 * 1024
	argonIterations   = 2
	argonParallelism  = 1
	argonSaltLen      = 16
	argonKeyLen       = 32
	argonMaxMemoryKiB = 1 << 20 // Obergrenze gegen ungebundenen Aufwand aus korrupten/boeswilligen Hash-Strings
	argonMaxIterations = 10
	argonMaxParallelism = 16
)

func hashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, argonIterations, argonMemoryKiB, argonParallelism, argonKeyLen)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemoryKiB, argonIterations, argonParallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key)), nil
}

func verifyPassword(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid hash format")
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false, errors.New("invalid or unsupported hash version")
	}
	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false, errors.New("invalid hash parameters")
	}
	if memory == 0 || memory > argonMaxMemoryKiB || iterations == 0 || iterations > argonMaxIterations || parallelism == 0 || int(parallelism) > argonMaxParallelism {
		return false, errors.New("hash parameters out of bounds")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) == 0 || len(salt) > 64 {
		return false, errors.New("invalid salt")
	}
	storedKey, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(storedKey) == 0 || len(storedKey) > 128 {
		return false, errors.New("invalid key")
	}
	computedKey := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(storedKey)))
	return subtle.ConstantTimeCompare(computedKey, storedKey) == 1, nil
}

type User struct {
	Email        string
	PasswordHash string
}

type csrfEntry struct {
	token   string
	expires time.Time
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
	return &Store{db: db, sessions: make(map[string]string), csrf: make(map[string]csrfEntry)}, nil
}

var (
	errAccountExists      = errors.New("Konto mit dieser E-Mail existiert bereits.")
	errInvalidCredentials = errors.New("E-Mail oder Passwort ist falsch.")
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
	err := s.db.QueryRow("SELECT email, password_hash FROM users WHERE email = ? COLLATE NOCASE", email).Scan(&storedEmail, &hash)
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
	return &User{Email: storedEmail, PasswordHash: hash}, nil
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
	s.mu.Lock()
	for tok, e := range s.sessions {
		if e == key {
			delete(s.sessions, tok)
		}
	}
	s.mu.Unlock()
	return nil
}

// NewCSRFToken erzeugt eine neue Cookie-ID und ein zugehoeriges Token, 30 Minuten gueltig.
func (s *Store) NewCSRFToken() (cookieID, token string, err error) {
	cookieID, err = randomToken()
	if err != nil {
		return "", "", err
	}
	token, err = randomToken()
	if err != nil {
		return "", "", err
	}
	s.mu.Lock()
	s.csrf[cookieID] = csrfEntry{token: token, expires: time.Now().Add(csrfTokenTTL)}
	s.mu.Unlock()
	return cookieID, token, nil
}

// CSRFTokenFor liefert das noch gueltige Token zu einer bestehenden Cookie-ID, falls vorhanden.
func (s *Store) CSRFTokenFor(cookieID string) (string, bool) {
	s.mu.Lock()
	e, ok := s.csrf[cookieID]
	s.mu.Unlock()
	if !ok || time.Now().After(e.expires) {
		return "", false
	}
	return e.token, true
}

// ValidateCSRF prueft Cookie-ID gegen eingereichtes Token: existent, nicht abgelaufen, zeitkonstant gleich.
func (s *Store) ValidateCSRF(cookieID, token string) bool {
	if cookieID == "" || token == "" {
		return false
	}
	s.mu.Lock()
	e, ok := s.csrf[cookieID]
	s.mu.Unlock()
	if !ok || time.Now().After(e.expires) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(e.token), []byte(token)) == 1
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func validateRegistration(email, password, confirm string) string {
	email = strings.TrimSpace(email)
	if email == "" || !emailRe.MatchString(email) {
		return "Bitte eine gueltige E-Mail-Adresse angeben."
	}
	if len(password) < 8 {
		return "Passwort muss mindestens 8 Zeichen haben."
	}
	if password != confirm {
		return "Passwoerter stimmen nicht ueberein."
	}
	return ""
}

const baseCSS = `
body{background:#0a0a0a;color:#f2f2f2;font-family:Arial,Helvetica,sans-serif;margin:0;padding:0}
header{display:flex;justify-content:space-between;align-items:center;padding:16px 24px;border-bottom:1px solid #222}
header a{color:#d4af37;text-decoration:none;font-weight:bold;letter-spacing:1px;margin-left:24px}
.brand{color:#d4af37;font-weight:bold;letter-spacing:1px}
main{max-width:640px;margin:0 auto;padding:32px 24px}
h1{color:#d4af37;text-align:center}
.tagline{color:#d4af37;text-align:center;font-weight:bold;text-transform:uppercase}
.points p{color:#fff;text-align:center;font-weight:bold;text-transform:uppercase;margin:4px 0}
.fliesstext{background:#141414;padding:24px;white-space:pre-line}
.ende{color:#d4af37;text-align:center;font-weight:bold;margin-top:24px}
form{display:flex;flex-direction:column;gap:12px;max-width:360px;margin:32px auto}
input{padding:10px;background:#111;border:1px solid #333;color:#fff}
button{padding:10px;background:#d4af37;border:none;color:#000;font-weight:bold;cursor:pointer}
.error{color:#ff6b6b;text-align:center}
.ok{color:#8fe388;text-align:center}
.account-bar{text-align:center;margin-top:24px}
.nav-btn{background:none;border:none;color:#d4af37;text-decoration:none;font-weight:bold;letter-spacing:1px;margin-left:24px;padding:0;cursor:pointer;font-size:inherit;font-family:inherit}
`

var tpl = template.Must(template.New("base").Parse(`
{{define "layout"}}
<!doctype html>
<html lang="de">
<head><meta charset="utf-8"><title>RAP-LAB</title><style>` + baseCSS + `</style></head>
<body>
<header>
  <span class="brand">RAP-LAB</span>
  <nav>
    <a href="#">LAB-STUDIO</a>
    <a href="#">LAB-CONNECT</a>
    {{if .LoggedIn}}
      <form method="post" action="/logout" style="display:inline">
        <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
        <button type="submit" class="nav-btn">Logout</button>
      </form>
    {{else}}
      <a href="/register">Registrieren</a>
      <a href="/login">Login</a>
    {{end}}
  </nav>
</header>
<main>{{template "content" .}}</main>
</body>
</html>
{{end}}

{{define "content-home"}}
<h1>RAP-LAB</h1>
<p class="tagline">Ein Ort fuer das, was zwischen Idee und Musik passiert.</p>
<div class="points">
<p>Musik entwickeln.</p>
<p>Menschen finden.</p>
<p>Zusammenarbeiten.</p>
<p>Weiterkommen.</p>
</div>
<div class="fliesstext">Denn mal ehrlich:

Du kannst den naechsten Song auch wieder in deinen Handy-Notizen schreiben.
Den Producer kannst du wieder ueber fuenf Ecken suchen.
Den Mixer wieder in irgendeiner Instagram-Story finden.
Den Videografen wieder irgendwo in einem Discord-Server suchen.
Und deine Brancheninformationen kannst du dir wieder aus zehn verschiedenen Quellen zusammensuchen.
Oder du hoerst irgendwann auf, deinen Musikprozess aus Einzelteilen zusammenzubauen.</div>
<p class="ende">RAP-LAB.</p>
{{if .LoggedIn}}
<div class="account-bar">
  Angemeldet als {{.Email}} &middot;
  <form method="post" action="/delete-account" style="display:inline">
    <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
    <button type="submit" onclick="return confirm('Account wirklich endgueltig loeschen?')">Account loeschen</button>
  </form>
</div>
{{end}}
{{end}}

{{define "content-register"}}
<h1>Registrieren</h1>
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
<form method="post" action="/register">
  <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
  <input type="email" name="email" placeholder="E-Mail" value="{{.Email}}" required>
  <input type="password" name="password" placeholder="Passwort (min. 8 Zeichen)" required>
  <input type="password" name="confirm" placeholder="Passwort wiederholen" required>
  <button type="submit">Konto erstellen</button>
</form>
{{end}}

{{define "content-login"}}
<h1>Login</h1>
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
<form method="post" action="/login">
  <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
  <input type="email" name="email" placeholder="E-Mail" value="{{.Email}}" required>
  <input type="password" name="password" placeholder="Passwort" required>
  <button type="submit">Anmelden</button>
</form>
{{end}}
`))

type pageData struct {
	LoggedIn  bool
	Email     string
	Error     string
	CSRFToken string
}

const csrfCookieName = "csrf_id"

// csrfTokenFor liefert ein gueltiges CSRF-Token fuer diesen Request: wiederverwendet ein
// bestehendes gueltiges Cookie, erzeugt sonst ein neues und setzt das Cookie.
func csrfTokenFor(w http.ResponseWriter, r *http.Request, store *Store) string {
	if c, err := r.Cookie(csrfCookieName); err == nil {
		if token, ok := store.CSRFTokenFor(c.Value); ok {
			return token
		}
	}
	cookieID, token, err := store.NewCSRFToken()
	if err != nil {
		log.Printf("csrf token generation failed: %v", err)
		return ""
	}
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    cookieID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(csrfTokenTTL.Seconds()),
	})
	return token
}

// checkCSRF validiert das csrf_token-Formularfeld gegen das csrf_id-Cookie. Bei Fehlschlag
// wird 403 geschrieben und false zurueckgegeben; der Aufrufer darf die Mutation dann nicht ausfuehren.
func checkCSRF(w http.ResponseWriter, r *http.Request, store *Store) bool {
	cookieID := ""
	if c, err := r.Cookie(csrfCookieName); err == nil {
		cookieID = c.Value
	}
	token := r.FormValue("csrf_token")
	if !store.ValidateCSRF(cookieID, token) {
		http.Error(w, "invalid or missing csrf token", http.StatusForbidden)
		return false
	}
	return true
}

func currentUser(r *http.Request, store *Store) (*User, string, bool) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return nil, "", false
	}
	u, ok := store.UserByToken(c.Value)
	if !ok {
		return nil, "", false
	}
	return u, c.Value, true
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

func main() {
	store, err := newStore("data/raplab.sqlite")
	if err != nil {
		log.Fatalf("store init failed: %v", err)
	}
	defer store.db.Close()
	limiter := newRateLimiter()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		u, _, loggedIn := currentUser(r, store)
		data := pageData{LoggedIn: loggedIn}
		if loggedIn {
			data.Email = u.Email
		}
		render(w, r, store, "content-home", data)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, r, store, "content-register", pageData{})
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !checkCSRF(w, r, store) {
			return
		}
		if !allowOrTooManyRequests(w, limiter, "register:"+clientIP(r)) {
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirm := r.FormValue("confirm")
		if msg := validateRegistration(email, password, confirm); msg != "" {
			render(w, r, store, "content-register", pageData{Error: msg, Email: email})
			return
		}
		if err := store.CreateUser(email, password); err != nil {
			render(w, r, store, "content-register", pageData{Error: err.Error(), Email: email})
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, r, store, "content-login", pageData{})
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !checkCSRF(w, r, store) {
			return
		}
		if !allowOrTooManyRequests(w, limiter, "login:"+clientIP(r)) {
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		u, err := store.Authenticate(email, password)
		if err != nil {
			render(w, r, store, "content-login", pageData{Error: err.Error(), Email: email})
			return
		}
		token, err := store.CreateSession(u.Email)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !checkCSRF(w, r, store) {
			return
		}
		if c, err := r.Cookie("session_token"); err == nil {
			store.DeleteSession(c.Value)
		}
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("/delete-account", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !checkCSRF(w, r, store) {
			return
		}
		u, token, loggedIn := currentUser(r, store)
		if !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if err := store.DeleteUser(u.Email); err != nil {
			log.Printf("delete-account error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		store.DeleteSession(token)
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("RAP-LAB laeuft auf http://127.0.0.1:8090")
	log.Fatal(http.ListenAndServe("127.0.0.1:8090", mux))
}
