package ports

import (
	"forum/internal/pkg/e3r"
	"forum/internal/pkg/httphelper"
	"forum/internal/pkg/sessions"
	"forum/internal/service"
	"net/http"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	session, err := h.svc.SignUp(r.Context(), service.SignUpInput{
		Email:    r.FormValue("email"),
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		e3r.ErrorEncoder(err, w, false)
		return
	}

	sessions.Set(w, session.Token, *session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) SignUpPage(w http.ResponseWriter, r *http.Request) {
	httphelper.Render(w, http.StatusOK, "sign-up", httphelper.GetTmplData(nil, false))
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	session, err := h.svc.SignIn(r.Context(), service.SignInInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		e3r.ErrorEncoder(err, w, false)
		return
	}

	sessions.Set(w, session.Token, *session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) SignInPage(w http.ResponseWriter, r *http.Request) {
	httphelper.Render(w, http.StatusOK, "sign-in", httphelper.GetTmplData(nil, false))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/authn/sign-in", http.StatusSeeOther)
		return
	}

	if err := h.svc.Logout(r.Context(), service.LogOutInput{
		Token: cookie.Value,
	}); err != nil {
		e3r.ErrorEncoder(err, w, true)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
