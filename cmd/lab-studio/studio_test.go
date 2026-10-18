package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"testing"
)

// --- Helfer fuer die Studio-API ---------------------------------------------------------

// loggedIn legt ein bestaetigtes Konto an und meldet einen neuen Test-Browser damit an.
func loggedIn(t *testing.T, store *Store, email string) *browser {
	t.Helper()
	newVerifiedUser(t, store, email, "sicheres1")
	web := newBrowser(t, store)
	login := url.Values{"email": {email}, "password": {"sicheres1"}}
	expectContains(t, web.submit("/login", "/login", login), "Angemeldet als")
	return web
}

var metaTokenRe = regexp.MustCompile(`name="csrf-token" content="([0-9a-f]+)"`)

// studioToken oeffnet /studio und liest das CSRF-Token aus der Seite.
func (b *browser) studioToken() string {
	b.t.Helper()
	m := metaTokenRe.FindStringSubmatch(b.get("/studio"))
	if m == nil {
		b.t.Fatal("studio page has no csrf token")
	}
	return m[1]
}

// call schickt eine JSON-Anfrage an die Studio-API und liefert Status und Antwort.
func (b *browser) call(method, path, token string, body any) (int, string) {
	b.t.Helper()
	var payload *bytes.Reader
	if body == nil {
		payload = bytes.NewReader(nil)
	} else {
		data, _ := json.Marshal(body)
		payload = bytes.NewReader(data)
	}
	req, _ := http.NewRequest(method, b.base+path, payload)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("X-CSRF-Token", token)
	}
	resp, err := b.client.Do(req)
	if err != nil {
		b.t.Fatalf("%s %s: %v", method, path, err)
	}
	return resp.StatusCode, readBody(resp)
}

func expectStatus(t *testing.T, got, want int, body string) {
	t.Helper()
	if got != want {
		t.Fatalf("expected status %d, got %d: %s", want, got, body)
	}
}

func decode[T any](t *testing.T, body string) T {
	t.Helper()
	var v T
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		t.Fatalf("decode %q: %v", body, err)
	}
	return v
}

type workspaceResponse struct {
	View   View    `json:"view"`
	Fields []Field `json:"fields"`
}

func (b *browser) workspace() workspaceResponse {
	b.t.Helper()
	status, body := b.call("GET", "/studio/api/workspace", "", nil)
	expectStatus(b.t, status, http.StatusOK, body)
	return decode[workspaceResponse](b.t, body)
}

func (b *browser) createSong(token, title string) Song {
	b.t.Helper()
	status, body := b.call("POST", "/studio/api/songs", token, map[string]string{"title": title})
	expectStatus(b.t, status, http.StatusCreated, body)
	return decode[Song](b.t, body)
}

func (b *browser) createField(token string, f Field) Field {
	b.t.Helper()
	status, body := b.call("POST", "/studio/api/fields", token, f)
	expectStatus(b.t, status, http.StatusCreated, body)
	return decode[Field](b.t, body)
}

// --- Zugang -----------------------------------------------------------------------------

func TestStudioPageRedirectsToLogin(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	expectContains(t, web.get("/studio"), "<h1>Anmelden</h1>")
}

func TestStudioAPIRequiresLogin(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	status, body := web.call("GET", "/studio/api/songs", "", nil)
	expectStatus(t, status, http.StatusUnauthorized, body)
}

func TestStudioPageLoadsItsFiles(t *testing.T) {
	web := loggedIn(t, newTestStore(t), "a@example.com")
	page := web.get("/studio")
	expectContains(t, page, `<script src="/static/studio.js">`)
	expectContains(t, web.get("/static/studio.js"), "async function start()")
	expectContains(t, web.get("/static/studio.css"), "#viewport")
	expectContains(t, web.get("/static/brand.css"), "--neon: #ACE182")
	expectContains(t, web.get("/"), `href="/static/brand.css"`)
}

func TestNavLinksToStudio(t *testing.T) {
	web := newBrowser(t, newTestStore(t))
	expectContains(t, web.get("/"), `href="/studio">LAB-STUDIO<`)
}

func TestStudioWriteNeedsCSRFToken(t *testing.T) {
	web := loggedIn(t, newTestStore(t), "a@example.com")
	song := map[string]string{"title": "Ohne Token"}

	status, body := web.call("POST", "/studio/api/songs", "", song)
	expectStatus(t, status, http.StatusForbidden, body)

	status, body = web.call("POST", "/studio/api/songs", "falsch", song)
	expectStatus(t, status, http.StatusForbidden, body)

	status, body = web.call("GET", "/studio/api/songs", "", nil)
	expectStatus(t, status, http.StatusOK, body)
	if songs := decode[[]Song](t, body); len(songs) != 0 {
		t.Fatalf("song was created without valid token: %v", songs)
	}
}

