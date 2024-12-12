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
		user        = ctx.Value(myKey).(User)
		categoryIDs = r.URL.Query()["categories"]
	)

	myLikedPosts, err := h.svc.GetMyLikedPosts(
		ctx,
		service.GetPostsListInput{
			UserID:      user.ID,
			CategoryIDs: categoryIDs,
		},
	)
	if err != nil {
		e3r.ErrorEncoder(err, w, user.IsAuthN)
		return
	}

	httphelper.Render(w, http.StatusOK, "home", httphelper.GetTmplData(myLikedPosts, user.IsAuthN))
}
