package ports

import (
	"net/http"
	"strconv"

	"forum/internal/pkg/e3r"
	"forum/internal/service"
)

func (h *Handler) UpgradeRoleReject(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           = r.Context()
		user          = getUserData(ctx)
		waitingUserID = r.URL.Query().Get("waiting_user_id")
		userID        = strconv.Itoa(int(user.ID))
	)

	if err := h.svc.UpgradeRoleReject(ctx, service.UpgradeRoleRejectInput{
		WaitingUserID: waitingUserID,
		UserID:        userID,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
