package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"rap-lab/cmd/lab-studio/internal/auth"
)

type User struct {
	Email        string
	PasswordHash string
}

// Tabellen fuer Einmal-Links per Mail. Nur diese festen Namen gehen in SQL-Strings ein.
const (
	verificationTable = "email_verifications"
	resetTable        = "password_resets"
	verificationTTL   = 24 * time.Hour
	resetTTL          = 1 * time.Hour
)

var (
	ErrAccountExists      = errors.New("Konto mit dieser E-Mail existiert bereits.")
	ErrInvalidCredentials = errors.New("E-Mail oder Passwort ist falsch.")
	ErrNotVerified        = errors.New("Bitte bestaetige zuerst deine E-Mail-Adresse.")
	ErrInvalidToken       = errors.New("Der Link ist ungueltig oder abgelaufen.")
)

func (s *Store) CreateUser(email, password string) error {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, hash)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrAccountExists
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
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	ok, err := auth.VerifyPassword(hash, password)
	if err != nil {
		log.Printf("verifyPassword error for %s: %v", storedEmail, err)
		return nil, ErrInvalidCredentials
	}
	if !ok {
		return nil, ErrInvalidCredentials
	}
	// Erst nach korrektem Passwort pruefen, damit der Hinweis nichts ueber fremde Konten verraet.
	if !verified {
		return nil, ErrNotVerified
	}
	return &User{Email: storedEmail, PasswordHash: hash}, nil
}

// createToken erzeugt einen Einmal-Token fuer email in table und ersetzt aeltere derselben Adresse.
// In der DB liegt nur der SHA-256-Hash; der Klartext-Token geht ausschliesslich in den Link.
func (s *Store) createToken(table, email string, ttl time.Duration) (string, error) {
	token, err := auth.RandomToken()
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
		return "", ErrInvalidToken
	}
	var email string
	err := s.db.QueryRow("DELETE FROM "+table+" WHERE token_hash = ? AND expires_at > ? RETURNING email",
		hashToken(token), time.Now().Unix()).Scan(&email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidToken
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
	hash, err := auth.HashPassword(newPassword)
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
	token, err := auth.RandomToken()
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
