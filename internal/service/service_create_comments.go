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

type DeleteCommentInput struct {
	CommentID	string
	UserID	int64
}

type EditCommentInput struct {
	CommentID	string
	Content 	string
	UserID		int64
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


func (s *service) DeleteComment(ctx context.Context, input DeleteCommentInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	comment, err := s.comments.GetOne(ctx, input.CommentID); 
	if err != nil {
		return err
	}

	if comment.AuthorID != strconv.Itoa(int(input.UserID)) {
		return domain.ErrForbidden
	}

	if err := s.comments.Delete(ctx, input.CommentID); err != nil {
		return err
	}

	return nil
}


func (i DeleteCommentInput) validate() error {
	if i.CommentID == "" {
		return domain.ErrInvalidPostID
	}

	if i.UserID <= 0 {
		return domain.ErrInvalidUserID
	}

	return nil
}


func (s *service) EditComment(ctx context.Context, input EditCommentInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	comment, err := s.comments.GetOne(ctx, input.CommentID); 
	if err != nil {
		return err
	}

	if comment.AuthorID != strconv.Itoa(int(input.UserID)) {
		return domain.ErrForbidden
	}

	if err := s.comments.Edit(ctx, domain.EditCommentInput{
		CommentID: input.CommentID,
		Content: input.Content}); 
		err != nil {
		return err
	}

	return nil
}


func (i EditCommentInput) validate() error {
	if i.CommentID == "" {
		return domain.ErrInvalidPostID
	}

	if i.UserID <= 0 {
		return domain.ErrInvalidUserID
	}

	i.Content = strings.ReplaceAll(i.Content, "ㅤ", "")

	i.Content = strings.TrimSpace(i.Content)
	if i.Content == "" {
		return domain.ErrInvalidContent
	}

	return nil
}