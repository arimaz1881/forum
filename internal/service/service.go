package service

import (
	"context"

	"forum/internal/domain"
)

type Service interface {
	CreatePost(context.Context, CreatePostInput) (int64, error)

	GetPostsOne(ctx context.Context, input GetPostOneInput) (*GetPostOneResponse, error)
	GetCatigories(context.Context) ([]domain.Categoria, error)

	GetPostsList(ctx context.Context, categories []string) ([]domain.PostView, error)

	GetUserByToken(ctx context.Context, token string) (int64, error)

	CreateComment(ctx context.Context, input CreateCommentInput) error

	SignIn(ctx context.Context, input SignInInput) (*SignInResponse, error)
	SignUp(ctx context.Context, input SignUpInput) (*SignUpResponse, error)
	Logout(ctx context.Context, input LogOutInput) error

	PostReaction(ctx context.Context, input PostReactionInput) error
	CommentReaction(ctx context.Context, input CommentReactionInput) error

	GetMyCreatedPosts(ctx context.Context, input GetPostsListInput) ([]domain.PostView, error)
	GetMyLikedPosts(ctx context.Context, input GetPostsListInput) ([]domain.PostView, error)
}

type service struct {
	users            domain.UsersRepository
	posts            domain.PostsRepository
	categories       domain.CategoriesRepository
	postCategories   domain.PostCategoriesRepository
	postReactions    domain.PostReactionsRepository
	sessions         domain.SessionsRepository
	comments         domain.CommentsRepository
	commentReactions domain.CommentReactionsRepository
}

func NewService(
	users domain.UsersRepository,
	posts domain.PostsRepository,
	categories domain.CategoriesRepository,
	postCategories domain.PostCategoriesRepository,
	postReactions domain.PostReactionsRepository,
	sessions domain.SessionsRepository,
	comments domain.CommentsRepository,
	commentReactions domain.CommentReactionsRepository,
) Service {
	return &service{
		users:            users,
		posts:            posts,
		categories:       categories,
		postCategories:   postCategories,
		postReactions:    postReactions,
		sessions:         sessions,
		comments:         comments,
		commentReactions: commentReactions,
	}
}
