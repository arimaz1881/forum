package service

import (
	"context"
	"forum/internal/domain"
)

func (s *service) GetWaitlistUsers(ctx context.Context, userID int64) ([]domain.User, error) {
	if userID == 0 {
        return nil, domain.ErrInvalidUserID
    }
	if userID != 1 {
		return nil, domain.ErrForbidden
	}
	usersList, err := s.users.GetWaitlist(ctx)
	if err!= nil {
        return nil, err
    }
    return usersList, nil
}
