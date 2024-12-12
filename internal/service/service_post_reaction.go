package service

import (
	"context"
	"database/sql"
	"errors"

	"forum/internal/domain"
)

type PostReactionInput struct {
	PostID string
	UserID int64
	Action string
}

func (s *service) PostReaction(ctx context.Context, input PostReactionInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	if _, err := s.posts.GetOne(ctx, input.PostID); err != nil {
		return err
	}

	reaction, err := s.postReactions.GetOne(ctx, domain.GetOnePostReactionInput{
		UserID: input.UserID,
		PostID: input.PostID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	newAction := ""

	if reaction.Action != input.Action {
		newAction = input.Action
	}

	if err := s.postReactions.Create(ctx, domain.CreatePostReactionInput{
		UserID: input.UserID,
		PostID: input.PostID,
		Action: newAction,
	}); err != nil {
		return err
	}

	return nil
}

func (i PostReactionInput) validate() error {
	if i.PostID == "" {
		return domain.ErrInvalidPostID
	}

	if i.UserID <= 0 {
		return domain.ErrInvalidUserID
	}

	if i.Action != domain.ActionTypeLike && i.Action != domain.ActionTypeDislike {
		return domain.ErrInvalidPostAction
	}

	return nil
}
