package ports

import (
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
	"forum/internal/service"
)

func (h *Handler) GetMyComments(w http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		user        = getUserData(ctx)
		categoryIDs = r.URL.Query()["categories"]
	)

	myCreatedPosts, err := h.svc.GetMyCommentsList(
		ctx,
		service.GetPostsListInput{
			UserID:      user.ID,
			CategoryIDs: categoryIDs,
		},
	)
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	httphelper.Render(w, http.StatusOK, "comments", httphelper.GetTmplData(myCreatedPosts, user))
}
