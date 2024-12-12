package ports

import (
	"fmt"
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/service"
)

func (h *Handler) CreatePostAction(w http.ResponseWriter, r *http.Request) {
	var (
		postID = r.URL.Query().Get("id")
		action = r.FormValue("action")
		ctx    = r.Context()
	)

	user := ctx.Value(myKey).(User)

	if err := h.svc.PostReaction(ctx, service.PostReactionInput{
		PostID: postID,
		UserID: user.ID,
		Action: action,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user.IsAuthN)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%s", postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}
