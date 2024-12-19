package ports

import (
	"fmt"
	"net/http"

	"forum/internal/pkg/e3r"
	"forum/internal/service"
)

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		postID = r.URL.Query().Get("id")
		user   = getUserData(ctx)
	)

	if err := h.svc.CreateComment(ctx, service.CreateCommentInput{
		Content: r.FormValue("content"),
		UserID:  user.ID,
		PostID:  postID,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%s", postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {

	var (
		commentID = r.URL.Query().Get("comment_id")
		postID = r.URL.Query().Get("post_id")
		ctx    = r.Context()
		user   = getUserData(ctx)
	)

	err := h.svc.DeleteComment(ctx, service.DeleteCommentInput{
		CommentID: commentID,
		UserID: user.ID})
	if err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%s", postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) EditCommennt(w http.ResponseWriter, r *http.Request) {
	var (
		commentID = r.URL.Query().Get("comment_id")
		postID = r.URL.Query().Get("post_id")
		ctx    = r.Context()
		user   = getUserData(ctx)
	)

	if err := h.svc.EditComment(ctx, service.EditCommentInput{
		CommentID: 	commentID,
		Content:    r.FormValue("edit-comment-content"),
		UserID:     user.ID,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	url := fmt.Sprintf("/posts/view?id=%s", postID)

	http.Redirect(w, r, url, http.StatusSeeOther)
}
