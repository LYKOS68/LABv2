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

// --- Test-Browser: spricht ueber HTTP mit einem Testserver -------------------------------

type browser struct {
	t      *testing.T
	base   string
	client *http.Client
	mail   *bytes.Buffer // Mailer ohne SMTP_HOST loggt die Links hierher
}

func newBrowser(t *testing.T, store *Store) *browser {
	t.Helper()
	mailer := &Mailer{baseURL: "http://test"}
	srv := httptest.NewServer(newMux(store, newRateLimiter(), mailer))
	t.Cleanup(srv.Close)

	mail := &bytes.Buffer{}
	log.SetOutput(mail)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	jar, _ := cookiejar.New(nil)
	return &browser{t: t, base: srv.URL, client: &http.Client{Jar: jar}, mail: mail}
}

func (b *browser) get(path string) string {
	b.t.Helper()
	resp, err := b.client.Get(b.base + path)
	if err != nil {
		b.t.Fatalf("GET %s: %v", path, err)
	}
	return readBody(resp)
}

// submit oeffnet formPage (dort steht das CSRF-Token) und schickt das Formular an action.
func (b *browser) submit(formPage, action string, form url.Values) string {
	b.t.Helper()
	form.Set("csrf_token", csrfTokenIn(b.t, b.get(formPage)))
	resp, err := b.client.PostForm(b.base+action, form)
	if err != nil {
		b.t.Fatalf("POST %s: %v", action, err)
	}
	return readBody(resp)
}

// mailedLink liefert den zuletzt verschickten Link zu path (z.B. "/verify") und seinen Token.
func (b *browser) mailedLink(path string) (link, token string) {
	b.t.Helper()
	re := regexp.MustCompile(`http://test(` + regexp.QuoteMeta(path) + `\?token=([0-9a-f]+))`)
	found := re.FindAllStringSubmatch(b.mail.String(), -1)
	if len(found) == 0 {
		b.t.Fatalf("no %s link was mailed", path)
	}
	last := found[len(found)-1]
	return last[1], last[2]
}

