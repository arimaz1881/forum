package service

import (
	"context"
	"time"
)

func (s *service) GetUserByToken(ctx context.Context, token string) (int64, error) {
	session, err := s.sessions.GetOne(ctx, token)
	if err != nil {
		return 0, err
	}

	var (
		now     = time.Now().UTC()
		expired = now.After(session.ExpresAt)
	)

	if expired {
		return 0, nil
	}

	return session.UserID, nil
}
