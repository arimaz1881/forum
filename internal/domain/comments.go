package domain

import (
	"context"
	"forum/internal/pkg/e3r"
	"time"
)

type Comment struct {
	ID           string
	Author       string
	CreatedAt    time.Time
	CreatedAtStr string
	Content      string
	PostID       string
	Likes        *Action
	Dislikes     *Action
}

type CommentsRepository interface {
	Create(ctx context.Context, input CreateCommentInput) error
	GetOne(ctx context.Context, commentID string) (*Comment, error)
	GetList(ctx context.Context, postID string) ([]Comment, error)
}

type CreateCommentInput struct {
	PostID  string
	UserID  int64
	Content string
}

var (
	ErrInvalidCommentID      = e3r.BadRequest("invalid comment_id")
	ErrInvalidCommentContent = e3r.BadRequest("invalid comment content")
	ErrCommentNotFound       = e3r.NotFound("comment not found")
)
