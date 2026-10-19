package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"
)

// Alles rund ums Konto: Registrieren mit Mail-Bestaetigung, Anmelden, Passwort vergessen,
// Konto loeschen - samt Sitzungen, CSRF-Schutz, Rate-Limit, Passwort-Hashing und Mailversand.

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

type User struct {
	Email        string
	PasswordHash string
}

func createAuthTables(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		email TEXT NOT NULL UNIQUE COLLATE NOCASE,
		password_hash TEXT NOT NULL
	)`); err != nil {
		return err
	}
	if err := migrateVerifiedColumn(db); err != nil {
		return err
	}
	for _, table := range []string{verificationTable, resetTable} {
		if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS ` + table + ` (
			token_hash TEXT PRIMARY KEY,
			email TEXT NOT NULL COLLATE NOCASE,
			expires_at INTEGER NOT NULL
		)`); err != nil {
			return err
		}
	}
	return nil
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

// ── Handler ─────────────────────────────────────────────────────────────────────────────────────

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

// sendVerification erzeugt einen Token und verschickt den Link. Fehler werden nur geloggt:
// das Konto existiert bereits, und der Link kann ueber /resend-verification neu angefordert werden.
func sendVerification(store *Store, mailer *Mailer, email string) {
	token, err := store.CreateVerification(email)
	if err != nil {
		log.Printf("create verification for %s: %v", email, err)
		return
	}
	if err := mailer.SendVerificationLink(email, token); err != nil {
		log.Printf("send verification mail to %s: %v", email, err)
	}
}

// forgotPasswordHandler antwortet immer gleich, egal ob die E-Mail existiert.
func forgotPasswordHandler(store *Store, limiter *RateLimiter, mailer *Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkCSRF(w, r, store) {
			return
		}
		if !allowOrTooManyRequests(w, limiter, "forgot:"+clientIP(r)) {
			return
		}
		email := r.FormValue("email")
		if store.AccountExists(email) {
			token, err := store.CreatePasswordReset(email)
			if err != nil {
				log.Printf("create password reset for %s: %v", email, err)
			} else if err := mailer.SendPasswordResetLink(email, token); err != nil {
				log.Printf("send password reset mail to %s: %v", email, err)
			}
		}
		render(w, r, store, "content-forgot-sent", pageData{Email: email})
	}
}

// pageHandler zeigt eine Seite ohne eigene Daten, z.B. ein leeres Formular.
func pageHandler(store *Store, contentName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, store, contentName, pageData{})
	}
}

// resetPasswordPageHandler zeigt nur das Formular; der Token wird erst beim POST mit
// gueltigem neuem Passwort verbraucht (Mail-Scanner oeffnen Links vorab).
func resetPasswordPageHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, store, "content-reset", pageData{Token: r.URL.Query().Get("token")})
	}
}

func resetPasswordHandler(store *Store, limiter *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkCSRF(w, r, store) {
			return
		}
		if !allowOrTooManyRequests(w, limiter, "reset:"+clientIP(r)) {
			return
		}
		token := r.FormValue("token")
		if msg := validateNewPassword(r.FormValue("password"), r.FormValue("confirm")); msg != "" {
			render(w, r, store, "content-reset", pageData{Error: msg, Token: token})
			return
		}
		if err := store.ResetPassword(token, r.FormValue("password")); err != nil {
			if err != errInvalidToken {
				log.Printf("reset password: %v", err)
			}
			render(w, r, store, "content-reset-result", pageData{Error: errInvalidToken.Error()})
			return
		}
		render(w, r, store, "content-reset-result", pageData{})
	}
}

func registerHandler(store *Store, limiter *RateLimiter, mailer *Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		sendVerification(store, mailer, email)
		render(w, r, store, "content-verify-sent", pageData{Email: email})
	}
}

func verifyHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := store.VerifyEmail(r.URL.Query().Get("token")); err != nil {
			if err != errInvalidToken {
				log.Printf("verify email: %v", err)
			}
			render(w, r, store, "content-verify-result", pageData{Error: errInvalidToken.Error()})
			return
		}
		render(w, r, store, "content-verify-result", pageData{})
	}
}

// resendVerificationHandler antwortet immer gleich, egal ob die E-Mail existiert, damit
// sich darueber nicht herausfinden laesst, welche Adressen registriert sind.
func resendVerificationHandler(store *Store, limiter *RateLimiter, mailer *Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkCSRF(w, r, store) {
			return
		}
		if !allowOrTooManyRequests(w, limiter, "resend:"+clientIP(r)) {
			return
		}
		email := r.FormValue("email")
		if store.IsUnverified(email) {
			sendVerification(store, mailer, email)
		}
		render(w, r, store, "content-verify-sent", pageData{Email: email})
	}
}

