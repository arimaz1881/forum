package domain

import (
	"context"
	"forum/internal/pkg/e3r"
	"time"
)

type CommentsList struct {
	CommentID      string
	PostID         string
	CommentContent string
	CommentDate    time.Time
	CommentDateStr string
	PostTitle      string
	PostContent    string
	PostDate       time.Time
	PostDateStr    string
	PostAuthor     string
}

type Comment struct {
	ID           string
	Author       string
	AuthorID	 string
	CreatedAt    time.Time
	CreatedAtStr string
	Content      string
	PostID       string
	Likes        *Action
	Dislikes     *Action
}

type CommentsRepository interface {
	Create(ctx context.Context, input CreateCommentInput) error
	Delete(ctx context.Context, commentID string) error
	Edit(ctx context.Context, input EditCommentInput) error
	GetOne(ctx context.Context, commentID string) (*Comment, error)
	GetList(ctx context.Context, postID string) ([]Comment, error)
	GetMyCommentsList(ctx context.Context, userID int64) ([]CommentsList, error)
}

type CreateCommentInput struct {
	PostID  string
	UserID  int64
	Content string
}

type EditCommentInput struct {
	CommentID  	string
	UserID  	int64
	Content 	string
}

var (
	ErrInvalidCommentID      = e3r.BadRequest("invalid comment_id")
	ErrInvalidCommentContent = e3r.BadRequest("invalid comment content")
	ErrCommentNotFound       = e3r.NotFound("comment not found")
)
