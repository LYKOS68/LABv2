package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"time"
)

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