func loginHandler(store *Store, limiter *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			render(w, r, store, "content-login", pageData{Error: err.Error(), Email: email, ShowResend: err == errNotVerified})
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
	}
}

func logoutHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkCSRF(w, r, store) {
			return
		}
		if c, err := r.Cookie("session_token"); err == nil {
			store.DeleteSession(c.Value)
		}
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func deleteAccountHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// ── Konten und Einmal-Links ─────────────────────────────────────────────────────────────────────

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

// ── Eingaben pruefen ────────────────────────────────────────────────────────────────────────────

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func validateRegistration(email, password, confirm string) string {
	email = strings.TrimSpace(email)
	if email == "" || !emailRe.MatchString(email) {
		return "Bitte eine gueltige E-Mail-Adresse angeben."
	}
	return validateNewPassword(password, confirm)
}

func validateNewPassword(password, confirm string) string {
	if len(password) < 8 {
		return "Passwort muss mindestens 8 Zeichen haben."
	}
	if password != confirm {
		return "Passwoerter stimmen nicht ueberein."
	}
	return ""
}

// ── Passwort-Hashing (Argon2id) ─────────────────────────────────────────────────────────────────

const (
	argonMemoryKiB      = 19 * 1024
	argonIterations     = 2
	argonParallelism    = 1
	argonSaltLen        = 16
	argonKeyLen         = 32
	argonMaxMemoryKiB   = 1 << 20 // Obergrenze gegen ungebundenen Aufwand aus korrupten/boeswilligen Hash-Strings
	argonMaxIterations  = 10
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

// ── CSRF ────────────────────────────────────────────────────────────────────────────────────────

const csrfTokenTTL = 30 * time.Minute
const csrfCookieName = "csrf_id"

type csrfEntry struct {
	token   string
	expires time.Time
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

// ── Rate-Limit ──────────────────────────────────────────────────────────────────────────────────

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

// ── Mail ────────────────────────────────────────────────────────────────────────────────────────

// Mailer verschickt Bestaetigungslinks per SMTP. Ohne SMTP_HOST wird der Link nur in die
// Konsole geschrieben, damit der Ablauf lokal ohne Mail-Dienst testbar bleibt.
type Mailer struct {
	host, port, user, pass, from string
	baseURL                      string
}

func newMailerFromEnv() *Mailer {
	m := &Mailer{
		host:    os.Getenv("SMTP_HOST"),
		port:    os.Getenv("SMTP_PORT"),
		user:    os.Getenv("SMTP_USER"),
		pass:    os.Getenv("SMTP_PASS"),
		from:    os.Getenv("SMTP_FROM"),
		baseURL: os.Getenv("BASE_URL"),
	}
	if m.port == "" {
		m.port = "587"
	}
	if m.from == "" {
		m.from = m.user
	}
	if m.baseURL == "" {
		m.baseURL = "http://127.0.0.1:8090"
	}
	return m
}

func (m *Mailer) SendVerificationLink(email, token string) error {
	return m.sendLink(email, "E-Mail-Adresse bestaetigen", "/verify",
		"bitte bestaetige deine E-Mail-Adresse fuer RAP-LAB:",
		"Der Link ist 24 Stunden gueltig. Wenn du dich nicht registriert hast, ignoriere diese Mail.", token)
}

func (m *Mailer) SendPasswordResetLink(email, token string) error {
	return m.sendLink(email, "Passwort zuruecksetzen", "/reset-password",
		"ueber diesen Link kannst du ein neues Passwort fuer RAP-LAB setzen:",
		"Der Link ist 1 Stunde gueltig. Wenn du das nicht angefordert hast, ignoriere diese Mail - dein Passwort bleibt unveraendert.", token)
}

func (m *Mailer) sendLink(email, subject, path, intro, outro, token string) error {
	link := m.baseURL + path + "?token=" + url.QueryEscape(token)
	if m.host == "" {
		log.Printf("[mail] kein SMTP_HOST gesetzt - %s fuer %s: %s", subject, email, link)
		return nil
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: RAP-LAB: %s\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n\r\n"+
		"Hallo,\r\n\r\n%s\r\n\r\n%s\r\n\r\n%s\r\n",
		m.from, email, subject, intro, link, outro)
	var auth smtp.Auth
	if m.user != "" {
		auth = smtp.PlainAuth("", m.user, m.pass, m.host)
	}
	return smtp.SendMail(m.host+":"+m.port, auth, m.from, []string{email}, []byte(msg))
}
