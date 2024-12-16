package service

import (
	"context"
	"strings"
	"time"

	"forum/internal/domain"
	"forum/internal/pkg/e3r"

	"golang.org/x/crypto/bcrypt"
)

/*
- [ ] validate data
- [ ] hash password
- [ ] check user
- [ ] save session
*/

type (
	SignInInput struct {
		Email    string
		Password string
	}

	SignInResponse struct {
		Token     string
		ExpiresAt *time.Time
	}
)

func (s *service) SignIn(ctx context.Context, input SignInInput) (*SignInResponse, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}

	user, err := s.users.GetOne(ctx, domain.GetUserInput{
		Email: &input.Email,
	})
	if err != nil {
		return nil, err
	}

	if user.HashedPassword == "" {
		return nil, e3r.BadRequest("Incorrect email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password)); err != nil {
		return nil, e3r.BadRequest("Incorrect email or password")
	}

	if err := s.sessions.Close(ctx, domain.CloseSessionInput{
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC(),
	}); err != nil {
		return nil, err
	}

	token, expiresAt, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &SignInResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (i *SignInInput) validate() error {
	i.Email = strings.TrimSpace(i.Email)
	i.Email = strings.ToLower(i.Email)
	if ok := isValidEmail(i.Email); !ok {
		return domain.ErrInvalidEmail
	}

	if len(i.Password) < 8 {
		return domain.ErrInvalidPasswordLen
	}

	if len(i.Password) > 72 {
		return domain.ErrPasswordTooLong
	}

	i.Password = strings.TrimSpace(i.Password)
	password := strings.ReplaceAll(i.Password, "ㅤ", "")
	if password == "" {
		return domain.ErrInvalidPassword
	}

	return nil
}
