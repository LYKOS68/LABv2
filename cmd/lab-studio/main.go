package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email        string
	PasswordHash string
}

type Store struct {
	mu       sync.Mutex
	users    map[string]*User  // key: lowercase email
	sessions map[string]string // token -> lowercase email
}

func newStore() *Store {
	return &Store{
		users:    make(map[string]*User),
		sessions: make(map[string]string),
	}
}

var (
	errAccountExists      = errors.New("Konto mit dieser E-Mail existiert bereits.")
	errInvalidCredentials = errors.New("E-Mail oder Passwort ist falsch.")
)

func (s *Store) CreateUser(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	key := strings.ToLower(email)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[key]; exists {
		return errAccountExists
	}
	s.users[key] = &User{Email: email, PasswordHash: string(hash)}
	return nil
}

func (s *Store) Authenticate(email, password string) (*User, error) {
	key := strings.ToLower(email)
	s.mu.Lock()
	u, exists := s.users[key]
	s.mu.Unlock()
	if !exists {
		return nil, errInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errInvalidCredentials
	}
	return u, nil
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
	s.mu.Lock()
	u, exists := s.users[email]
	s.mu.Unlock()
	return u, exists
}

func (s *Store) DeleteSession(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func (s *Store) DeleteUser(email string) {
	key := strings.ToLower(email)
	s.mu.Lock()
	delete(s.users, key)
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
      <a href="/logout">Logout</a>
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
    <button type="submit" onclick="return confirm('Account wirklich endgueltig loeschen?')">Account loeschen</button>
  </form>
</div>
{{end}}
{{end}}

{{define "content-register"}}
<h1>Registrieren</h1>
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
<form method="post" action="/register">
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
  <input type="email" name="email" placeholder="E-Mail" value="{{.Email}}" required>
  <input type="password" name="password" placeholder="Passwort" required>
  <button type="submit">Anmelden</button>
</form>
{{end}}
`))

type pageData struct {
	LoggedIn bool
	Email    string
	Error    string
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

func render(w http.ResponseWriter, contentName string, data pageData) {
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
	store := newStore()
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
		render(w, "content-home", data)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, "content-register", pageData{})
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirm := r.FormValue("confirm")
		if msg := validateRegistration(email, password, confirm); msg != "" {
			render(w, "content-register", pageData{Error: msg, Email: email})
			return
		}
		if err := store.CreateUser(email, password); err != nil {
			render(w, "content-register", pageData{Error: err.Error(), Email: email})
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, "content-login", pageData{})
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		u, err := store.Authenticate(email, password)
		if err != nil {
			render(w, "content-login", pageData{Error: err.Error(), Email: email})
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
		u, token, loggedIn := currentUser(r, store)
		if !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		store.DeleteUser(u.Email)
		store.DeleteSession(token)
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("RAP-LAB laeuft auf http://127.0.0.1:8090")
	log.Fatal(http.ListenAndServe("127.0.0.1:8090", mux))
}
