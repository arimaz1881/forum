package domain
// методы в adapters унаследованные из моделей в domain нужны для обращения к внешним сервисам

import (
	"context"
	"forum/internal/pkg/e3r"
	"time"
)

type Post struct {
	ID           int64
	Title        string
	Content      string
	UserID       string
	CreatedAt    time.Time
	CreatedAtStr string
	FileKey		 string
}

type PostView struct {
	ID           int64
	Title        string
	Author       string
	CreatedAtStr string
	CreatedAt    time.Time
	Categories   []*Categoria
}

type PostsRepository interface {
	Create(ctx context.Context, input CreatePostInput) (int64, error)
	GetOne(ctx context.Context, postID string) (*Post, error)
	GetList(ctx context.Context) ([]PostView, error)
	GetCreatedList(ctx context.Context, userID int64) ([]PostView, error)
	GetLikedList(ctx context.Context, userID int64, action string) ([]PostView, error)
}

type CreatePostInput struct {
	Title   string
	Content string
	UserID  int64
	FileKey string
}

var (
	ErrInvalidPostID  = e3r.BadRequest("invalid post_id")
	ErrInvalidTitle   = e3r.BadRequest("invalid title")
	ErrInvalidContent = e3r.BadRequest("invalid content")
	ErrPostNotFound   = e3r.NotFound("post not found")
	ErrInvalidFile	  = e3r.BadRequest("invalid file")
)
