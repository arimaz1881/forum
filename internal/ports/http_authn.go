package ports

import (
	"forum/internal/domain"
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
		e3r.ErrorEncoder(err, w, httphelper.User{
			IsAuthN: false,
			Role:    domain.RoleGuest,
		})
		return
	}

	sessions.Set(w, session.Token, *session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) SignUpPage(w http.ResponseWriter, r *http.Request) {
	httphelper.Render(w, http.StatusOK, "sign-up", httphelper.GetTmplData(nil, httphelper.User{
		IsAuthN: false,
		Role:    domain.RoleGuest,
	}))
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	session, err := h.svc.SignIn(r.Context(), service.SignInInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		e3r.ErrorEncoder(err, w, httphelper.User{
			IsAuthN: false,
			Role:    domain.RoleGuest,
		})
		return
	}

	sessions.Set(w, session.Token, *session.ExpiresAt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) SignInPage(w http.ResponseWriter, r *http.Request) {
	httphelper.Render(w, http.StatusOK, "sign-in", httphelper.GetTmplData(nil, httphelper.User{
		IsAuthN: false,
		Role:    domain.RoleGuest,
	}))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/authn/sign-in", http.StatusSeeOther)
		return
	}

	var (
		ctx  = r.Context()
		user = getUserData(ctx)
	)

	if err := h.svc.Logout(ctx, service.LogOutInput{
		Token: cookie.Value,
	}); err != nil {
		e3r.ErrorEncoder(err, w, user)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
