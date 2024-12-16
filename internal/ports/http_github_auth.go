package ports

import (
	"fmt"
	"forum/internal/domain"
	"forum/internal/pkg/e3r"
	"forum/internal/pkg/sessions"
	"net/http"
	"net/url"

	_ "github.com/mattn/go-sqlite3"
)

func (h *Handler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	// GitHub OAuth authorization URL
	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		domain.GithubClientID, url.QueryEscape(domain.GithubRedirectURI),
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	// Get the "code" from the query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}
	session, err := h.svc.GithubAuth(ctx, code)
	if err != nil {
		e3r.ErrorEncoder(err, w, false)
		return
	}

	sessions.Set(w, session.Token, *session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
