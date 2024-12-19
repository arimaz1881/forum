package service

import (
	"context"
	"forum/internal/domain"
	"strconv"
)

type UpgradeRoleRejectInput struct {
	UserID        string
	WaitingUserID string
}

func (s *service) UpgradeRoleReject(ctx context.Context, input UpgradeRoleRejectInput) error {
	if err := input.validate(); err != nil {
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
	})
}

// TODO: изменить тексты ошибок
func (s *service) allowedRejectUserRole(ctx context.Context, input UpgradeRoleRejectInput) error {
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

func (i *UpgradeRoleRejectInput) validate() error {
	if _, err := strconv.Atoi(i.WaitingUserID); err != nil {
		return domain.ErrInvalidWaitingUserID
	}

	if _, err := strconv.Atoi(i.WaitingUserID); err != nil {
		return domain.ErrInvalidUserID
	}

	return nil
}
