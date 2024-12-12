package service

import (
	"context"
	"strings"

	"forum/internal/domain"
)

// TODO: add comment for block of code

type CreatePostInput struct {
	Title      string
	Content    string
	Categories []string
	UserID     int64
}

func (s *service) CreatePost(ctx context.Context, input CreatePostInput) (int64, error) {
	if err := input.validate(); err != nil {
		return 0, err
	}

	// check exist category IDs
	for _, categoreiaID := range input.Categories {
		if _, err := s.categories.GetOne(ctx, categoreiaID); err != nil {
			return 0, domain.ErrInvalidCategoryID
		}
	}

	postID, err := s.posts.Create(ctx, domain.CreatePostInput{
		Title:   input.Title,
		Content: input.Content,
		UserID:  input.UserID,
	})
	if err != nil {
		return 0, nil
	}

	for _, categoriaID := range input.Categories {
		if err := s.postCategories.Create(ctx, domain.CreateCategoriaInput{
			PostID:      postID,
			CategoriaID: categoriaID,
		}); err != nil {
			return 0, err
		}
	}

	return postID, nil
}


func (i CreatePostInput) validate() error {
	if i.UserID == 0 {
		return domain.ErrInvalidToken
	}

	if len(i.Categories) == 0 {
		return domain.ErrInvalidCategory
	}

	i.Title = strings.ReplaceAll(i.Title, "ㅤ", "")

	i.Title = strings.TrimSpace(i.Title)
	if i.Title == "" {
		return domain.ErrInvalidContent
	}

	i.Content = strings.ReplaceAll(i.Content, "ㅤ", "")

	i.Content = strings.TrimSpace(i.Content)
	if i.Content == "" {
		return domain.ErrInvalidContent
	}

	return nil
}

