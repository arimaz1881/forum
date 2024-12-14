package ports

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/pkg/e3r"
	"log"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu            sync.Mutex
	requestCounts map[string]*ClientRequestInfo
	rateLimit     int           // Number of requests allowed
	timeWindow    time.Duration // Time window for the rate limit
}

type ClientRequestInfo struct {
	requestCount int
	lastReset    time.Time
}

func NewRateLimiter(rateLimit int, timeWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		requestCounts: make(map[string]*ClientRequestInfo),
		rateLimit:     rateLimit,
		timeWindow:    timeWindow,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	clientInfo, exists := rl.requestCounts[ip]
	if !exists {
		rl.requestCounts[ip] = &ClientRequestInfo{
			requestCount: 1,
			lastReset:    time.Now(),
		}
		return true
	}

	// Reset the count if the time window has passed
	if time.Since(clientInfo.lastReset) > rl.timeWindow {
		clientInfo.requestCount = 1
		clientInfo.lastReset = time.Now()
		return true
	}

	// Check if the request count exceeds the limit
	if clientInfo.requestCount < rl.rateLimit {
		clientInfo.requestCount++
		return true
	}

	return false
}

func (h *Handler) rateLimit(rl *RateLimiter, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr // Use client's IP as the key
		if !rl.Allow(clientIP) {
			e3r.ErrorEncoder(e3r.TooManyRequests("too many requests"), w, false)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) recoverPanic(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				fmt.Fprint(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuthN(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/authn/sign-in", http.StatusSeeOther)
			return
		}

		userID, err := h.svc.GetUserByToken(ctx, cookie.Value)
		if err != nil {
			log.Printf("user by token: %s\n", err)
			return
		}

		if userID == 0 {
			http.Redirect(w, r, "/authn/sign-in", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type key string

const myKey key = "myKey"

type User struct {
	ID      int64
	IsAuthN bool
}

func (h *Handler) withContext(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie("token")
		if err != nil {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, err := h.svc.GetUserByToken(ctx, cookie.Value)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/authn/sign-in", http.StatusSeeOther)
			return
		}

		ctx = context.WithValue(ctx, myKey, User{
			ID:      userID,
			IsAuthN: userID != 0,
		})

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
