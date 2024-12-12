package ports

import (
	"fmt"
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/service"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		postID = r.URL.Query().Get("id")
		user   = ctx.Value(myKey).(User)
	)

	if err := h.svc.CreateComment(ctx, service.CreateCommentInput{
		Content: r.FormValue("content"),
		UserID:  user.ID,
		PostID:  postID,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user.IsAuthN)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%s", postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}
