package service

import (
	"context"
	"forum/internal/domain"
)

func (s *service) GetWaitlistUsers(ctx context.Context) ([]domain.User, error) {
	return s.users.GetWaitlist(ctx)
}
