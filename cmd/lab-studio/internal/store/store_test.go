package store

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"rap-lab/cmd/lab-studio/internal/auth"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("newStore: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestUnverifiedAccountCannotLogIn(t *testing.T) {
	store := newTestStore(t)
	if err := store.CreateUser("a@example.com", "sicheres1"); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "sicheres1"); err != ErrNotVerified {
		t.Fatalf("expected ErrNotVerified, got %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "falsch123"); err != ErrInvalidCredentials {
		t.Fatalf("wrong password on unverified account: expected ErrInvalidCredentials, got %v", err)
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
		if err := store.VerifyEmail(bad); err != ErrInvalidToken {
			t.Errorf("token %q: expected ErrInvalidToken, got %v", bad, err)
		}
	}
	if !store.IsUnverified("a@example.com") {
		t.Fatal("account became verified through an invalid token")
	}

	// Abgelaufener Token: Ablaufzeit in die Vergangenheit setzen.
	if _, err := store.db.Exec("UPDATE email_verifications SET expires_at = ?", time.Now().Add(-time.Minute).Unix()); err != nil {
		t.Fatalf("expire token: %v", err)
	}
	if err := store.VerifyEmail(token); err != ErrInvalidToken {
		t.Fatalf("expired token: expected ErrInvalidToken, got %v", err)
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
	if err := store.VerifyEmail(token); err != ErrInvalidToken {
		t.Fatalf("reused token: expected ErrInvalidToken, got %v", err)
	}
}

func TestResendReplacesOldToken(t *testing.T) {
	store := newTestStore(t)
	if err := store.CreateUser("a@example.com", "sicheres1"); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	old, _ := store.CreateVerification("a@example.com")
	fresh, _ := store.CreateVerification("a@example.com")
	if err := store.VerifyEmail(old); err != ErrInvalidToken {
		t.Fatalf("old token after resend: expected ErrInvalidToken, got %v", err)
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
	hash, _ := auth.HashPassword("sicheres1")
	if _, err := old.Exec("INSERT INTO users VALUES (?, ?)", "alt@example.com", hash); err != nil {
		t.Fatalf("insert: %v", err)
	}
	old.Close()

	store, err := Open(path)
	if err != nil {
		t.Fatalf("newStore on old schema: %v", err)
	}
	defer store.Close()
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
		if err := store.ResetPassword(bad, "neuespasswort"); err != ErrInvalidToken {
			t.Errorf("token %q: expected ErrInvalidToken, got %v", bad, err)
		}
	}
	if _, err := store.db.Exec("UPDATE password_resets SET expires_at = ?", time.Now().Add(-time.Minute).Unix()); err != nil {
		t.Fatalf("expire token: %v", err)
	}
	if err := store.ResetPassword(token, "neuespasswort"); err != ErrInvalidToken {
		t.Fatalf("expired token: expected ErrInvalidToken, got %v", err)
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

	if err := store.ResetPassword(old, "neuespasswort"); err != ErrInvalidToken {
		t.Fatalf("superseded token: expected ErrInvalidToken, got %v", err)
	}
	if err := store.ResetPassword(token, "neuespasswort"); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "altpasswort"); err != ErrInvalidCredentials {
		t.Fatalf("old password: expected ErrInvalidCredentials, got %v", err)
	}
	if _, err := store.Authenticate("a@example.com", "neuespasswort"); err != nil {
		t.Fatalf("new password: %v", err)
	}
	if _, ok := store.UserByToken(session); ok {
		t.Fatal("session from before the reset is still active")
	}
	if err := store.ResetPassword(token, "drittespasswort"); err != ErrInvalidToken {
		t.Fatalf("reused token: expected ErrInvalidToken, got %v", err)
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

func TestStudioDataSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "restart.sqlite")
	first, err := Open(path)
	if err != nil {
		t.Fatalf("newStore: %v", err)
	}
	song, _ := first.CreateSong("a@example.com", "Bleibt")
	first.SaveSongText("a@example.com", song.ID, "Text")
	first.CreateField("a@example.com", Field{Type: "note", W: 240, H: 160, Text: "Notiz"})
	first.SaveView("a@example.com", View{X: 3, Y: 4, Scale: 2})
	first.Close()

	second, err := Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer second.Close()
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
