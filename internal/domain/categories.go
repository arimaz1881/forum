package domain

import (
	"context"
	"forum/internal/pkg/e3r"
)

type Categoria struct {
	ID    string
	Title string
}

type CategoriesRepository interface {
	GetOne(ctx context.Context, categoriaID string) (*Categoria, error)
	GetMany(ctx context.Context) ([]Categoria, error)
}

var (
	ErrInvalidCategory   = e3r.BadRequest("invalid category")
	ErrInvalidCategoryID = e3r.BadRequest("invalid category_id")
	ErrCategoryNotFound  = e3r.NotFound("category not found")
)
