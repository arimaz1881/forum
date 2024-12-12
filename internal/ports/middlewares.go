package ports

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
)

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