// --- Leeres Studio und Grundfunktionen ----------------------------------------------------

func TestNewAccountStudioIsEmpty(t *testing.T) {
	web := loggedIn(t, newTestStore(t), "a@example.com")

	status, body := web.call("GET", "/studio/api/songs", "", nil)
	expectStatus(t, status, http.StatusOK, body)
	if songs := decode[[]Song](t, body); len(songs) != 0 {
		t.Fatalf("new account has songs: %v", songs)
	}

	ws := web.workspace()
	if len(ws.Fields) != 0 {
		t.Fatalf("new account has fields: %v", ws.Fields)
	}
	if ws.View.Scale != 1 {
		t.Fatalf("new account view scale = %v, want 1", ws.View.Scale)
	}
}

func TestStudioSongRoundTrip(t *testing.T) {
	web := loggedIn(t, newTestStore(t), "a@example.com")
	token := web.studioToken()
	song := web.createSong(token, "Erster Song")

	status, body := web.call("PUT", "/studio/api/songs/"+itoa(song.ID), token, map[string]string{"text": "Zeile 1\nZeile 2"})
	expectStatus(t, status, http.StatusOK, body)

	status, body = web.call("GET", "/studio/api/songs/"+itoa(song.ID), "", nil)
	expectStatus(t, status, http.StatusOK, body)
	if got := decode[Song](t, body).Text; got != "Zeile 1\nZeile 2" {
		t.Fatalf("song text = %q", got)
	}
}

func TestStudioFieldPatchKeepsIDAndType(t *testing.T) {
	web := loggedIn(t, newTestStore(t), "a@example.com")
	token := web.studioToken()
	note := web.createField(token, Field{Type: "note", X: 10, Y: 20})

	changes := map[string]any{"id": 999, "type": "songs", "x": 50, "text": "Idee"}
	status, body := web.call("PATCH", "/studio/api/fields/"+itoa(note.ID), token, changes)
	expectStatus(t, status, http.StatusOK, body)

	got := web.workspace().Fields[0]
	if got.ID != note.ID || got.Type != "note" {
		t.Fatalf("patch changed id/type: %+v", got)
	}
	if got.X != 50 || got.Y != 20 || got.Text != "Idee" {
		t.Fatalf("patch result wrong: %+v", got)
	}
}

// --- Ein Konto sieht nur seine eigenen Daten ----------------------------------------------

func TestStudioOtherAccountCannotTouchSongs(t *testing.T) {
	store := newTestStore(t)
	alice := loggedIn(t, store, "alice@example.com")
	aliceToken := alice.studioToken()
	song := alice.createSong(aliceToken, "Privat")
	alice.call("PUT", "/studio/api/songs/"+itoa(song.ID), aliceToken, map[string]string{"text": "geheim"})

	bob := loggedIn(t, store, "bob@example.com")
	bobToken := bob.studioToken()
	songPath := "/studio/api/songs/" + itoa(song.ID)

	status, body := bob.call("GET", songPath, "", nil)
	expectStatus(t, status, http.StatusNotFound, body)

	status, body = bob.call("PUT", songPath, bobToken, map[string]string{"text": "ueberschrieben"})
	expectStatus(t, status, http.StatusNotFound, body)

	status, body = bob.call("GET", "/studio/api/songs", "", nil)
	expectStatus(t, status, http.StatusOK, body)
	if songs := decode[[]Song](t, body); len(songs) != 0 {
		t.Fatalf("bob sees alice's songs: %v", songs)
	}

	stored, err := store.GetSong("alice@example.com", song.ID)
	if err != nil || stored.Text != "geheim" {
		t.Fatalf("alice's song changed: %+v, %v", stored, err)
	}
}

