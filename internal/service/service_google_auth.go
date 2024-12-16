package service

import (
	"context"
	"forum/internal/domain"
	"strings"
	"time"
)

type GoogleUserInfo struct {
	GoogleID string		`json:"id"`
	Email    string		`json:"email"`
	Login    string		`json:"name"`
}

func (s *service) GoogleAuth(ctx context.Context, input *GoogleUserInfo) (*SignInResponse, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}
	userID, err := s.users.OAuthFindOrCreateUser(ctx, domain.GoogleAuthInput{
		Provider: "google",
		OAuthID: input.GoogleID,
		Email: input.Email,
		Login: input.Login,
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
