package ports

import (
	"fmt"
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/service"
)

func (h *Handler) CreateCommentReaction(w http.ResponseWriter, r *http.Request) {
	var (
		commentID = r.URL.Query().Get("comment_id")
		postID    = r.URL.Query().Get("post_id")
		action    = r.FormValue("action")
		ctx       = r.Context()
		user      = getUserData(ctx)
	)

	if err := h.svc.CommentReaction(ctx, service.CommentReactionInput{
		PostID:    postID,
		CommentID: commentID,
		UserID:    user.ID,
		Action:    action,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%s", postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}
