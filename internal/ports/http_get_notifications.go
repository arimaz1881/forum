package ports

import (
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
	"forum/internal/service"

)

func (h *Handler) GetNotificationsList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		user        = getUserData(ctx)
	)

	notificationsList, err := h.svc.GetNotificationsList(
		ctx,
		user.ID,
	)
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	httphelper.Render(w, http.StatusOK, "notification", httphelper.GetTmplData(notificationsList, user))
}


func (h *Handler) NotificationLook(w http.ResponseWriter, r *http.Request) {
	var (
        ctx         = r.Context()
        user        = getUserData(ctx)
        notificationID = r.URL.Query().Get("notification_id")
    )

    err := h.svc.LookNotification(
        ctx,
        service.LookNotificationInput{
		UserID:	user.ID,
        NotificationID: notificationID,
	})
    if err != nil {
        e3r.ErrorEncoder(err, w, user)
        return
    }

    http.Redirect(w, r, "/notifications", http.StatusSeeOther)
}