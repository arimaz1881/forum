package adapters

import (
	"context"
	"database/sql"

	"forum/internal/domain"
)

type CommentReactionsRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.CommentReactionsRepository = (*CommentReactionsRepositorySqlite3)(nil)

func NewCommentReactionsRepositorySqlite3(db *sql.DB) *CommentReactionsRepositorySqlite3 {
	return &CommentReactionsRepositorySqlite3{
		db: db,
	}
}

const createCommentReaction = `
INSERT INTO
    comment_reactions (
        comment_id,
        user_id,
        action
    )
VALUES
    (?, ?, ?)
ON CONFLICT (comment_id, user_id)
DO UPDATE SET
    action = EXCLUDED.action;
`

func (q *CommentReactionsRepositorySqlite3) Create(ctx context.Context, input domain.CreateCommentReactionInput) error {
	_, err := q.db.ExecContext(ctx, createCommentReaction, input.CommentID, input.UserID, input.Action)
	return err
}

const getCommentReactionsOne = `
SELECT
  comment_id,
  user_id,
  action
FROM
  comment_reactions
WHERE
	user_id = ?
    and comment_id = ?
LIMIT 1;
`

func (q *CommentReactionsRepositorySqlite3) GetOne(ctx context.Context, input domain.GetOneCommentReactionInput) (*domain.CommentReaction, error) {
	row := q.db.QueryRowContext(ctx, getCommentReactionsOne, input.UserID, input.CommentID)

	var CommentReaction domain.CommentReaction

	err := row.Scan(
		&CommentReaction.CommentID,
		&CommentReaction.UserID,
		&CommentReaction.Action,
	)

	return &CommentReaction, err
}

const getCommentReactionsMany = `
SELECT
  comment_id,
  action,
  user_id
FROM
  comment_reactions
WHERE
	comment_id = ?
and 
	action = ?;
`

func (q *CommentReactionsRepositorySqlite3) GetMany(ctx context.Context, input domain.GetManyCommentReactionInput) ([]domain.CommentReaction, error) {
	rows, err := q.db.QueryContext(ctx, getCommentReactionsMany, input.CommentID, input.Action)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentReactions []domain.CommentReaction

	for rows.Next() {
		var commentReaction domain.CommentReaction
		err := rows.Scan(
			&commentReaction.CommentID,
			&commentReaction.Action,
			&commentReaction.UserID,
		)
		if err != nil {
			return nil, err
		}
		commentReactions = append(commentReactions, commentReaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return commentReactions, nil
}