func readBody(resp *http.Response) string {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

var csrfFieldRe = regexp.MustCompile(`name="csrf_token" value="([0-9a-f]+)"`)

func csrfTokenIn(t *testing.T, page string) string {
	t.Helper()
	m := csrfFieldRe.FindStringSubmatch(page)
	if m == nil {
		t.Fatal("page has no csrf token")
	}
	return m[1]
}

func expectContains(t *testing.T, page, want string) {
	t.Helper()
	if !strings.Contains(page, want) {
		t.Fatalf("expected %q in page:\n%s", want, page)
	}
}

func expectMissing(t *testing.T, page, unwanted string) {
	t.Helper()
	if strings.Contains(page, unwanted) {
		t.Fatalf("did not expect %q in page:\n%s", unwanted, page)
	}
}

// --- Navigation ---------------------------------------------------------------------------

func TestNavShowsOnlyAnmelden(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	home := web.get("/")
	nav := home[strings.Index(home, "<nav>"):strings.Index(home, "</nav>")]

	expectContains(t, nav, `href="/login">Anmelden<`)
	expectMissing(t, nav, "/register")
	expectMissing(t, nav, "Login<")
}

func TestLoginPageLinksToRegister(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	expectContains(t, web.get("/login"), `href="/register">Registrieren<`)
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

// Ende-zu-Ende: Registrieren -> Login abgelehnt -> Link aus der Mail oeffnen -> Login klappt.
func TestRegistrationNeedsEmailLink(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	login := url.Values{"email": {"neu@example.com"}, "password": {"sicheres1"}}
	register := url.Values{"email": {"neu@example.com"}, "password": {"sicheres1"}, "confirm": {"sicheres1"}}

	page := web.submit("/register", "/register", register)
	expectContains(t, page, "Bestaetigungslink an neu@example.com")

	page = web.submit("/login", "/login", login)
	expectContains(t, page, errNotVerified.Error())
	expectContains(t, page, "/resend-verification")

	link, _ := web.mailedLink("/verify")
	expectContains(t, web.get(link), "ist bestaetigt")

	page = web.submit("/login", "/login", login)
	expectContains(t, page, "Angemeldet als neu@example.com")
}

func newVerifiedUser(t *testing.T, store *Store, email, password string) {
	t.Helper()
	if err := store.CreateUser(email, password); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	token, _ := store.CreateVerification(email)
	if err := store.VerifyEmail(token); err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
}

func TestResetPasswordRejectsInvalidTokens(t *testing.T) {
	store := newTestStore(t)
	newVerifiedUser(t, store, "a@example.com", "altpasswort")
	token, err := store.CreatePasswordReset("a@example.com")
	if err != nil {
		t.Fatalf("CreatePasswordReset: %v", err)
	}
	verifyToken, _ := store.CreateVerification("a@example.com")

	// Bestaetigungs-Token darf nicht als Reset-Token durchgehen (getrennte Tabellen).
	for _, bad := range []string{"", "falsch", token + "x", hashToken(token), verifyToken} {
		if err := store.ResetPassword(bad, "neuespasswort"); err != errInvalidToken {
			t.Errorf("token %q: expected errInvalidToken, got %v", bad, err)
		}
	}
	if _, err := store.db.Exec("UPDATE password_resets SET expires_at = ?", time.Now().Add(-time.Minute).Unix()); err != nil {
		t.Fatalf("expire token: %v", err)
	}
	if err := store.ResetPassword(token, "neuespasswort"); err != errInvalidToken {
		t.Fatalf("expired token: expected errInvalidToken, got %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "altpasswort"); err != nil {
		t.Fatalf("old password must still work after rejected resets: %v", err)
	}
}

func TestResetPasswordWorksOnceAndEndsSessions(t *testing.T) {
	store := newTestStore(t)
	newVerifiedUser(t, store, "a@example.com", "altpasswort")
	session, _ := store.CreateSession("a@example.com")
	old, _ := store.CreatePasswordReset("a@example.com")
	token, _ := store.CreatePasswordReset("a@example.com")

	if err := store.ResetPassword(old, "neuespasswort"); err != errInvalidToken {
		t.Fatalf("superseded token: expected errInvalidToken, got %v", err)
	}
	if err := store.ResetPassword(token, "neuespasswort"); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "altpasswort"); err != errInvalidCredentials {
		t.Fatalf("old password: expected errInvalidCredentials, got %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "neuespasswort"); err != nil {
		t.Fatalf("new password: %v", err)
	}
	if _, ok := store.UserByToken(session); ok {
		t.Fatal("session from before the reset is still active")
	}
	if err := store.ResetPassword(token, "drittespasswort"); err != errInvalidToken {
		t.Fatalf("reused token: expected errInvalidToken, got %v", err)
	}
}

func TestResetPasswordVerifiesUnverifiedAccount(t *testing.T) {
	store := newTestStore(t)
	if err := store.CreateUser("a@example.com", "altpasswort"); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	token, _ := store.CreatePasswordReset("a@example.com")
	if err := store.ResetPassword(token, "neuespasswort"); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "neuespasswort"); err != nil {
		t.Fatalf("login after reset of unverified account: %v", err)
	}
}

func TestLoginPageLinksToForgotPassword(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	expectContains(t, web.get("/login"), `href="/forgot-password"`)
}

// Fremde Adresse: gleiche Antwort wie bei echter, aber es geht keine Mail raus.
func TestForgotPasswordHidesUnknownEmail(t *testing.T) {
	store := newTestStore(t)
	newVerifiedUser(t, store, "a@example.com", "altpasswort")
	web := newBrowser(t, store)

	known := web.submit("/forgot-password", "/forgot-password", url.Values{"email": {"a@example.com"}})
	unknown := web.submit("/forgot-password", "/forgot-password", url.Values{"email": {"b@example.com"}})

	if strings.ReplaceAll(unknown, "b@example.com", "a@example.com") != known {
		t.Fatal("answer for unknown email differs from answer for known email")
	}
	expectMissing(t, web.mail.String(), "b@example.com")
}

// Ende-zu-Ende: Link anfordern -> zu kurzes Passwort kostet den Link nicht -> neues Passwort -> Login.
func TestForgotPasswordFlow(t *testing.T) {
	store := newTestStore(t)
	newVerifiedUser(t, store, "a@example.com", "altpasswort")
	web := newBrowser(t, store)

	web.submit("/forgot-password", "/forgot-password", url.Values{"email": {"a@example.com"}})
	link, token := web.mailedLink("/reset-password")

	tooShort := url.Values{"token": {token}, "password": {"kurz"}, "confirm": {"kurz"}}
	expectContains(t, web.submit(link, "/reset-password", tooShort), "mindestens 8 Zeichen")

	valid := url.Values{"token": {token}, "password": {"neuespasswort"}, "confirm": {"neuespasswort"}}
	expectContains(t, web.submit(link, "/reset-password", valid), "Passwort wurde geaendert")

	login := url.Values{"email": {"a@example.com"}, "password": {"neuespasswort"}}
	expectContains(t, web.submit("/login", "/login", login), "Angemeldet als a@example.com")
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
