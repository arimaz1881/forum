package sessions

import (
	"net/http"
	"time"
)

func Set(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiresAt,
		MaxAge:   3600,
		Path:     "/",
		HttpOnly: true,
	})
}

func Close(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	sessionToken := c.Value

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().UTC(),
		HttpOnly: false,
		MaxAge:   -1,
	})

	return sessionToken, nil
}
