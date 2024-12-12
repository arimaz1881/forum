package service

import (
	"context"
	"regexp"
	"strings"
	"time"

	"forum/internal/domain"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
- [x] validate data
- [x] hash password
- [x] create user
- [x] save session
*/

type (
	SignUpInput struct {
		Email    string
		Login    string
		Password string
	}
	SignUpResponse struct {
		Token     string
		ExpiresAt *time.Time
	}
)

func (s *service) SignUp(ctx context.Context, input SignUpInput) (*SignUpResponse, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 6)
	if err != nil {
		return nil, err
	}

	userID, err := s.users.Create(ctx, domain.CreateUserInput{
		Email:          input.Email,
		Login:          input.Login,
		HashedPassword: string(hashedPasswordBytes),
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

	return &SignUpResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *service) createSession(ctx context.Context, userID int64) (string, *time.Time, error) {
	newToken, err := uuid.NewV4()
	if err != nil {
		return "", nil, err
	}

	var (
		expiresAt = time.Now().UTC().Add(time.Hour * 24)
		tokenStr  = newToken.String()
	)

	if err := s.sessions.Create(ctx, domain.CreateSessionInput{
		UserID:   userID,
		Token:    tokenStr,
		ExpresAt: expiresAt,
	}); err != nil {
		return "", nil, err
	}

	return tokenStr, &expiresAt, nil
}

func (i *SignUpInput) validate() error {
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

	// TODO: check valid password
	i.Password = strings.TrimSpace(i.Password)
	password := strings.ReplaceAll(i.Password, "ㅤ", "")
	if password == "" {
		return domain.ErrInvalidPassword
	}

	i.Login = strings.TrimSpace(i.Login)
	login := strings.ReplaceAll(i.Password, "ㅤ", "")
	if login == "" {
		return domain.ErrInvalidLogin
	}

	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
