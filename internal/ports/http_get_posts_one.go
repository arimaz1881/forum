package ports

import (
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
	"forum/internal/service"
)

func (h *Handler) GetPostsOne(w http.ResponseWriter, r *http.Request) {
	var (
		postID = r.URL.Query().Get("id")
		ctx    = r.Context()
		user   = getUserData(ctx)
	)

	post, err := h.svc.GetPostsOne(ctx, service.GetPostOneInput{
		PostID: postID,
		UserID: user.ID,
	},
	)
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}
	

	httphelper.Render(w, http.StatusOK, "view", httphelper.GetTmplData(post, user))
}
