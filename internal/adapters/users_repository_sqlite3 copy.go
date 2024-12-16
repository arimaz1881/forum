package adapters

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"forum/internal/domain"
)

type UsersRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.UsersRepository = (*UsersRepositorySqlite3)(nil)

func NewUsersRepositorySqlite3(db *sql.DB) *UsersRepositorySqlite3 {
	return &UsersRepositorySqlite3{
		db: db,
	}
}

const createUser = `-- name: Create :exec
INSERT INTO
  users (email, login, hashed_password)
VALUES
  (?, ?, ?) RETURNING id
`

func (q *UsersRepositorySqlite3) Create(ctx context.Context, input domain.CreateUserInput) (int64, error) {
	row := q.db.QueryRowContext(ctx, createUser, input.Email, input.Login, input.HashedPassword)
	var userID int64
	err := row.Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return 0, domain.ErrUserExists
		}
		return 0, err
	}

	return userID, err
}

const getUsersOne = `
SELECT
  id, login, email, coalesce(hashed_password, '')
FROM
  users
WHERE
  (
    ?1 is not null
    and id = ?1
  )
  OR (
    ?2 is not null
    and email = ?2
  )
`

func (q *UsersRepositorySqlite3) GetOne(ctx context.Context, input domain.GetUserInput) (*domain.User, error) {
	row := q.db.QueryRowContext(ctx, getUsersOne, input.UserID, input.Email)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.HashedPassword,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	return &user, err
}

const OAuthGetUser = `
SELECT 
  id 
FROM
  users
WHERE
  oauth_provider = ? AND oauth_id = ?
`

const OAuthCreateUser = `
INSERT INTO
  users (oauth_provider, oauth_id, email, login)
VALUES
  (?, ?, ?, ?)
`

func (q *UsersRepositorySqlite3) OAuthFindOrCreateUser(ctx context.Context, input domain.GoogleAuthInput) (int64, error) {
	var userID int64
	err := q.db.QueryRowContext(ctx, OAuthGetUser, input.Provider, input.OAuthID).Scan(&userID)
	if err == sql.ErrNoRows {
		// Insert new user
		result, err := q.db.ExecContext(ctx, OAuthCreateUser, input.Provider, input.OAuthID, input.Email, input.Login)
		if err != nil {
			return 0, err
		}
		userID, err = result.LastInsertId()
		return userID, err
	}
	return userID, err
}
