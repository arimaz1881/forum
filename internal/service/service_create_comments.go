package service

import (
	"context"
	"forum/internal/domain"
	"strconv"
	"strings"
)

type CreateCommentInput struct {
	PostID  string
	UserID  int64
	Content string
}

func (s *service) CreateComment(ctx context.Context, input CreateCommentInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	if _, err := s.posts.GetOne(ctx, input.PostID); err != nil {
		return domain.ErrPostNotFound
	}

	return s.comments.Create(ctx, domain.CreateCommentInput{
		PostID:  input.PostID,
		UserID:  input.UserID,
		Content: input.Content,
	})
}

func (i CreateCommentInput) validate() error {
	if _, err := strconv.Atoi(i.PostID); err != nil {
		return domain.ErrInvalidPostID
	}
	if i.PostID == "" {
		return domain.ErrInvalidPostID
	}
	if i.UserID == 0 {
		return domain.ErrInvalidUserID
	}
	i.Content = strings.TrimSpace(i.Content)
	if i.Content == "" {
		return domain.ErrInvalidCommentContent
	}

	return nil
}
