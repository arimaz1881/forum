package domain

import (
	"context"
	"forum/internal/pkg/e3r"
)

var (
	ActionTypeLike    = "like"
	ActionTypeDislike = "dislike"
)

type PostReaction struct {
	PostID int64
	Action string
	UserID int64
}

type Action struct {
	Actions  int
	Actioned bool
}

type PostReactionsRepository interface {
	Create(ctx context.Context, input CreatePostReactionInput) error
	GetOne(ctx context.Context, input GetOnePostReactionInput) (*PostReaction, error)
	GetMany(ctx context.Context, input GetManyPostReactionInput) ([]PostReaction, error)
}

type CreatePostReactionInput struct {
	PostID string
	Action string
	UserID int64
}

type GetOnePostReactionInput struct {
	PostID string
	UserID int64
}

type GetManyPostReactionInput struct {
	PostID string
	Action string
}

var ErrInvalidPostAction = e3r.BadRequest("invalid post action type")
