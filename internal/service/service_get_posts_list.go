package service

import (
	"context"

	"forum/internal/domain"
)

/*
TODO: repository decorator check error
*/

/*
- validate
- change postID -> int64
- get post
- get categories
- get author
// TODO:
- get reaction
- get comments
*/

func (s *service) GetPostsList(ctx context.Context, categoryIDs []string) ([]domain.PostView, error) {
	postsList, err := s.posts.GetList(ctx)
	if err != nil {
		return nil, err
	}

	for i, post := range postsList {
		categories, err := s.getCatigories(ctx, post.ID)
		if err != nil {
			return nil, err
		}

		postsList[i].Categories = categories
	}

	return filterPostsByCategories(postsList, categoryIDs), nil
}

func (s *service) getCatigories(ctx context.Context, postID int64) ([]*domain.Categoria, error) {
	categoryIDs, err := s.postCategories.GetMany(ctx, postID)
	if err != nil {
		return nil, err
	}

	categories := make([]*domain.Categoria, 0, len(categoryIDs))

	for _, categoriaID := range categoryIDs {
		categoria, err := s.categories.GetOne(ctx, categoriaID)
		if err != nil {
			return nil, err
		}

		categories = append(categories, categoria)
	}

	return categories, nil
}

func filterPostsByCategories(posts []domain.PostView, categoryIDs []string) []domain.PostView {
	if len(categoryIDs) == 0 {
		return posts
	}

	var filteredPosts []domain.PostView

	categoryIDMap := make(map[string]struct{})
	for _, id := range categoryIDs {
		categoryIDMap[id] = struct{}{}
	}

	for _, post := range posts {
		found := false
		for _, cat := range post.Categories {
			if _, exists := categoryIDMap[cat.ID]; exists {
				found = true
				break
			}
		}

		if found {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}
