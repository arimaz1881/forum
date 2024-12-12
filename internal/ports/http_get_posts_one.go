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
		userID int64
	)

	user, ok := ctx.Value(myKey).(User)
	if ok {
		userID = user.ID
	}

	post, err := h.svc.GetPostsOne(ctx, service.GetPostOneInput{
		PostID: postID,
		UserID: userID,
	},
	)
	if err != nil {
		e3r.ErrorEncoder(err, w, user.IsAuthN)
		return
	}

	httphelper.Render(w, http.StatusOK, "view", httphelper.GetTmplData(post, user.IsAuthN))
}
