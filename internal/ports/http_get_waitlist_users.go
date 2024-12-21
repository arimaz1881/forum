package ports

import (
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
)

func (h *Handler) GetWaitlistUsers(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = getUserData(ctx)
	)

	usersWaitlist, err := h.svc.GetWaitlistUsers(ctx, user.ID)
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	httphelper.Render(w, http.StatusOK, "users-waitlist", httphelper.GetTmplData(usersWaitlist, user))
}
