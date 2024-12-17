package service


import (
	"context"

	"forum/internal/domain"
)


func (s *service) GetMyCommentsList(ctx context.Context, input GetPostsListInput) ([]domain.CommentsList, error) {
	commentsList, err := s.comments.GetMyCommentsList(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return commentsList, nil
}