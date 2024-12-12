package ports

import (
	"fmt"
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
	"forum/internal/service"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = ctx.Value(myKey).(User)
	)

	r.ParseForm()

	id, err := h.svc.CreatePost(ctx, service.CreatePostInput{
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		Categories: r.Form["catigoria"],
		UserID:     user.ID,
	})
	if err != nil {
		e3r.ErrorEncoder(err, w, user.IsAuthN)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%d", id)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	catigories, err := h.svc.GetCatigories(r.Context())
	if err != nil {
		e3r.ErrorEncoder(err, w, true)
		return
	}

	httphelper.Render(w, http.StatusOK, "create", httphelper.GetTmplData(catigories, true))
}
