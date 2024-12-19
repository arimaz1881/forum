package ports

import (
	"fmt"
	"net/http"

	"forum/internal/domain"
	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
	"forum/internal/service"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		imgMaxSize int64 = 32 << 20
		ctx              = r.Context()
		user             = getUserData(ctx)
	)

	err := r.ParseMultipartForm(imgMaxSize)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}

	file, err := httphelper.FileFromForm(r, "img")
	if err != nil {
		e3r.ErrorEncoder(domain.ErrInvalidFile, w, user)
		return
	}
	if file != nil {
		defer file.FileReader.Close()
	}

	id, err := h.svc.CreatePost(ctx, service.CreatePostInput{
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		Categories: r.Form["catigoria"],
		File:       file,
		UserID:     user.ID,
	})
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%d", id)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = getUserData(ctx)
	)

	catigories, err := h.svc.GetCatigories(r.Context())
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	httphelper.Render(w, http.StatusOK, "create", httphelper.GetTmplData(catigories, user))
}


func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	var (
		postID = r.URL.Query().Get("post_id")
		ctx    = r.Context()
		user   = getUserData(ctx)
	)


	err := h.svc.DeletePost(ctx, service.DeletePostInput{
		PostID: postID,
		UserID: user.ID})
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	var (
		postID = r.URL.Query().Get("post_id")
		ctx    = r.Context()
		user   = getUserData(ctx)
	)

	if err := h.svc.EditPost(ctx, service.EditPostInput{
		PostID: 	postID,
		Content:    r.FormValue("edit-post-content"),
		UserID:     user.ID,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%s", postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}