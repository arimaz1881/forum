package adapters

import (
	"context"
	"database/sql"
	"errors"

	"forum/internal/domain"
)

type CommentsRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.CommentsRepository = (*CommentsRepositorySqlite3)(nil)

func NewCommentsRepositorySqlite3(db *sql.DB) *CommentsRepositorySqlite3 {
	return &CommentsRepositorySqlite3{
		db: db,
	}
}

const createComment = `
INSERT INTO
	comments (
        post_id,
        content,
        user_id
    )
VALUES
    (?, ?, ?);
`

func (q *CommentsRepositorySqlite3) Create(ctx context.Context, input domain.CreateCommentInput) error {
	_, err := q.db.ExecContext(ctx, createComment, input.PostID, input.Content, input.UserID)
	return err
}

const getCommentsList = `
SELECT
    c.id,
	c.post_id,
    c.content,
    u.login,
    c.created_at
FROM
    comments c
JOIN
	users u
	on u.id = c.user_id
WHERE 
	c.post_id = ?
ORDER BY
    c.created_at desc
`

func (q *CommentsRepositorySqlite3) GetList(ctx context.Context, postID string) ([]domain.Comment, error) {
	rows, err := q.db.QueryContext(ctx, getCommentsList, postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var commentsList []domain.Comment

	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.Author,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}
		comment.CreatedAtStr = comment.CreatedAt.Format("2006-01-02 15:04:05")
		commentsList = append(commentsList, comment)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commentsList, nil
}

const getCommentsOne = `
SELECT
    c.id,
	c.post_id,
    c.content,
    u.login,
    c.created_at
FROM
    comments c
JOIN
	users u
	on u.id = c.user_id
WHERE 
	c.id = ?
LIMIT
	1
`

func (q *CommentsRepositorySqlite3) GetOne(ctx context.Context, commentID string) (*domain.Comment, error) {
	rows := q.db.QueryRowContext(ctx, getCommentsOne, commentID)
	var comment domain.Comment
	err := rows.Scan(
		&comment.ID,
		&comment.PostID,
		&comment.Content,
		&comment.Author,
		&comment.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCommentNotFound
	}

	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &comment, nil
}
