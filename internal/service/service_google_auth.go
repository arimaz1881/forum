package service

import (
	"context"
	"encoding/json"
	"fmt"
	"forum/internal/domain"
	"io" 
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GoogleUserInfo struct {
	GoogleID string `json:"id"`
	Email    string `json:"email"`
	Login    string `json:"name"`
}

func (s *service) GoogleAuth(ctx context.Context, code string) (*SignInResponse, error) {
	token, err := exchangeCodeForToken(code)

	if err != nil {
		return nil, err
	}

	// Fetch user info
	googleUserInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		return nil, err
	}

	if err := googleUserInfo.validate(); err != nil {
		return nil, err
	}
	userID, err := s.users.OAuthFindOrCreateUser(ctx, domain.GoogleAuthInput{
		Provider: "google",
		OAuthID:  googleUserInfo.GoogleID,
		Email:    googleUserInfo.Email,
		Login:    googleUserInfo.Login,
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

func (i *GoogleUserInfo) validate() error {
	i.Email = strings.TrimSpace(i.Email)
	i.Email = strings.ToLower(i.Email)
	if ok := isValidEmail(i.Email); !ok {
		return domain.ErrInvalidEmail
	}

	return nil
}

func exchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", domain.GoogleClientID)
	data.Set("client_secret", domain.GoogleClientSecret)
	data.Set("redirect_uri", domain.GoogleRedirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	resp, err := http.PostForm(domain.GoogleTokenURL, data)
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

func fetchGoogleUserInfo(token string) (*GoogleUserInfo, error) {
	req, _ := http.NewRequest("GET", domain.GoogleUserInfoURL, nil)
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

	userInfoGoogle := &GoogleUserInfo{}
	err = json.Unmarshal(body, &userInfoGoogle)
	if err != nil {
		fmt.Printf("callbackgithub: failed unmarshal userInfo: %s\n", err.Error())
		return nil, err
	}

	return userInfoGoogle, nil
}
