package ports

import (
	"fmt"
	"net/http"
	"time"

	"forum/internal/service"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
	Auth    bool
}

func (h *Handler) Routes() []Route {
	return []Route{
		{
			Path:    "/",
			Handler: h.GetPostsList,
			Method:  http.MethodGet,
			Auth:    false,
		},
		{
			Path:    "/posts/my-posts",
			Method:  http.MethodGet,
			Handler: h.GetMyCreatedPosts,
			Auth:    true,
		},
		{
			Path:    "/posts/my-liked",
			Method:  http.MethodGet,
			Handler: h.GetMyLikedPosts,
			Auth:    true,
		},
		{
			Path:    "/posts/create",
			Handler: h.CreatePost,
			Method:  http.MethodPost,
			Auth:    true,
		},
		{
			Path:    "/posts/create",
			Handler: h.CreatePostPage,
			Method:  http.MethodGet,
			Auth:    true,
		},
		{
			Path:    "/posts/view",
			Method:  http.MethodGet,
			Handler: h.GetPostsOne,
			Auth:    false,
		},
		{
			Path:    "/authn/sign-up",
			Method:  http.MethodPost,
			Handler: h.SignUp,
			Auth:    false,
		},
		{
			Path:    "/authn/sign-up",
			Method:  http.MethodGet,
			Handler: h.SignUpPage,
			Auth:    false,
		},
		{
			Path:    "/authn/sign-in",
			Method:  http.MethodPost,
			Handler: h.SignIn,
			Auth:    false,
		},
		{
			Path:    "/authn/sign-in",
			Method:  http.MethodGet,
			Handler: h.SignInPage,
			Auth:    false,
		},
		{
			Path:    "/logout",
			Method:  http.MethodPost,
			Handler: h.Logout,
			Auth:    true,
		},
		{
			Path:    "/posts/reaction",
			Method:  http.MethodPost,
			Handler: h.CreatePostAction,
			Auth:    true,
		},
		{
			Path:    "/comments",
			Method:  http.MethodPost,
			Handler: h.CreateComment,
			Auth:    true,
		},
		{
			Path:    "/comments/reaction",
			Method:  http.MethodPost,
			Handler: h.CreateCommentReaction,
			Auth:    true,
		},
		{
			Path:    "/auth/google/login",
			Method:  http.MethodGet,
			Handler: h.GoogleLogin,
			Auth:    false,
		},
		{
			Path:    "/auth/google/callback",
			Method:  http.MethodGet,
			Handler: h.GoogleCallback,
			Auth:    false,
		},
		{
			Path:    "/auth/github/login",
			Method:  http.MethodGet,
			Handler: h.GitHubLogin,
			Auth:    false,
		},
		{
			Path:    "/auth/github/callback",
			Method:  http.MethodGet,
			Handler: h.GitHubCallback,
			Auth:    false,
		},
	}
}

func (h *Handler) InitRouters() http.Handler {
	rateLimiter := NewRateLimiter(30, 30*time.Second) // 30 requests per 30 sec

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	uploadServer := http.FileServer(http.Dir("./uploads/"))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads", uploadServer))

	for _, route := range h.Routes() {
		handler := route.Handler

		if route.Auth {
			handler = h.requireAuthN(handler)
		}

		pathWMethod := fmt.Sprintf("%s %s", route.Method, route.Path)

		mux.Handle(pathWMethod, h.rateLimit(rateLimiter, (h.withContext(h.recoverPanic(handler)))))
	}

	return mux
}
