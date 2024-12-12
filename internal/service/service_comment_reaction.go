package service

import (
	"context"
	"database/sql"
	"errors"

	"forum/internal/domain"
)

type CommentReactionInput struct {
	PostID    string
	CommentID string
	UserID    int64
	Action    string
}

func (s *service) CommentReaction(ctx context.Context, input CommentReactionInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	if _, err := s.posts.GetOne(ctx, input.PostID); err != nil {
		return err
	}

	comment, err := s.comments.GetOne(ctx, input.CommentID)
	if err != nil {
		return err
	}

	if comment.PostID != input.PostID {
		return domain.ErrPostNotFound
	}

	reaction, err := s.commentReactions.GetOne(ctx, domain.GetOneCommentReactionInput{
		UserID:    input.UserID,
		CommentID: input.CommentID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	newAction := ""

	if reaction.Action != input.Action {
		newAction = input.Action
	}

	if err := s.commentReactions.Create(ctx, domain.CreateCommentReactionInput{
		CommentID: input.CommentID,
		UserID:    input.UserID,
		Action:    newAction,
	}); err != nil {
		return err
	}

	return nil
}

func (i CommentReactionInput) validate() error {
	if i.CommentID == "" {
		return domain.ErrInvalidCommentID
	}

	if i.UserID <= 0 {
		return domain.ErrInvalidUserID
	}

	if i.Action != domain.ActionTypeLike && i.Action != domain.ActionTypeDislike {
		return domain.ErrInvalidPostAction
	}

	return nil
}
