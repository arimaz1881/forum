package domain

import (
	"context"
	"forum/internal/pkg/e3r"
)

type PostCategoriesRepository interface {
	Create(ctx context.Context, input CreateCategoriaInput) error
	GetMany(ctx context.Context, postID int64) ([]string, error)
}

type CreateCategoriaInput struct {
	CategoriaID string
	PostID      int64
}

var (
	ErrPostCategoriesNotFound = e3r.NotFound("post catigories not found")
)
