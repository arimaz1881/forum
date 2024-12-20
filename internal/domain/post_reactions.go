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
	ID     int64
	PostID int64
	Action string
	UserID int64
}

type Action struct {
	Actions  int
	Actioned bool
}

type PostReactionsRepository interface {
	Create(ctx context.Context, input CreatePostReactionInput) (reactionID int64, err error)
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

type PostReactionList struct {
	LikedPostsList 	[]PostView
	DisikedPostsList	[]PostView
}

var ErrInvalidPostAction = e3r.BadRequest("invalid post action type")
