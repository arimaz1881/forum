package service

import (
	"context"

	"forum/internal/domain"
)

type GetPostsListInput struct {
	UserID      int64
	CategoryIDs []string
}

func (s *service) GetMyCreatedPosts(ctx context.Context, input GetPostsListInput) ([]domain.PostView, error) {
	if input.UserID == 0 {
		return nil, domain.ErrInvalidUserID
	}

	myCreatedPosts, err := s.posts.GetCreatedList(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	for i, post := range myCreatedPosts {
		categories, err := s.getCatigories(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		myCreatedPosts[i].Categories = categories
	}

	return filterPostsByCategories(myCreatedPosts, input.CategoryIDs), nil
}
