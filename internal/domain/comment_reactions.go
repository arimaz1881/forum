package domain

import "context"

type CommentReaction struct {
	CommentID int64
	Action    string
	UserID    int64
}

type CommentReactionsRepository interface {
	Create(ctx context.Context, input CreateCommentReactionInput) error
	GetOne(ctx context.Context, input GetOneCommentReactionInput) (*CommentReaction, error)
	GetMany(ctx context.Context, input GetManyCommentReactionInput) ([]CommentReaction, error)
}

type CreateCommentReactionInput struct {
	CommentID string
	Action    string
	UserID    int64
}

type GetOneCommentReactionInput struct {
	CommentID string
	UserID    int64
}

type GetManyCommentReactionInput struct {
	CommentID string
	Action    string
}
