package http

import (
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CORS wraps a handler allowing the configured origin and the methods used.
func CORS(allowOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// bucket is a per-client token bucket.
type bucket struct {
	tokens   float64
	lastSeen time.Time
}

// rateLimiter throttles each client IP to a fixed rate using token buckets.
type rateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*bucket
	rate     float64 // tokens refilled per second
	capacity float64 // max burst
}

func newRateLimiter(perMinute int) *rateLimiter {
	rl := &rateLimiter{
		clients:  make(map[string]*bucket),
		rate:     float64(perMinute) / 60.0,
		capacity: float64(perMinute),
	}
	go rl.cleanup()
	return rl
}

// cleanup periodically evicts idle clients to bound memory.
func (rl *rateLimiter) cleanup() {
	for range time.Tick(3 * time.Minute) {
		rl.mu.Lock()
		for ip, b := range rl.clients {
			if time.Since(b.lastSeen) > 3*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// allow reports whether the given client may make a request now.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.clients[ip]
	if !ok {
		rl.clients[ip] = &bucket{tokens: rl.capacity - 1, lastSeen: now}
		return true
	}

	// Refill proportionally to elapsed time, capped at capacity.
	b.tokens = math.Min(rl.capacity, b.tokens+now.Sub(b.lastSeen).Seconds()*rl.rate)
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// RateLimit limits each client IP to perMinute requests per minute.
func RateLimit(perMinute int, next http.Handler) http.Handler {
	rl := newRateLimiter(perMinute)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(clientIP(r)) {
			w.Header().Set("Retry-After", "60")
			writeError(w, http.StatusTooManyRequests, "rate limit exceeded: max 60 requests per minute")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP extracts the originating client IP, honoring X-Forwarded-For.
func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := strings.IndexByte(fwd, ','); i >= 0 {
			return strings.TrimSpace(fwd[:i])
		}
		return strings.TrimSpace(fwd)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
