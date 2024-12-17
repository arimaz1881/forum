package adapters

import (
	"context"
	"database/sql"
	"errors"

	"forum/internal/domain"
)

type PostsRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.PostsRepository = (*PostsRepositorySqlite3)(nil)

func NewPostsRepositorySqlite3(db *sql.DB) *PostsRepositorySqlite3 {
	return &PostsRepositorySqlite3{
		db: db,
	}
}

const createPost = `
INSERT INTO
    posts (
        title,
        content,
        user_id,
		file_key
    )
VALUES
    (?, ?, ?, ?) RETURNING id
`

func (q *PostsRepositorySqlite3) Create(ctx context.Context, input domain.CreatePostInput) (int64, error) {
	row := q.db.QueryRowContext(ctx, createPost, input.Title, input.Content, input.UserID, input.FileKey)
	var postID int64
	err := row.Scan(&postID)
	return postID, err
}

const deletePost = `
DELETE FROM
    posts
WHERE
    id = ?
`

func (q *PostsRepositorySqlite3) Delete(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPostsList = `
SELECT
    p.id,
    p.title,
    u.login,
    p.created_at
FROM
    posts p
JOIN
	users u
	on u.id = p.user_id
ORDER BY
    p.created_at desc
`

func (q *PostsRepositorySqlite3) GetList(ctx context.Context) ([]domain.PostView, error) {
	rows, err := q.db.QueryContext(ctx, getPostsList)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var postsList []domain.PostView

	for rows.Next() {
		var post domain.PostView
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Author,
			&post.CreatedAt,
		); err != nil {
			return nil, err
		}
		post.CreatedAtStr = post.CreatedAt.Format("2006-01-02 15:04:05")
		postsList = append(postsList, post)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return postsList, nil
}

const getMyCreatedPostsList = `
SELECT
    p.id,
    p.title,
    u.login,
    p.created_at
FROM
    posts p
JOIN
	users u
	on u.id = p.user_id
WHERE
	p.user_id = ?
ORDER BY
    p.created_at desc
`

func (q *PostsRepositorySqlite3) GetCreatedList(ctx context.Context, userID int64) ([]domain.PostView, error) {
	rows, err := q.db.QueryContext(ctx, getMyCreatedPostsList, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var postsList []domain.PostView

	for rows.Next() {
		var post domain.PostView
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Author,
			&post.CreatedAt,
		); err != nil {
			return nil, err
		}
		post.CreatedAtStr = post.CreatedAt.Format("2006-01-02 15:04:05")
		postsList = append(postsList, post)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return postsList, nil
}

const getMyLikedPostsList = `
SELECT
    p.id,
    p.title,
    u.login,
    p.created_at
FROM
    posts p
JOIN
	users u
	on u.id = p.user_id
JOIN
	post_reactions pr
	on pr.post_id = p.id
WHERE
	pr.user_id = ?
and
	pr.action = ?
ORDER BY
    p.created_at desc
`

func (q *PostsRepositorySqlite3) GetLikedList(ctx context.Context, userID int64, action string) ([]domain.PostView, error) {
	rows, err := q.db.QueryContext(ctx, getMyLikedPostsList, userID, action)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var postsList []domain.PostView

	for rows.Next() {
		var post domain.PostView
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Author,
			&post.CreatedAt,
		); err != nil {
			return nil, err
		}
		post.CreatedAtStr = post.CreatedAt.Format("2006-01-02 15:04:05")
		postsList = append(postsList, post)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return postsList, nil
}

const getPostsOne = `
SELECT
    id,
    created_at,
    title,
    content,
    user_id,
	file_key
FROM
    posts
WHERE
    id = ?
LIMIT
    1
`

func (q *PostsRepositorySqlite3) GetOne(ctx context.Context, id string) (*domain.Post, error) {
	row := q.db.QueryRowContext(ctx, getPostsOne, id)

	var post domain.Post

	err := row.Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.FileKey,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPostNotFound
	}

	post.CreatedAtStr = post.CreatedAt.Format("2006-01-02 15:04:05")

	return &post, err
}
