package service

import (
	"context"
	"forum/internal/domain"
)

func (s *service) UpgradeRoleReject(ctx context.Context, input UpgradeRoleInput) error {
	id, err := input.validate();
	if err != nil {
		return err
	}


	if err := s.allowedRejectUserRole(ctx, input); err != nil {
		return err
	}

	newRole := domain.RoleUser
	newStatus := false

	return s.users.Update(ctx, domain.UpdateUserInput{
		Role:                 &newRole,
		ModeratorRoleRequest: &newStatus,
		UserID:               id,
	})
}
