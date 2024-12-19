package ports

import (
	"net/http"
	"strconv"

	"forum/internal/pkg/e3r"
)

func (h *Handler) SubmitRoleUpgrade(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		user   = getUserData(ctx)
		userID = strconv.Itoa(int(user.ID))
	)

	if err := h.svc.SubmitRoleUpgrade(ctx, userID); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
