package main

import "testing"

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
