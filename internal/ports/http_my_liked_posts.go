package ports

import (
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
	"forum/internal/service"
)

func (h *Handler) GetMyLikedPosts(w http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		user        = getUserData(ctx)
		categoryIDs = r.URL.Query()["categories"]
	)

	myActionPosts, err := h.svc.GetMyLikedPosts(
		ctx,
		service.GetPostsListInput{
			UserID:      user.ID,
			CategoryIDs: categoryIDs,
		},
		"like",
	)
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	httphelper.Render(w, http.StatusOK, "reacted", httphelper.GetTmplData(myActionPosts, user))

}

func (h *Handler) GetMyDislikedPosts(w http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		user        = getUserData(ctx)
		categoryIDs = r.URL.Query()["categories"]
	)

	myActionPosts, err := h.svc.GetMyLikedPosts(
		ctx,
		service.GetPostsListInput{
			UserID:      user.ID,
			CategoryIDs: categoryIDs,
		},
		"dislike",
	)
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	httphelper.Render(w, http.StatusOK, "reacted", httphelper.GetTmplData(myActionPosts, user))

}
