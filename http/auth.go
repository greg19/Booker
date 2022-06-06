package http

import (
	"booker/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (s *server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if !verifyForm(r, "username", "password") {
		renderError(w, r, http.StatusBadRequest)
		return
	}

	u, err := models.GetUserByUsername(s.db, r.Form.Get("username"))

	if err != nil {
		renderError(w, r, http.StatusInternalServerError)
		log.Println(err)
		return
	} else if u.Password != r.Form.Get("password") {
		addError(w, r, http.StatusBadRequest, "invalid username or password")
		renderTemplate(w, r, "login.html", nil)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(12 * time.Hour)

	if err = models.CreateSession(s.db, sessionToken, u.Id, expiresAt); err != nil {
		renderError(w, r, http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			renderError(w, r, http.StatusUnauthorized)
		} else {
			renderError(w, r, http.StatusBadRequest)
		}
		return
	}

	sessionToken := c.Value
	models.DeleteSession(s.db, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

// middleware to read user from session
func (s *server) readUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		c, err := r.Cookie("session_token")
		if err != nil {
			user = nil
		} else {
			sessionToken := c.Value
			session, err := models.GetSessionByToken(s.db, sessionToken)
			if err != nil {
				renderError(w, r, http.StatusInternalServerError)
			} else if session.IsExpired() {
				user = nil
				models.DeleteSession(s.db, sessionToken)
			} else {
				user, err = models.GetUserById(s.db, session.UserId)
			}
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
