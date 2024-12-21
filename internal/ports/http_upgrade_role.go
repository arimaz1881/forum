package ports

import (
	"forum/internal/pkg/e3r"
	"forum/internal/service"
	"net/http"
	"strconv"
)

func (h *Handler) UpgradeRoleApprove(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           = r.Context()
		user          = getUserData(ctx)
		waitingUserID = r.URL.Query().Get("waiting_user_id")
		userID        = strconv.Itoa(int(user.ID))
	)

	if err := h.svc.UpgradeRoleApprove(ctx, service.UpgradeRoleInput{
		WaitingUserID: waitingUserID,
		UserID:        userID,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}
	http.Redirect(w, r, "/users/roles/moderator-waitlist", http.StatusSeeOther)
}

func (h *Handler) UpgradeRoleReject(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           = r.Context()
		user          = getUserData(ctx)
		waitingUserID = r.URL.Query().Get("waiting_user_id")
		userID        = strconv.Itoa(int(user.ID))
	)

	if err := h.svc.UpgradeRoleReject(ctx, service.UpgradeRoleInput{
		WaitingUserID: waitingUserID,
		UserID:        userID,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}
	http.Redirect(w, r, "/users/roles/moderator-waitlist", http.StatusSeeOther)
}
