package domain

import (
	"context"
	"forum/internal/pkg/e3r"
	"time"
)

type Notification struct {
	ID        string
	CreatedAt time.Time
	PostID    string
	PostTitle string
	AuthorID  int64
	Action    string
	Seen      bool
}

type CreateNotificationInput struct {
	PostID    string
	AuthorID  string
	Action    string
	ActionID  int64
	CommentID int64
	Seen      bool
}

type NotificationsRepository interface {
	Create(ctx context.Context, input CreateNotificationInput) error
	Delete(ctx context.Context, reactionID int64) error
	Look(ctx context.Context, id string) error
	GetList(ctx context.Context, userID int64) ([]Notification, error)
	GetOne(ctx context.Context, id string) (*Notification, error)
}

var ErrNotificationNotFound = e3r.NotFound("Notification not found")
