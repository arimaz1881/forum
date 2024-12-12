package service

import (
	"context"
	"time"

	"forum/internal/domain"

	"github.com/gofrs/uuid"
)

type LogOutInput struct {
	Token string
}

func (s *service) Logout(ctx context.Context, input LogOutInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	if err := s.sessions.Close(ctx, domain.CloseSessionInput{
		Token:     input.Token,
		ExpiresAt: time.Now().UTC(),
	}); err != nil {
		return err
	}

	return nil
}

func (i LogOutInput) validate() error {
	if i.Token == "" {
		return domain.ErrInvalidToken
	}

	if err := uuid.NamespaceDNS.Parse(i.Token); err != nil {
		return domain.ErrInvalidToken
	}

	return nil
}
