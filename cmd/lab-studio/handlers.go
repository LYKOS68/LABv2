package main

import (
	"log"
	"net/http"
)

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

func homeHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
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

func registerHandler(store *Store, limiter *RateLimiter, mailer *Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		sendVerification(store, mailer, email)
		render(w, r, store, "content-verify-sent", pageData{Email: email})
	}
}

func verifyHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
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
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
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
	}
}

func deleteAccountHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
