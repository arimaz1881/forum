package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"forum/internal/domain"
)

/*
TODO: repository decorator check error
*/

/*
	-[x] validate
	-[x] change postID -> int64
	-[x] get post
	-[x] get categories
	-[x] get author

// TODO:

	-[ ] get reaction
	-[ ] get comments
*/
type GetPostOneInput struct {
	PostID string
	UserID int64
}

type GetPostOneResponse struct {
	domain.Post
	Author     string
	Categories []*domain.Categoria
	Likes      *domain.Action
	Dislikes   *domain.Action
	Comments   []domain.Comment
}

func (s *service) GetPostsOne(ctx context.Context, input GetPostOneInput) (*GetPostOneResponse, error) {
	if _, err := strconv.Atoi(input.PostID); err != nil {
		return nil, domain.ErrInvalidPostID
	}

	post, err := s.posts.GetOne(ctx, input.PostID)
	if err != nil {
		return nil, err
	}

	categories, err := s.getCatigories(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	author, err := s.users.GetOne(ctx, domain.GetUserInput{
		UserID: &post.UserID,
	})
	if err != nil {
		return nil, err
	}

	likes, err := s.getActions(ctx, input, domain.ActionTypeLike)
	if err != nil {
		return nil, err
	}

	dislikes, err := s.getActions(ctx, input, domain.ActionTypeDislike)
	if err != nil {
		return nil, err
	}

	comments, err := s.comments.GetList(ctx, input.PostID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	for i := range comments {
		commentLikes, err := s.getCommentActions(ctx, input.UserID, comments[i].ID, domain.ActionTypeLike)
		if err != nil {
			return nil, err
		}
		commentDislikes, err := s.getCommentActions(ctx, input.UserID, comments[i].ID, domain.ActionTypeDislike)
		if err != nil {
			return nil, err
		}
		comments[i].Likes = commentLikes
		comments[i].Dislikes = commentDislikes
	}

	return &GetPostOneResponse{
		Post:       *post,
		Author:     author.Login,
		Categories: categories,
		Likes:      likes,
		Dislikes:   dislikes,
		Comments:   comments,
	}, nil
}

func (s *service) getActions(ctx context.Context, input GetPostOneInput, action string) (*domain.Action, error) {
	actions, err := s.postReactions.GetMany(ctx, domain.GetManyPostReactionInput{
		PostID: input.PostID,
		Action: action,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return &domain.Action{}, err
	}
	if err != nil {
		return nil, err
	}

	return &domain.Action{
		Actions:  len(actions),
		Actioned: isActioned(actions, input.UserID),
	}, nil
}

func isActioned(actions []domain.PostReaction, userID int64) bool {
	for _, action := range actions {
		if action.UserID == userID && action.Action != "" {
			return true
		}
	}

	return false
}

func (s *service) getCommentActions(ctx context.Context, userID int64, commentID, action string) (*domain.Action, error) {
	actions, err := s.commentReactions.GetMany(ctx, domain.GetManyCommentReactionInput{
		CommentID: commentID,
		Action:    action,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return &domain.Action{}, err
	}
	if err != nil {
		return nil, err
	}

	return &domain.Action{
		Actions:  len(actions),
		Actioned: isCommentActioned(actions, userID),
	}, nil
}

func isCommentActioned(actions []domain.CommentReaction, userID int64) bool {
	for _, action := range actions {
		if action.UserID == userID && action.Action != "" {
			return true
		}
	}

	return false
}
