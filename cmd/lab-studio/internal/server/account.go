package server

import (
	"log"
	"net/http"

	"rap-lab/cmd/lab-studio/internal/auth"
	"rap-lab/cmd/lab-studio/internal/mail"
	"rap-lab/cmd/lab-studio/internal/store"
)

func currentUser(r *http.Request, st *store.Store) (*store.User, string, bool) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return nil, "", false
	}
	u, ok := st.UserByToken(c.Value)
	if !ok {
		return nil, "", false
	}
	return u, c.Value, true
}

func homeHandler(st *store.Store, csrf *auth.CSRF) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		u, _, loggedIn := currentUser(r, st)
		data := pageData{LoggedIn: loggedIn}
		if loggedIn {
			data.Email = u.Email
		}
		render(w, r, csrf, "content-home", data)
	}
}

// sendVerification erzeugt einen Token und verschickt den Link. Fehler werden nur geloggt:
// das Konto existiert bereits, und der Link kann ueber /resend-verification neu angefordert werden.
func sendVerification(st *store.Store, mailer *mail.Mailer, email string) {
	token, err := st.CreateVerification(email)
	if err != nil {
		log.Printf("create verification for %s: %v", email, err)
		return
	}
	if err := mailer.SendVerificationLink(email, token); err != nil {
		log.Printf("send verification mail to %s: %v", email, err)
	}
}

// forgotPasswordHandler antwortet immer gleich, egal ob die E-Mail existiert.
func forgotPasswordHandler(st *store.Store, csrf *auth.CSRF, limiter *auth.RateLimiter, mailer *mail.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, r, csrf, "content-forgot", pageData{})
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !csrf.Check(w, r) {
			return
		}
		if !limiter.AllowOrTooManyRequests(w, "forgot:"+auth.ClientIP(r)) {
			return
		}
		email := r.FormValue("email")
		if st.AccountExists(email) {
			token, err := st.CreatePasswordReset(email)
			if err != nil {
				log.Printf("create password reset for %s: %v", email, err)
			} else if err := mailer.SendPasswordResetLink(email, token); err != nil {
				log.Printf("send password reset mail to %s: %v", email, err)
			}
		}
		render(w, r, csrf, "content-forgot-sent", pageData{Email: email})
	}
}

// resetPasswordHandler zeigt beim GET nur das Formular; der Token wird erst beim POST
// mit gueltigem neuem Passwort verbraucht (Mail-Scanner oeffnen Links vorab).
func resetPasswordHandler(st *store.Store, csrf *auth.CSRF, limiter *auth.RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, r, csrf, "content-reset", pageData{Token: r.URL.Query().Get("token")})
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !csrf.Check(w, r) {
			return
		}
		if !limiter.AllowOrTooManyRequests(w, "reset:"+auth.ClientIP(r)) {
			return
		}
		token := r.FormValue("token")
		if msg := auth.ValidateNewPassword(r.FormValue("password"), r.FormValue("confirm")); msg != "" {
			render(w, r, csrf, "content-reset", pageData{Error: msg, Token: token})
			return
		}
		if err := st.ResetPassword(token, r.FormValue("password")); err != nil {
			if err != store.ErrInvalidToken {
				log.Printf("reset password: %v", err)
			}
			render(w, r, csrf, "content-reset-result", pageData{Error: store.ErrInvalidToken.Error()})
			return
		}
		render(w, r, csrf, "content-reset-result", pageData{})
	}
}

func registerHandler(st *store.Store, csrf *auth.CSRF, limiter *auth.RateLimiter, mailer *mail.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, r, csrf, "content-register", pageData{})
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !csrf.Check(w, r) {
			return
		}
		if !limiter.AllowOrTooManyRequests(w, "register:"+auth.ClientIP(r)) {
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirm := r.FormValue("confirm")
		if msg := auth.ValidateRegistration(email, password, confirm); msg != "" {
			render(w, r, csrf, "content-register", pageData{Error: msg, Email: email})
			return
		}
		if err := st.CreateUser(email, password); err != nil {
			render(w, r, csrf, "content-register", pageData{Error: err.Error(), Email: email})
			return
		}
		sendVerification(st, mailer, email)
		render(w, r, csrf, "content-verify-sent", pageData{Email: email})
	}
}

func verifyHandler(st *store.Store, csrf *auth.CSRF) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := st.VerifyEmail(r.URL.Query().Get("token")); err != nil {
			if err != store.ErrInvalidToken {
				log.Printf("verify email: %v", err)
			}
			render(w, r, csrf, "content-verify-result", pageData{Error: store.ErrInvalidToken.Error()})
			return
		}
		render(w, r, csrf, "content-verify-result", pageData{})
	}
}

// resendVerificationHandler antwortet immer gleich, egal ob die E-Mail existiert, damit
// sich darueber nicht herausfinden laesst, welche Adressen registriert sind.
func resendVerificationHandler(st *store.Store, csrf *auth.CSRF, limiter *auth.RateLimiter, mailer *mail.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !csrf.Check(w, r) {
			return
		}
		if !limiter.AllowOrTooManyRequests(w, "resend:"+auth.ClientIP(r)) {
			return
		}
		email := r.FormValue("email")
		if st.IsUnverified(email) {
			sendVerification(st, mailer, email)
		}
		render(w, r, csrf, "content-verify-sent", pageData{Email: email})
	}
}

func loginHandler(st *store.Store, csrf *auth.CSRF, limiter *auth.RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			render(w, r, csrf, "content-login", pageData{})
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !csrf.Check(w, r) {
			return
		}
		if !limiter.AllowOrTooManyRequests(w, "login:"+auth.ClientIP(r)) {
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		u, err := st.Authenticate(email, password)
		if err != nil {
			render(w, r, csrf, "content-login", pageData{Error: err.Error(), Email: email, ShowResend: err == store.ErrNotVerified})
			return
		}
		token, err := st.CreateSession(u.Email)
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

func logoutHandler(st *store.Store, csrf *auth.CSRF) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !csrf.Check(w, r) {
			return
		}
		if c, err := r.Cookie("session_token"); err == nil {
			st.DeleteSession(c.Value)
		}
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func deleteAccountHandler(st *store.Store, csrf *auth.CSRF) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !csrf.Check(w, r) {
			return
		}
		u, token, loggedIn := currentUser(r, st)
		if !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if err := st.DeleteUser(u.Email); err != nil {
			log.Printf("delete-account error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		st.DeleteSession(token)
		http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Path: "/", MaxAge: -1})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
