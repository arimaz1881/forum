package adapters

import (
	"context"
	"database/sql"
	"errors"

	"forum/internal/domain"
)

type CategoriesRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.CategoriesRepository = (*CategoriesRepositorySqlite3)(nil)

func NewCategoriesRepositorySqlite3(db *sql.DB) *CategoriesRepositorySqlite3 {
	return &CategoriesRepositorySqlite3{
		db: db,
	}
}

const getCriteriaOne = `
SELECT
  id,
  title
FROM
  categories
WHERE
  id = ?
LIMIT
  1
`

func (q *CategoriesRepositorySqlite3) GetOne(ctx context.Context, id string) (*domain.Categoria, error) {
	row := q.db.QueryRowContext(ctx, getCriteriaOne, id)
	var categoria domain.Categoria

	err := row.Scan(&categoria.ID, &categoria.Title)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCategoryNotFound
	}

	return &categoria, err
}

const getCriteriesMany = `
SELECT
  id,
  title
FROM
  categories
`

func (q *CategoriesRepositorySqlite3) GetMany(ctx context.Context) ([]domain.Categoria, error) {
	rows, err := q.db.QueryContext(ctx, getCriteriesMany)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []domain.Categoria

	for rows.Next() {
		var categoria domain.Categoria
		if err := rows.Scan(&categoria.ID, &categoria.Title); err != nil {
			return nil, err
		}
		categories = append(categories, categoria)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
