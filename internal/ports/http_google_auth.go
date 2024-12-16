package ports

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"forum/internal/pkg/e3r"
	"forum/internal/pkg/sessions"
	"forum/internal/service"
	"io"
	"net/http"
	"net/url"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	googleClientID     = "" //your client id
	googleClientSecret = "" //your cleant secret
	googleRedirectURI  = "https://localhost:8082/auth/google/callback"
	googleAuthURL      = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL     = "https://oauth2.googleapis.com/token"
	googleUserInfoURL  = "https://www.googleapis.com/oauth2/v2/userinfo"
)

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
		googleAuthURL, googleClientID, url.QueryEscape(googleRedirectURI), state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	// Verify state token
	state, err := r.Cookie("oauthstate")

	if err != nil || state.Value != r.URL.Query().Get("state") {
		e3r.ErrorEncoder(e3r.BadRequest("Invalid state parameter"), w, false)
		return
	}

	// Exchange authorization code for access token
	code := r.URL.Query().Get("code")

	token, err := exchangeCodeForToken(code)

	if err != nil {
		e3r.ErrorEncoder(e3r.Internal("Failed to exchange token "+err.Error()), w, false)
		return
	}

	// Fetch user info
	googleUserInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		e3r.ErrorEncoder(e3r.Internal("Failed to fetch user info "+err.Error()), w, false)
		return
	}

	session, err := h.svc.GoogleAuth(ctx, googleUserInfo)
	if err != nil {
		e3r.ErrorEncoder(err, w, false)
	}

	sessions.Set(w, session.Token, *session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func exchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("redirect_uri", googleRedirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	resp, err := http.PostForm(googleTokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["access_token"].(string), nil
}

func fetchGoogleUserInfo(token string) (*service.GoogleUserInfo, error) {
	req, _ := http.NewRequest("GET", googleUserInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("callbackGoogle: failed to read response body: %s\n", err.Error())
		return nil, err
	}

	userInfoGoogle := &service.GoogleUserInfo{}
	err = json.Unmarshal(body, &userInfoGoogle)
	if err != nil {
		fmt.Printf("callbackgithub: failed unmarshal userInfo: %s\n", err.Error())
		return nil, err
	}

	return userInfoGoogle, nil
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
