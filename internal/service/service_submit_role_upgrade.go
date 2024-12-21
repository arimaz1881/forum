package service

import (
	"context"
	"forum/internal/domain"
)

func (s *service) SubmitRoleUpgrade(ctx context.Context, userID string) error {
	user, err := s.users.GetOne(ctx, domain.GetUserInput{
		UserID: &userID,
	})
	if err != nil {
		return err
	}

	if user.Role != domain.RoleUser {
		return domain.ErrUnableToPromoteToModerator
	}

	newModeratorRoleRequest := true

	return s.users.Update(ctx, domain.UpdateUserInput{
		UserID:               user.ID,
		Role:                 &user.Role,
		ModeratorRoleRequest: &newModeratorRoleRequest,
	})
}
