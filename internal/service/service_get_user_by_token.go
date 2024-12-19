package service

import (
	"context"
	"forum/internal/domain"
	"strconv"
	"time"
)

type GetUserByTokenResponse struct {
	ID             int64
	Role           string
	CanSendRequest bool
	Login		   string
}

func (s *service) GetUserByToken(ctx context.Context, token string) (*GetUserByTokenResponse, error) {
	session, err := s.sessions.GetOne(ctx, token)
	if err != nil {
		return nil, err
	}

	var (
		now     = time.Now().UTC()
		expired = now.After(session.ExpresAt)
		userID  = strconv.Itoa(int(session.UserID))
	)

	user, err := s.users.GetOne(ctx, domain.GetUserInput{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}

	if expired {
		return &GetUserByTokenResponse{
			Role: user.Role,
		}, nil
	}

	return &GetUserByTokenResponse{
		ID:             session.UserID,
		Role:           user.Role,
		CanSendRequest: !user.ModeratorRoleRequest,
		Login:			user.Login,
	}, nil
}
