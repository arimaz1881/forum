package adapters

import (
	"context"
	"database/sql"

	"forum/internal/domain"
)

type PostCriteriesRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.PostCategoriesRepository = (*PostCriteriesRepositorySqlite3)(nil)

func NewPostCategoriesRepositorySqlite3(db *sql.DB) *PostCriteriesRepositorySqlite3 {
	return &PostCriteriesRepositorySqlite3{
		db: db,
	}
}

const createPostCatigories = `
INSERT INTO
  posts_categories (post_id, categoria_id)
VALUES
  (?, ?)
`

func (q *PostCriteriesRepositorySqlite3) Create(ctx context.Context, arg domain.CreateCategoriaInput) error {
	_, err := q.db.ExecContext(ctx, createPostCatigories, arg.PostID, arg.CategoriaID)
	return err
}

const getPostCatigoriesMany = `
SELECT
  categoria_id
FROM
  posts_categories
WHERE
  post_id = ?
`

func (q *PostCriteriesRepositorySqlite3) GetMany(ctx context.Context, postID int64) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getPostCatigoriesMany, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []string

	for rows.Next() {
		var i string
		if err := rows.Scan(&i); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
