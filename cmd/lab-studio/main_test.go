package main

import "testing"

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
