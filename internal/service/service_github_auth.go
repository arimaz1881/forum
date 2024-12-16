package service

import (
	"context"
	"encoding/json"
	"forum/internal/domain"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GitHubUser struct {
	GithubID int64  `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
}

func (s *service) GithubAuth(ctx context.Context, code string) (*SignInResponse, error) {
	// Exchange the authorization code for an access token
	token, err := getGitHubAccessToken(code)
	if err != nil {
		return nil, err
	}

	// Fetch user information using the access token
	user, err := getGitHubUserInfo(token)
	if err != nil {
		return nil, err
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	userID, err := s.users.OAuthFindOrCreateUser(ctx, domain.GoogleAuthInput{
		Provider: "github",
		OAuthID:  strconv.Itoa(int(user.GithubID)),
		Email:    user.Email,
		Login:    user.Login,
		Password: "",
	})
	if err != nil {
		return nil, err
	}

	if err := s.sessions.Close(ctx, domain.CloseSessionInput{
		UserID:    userID,
		ExpiresAt: time.Now().UTC(),
	}); err != nil {
		return nil, err
	}

	token, expiresAt, err := s.createSession(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &SignInResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func getGitHubAccessToken(code string) (string, error) {
	// Prepare the POST request to exchange the code for a token
	data := url.Values{}
	data.Set("client_id", domain.GithubClientID)
	data.Set("client_secret", domain.GithubClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", domain.GithubRedirectURI)

	resp, err := http.PostForm("https://github.com/login/oauth/access_token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract the access token
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}
	return values.Get("access_token"), nil
}

func getGitHubUserInfo(accessToken string) (*GitHubUser, error) {
	// Prepare the request to GitHub's user API
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (i *GitHubUser) validate() error {
	i.Login = strings.TrimSpace(i.Login)
	if len(i.Login) == 0 {
		return domain.ErrIncorrectCredentials
	}
	return nil
}
