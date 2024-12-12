package adapters

import (
	"context"
	"database/sql"

	"forum/internal/domain"
)

type PostReactionsRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.PostReactionsRepository = (*PostReactionsRepositorySqlite3)(nil)

func NewPostReactionsRepositorySqlite3(db *sql.DB) *PostReactionsRepositorySqlite3 {
	return &PostReactionsRepositorySqlite3{
		db: db,
	}
}

const createPostReaction = `
INSERT INTO
    post_reactions (
        post_id,
        user_id,
        action
    )
VALUES
    (?, ?, ?)
ON CONFLICT (post_id, user_id)
DO UPDATE SET
    action = EXCLUDED.action;
`

func (q *PostReactionsRepositorySqlite3) Create(ctx context.Context, input domain.CreatePostReactionInput) error {
	_, err := q.db.ExecContext(ctx, createPostReaction, input.PostID, input.UserID, input.Action)
	return err
}

const getPostReactionsOne = `
SELECT
  post_id,
  user_id,
  action
FROM
  post_reactions
WHERE
	user_id = ?
    and post_id = ?
LIMIT 1;
`

func (q *PostReactionsRepositorySqlite3) GetOne(ctx context.Context, input domain.GetOnePostReactionInput) (*domain.PostReaction, error) {
	row := q.db.QueryRowContext(ctx, getPostReactionsOne, input.UserID, input.PostID)

	var postReaction domain.PostReaction

	err := row.Scan(
		&postReaction.PostID,
		&postReaction.UserID,
		&postReaction.Action,
	)

	return &postReaction, err
}

const getPostReactionsMany = `
SELECT
  post_id,
  action,
  user_id
FROM
  post_reactions
WHERE
	post_id = ?
and 
	action = ?;
`

func (q *PostReactionsRepositorySqlite3) GetMany(ctx context.Context, input domain.GetManyPostReactionInput) ([]domain.PostReaction, error) {
	rows, err := q.db.QueryContext(ctx, getPostReactionsMany, input.PostID, input.Action)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postReactions []domain.PostReaction

	for rows.Next() {
		var postReaction domain.PostReaction
		err := rows.Scan(
			&postReaction.PostID,
			&postReaction.Action,
			&postReaction.UserID,
		)
		if err != nil {
			return nil, err
		}
		postReactions = append(postReactions, postReaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return postReactions, nil
}
