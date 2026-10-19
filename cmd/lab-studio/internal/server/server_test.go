package server

import (
	"bytes"
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

	"rap-lab/cmd/lab-studio/internal/auth"
	"rap-lab/cmd/lab-studio/internal/mail"
	"rap-lab/cmd/lab-studio/internal/store"
)

// --- Test-Browser: spricht ueber HTTP mit einem Testserver -------------------------------

type browser struct {
	t      *testing.T
	base   string
	client *http.Client
	mail   *bytes.Buffer // Mailer ohne SMTP_HOST loggt die Links hierher
}

func newBrowser(t *testing.T, st *store.Store) *browser {
	t.Helper()
	t.Setenv("SMTP_HOST", "")
	t.Setenv("BASE_URL", "http://test")
	srv := httptest.NewServer(NewMux(st, auth.NewCSRF(), auth.NewRateLimiter(), mail.NewFromEnv()))
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

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("newStore: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return st
}

// Ende-zu-Ende: Registrieren -> Login abgelehnt -> Link aus der Mail oeffnen -> Login klappt.

func TestRegistrationNeedsEmailLink(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	login := url.Values{"email": {"neu@example.com"}, "password": {"sicheres1"}}
	register := url.Values{"email": {"neu@example.com"}, "password": {"sicheres1"}, "confirm": {"sicheres1"}}

	page := web.submit("/register", "/register", register)
	expectContains(t, page, "Bestaetigungslink an neu@example.com")

	page = web.submit("/login", "/login", login)
	expectContains(t, page, store.ErrNotVerified.Error())
	expectContains(t, page, "/resend-verification")

	link, _ := web.mailedLink("/verify")
	expectContains(t, web.get(link), "ist bestaetigt")

	page = web.submit("/login", "/login", login)
	expectContains(t, page, "Angemeldet als neu@example.com")
}

func newVerifiedUser(t *testing.T, st *store.Store, email, password string) {
	t.Helper()
	if err := st.CreateUser(email, password); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	token, _ := st.CreateVerification(email)
	if err := st.VerifyEmail(token); err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
}

func TestLoginPageLinksToForgotPassword(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	expectContains(t, web.get("/login"), `href="/forgot-password"`)
}

// Fremde Adresse: gleiche Antwort wie bei echter, aber es geht keine Mail raus.

func TestForgotPasswordHidesUnknownEmail(t *testing.T) {
	st := newTestStore(t)
	newVerifiedUser(t, st, "a@example.com", "altpasswort")
	web := newBrowser(t, st)

	known := web.submit("/forgot-password", "/forgot-password", url.Values{"email": {"a@example.com"}})
	unknown := web.submit("/forgot-password", "/forgot-password", url.Values{"email": {"b@example.com"}})

	if strings.ReplaceAll(unknown, "b@example.com", "a@example.com") != known {
		t.Fatal("answer for unknown email differs from answer for known email")
	}
	expectMissing(t, web.mail.String(), "b@example.com")
}

// Ende-zu-Ende: Link anfordern -> zu kurzes Passwort kostet den Link nicht -> neues Passwort -> Login.

func TestForgotPasswordFlow(t *testing.T) {
	st := newTestStore(t)
	newVerifiedUser(t, st, "a@example.com", "altpasswort")
	web := newBrowser(t, st)

	web.submit("/forgot-password", "/forgot-password", url.Values{"email": {"a@example.com"}})
	link, token := web.mailedLink("/reset-password")

	tooShort := url.Values{"token": {token}, "password": {"kurz"}, "confirm": {"kurz"}}
	expectContains(t, web.submit(link, "/reset-password", tooShort), "mindestens 8 Zeichen")

	valid := url.Values{"token": {token}, "password": {"neuespasswort"}, "confirm": {"neuespasswort"}}
	expectContains(t, web.submit(link, "/reset-password", valid), "Passwort wurde geaendert")

	login := url.Values{"email": {"a@example.com"}, "password": {"neuespasswort"}}
	expectContains(t, web.submit("/login", "/login", login), "Angemeldet als a@example.com")
}

func TestHomeShowsBrandImages(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	home := web.get("/")
	images := []string{"krone.jpg", "schriftzug.png", "leder-hintergrund.jpg"}
	for _, name := range images {
		expectContains(t, home, "/static/img/"+name)
		resp, err := web.client.Get(web.base + "/static/img/" + name)
		if err != nil {
			t.Fatalf("GET %s: %v", name, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s: status %d", name, resp.StatusCode)
		}
	}
}
