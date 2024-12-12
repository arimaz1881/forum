package domain

import (
	"context"
	"time"

	"forum/internal/pkg/e3r"
)

type Session struct {
	UserID   int64
	Token    string
	ExpresAt time.Time
}

type SessionsRepository interface {
	Create(ctx context.Context, input CreateSessionInput) error
	GetOne(ctx context.Context, token string) (*Session, error)
	Close(ctx context.Context, input CloseSessionInput) error
}

type CreateSessionInput struct {
	UserID   int64
	Token    string
	ExpresAt time.Time
}

type CloseSessionInput struct {
	UserID    int64
	Token     string
	ExpiresAt time.Time
}

var (
	ErrInvalidToken    = e3r.BadRequest("invalid token")
	ErrSessionNotFound = e3r.NotFound("session not found")
)
