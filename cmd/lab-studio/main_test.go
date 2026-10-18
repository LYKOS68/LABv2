package main

import (
	"bytes"
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestNavShowsOnlyAnmeldenAndLoginLinksToRegister(t *testing.T) {
	store := newTestStore(t)

	home := httptest.NewRecorder()
	homeHandler(store)(home, httptest.NewRequest("GET", "/", nil))
	nav := home.Body.String()
	nav = nav[strings.Index(nav, "<nav>"):strings.Index(nav, "</nav>")]
	if !strings.Contains(nav, `href="/login">Anmelden<`) {
		t.Error("nav: missing Anmelden link to /login")
	}
	if strings.Contains(nav, "/register") || strings.Contains(nav, "Login<") {
		t.Errorf("nav: should only show Anmelden, got %q", nav)
	}

	login := httptest.NewRecorder()
	loginHandler(store, newRateLimiter())(login, httptest.NewRequest("GET", "/login", nil))
	if !strings.Contains(login.Body.String(), `href="/register">Registrieren<`) {
		t.Error("login page: missing Registrieren link")
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	store, err := newStore(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("newStore: %v", err)
	}
	t.Cleanup(func() { store.db.Close() })
	return store
}

func TestUnverifiedAccountCannotLogIn(t *testing.T) {
	store := newTestStore(t)
	if err := store.CreateUser("a@example.com", "sicheres1"); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "sicheres1"); err != errNotVerified {
		t.Fatalf("expected errNotVerified, got %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "falsch123"); err != errInvalidCredentials {
		t.Fatalf("wrong password on unverified account: expected errInvalidCredentials, got %v", err)
	}
}

func TestVerifyEmailRejectsInvalidTokens(t *testing.T) {
	store := newTestStore(t)
	if err := store.CreateUser("a@example.com", "sicheres1"); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	token, err := store.CreateVerification("a@example.com")
	if err != nil {
		t.Fatalf("CreateVerification: %v", err)
	}

	for _, bad := range []string{"", "falsch", token + "x", hashToken(token)} {
		if err := store.VerifyEmail(bad); err != errInvalidToken {
			t.Errorf("token %q: expected errInvalidToken, got %v", bad, err)
		}
	}
	if !store.IsUnverified("a@example.com") {
		t.Fatal("account became verified through an invalid token")
	}

	// Abgelaufener Token: Ablaufzeit in die Vergangenheit setzen.
	if _, err := store.db.Exec("UPDATE email_verifications SET expires_at = ?", time.Now().Add(-time.Minute).Unix()); err != nil {
		t.Fatalf("expire token: %v", err)
	}
	if err := store.VerifyEmail(token); err != errInvalidToken {
		t.Fatalf("expired token: expected errInvalidToken, got %v", err)
	}
	if !store.IsUnverified("a@example.com") {
		t.Fatal("account became verified through an expired token")
	}
}

func TestVerifyEmailTokenWorksExactlyOnce(t *testing.T) {
	store := newTestStore(t)
	if err := store.CreateUser("a@example.com", "sicheres1"); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	token, err := store.CreateVerification("a@example.com")
	if err != nil {
		t.Fatalf("CreateVerification: %v", err)
	}
	var stored string
	if err := store.db.QueryRow("SELECT token_hash FROM email_verifications").Scan(&stored); err != nil {
		t.Fatalf("read token_hash: %v", err)
	}
	if stored == token {
		t.Fatal("token stored in plain text")
	}
	if err := store.VerifyEmail(token); err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "sicheres1"); err != nil {
		t.Fatalf("login after verification: %v", err)
	}
	if err := store.VerifyEmail(token); err != errInvalidToken {
		t.Fatalf("reused token: expected errInvalidToken, got %v", err)
	}
}

func TestResendReplacesOldToken(t *testing.T) {
	store := newTestStore(t)
	if err := store.CreateUser("a@example.com", "sicheres1"); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	old, _ := store.CreateVerification("a@example.com")
	fresh, _ := store.CreateVerification("a@example.com")
	if err := store.VerifyEmail(old); err != errInvalidToken {
		t.Fatalf("old token after resend: expected errInvalidToken, got %v", err)
	}
	if err := store.VerifyEmail(fresh); err != nil {
		t.Fatalf("fresh token: %v", err)
	}
}

func TestMigrationMarksExistingAccountsVerified(t *testing.T) {
	path := filepath.Join(t.TempDir(), "old.sqlite")
	old, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if _, err := old.Exec(`CREATE TABLE users (email TEXT NOT NULL UNIQUE COLLATE NOCASE, password_hash TEXT NOT NULL)`); err != nil {
		t.Fatalf("create old schema: %v", err)
	}
	hash, _ := hashPassword("sicheres1")
	if _, err := old.Exec("INSERT INTO users VALUES (?, ?)", "alt@example.com", hash); err != nil {
		t.Fatalf("insert: %v", err)
	}
	old.Close()

	store, err := newStore(path)
	if err != nil {
		t.Fatalf("newStore on old schema: %v", err)
	}
	defer store.db.Close()
	if _, err := store.Authenticate("alt@example.com", "sicheres1"); err != nil {
		t.Fatalf("existing account locked out after migration: %v", err)
	}
	if err := store.CreateUser("neu@example.com", "sicheres1"); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if !store.IsUnverified("neu@example.com") {
		t.Fatal("new account after migration should start unverified")
	}
}

// Ende-zu-Ende ueber HTTP: Registrieren -> Login abgelehnt -> Link aus Konsolen-Log oeffnen -> Login klappt.
func TestRegistrationFlowRequiresEmailLink(t *testing.T) {
	store := newTestStore(t)
	limiter := newRateLimiter()
	mailer := &Mailer{baseURL: "http://test"}
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler(store))
	mux.HandleFunc("/register", registerHandler(store, limiter, mailer))
	mux.HandleFunc("/login", loginHandler(store, limiter))
	mux.HandleFunc("/verify", verifyHandler(store))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	tokenRe := regexp.MustCompile(`name="csrf_token" value="([0-9a-f]+)"`)
	post := func(path string, form url.Values) string {
		t.Helper()
		page, err := client.Get(srv.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		body, _ := io.ReadAll(page.Body)
		page.Body.Close()
		m := tokenRe.FindSubmatch(body)
		if m == nil {
			t.Fatalf("no csrf token on %s", path)
		}
		form.Set("csrf_token", string(m[1]))
		resp, err := client.PostForm(srv.URL+path, form)
		if err != nil {
			t.Fatalf("POST %s: %v", path, err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return string(b)
	}
	login := url.Values{"email": {"neu@example.com"}, "password": {"sicheres1"}}

	body := post("/register", url.Values{"email": {"neu@example.com"}, "password": {"sicheres1"}, "confirm": {"sicheres1"}})
	if !strings.Contains(body, "Bestaetigungslink an neu@example.com") {
		t.Fatalf("register: expected verify-sent page, got %q", body)
	}
	if body := post("/login", login); !strings.Contains(body, errNotVerified.Error()) || !strings.Contains(body, "/resend-verification") {
		t.Fatal("login before verification: expected rejection with resend button")
	}

	link := regexp.MustCompile(`http://test(/verify\?token=[0-9a-f]+)`).FindStringSubmatch(logBuf.String())
	if link == nil {
		t.Fatalf("no verification link in log: %q", logBuf.String())
	}
	resp, err := client.Get(srv.URL + link[1])
	if err != nil {
		t.Fatalf("GET verify: %v", err)
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if !strings.Contains(string(b), "ist bestaetigt") {
		t.Fatalf("verify: expected success page, got %q", b)
	}

	if body := post("/login", login); !strings.Contains(body, "Angemeldet als neu@example.com") {
		t.Fatalf("login after verification: expected logged-in home, got %q", body)
	}
}

func TestRateLimiterAllowsUpToLimitThenBlocks(t *testing.T) {
	rl := newRateLimiter()
	for i := 0; i < rateLimitAttempts; i++ {
		allowed, _ := rl.Allow("test-key")
		if !allowed {
			t.Fatalf("attempt %d: expected allowed, got blocked", i+1)
		}
	}
	allowed, retryAfter := rl.Allow("test-key")
	if allowed {
		t.Fatal("attempt beyond limit was allowed")
	}
	if retryAfter <= 0 {
		t.Fatalf("expected positive retryAfter, got %v", retryAfter)
	}
}

func TestRateLimiterKeysAreIndependent(t *testing.T) {
	rl := newRateLimiter()
	for i := 0; i < rateLimitAttempts; i++ {
		if allowed, _ := rl.Allow("key-a"); !allowed {
			t.Fatalf("key-a attempt %d unexpectedly blocked", i+1)
		}
	}
	if allowed, _ := rl.Allow("key-a"); allowed {
		t.Fatal("key-a should be blocked after exhausting its limit")
	}
	if allowed, _ := rl.Allow("key-b"); !allowed {
		t.Fatal("key-b should be unaffected by key-a's limit")
	}
}

func TestHashAndVerifyPasswordRoundTrip(t *testing.T) {
	hash, err := hashPassword("sicheres1")
	if err != nil {
		t.Fatalf("hashPassword: %v", err)
	}
	ok, err := verifyPassword(hash, "sicheres1")
	if err != nil {
		t.Fatalf("verifyPassword: %v", err)
	}
	if !ok {
		t.Fatal("correct password rejected")
	}
	ok, err = verifyPassword(hash, "falsch")
	if err != nil {
		t.Fatalf("verifyPassword with wrong password returned error: %v", err)
	}
	if ok {
		t.Fatal("wrong password accepted")
	}
}

func TestVerifyPasswordRejectsCorruptHashWithoutPanic(t *testing.T) {
	corrupt := []string{
		"GARBAGE;NOT;A;HASH",
		"$argon2id$v=19$m=19456,t=2,p=1$short$short",
		"$argon2id$v=1$m=19456,t=2,p=1$c29tZXNhbHQ$c29tZWhhc2g",
		"$argon2id$v=19$m=99999999999,t=2,p=1$c29tZXNhbHQ$c29tZWhhc2g",
		"",
	}
	for _, h := range corrupt {
		ok, err := verifyPassword(h, "irrelevant")
		if err == nil {
			t.Errorf("corrupt hash %q: expected error, got nil (ok=%v)", h, ok)
		}
		if ok {
			t.Errorf("corrupt hash %q: verifyPassword returned true", h)
		}
	}
}
