package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/internal/domain"
)

type SessionsRepositorySqlite3 struct {
	db *sql.DB
}

var _ domain.SessionsRepository = (*SessionsRepositorySqlite3)(nil)

func NewSessionsRepositorySqlite3(db *sql.DB) *SessionsRepositorySqlite3 {
	return &SessionsRepositorySqlite3{
		db: db,
	}
}

const createSession = `
INSERT INTO
  sessions (user_id, token, expires_at)
VALUES
  (?, ?, ?)
`

func (q *SessionsRepositorySqlite3) Create(ctx context.Context, input domain.CreateSessionInput) error {
	if _, err := q.db.ExecContext(
		ctx,
		createSession,
		input.UserID,
		input.Token,
		input.ExpresAt,
	); err != nil {
		return err
	}

	return nil
}

const getSessionsOne = `
SELECT
  user_id, token, expires_at
FROM
  sessions
WHERE
	token = ?
`

func (q *SessionsRepositorySqlite3) GetOne(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	var expiresAtStr string

	if err := q.db.QueryRowContext(
		ctx,
		getSessionsOne,
		token,
	).Scan(
		&session.UserID,
		&session.Token,
		&expiresAtStr,
	); err != nil {
		return nil, err
	}

	expiresAt, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", expiresAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expires_at: %v", err)
	}

	session.ExpresAt = expiresAt
	return &session, nil
}

const closeSession = `
UPDATE sessions 
set
  expires_at = ?	
WHERE
  token = ? OR user_id = ?
`

func (q *SessionsRepositorySqlite3) Close(ctx context.Context, input domain.CloseSessionInput) error {
	_, err := q.db.ExecContext(ctx, closeSession, input.ExpiresAt, input.Token, input.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrSessionNotFound
	}

	return err
}
