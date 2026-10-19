package auth

import (
	"crypto/subtle"
	"log"
	"net/http"
	"sync"
	"time"
)

const csrfTokenTTL = 30 * time.Minute
const csrfCookieName = "csrf_id"

type csrfEntry struct {
	token   string
	expires time.Time
}

// CSRF haelt die ausgegebenen CSRF-Tokens im Speicher (nicht persistiert), je Cookie-ID eins.
type CSRF struct {
	mu      sync.Mutex
	entries map[string]csrfEntry // cookie-id -> Token+Ablauf
}

func NewCSRF() *CSRF {
	return &CSRF{entries: make(map[string]csrfEntry)}
}

// newToken erzeugt eine neue Cookie-ID und ein zugehoeriges Token, 30 Minuten gueltig.
func (c *CSRF) newToken() (cookieID, token string, err error) {
	cookieID, err = RandomToken()
	if err != nil {
		return "", "", err
	}
	token, err = RandomToken()
	if err != nil {
		return "", "", err
	}
	c.mu.Lock()
	c.entries[cookieID] = csrfEntry{token: token, expires: time.Now().Add(csrfTokenTTL)}
	c.mu.Unlock()
	return cookieID, token, nil
}

// tokenFor liefert das noch gueltige Token zu einer bestehenden Cookie-ID, falls vorhanden.
func (c *CSRF) tokenFor(cookieID string) (string, bool) {
	c.mu.Lock()
	e, ok := c.entries[cookieID]
	c.mu.Unlock()
	if !ok || time.Now().After(e.expires) {
		return "", false
	}
	return e.token, true
}

// validate prueft Cookie-ID gegen eingereichtes Token: existent, nicht abgelaufen, zeitkonstant gleich.
func (c *CSRF) validate(cookieID, token string) bool {
	if cookieID == "" || token == "" {
		return false
	}
	c.mu.Lock()
	e, ok := c.entries[cookieID]
	c.mu.Unlock()
	if !ok || time.Now().After(e.expires) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(e.token), []byte(token)) == 1
}

func cookieID(r *http.Request) string {
	if c, err := r.Cookie(csrfCookieName); err == nil {
		return c.Value
	}
	return ""
}

// Token liefert ein gueltiges CSRF-Token fuer diesen Request: wiederverwendet ein
// bestehendes gueltiges Cookie, erzeugt sonst ein neues und setzt das Cookie.
func (c *CSRF) Token(w http.ResponseWriter, r *http.Request) string {
	if token, ok := c.tokenFor(cookieID(r)); ok {
		return token
	}
	id, token, err := c.newToken()
	if err != nil {
		log.Printf("csrf token generation failed: %v", err)
		return ""
	}
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(csrfTokenTTL.Seconds()),
	})
	return token
}

// Check validiert das csrf_token-Formularfeld gegen das csrf_id-Cookie. Bei Fehlschlag
// wird 403 geschrieben und false zurueckgegeben; der Aufrufer darf die Mutation dann nicht ausfuehren.
func (c *CSRF) Check(w http.ResponseWriter, r *http.Request) bool {
	if !c.validate(cookieID(r), r.FormValue("csrf_token")) {
		http.Error(w, "invalid or missing csrf token", http.StatusForbidden)
		return false
	}
	return true
}

// ValidHeader prueft das Token aus dem Request-Header header gegen das csrf_id-Cookie.
func (c *CSRF) ValidHeader(r *http.Request, header string) bool {
	return c.validate(cookieID(r), r.Header.Get(header))
}
