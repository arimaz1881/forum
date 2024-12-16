package ports

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"forum/internal/domain"
	"forum/internal/pkg/e3r"
	"forum/internal/pkg/sessions"
	"net/http"
	"net/url"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const ()

func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate a state token for CSRF protection
	state := generateRandomString(16)
	http.SetCookie(w, &http.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: time.Now().Add(10 * time.Minute),
	})

	// Redirect to Google's OAuth2 consent page
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile&state=%s",
		domain.GoogleAuthURL, domain.GoogleClientID, url.QueryEscape(domain.GoogleRedirectURI), state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	// Verify state token
	state, err := r.Cookie("oauthstate")

	if err != nil || state.Value != r.URL.Query().Get("state") {
		e3r.ErrorEncoder(e3r.BadRequest("Invalid state parameter"), w, false)
		return
	}

	// Exchange authorization code for access token
	code := r.URL.Query().Get("code")

	session, err := h.svc.GoogleAuth(ctx, code)

	// session, err := h.svc.GoogleAuth(ctx, googleUserInfo)
	if err != nil {
		e3r.ErrorEncoder(err, w, false)
		return
	}

	sessions.Set(w, session.Token, *session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
