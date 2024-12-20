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

	post, err := s.posts.GetOne(ctx, input.PostID)
	if err != nil {
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

	reactionID, err := s.postReactions.Create(ctx, domain.CreatePostReactionInput{
		UserID: input.UserID,
		PostID: input.PostID,
		Action: newAction,
	})
	if err != nil {
		return err
	}

	if newAction != "" {
		if err = s.notifications.Create(ctx, domain.CreateNotificationInput{
			PostID:   input.PostID,
			AuthorID: post.UserID,
			Action:   newAction,
			ActionID: reactionID,
			Seen:     false,
		}); err != nil {
			return err
		}
	} else {
		if err = s.notifications.Delete(ctx, reactionID); err!= nil {
            return err
        }
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
