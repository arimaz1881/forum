package service

import (
	"context"

	"forum/internal/domain"
)

type LookNotificationInput struct {
	UserID         int64
	NotificationID string
}

func (s *service) GetNotificationsList(ctx context.Context, userID int64) ([]domain.Notification, error) {
	notifications, err := s.notifications.GetList(ctx, userID)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *service) LookNotification(ctx context.Context, input LookNotificationInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	notif, err := s.notifications.GetOne(ctx, input.NotificationID)
	if err != nil {
		return domain.ErrNotificationNotFound
	}

	if notif.AuthorID != input.UserID {
		return domain.ErrForbidden
	}

	return s.notifications.Look(ctx, input.NotificationID)
}

func (i LookNotificationInput) validate() error {
	if i.UserID <= 0 {
		return domain.ErrInvalidUserID
	}

	if i.NotificationID == "" {
		return domain.ErrInvalidPostID
	}

	return nil
}
