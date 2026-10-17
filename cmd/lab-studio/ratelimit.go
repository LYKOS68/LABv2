package main

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	rateLimitAttempts = 5
	rateLimitWindow   = 10 * time.Minute
	rateLimitMaxKeys  = 10000 // Obergrenze gegen unbegrenztes Map-Wachstum
)

type rateLimitEntry struct {
	count       int
	windowStart time.Time
}

// RateLimiter begrenzt Versuche pro Schluessel (z.B. "login:"+IP) auf N pro Zeitfenster.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
}

func newRateLimiter() *RateLimiter {
	return &RateLimiter{entries: make(map[string]*rateLimitEntry)}
}

// Allow prueft und zaehlt einen Versuch fuer key. allowed=false liefert zusaetzlich die
// verbleibende Wartezeit bis zum naechsten erlaubten Versuch.
func (rl *RateLimiter) Allow(key string) (allowed bool, retryAfter time.Duration) {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for k, e := range rl.entries {
		if now.Sub(e.windowStart) >= rateLimitWindow {
			delete(rl.entries, k)
		}
	}

	e, ok := rl.entries[key]
	if !ok {
		if len(rl.entries) >= rateLimitMaxKeys {
			return false, rateLimitWindow
		}
		rl.entries[key] = &rateLimitEntry{count: 1, windowStart: now}
		return true, 0
	}
	if e.count >= rateLimitAttempts {
		return false, rateLimitWindow - now.Sub(e.windowStart)
	}
	e.count++
	return true, 0
}

// allowOrTooManyRequests prueft das Rate-Limit fuer key. Bei Ueberschreitung wird 429 mit
// Retry-After geschrieben und false zurueckgegeben; der Aufrufer darf dann nicht fortfahren.
func allowOrTooManyRequests(w http.ResponseWriter, limiter *RateLimiter, key string) bool {
	allowed, retryAfter := limiter.Allow(key)
	if !allowed {
		w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
		http.Error(w, "too many attempts, try again later", http.StatusTooManyRequests)
		return false
	}
	return true
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
