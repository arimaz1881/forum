package domain

import (
	"context"
	"forum/internal/pkg/e3r"
)

type User struct {
	ID             int64
	Login          string
	Email          string
	HashedPassword string
}

type UsersRepository interface {
	Create(ctx context.Context, input CreateUserInput) (int64, error)
	GetOne(ctx context.Context, input GetUserInput) (*User, error)
}

type CreateUserInput struct {
	Email          string
	Login          string
	HashedPassword string
}

type GetUserInput struct {
	UserID *string
	Email  *string
	Login  *string
}

var (
	ErrInvalidPasswordLen   = e3r.BadRequest("the length of the coat should be longer")
	ErrPasswordTooLong      = e3r.BadRequest("password too long")
	ErrInvalidUserID        = e3r.BadRequest("invalid user_id")
	ErrInvalidEmail         = e3r.BadRequest("invalid email")
	ErrInvalidPassword      = e3r.BadRequest("invalid password")
	ErrInvalidLogin         = e3r.BadRequest("invalid login")
	ErrIncorrectCredentials = e3r.BadRequest("incorrect credentials")
	ErrUserExists           = e3r.BadRequest("user exist")
	ErrUserNotFound         = e3r.NotFound("user not found")
)
