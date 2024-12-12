package service

import (
	"context"
	"forum/internal/domain"
)

func (s *service) GetCatigories(ctx context.Context) ([]domain.Categoria, error) {
	return s.categories.GetMany(ctx)
}
