package service

import (
	"context"

	"forum/internal/domain"
)

type Service interface {
	CreatePost(context.Context, CreatePostInput) (int64, error)
	DeletePost(ctx context.Context, input DeletePostInput) error
	EditPost(ctx context.Context, input EditPostInput) error

	GetPostsOne(ctx context.Context, input GetPostOneInput) (*GetPostOneResponse, error)
	GetCatigories(context.Context) ([]domain.Categoria, error)

	GetPostsList(ctx context.Context, categories []string) ([]domain.PostView, error)

	GetUserByToken(ctx context.Context, token string) (*GetUserByTokenResponse, error)

	CreateComment(ctx context.Context, input CreateCommentInput) error
	DeleteComment(ctx context.Context, input DeleteCommentInput) error
	EditComment(ctx context.Context, input EditCommentInput) error


	SignIn(ctx context.Context, input SignInInput) (*SignInResponse, error)
	SignUp(ctx context.Context, input SignUpInput) (*SignUpResponse, error)
	Logout(ctx context.Context, input LogOutInput) error

	PostReaction(ctx context.Context, input PostReactionInput) error
	CommentReaction(ctx context.Context, input CommentReactionInput) error

	GetNotificationsList(ctx context.Context, userID int64) ([]domain.Notification, error)
	LookNotification(ctx context.Context, input LookNotificationInput) error

	GetMyCreatedPosts(ctx context.Context, input GetPostsListInput) ([]domain.PostView, error)
	GetMyLikedPosts(ctx context.Context, input GetPostsListInput, action string) ([]domain.PostView, error)
	GetMyCommentsList(ctx context.Context, input GetPostsListInput) ([]domain.CommentsList, error)

	GoogleAuth(ctx context.Context, code string) (*SignInResponse, error)
	GithubAuth(ctx context.Context, code string) (*SignInResponse, error)

	// moderator
	SubmitRoleUpgrade(ctx context.Context, userID string) error
	GetWaitlistUsers(ctx context.Context, userID int64) ([]domain.User, error)
	UpgradeRoleApprove(ctx context.Context, input UpgradeRoleInput) error
	UpgradeRoleReject(ctx context.Context, input UpgradeRoleInput) error
	GetModerators(ctx context.Context, userID int64) ([]domain.User, error)
}

type service struct {
	users      		 domain.UsersRepository
	posts      		 domain.PostsRepository
	categories 		 domain.CategoriesRepository
	notifications	 domain.NotificationsRepository
	postCategories   domain.PostCategoriesRepository
	postReactions    domain.PostReactionsRepository
	sessions         domain.SessionsRepository
	comments         domain.CommentsRepository
	commentReactions domain.CommentReactionsRepository
	fileStorage      string
}

func NewService(
	users domain.UsersRepository,
	posts domain.PostsRepository,
	categories domain.CategoriesRepository,
	notifications domain.NotificationsRepository,
	postCategories domain.PostCategoriesRepository,
	postReactions domain.PostReactionsRepository,
	sessions domain.SessionsRepository,
	comments domain.CommentsRepository,
	commentReactions domain.CommentReactionsRepository,
	fileStorage string,
) Service {
	return &service{
		users:            users,
		posts:            posts,
		categories:       categories,
		notifications:	  notifications,
		postCategories:   postCategories,
		postReactions:    postReactions,
		sessions:         sessions,
		comments:         comments,
		commentReactions: commentReactions,
		fileStorage:      fileStorage,
	}
}
