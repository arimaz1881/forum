package service

import (
	"context"
	"forum/internal/domain"
)

func (s *service) GetMyLikedPosts(ctx context.Context, input GetPostsListInput) ([]domain.PostView, error) {
	if input.UserID == 0 {
		return nil, domain.ErrInvalidUserID
	}

	myLikedPosts, err := s.posts.GetLikedList(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	for i, post := range myLikedPosts {
		categories, err := s.getCatigories(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		myLikedPosts[i].Categories = categories
	}

	return filterPostsByCategories(myLikedPosts, input.CategoryIDs), nil
}
