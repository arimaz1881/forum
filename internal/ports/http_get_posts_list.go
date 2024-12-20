package ports

import (
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
)

func (h *Handler) GetPostsList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		categoryIDs = r.URL.Query()["categories"]
		user        = getUserData(ctx)
	)

	if r.URL.Path != "/" {
		e3r.ErrorEncoder(e3r.NotFound("not found"), w, user)
		return
	}

	postsList, err := h.svc.GetPostsList(
		ctx,
		categoryIDs,
	)
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}
	httphelper.Render(w, http.StatusOK, "home", httphelper.GetTmplData(postsList, user))
}
