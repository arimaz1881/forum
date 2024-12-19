package domain

import (
	"context"
	"forum/internal/pkg/e3r"
)

const (
	GoogleClientID     = "" //your client id
	GoogleClientSecret = "" //your client secret
	GoogleRedirectURI  = "https://localhost:8082/auth/google/callback"
	GoogleAuthURL      = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL     = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL  = "https://www.googleapis.com/oauth2/v2/userinfo"
	GithubClientID     = "" //your client id
	GithubClientSecret = "" //your client secret
	GithubRedirectURI  = "https://localhost:8082/auth/github/callback"
)

const (
	RoleGuest     = "guest"
	RoleUser      = "user"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

type User struct {
	ID                   int64
	Role                 string
	Login                string
	Email                string
	HashedPassword       string
	ModeratorRoleRequest bool
}

type UsersRepository interface {
	Create(ctx context.Context, input CreateUserInput) (int64, error)
	Update(ctx context.Context, input UpdateUserInput) error
	GetOne(ctx context.Context, input GetUserInput) (*User, error)
	OAuthFindOrCreateUser(ctx context.Context, input GoogleAuthInput) (int64, error)
	GetWaitlist(ctx context.Context) ([]User, error)
}

type CreateUserInput struct {
	Email          string
	Role           string
	Login          string
	HashedPassword string
}

type GetUserInput struct {
	UserID *string
	Email  *string
	Login  *string
}

type UpdateUserInput struct {
	UserID               int64
	Role                 *string
	ModeratorRoleRequest *bool
}

type GoogleAuthInput struct {
	Provider string
	OAuthID  string
	Email    string
	Login    string
	Password string
	Role	 string
}

var (
	ErrInvalidPasswordLen         = e3r.BadRequest("the length of the coat should be longer")
	ErrPasswordTooLong            = e3r.BadRequest("password too long")
	ErrInvalidUserID              = e3r.BadRequest("invalid user_id")
	ErrInvalidWaitingUserID       = e3r.BadRequest("invalid waiting_user_id")
	ErrInvalidEmail               = e3r.BadRequest("invalid email")
	ErrInvalidPassword            = e3r.BadRequest("invalid password")
	ErrInvalidLogin               = e3r.BadRequest("invalid login")
	ErrIncorrectCredentials       = e3r.BadRequest("incorrect credentials")
	ErrUserExists                 = e3r.BadRequest("user exists")
	ErrUnableToPromoteToModerator = e3r.Forbidden("unable to promote to moderator")
	ErrUserNotFound               = e3r.NotFound("user not found")
	ErrUserRoleForbidden          = e3r.Forbidden("user role forbidden")
	ErrUserNotRequestedChangeRole = e3r.BadRequest("user not requested change role")
	ErrAdminRoleForbidden         = e3r.Forbidden("admin role forbidden")
)
