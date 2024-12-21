package ports

import (
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
)

func (h *Handler) GetModerators(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = getUserData(ctx)
	)

	moderatorslist, err := h.svc.GetModerators(ctx, user.ID)
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	httphelper.Render(w, http.StatusOK, "moderators-list", httphelper.GetTmplData(moderatorslist, user))
}
