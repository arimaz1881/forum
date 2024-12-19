package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
  users (email, login, hashed_password, role)
VALUES
  (?, ?, ?, ?) RETURNING id
`

func (q *UsersRepositorySqlite3) Create(ctx context.Context, input domain.CreateUserInput) (int64, error) {
	row := q.db.QueryRowContext(ctx, createUser, input.Email, input.Login, input.HashedPassword, input.Role)
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
  id, role, login, coalesce(email, ''), coalesce(hashed_password, ''), moderator_role_request
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
		&user.Role,
		&user.Login,
		&user.Email,
		&user.HashedPassword,
		&user.ModeratorRoleRequest,
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
  users (oauth_provider, oauth_id, email, login, role)
VALUES
  (?, ?, ?, ?, ?)
`

func (q *UsersRepositorySqlite3) OAuthFindOrCreateUser(ctx context.Context, input domain.GoogleAuthInput) (int64, error) {
	var userID int64
	err := q.db.QueryRowContext(ctx, OAuthGetUser, input.Provider, input.OAuthID).Scan(&userID)
	if err == sql.ErrNoRows {
		// Insert new user
		result, err := q.db.ExecContext(ctx, OAuthCreateUser, input.Provider, input.OAuthID, input.Email, input.Login, input.Role)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint") {
				return 0, domain.ErrUserExists
			}
			return 0, err
		}
		userID, err = result.LastInsertId()
		return userID, err
	}
	return userID, err
}

const getUsersWaitlist = `
SELECT
  id, login 
FROM
  users
WHERE
	moderator_role_request = ? AND role = ?
`

func (q *UsersRepositorySqlite3) GetWaitlist(ctx context.Context) ([]domain.User, error) {
	rows, err := q.db.QueryContext(ctx, getUsersWaitlist, true, domain.RoleUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Login,
		)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		users = append(users, user)
	}

	return users, err
}

const updateUserQuery = `
UPDATE users
SET 
    %s
WHERE id = ?
`

func (q *UsersRepositorySqlite3) Update(ctx context.Context, input domain.UpdateUserInput) error {
	var (
		setClauses []string
		args       []interface{}
	)

	if input.Role != nil {
		setClauses = append(setClauses, "role = ?")
		args = append(args, *input.Role)
	}
	if input.ModeratorRoleRequest != nil {
		setClauses = append(setClauses, "moderator_role_request = ?")
		args = append(args, *input.ModeratorRoleRequest)
	}

	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	// Формируем итоговый запрос
	query := fmt.Sprintf(updateUserQuery, strings.Join(setClauses, ", "))

	// Добавляем ID в конец аргументов
	args = append(args, input.UserID)

	_, err := q.db.ExecContext(ctx, query, args...)
	return err
}
