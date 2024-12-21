package service

import (
	"context"
	"forum/internal/domain"
	"strconv"
)

type UpgradeRoleInput struct {
	UserID        string
	WaitingUserID string
}

func (s *service) UpgradeRoleApprove(ctx context.Context, input UpgradeRoleInput) error {
	id, err := input.validate();
	if err != nil {
		return err
	}

	if err := s.allowedRejectUserRole(ctx, input); err != nil {
		return err
	}

	newRole := domain.RoleModerator
	newStatus := false

	return s.users.Update(ctx, domain.UpdateUserInput{
		Role:                 &newRole,
		ModeratorRoleRequest: &newStatus,
		UserID:               id,
	})
}

// TODO: изменить тексты ошибок
func (s *service) allowedRejectUserRole(ctx context.Context, input UpgradeRoleInput) error {
	waitingUser, err := s.users.GetOne(ctx, domain.GetUserInput{
		UserID: &input.WaitingUserID,
	})
	if err != nil {
		return err
	}

	if waitingUser.Role != domain.RoleUser {
		return domain.ErrUserRoleForbidden
	}

	if !waitingUser.ModeratorRoleRequest {
		return domain.ErrUserNotRequestedChangeRole
	}

	user, err := s.users.GetOne(ctx, domain.GetUserInput{
		UserID: &input.UserID,
	})
	if err != nil {
		return err
	}

	if user.Role != domain.RoleAdmin {
		return domain.ErrAdminRoleForbidden
	}

	return nil
}

func (i *UpgradeRoleInput) validate() (int64, error) {

	id, err := strconv.Atoi(i.WaitingUserID);
	if err != nil {
		return 0, domain.ErrInvalidWaitingUserID
	}

	if id <= 0 {
        return 0, domain.ErrInvalidWaitingUserID
    }


	return int64(id), nil
}