func TestStudioOtherAccountCannotTouchFields(t *testing.T) {
	store := newTestStore(t)
	alice := loggedIn(t, store, "alice@example.com")
	aliceToken := alice.studioToken()
	song := alice.createSong(aliceToken, "Privat")
	field := alice.createField(aliceToken, Field{Type: "note", X: 1, Y: 2})

	bob := loggedIn(t, store, "bob@example.com")
	bobToken := bob.studioToken()

	status, body := bob.call("PATCH", "/studio/api/fields/"+itoa(field.ID), bobToken, map[string]any{"x": 500})
	expectStatus(t, status, http.StatusNotFound, body)

	// Bob darf in seinem eigenen SONGS-Feld nicht auf Alices Song zeigen.
	status, body = bob.call("POST", "/studio/api/fields", bobToken, Field{Type: "songs", SongID: song.ID})
	expectStatus(t, status, http.StatusBadRequest, body)

	if fields := bob.workspace().Fields; len(fields) != 0 {
		t.Fatalf("bob sees fields: %v", fields)
	}
	if x := alice.workspace().Fields[0].X; x != 1 {
		t.Fatalf("alice's field moved to x=%v", x)
	}
}

// --- Ungueltige Eingaben ----------------------------------------------------------------

func TestStudioRejectsInvalidInput(t *testing.T) {
	web := loggedIn(t, newTestStore(t), "a@example.com")
	token := web.studioToken()
	note := web.createField(token, Field{Type: "note"})

	cases := []struct {
		name, method, path string
		body               any
	}{
		{"unknown field type", "POST", "/studio/api/fields", Field{Type: "bogus"}},
		{"zero size", "PATCH", "/studio/api/fields/" + itoa(note.ID), map[string]any{"w": 0}},
		{"empty title", "POST", "/studio/api/songs", map[string]string{"title": "   "}},
		{"broken json", "POST", "/studio/api/songs", "kein objekt"},
	}
	for _, c := range cases {
		status, body := web.call(c.method, c.path, token, c.body)
		if status != http.StatusBadRequest {
			t.Errorf("%s: expected 400, got %d: %s", c.name, status, body)
		}
	}

	ws := web.workspace()
	if len(ws.Fields) != 1 || ws.Fields[0].W != 240 {
		t.Fatalf("invalid input changed stored fields: %+v", ws.Fields)
	}
}

func TestStudioViewScaleIsClamped(t *testing.T) {
	web := loggedIn(t, newTestStore(t), "a@example.com")
	token := web.studioToken()

	web.call("PUT", "/studio/api/workspace/view", token, View{X: 5, Y: 6, Scale: 99})
	if got := web.workspace().View; got.Scale != maxScale || got.X != 5 {
		t.Fatalf("view after scale 99 = %+v", got)
	}

	web.call("PUT", "/studio/api/workspace/view", token, View{Scale: 0.01})
	if got := web.workspace().View.Scale; got != minScale {
		t.Fatalf("scale after 0.01 = %v", got)
	}
}

// --- Speicherung --------------------------------------------------------------------------

func TestStudioDataSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "restart.sqlite")
	first, err := newStore(path)
	if err != nil {
		t.Fatalf("newStore: %v", err)
	}
	song, _ := first.CreateSong("a@example.com", "Bleibt")
	first.SaveSongText("a@example.com", song.ID, "Text")
	first.CreateField("a@example.com", Field{Type: "note", W: 240, H: 160, Text: "Notiz"})
	first.SaveView("a@example.com", View{X: 3, Y: 4, Scale: 2})
	first.db.Close()

	second, err := newStore(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer second.db.Close()
	if got, _ := second.GetSong("a@example.com", song.ID); got.Text != "Text" {
		t.Fatalf("song after restart: %+v", got)
	}
	if fields, _ := second.ListFields("a@example.com"); len(fields) != 1 || fields[0].Text != "Notiz" {
		t.Fatalf("fields after restart: %+v", fields)
	}
	if view, _ := second.LoadView("a@example.com"); view.Scale != 2 {
		t.Fatalf("view after restart: %+v", view)
	}
}

func TestDeleteUserRemovesStudioData(t *testing.T) {
	store := newTestStore(t)
	newVerifiedUser(t, store, "a@example.com", "sicheres1")
	store.CreateSong("a@example.com", "Weg")
	store.CreateField("a@example.com", Field{Type: "note", W: 1, H: 1})

	if err := store.DeleteUser("a@example.com"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	newVerifiedUser(t, store, "a@example.com", "sicheres1")
	if songs, _ := store.ListSongs("a@example.com"); len(songs) != 0 {
		t.Fatalf("songs survived account deletion: %v", songs)
	}
	if fields, _ := store.ListFields("a@example.com"); len(fields) != 0 {
		t.Fatalf("fields survived account deletion: %v", fields)
	}
}

func itoa(n int64) string {
	data, _ := json.Marshal(n)
	return string(data)
}
