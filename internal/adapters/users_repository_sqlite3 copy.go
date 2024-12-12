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
  id, login, email, hashed_password
FROM
  users
WHERE
  (
    ?1 is not null
    and id = ?1
  )
  or (
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
